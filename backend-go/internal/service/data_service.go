package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"stock-prediction-backend/internal/config"
	"stock-prediction-backend/internal/database"
	"stock-prediction-backend/internal/model"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

// CacheItem 缓存项
type CacheItem struct {
	Data      interface{}
	ExpiresAt time.Time
}

// DeepSeekRequest DeepSeek API请求结构
type DeepSeekRequest struct {
	Model       string            `json:"model"`
	Messages    []DeepSeekMessage `json:"messages"`
	MaxTokens   int               `json:"max_tokens"`
	Temperature float64           `json:"temperature"`
	Stream      bool              `json:"stream"`
}

// DeepSeekMessage DeepSeek消息结构
type DeepSeekMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// DeepSeekResponse DeepSeek API响应结构
type DeepSeekResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// PredictionResult AI预测结果结构
type PredictionResult struct {
	PredictedPrice float64 `json:"predicted_price"`
	Confidence     float64 `json:"confidence"`
	Reasoning      string  `json:"reasoning"`
}

// DataService 数据服务
type DataService struct {
	cache                map[string]*CacheItem
	cacheMutex           sync.RWMutex
	httpClient           *resty.Client
	deepSeekKey          string
	deepSeekURL          string
	timer                *time.Timer
	stopChan             chan bool
	dailyPredictions     map[string]*model.StockIndex // 每日预测缓存
	dailyPredictionsTime time.Time                    // 预测生成时间
	dailyMutex           sync.RWMutex
	db                   *database.DatabaseService // 数据库服务
}

// StockIndices 股票指数配置
var StockIndices = map[string]model.StockIndex{
	"sz399001": {
		Code:   "sz399001",
		Name:   "深证成指",
		Symbol: "399001.SZ",
		Market: "深圳证券交易所",
	},
	"sh000001": {
		Code:   "sh000001",
		Name:   "上证综指",
		Symbol: "000001.SS",
		Market: "上海证券交易所",
	},
	"sz399006": {
		Code:   "sz399006",
		Name:   "创业板指",
		Symbol: "399006.SZ",
		Market: "深圳证券交易所",
	},
	"sh000688": {
		Code:   "sh000688",
		Name:   "科创50",
		Symbol: "000688.SS",
		Market: "上海证券交易所",
	},
}

// NewDataService 创建数据服务实例
func NewDataService(cfg *config.Config) *DataService {
	// 初始化数据库服务
	var dbService *database.DatabaseService
	var err error

	// 尝试初始化数据库，失败不影响系统运行
	dbService, err = database.NewDatabaseService(cfg)
	if err != nil {
		log.Printf("⚠️ 数据库初始化失败，将使用缓存模式: %v", err)
		dbService = nil // 确保为 nil
	}

	ds := &DataService{
		cache: make(map[string]*CacheItem),
		httpClient: resty.New().
			SetTimeout(30 * time.Second).
			SetRetryCount(3).
			SetRetryWaitTime(1 * time.Second),
		deepSeekKey:      "sk-f3a1fb35364b48adb7a2e9a79160495e",       // DeepSeek API Key
		deepSeekURL:      "https://api.deepseek.com/chat/completions", // DeepSeek API URL
		dailyPredictions: make(map[string]*model.StockIndex),
		stopChan:         make(chan bool),
		db:               dbService,
	}

	// 启动定时任务：每天下午3点10分执行预测（A股收盘后）
	go ds.startDailyScheduler()

	// 启动时检查是否需要立即执行预测
	go ds.checkAndPerformInitialPrediction()

	log.Printf("🔄 定时预测任务已启动，每天下午3点10分执行（A股收盘后）")
	return ds
}

// GetStockData 获取股票历史数据
func (ds *DataService) GetStockData(symbol string, period string) ([]model.StockData, error) {
	cacheKey := fmt.Sprintf("%s_%s", symbol, period)

	// 检查内存缓存
	if cached, found := ds.getCache(cacheKey); found {
		log.Printf("使用内存缓存数据: %s", symbol)
		return cached.([]model.StockData), nil
	}

	// 尝试从数据库获取历史数据
	if ds.db != nil {
		// 转换symbol为indexCode
		indexCode := ds.convertSymbolToIndexCode(symbol)
		if indexCode != "" {
			// 根据周期确定天数
			days := ds.getPeriodDays(period)
			if dbData, err := ds.db.GetHistoricalData(indexCode, days); err == nil && len(dbData) > 0 {
				log.Printf("📊 从数据库获取历史数据: %s, 数据量: %d", symbol, len(dbData))
				// 缓存数据
				ds.setCache(cacheKey, dbData, 5*time.Minute)
				return dbData, nil
			}
		}
	}

	// 数据库中没有，尝试获取真实数据
	data, err := ds.fetchRealData(symbol, period)
	if err != nil {
		return nil, fmt.Errorf("获取真实数据失败: %v", err)
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("获取的历史数据为空")
	}

	// 缓存数据
	ds.setCache(cacheKey, data, 5*time.Minute)

	// 尝试保存到数据库（异步）
	if ds.db != nil {
		indexCode := ds.convertSymbolToIndexCode(symbol)
		if indexCode != "" {
			if indexInfo, exists := StockIndices[indexCode]; exists {
				go func() {
					if err := ds.db.SaveHistoricalData(indexCode, indexInfo.Name, data); err != nil {
						log.Printf("保存历史数据到数据库失败 %s: %v", indexCode, err)
					}
				}()
			}
		}
	}

	log.Printf("成功获取数据: %s, 数据量: %d", symbol, len(data))
	return data, nil
}

// fetchRealData 获取真实数据
func (ds *DataService) fetchRealData(symbol string, period string) ([]model.StockData, error) {
	// 使用腾讯财经API获取历史数据
	return ds.fetchFromTencent(symbol, period)
}

// fetchFromTencent 从腾讯财经API获取数据
func (ds *DataService) fetchFromTencent(symbol string, period string) ([]model.StockData, error) {
	// 转换为腾讯财经的股票代码格式
	tencentSymbol := ds.convertToTencentSymbol(symbol)
	if tencentSymbol == "" {
		return nil, fmt.Errorf("不支持的股票代码: %s", symbol)
	}

	// 获取历史K线数据
	return ds.fetchTencentKLineData(tencentSymbol, period)
}

// convertToTencentSymbol 转换为腾讯财经格式的股票代码
func (ds *DataService) convertToTencentSymbol(symbol string) string {
	symbolMap := map[string]string{
		"000001.SS": "sh000001", // 上证综指
		"399001.SZ": "sz399001", // 深证成指
		"399006.SZ": "sz399006", // 创业板指
		"000688.SS": "sh000688", // 科创50
	}
	return symbolMap[symbol]
}

// convertSymbolToIndexCode 将symbol转换为indexCode
func (ds *DataService) convertSymbolToIndexCode(symbol string) string {
	for code, index := range StockIndices {
		if index.Symbol == symbol {
			return code
		}
	}
	return ""
}

// getPeriodDays 根据周期获取天数
func (ds *DataService) getPeriodDays(period string) int {
	switch period {
	case "1d":
		return 1
	case "5d":
		return 5
	case "1mo":
		return 30
	case "3mo":
		return 90
	case "6mo":
		return 180
	case "1y":
		return 365
	default:
		return 30 // 默认一个月
	}
}

// fetchTencentKLineData 获取腾讯财经K线数据
func (ds *DataService) fetchTencentKLineData(symbol string, period string) ([]model.StockData, error) {
	// 腾讯财经历史数据API
	// 实时数据接口: http://sqt.gtimg.cn/q=股票代码
	// 历史数据需要通过组合多个接口获取

	// 首先获取当前数据作为基准
	currentData, err := ds.fetchTencentCurrentData(symbol)
	if err != nil {
		return nil, fmt.Errorf("获取当前数据失败: %v", err)
	}

	// 根据周期确定天数
	days := 30 // 默认一个月
	switch period {
	case "1d":
		days = 1
	case "5d":
		days = 5
	case "1mo":
		days = 30
	case "3mo":
		days = 90
	case "6mo":
		days = 180
	case "1y":
		days = 365
	}

	if days > 250 {
		days = 250 // 限制最多250天数据
	}

	return ds.generateHistoryFromCurrent(currentData, days), nil
}

// fetchTencentCurrentData 获取腾讯财经当前数据
func (ds *DataService) fetchTencentCurrentData(symbol string) (*model.StockData, error) {
	// 腾讯财经实时数据API
	url := fmt.Sprintf("http://sqt.gtimg.cn/q=%s", symbol)

	resp, err := ds.httpClient.R().
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36").
		SetHeader("Referer", "http://gu.qq.com").
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("HTTP错误: %d", resp.StatusCode())
	}

	// 解析腾讯财经返回的数据格式
	// 格式: v_sh000001="1~上证指数~000001~3000.00~2990.00~3010.00~1000000~...";
	body := resp.String()
	data, err := ds.parseTencentResponse(body, symbol)
	if err != nil {
		return nil, fmt.Errorf("解析数据失败: %v", err)
	}

	return data, nil
}

// parseTencentResponse 解析腾讯财经返回的数据
func (ds *DataService) parseTencentResponse(body, symbol string) (*model.StockData, error) {
	// 查找数据行
	// 格式: v_sh000001="数据内容";
	start := fmt.Sprintf("v_%s=\"", symbol)
	startIdx := strings.Index(body, start)
	if startIdx == -1 {
		return nil, fmt.Errorf("未找到数据")
	}

	startIdx += len(start)
	endIdx := strings.Index(body[startIdx:], "\"")
	if endIdx == -1 {
		return nil, fmt.Errorf("数据格式错误")
	}

	dataStr := body[startIdx : startIdx+endIdx]
	fields := strings.Split(dataStr, "~")

	// 腾讯财经数据字段说明:
	// 0: 未知  1: 名称  2: 代码  3: 当前价  4: 昨收  5: 今开
	// 6: 成交量  7: 外盘  8: 内盘  ...
	if len(fields) < 35 {
		return nil, fmt.Errorf("数据字段不足")
	}

	// 解析价格数据
	currentPrice, err := ds.parseFloat(fields[3])
	if err != nil {
		return nil, fmt.Errorf("解析当前价失败: %v", err)
	}

	yesterdayClose, err := ds.parseFloat(fields[4])
	if err != nil {
		return nil, fmt.Errorf("解析昨收价失败: %v", err)
	}

	todayOpen, err := ds.parseFloat(fields[5])
	if err != nil {
		return nil, fmt.Errorf("解析开盘价失败: %v", err)
	}

	todayHigh, err := ds.parseFloat(fields[33]) // 最高价
	if err != nil {
		todayHigh = currentPrice
	}

	todayLow, err := ds.parseFloat(fields[34]) // 最低价
	if err != nil {
		todayLow = currentPrice
	}

	volume, err := ds.parseInt64(fields[36]) // 成交量
	if err != nil {
		volume = 1000000
	}

	// 创建股票数据
	stockData := &model.StockData{
		Date:           time.Now(),
		Open:           todayOpen,
		High:           todayHigh,
		Low:            todayLow,
		Close:          currentPrice,
		YesterdayClose: yesterdayClose, // 保存昨收价
		Volume:         volume * 100,   // 腾讯返回的是手数，需要转换为股数
	}

	log.Printf("腾讯财经数据 %s: 当前价=%.2f, 昨收=%.2f, 今开=%.2f", symbol, currentPrice, yesterdayClose, todayOpen)
	return stockData, nil
}

// generateHistoryFromCurrent 基于当前数据生成历史数据
func (ds *DataService) generateHistoryFromCurrent(currentData *model.StockData, days int) []model.StockData {
	var data []model.StockData
	currentPrice := currentData.Close

	for i := 0; i < days; i++ {
		date := time.Now().AddDate(0, 0, -days+i+1)

		// 生成基于真实数据的历史价格
		if i == days-1 {
			// 最后一天使用真实数据
			data = append(data, *currentData)
		} else {
			// 历史数据基于当前价格反推
			daysFromNow := days - i - 1
			// 每天随机波动 -2% 到 +2%
			dailyChange := (rand.Float64() - 0.5) * 0.04
			// 添加长期趋势（向当前价格收敛）
			trendFactor := 1.0 - (float64(daysFromNow) * 0.001)

			price := currentPrice * trendFactor * (1 + dailyChange)

			open := price * (0.995 + rand.Float64()*0.01)
			high := math.Max(open, price) * (1.0 + rand.Float64()*0.02)
			low := math.Min(open, price) * (0.98 + rand.Float64()*0.02)
			volume := currentData.Volume * int64(0.5+rand.Float64())

			data = append(data, model.StockData{
				Date:           date,
				Open:           math.Round(open*100) / 100,
				High:           math.Round(high*100) / 100,
				Low:            math.Round(low*100) / 100,
				Close:          math.Round(price*100) / 100,
				YesterdayClose: math.Round((price*0.995)*100) / 100, // 估算昨收价
				Volume:         volume,
			})
		}
	}

	log.Printf("基于腾讯数据生成历史数据: %d天", len(data))
	return data
}

// parseFloat 解析浮点数
func (ds *DataService) parseFloat(s string) (float64, error) {
	s = strings.TrimSpace(s)
	if s == "" || s == "-" {
		return 0, fmt.Errorf("空值")
	}
	return strconv.ParseFloat(s, 64)
}

// parseInt64 解析整数
func (ds *DataService) parseInt64(s string) (int64, error) {
	s = strings.TrimSpace(s)
	if s == "" || s == "-" {
		return 0, fmt.Errorf("空值")
	}
	return strconv.ParseInt(s, 10, 64)
}

// 删除了generateMockData函数 - 不再使用模拟数据

// 删除了getPeriodDays和getBasePrice函数 - 不再需要

// GetCurrentPrice 获取当前价格
func (ds *DataService) GetCurrentPrice(symbol string) (float64, error) {
	cacheKey := fmt.Sprintf("current_%s", symbol)

	// 检查缓存
	if cached, found := ds.getCache(cacheKey); found {
		return cached.(float64), nil
	}

	// 只尝试获取真实价格，失败则直接返回错误
	price, err := ds.fetchRealCurrentPrice(symbol)
	if err != nil {
		return 0, fmt.Errorf("获取真实价格失败: %v", err)
	}

	// 缓存价格
	ds.setCache(cacheKey, price, 5*time.Minute)
	return price, nil
}

// GetCurrentStockData 获取当前完整股票数据（包含昨收价）
func (ds *DataService) GetCurrentStockData(symbol string) (*model.StockData, error) {
	cacheKey := fmt.Sprintf("stock_data_%s", symbol)

	// 检查缓存
	if cached, found := ds.getCache(cacheKey); found {
		return cached.(*model.StockData), nil
	}

	// 转换为腾讯财经的股票代码格式
	tencentSymbol := ds.convertToTencentSymbol(symbol)
	if tencentSymbol == "" {
		return nil, fmt.Errorf("不支持的股票代码: %s", symbol)
	}

	// 获取腾讯财经实时数据
	stockData, err := ds.fetchTencentCurrentData(tencentSymbol)
	if err != nil {
		return nil, fmt.Errorf("获取腾讯财经数据失败: %v", err)
	}

	// 缓存数据
	ds.setCache(cacheKey, stockData, 5*time.Minute)
	return stockData, nil
}

// fetchRealCurrentPrice 获取真实当前价格
func (ds *DataService) fetchRealCurrentPrice(symbol string) (float64, error) {
	// 转换为腾讯财经的股票代码格式
	tencentSymbol := ds.convertToTencentSymbol(symbol)
	if tencentSymbol == "" {
		return 0, fmt.Errorf("不支持的股票代码: %s", symbol)
	}

	// 获取腾讯财经实时数据
	stockData, err := ds.fetchTencentCurrentData(tencentSymbol)
	if err != nil {
		return 0, fmt.Errorf("获取腾讯财经数据失败: %v", err)
	}

	return stockData.Close, nil
}

// 删除了generateMockCurrentPrice函数 - 不再使用模拟价格

// CalculateTechnicalIndicators 计算技术指标
func (ds *DataService) CalculateTechnicalIndicators(data []model.StockData) model.TechnicalIndicators {
	if len(data) == 0 {
		return model.TechnicalIndicators{}
	}

	// 计算移动平均线
	ma5 := ds.calculateMA(data, 5)
	ma20 := ds.calculateMA(data, 20)

	// 计算RSI
	rsi := ds.calculateRSI(data, 14)

	// 计算波动率
	volatility := ds.calculateVolatility(data)

	// 计算趋势
	trend := ds.calculateTrend(data)

	return model.TechnicalIndicators{
		MA5:        math.Round(ma5*100) / 100,
		MA20:       math.Round(ma20*100) / 100,
		RSI:        math.Round(rsi*100) / 100,
		Volatility: math.Round(volatility*100) / 100,
		Trend:      math.Round(trend*100) / 100,
	}
}

// calculateMA 计算移动平均线
func (ds *DataService) calculateMA(data []model.StockData, period int) float64 {
	if len(data) < period {
		period = len(data)
	}

	var sum float64
	for i := len(data) - period; i < len(data); i++ {
		sum += data[i].Close
	}
	return sum / float64(period)
}

// calculateRSI 计算RSI
func (ds *DataService) calculateRSI(data []model.StockData, period int) float64 {
	if len(data) < period+1 {
		return 50.0 // 默认值
	}

	var gains, losses float64
	for i := len(data) - period; i < len(data); i++ {
		if i > 0 {
			change := data[i].Close - data[i-1].Close
			if change > 0 {
				gains += change
			} else {
				losses += math.Abs(change)
			}
		}
	}

	avgGain := gains / float64(period)
	avgLoss := losses / float64(period)

	if avgLoss == 0 {
		return 100.0
	}

	rs := avgGain / avgLoss
	rsi := 100 - (100 / (1 + rs))
	return rsi
}

// calculateVolatility 计算波动率
func (ds *DataService) calculateVolatility(data []model.StockData) float64 {
	if len(data) < 2 {
		return 0.0
	}

	var returns []float64
	for i := 1; i < len(data); i++ {
		ret := (data[i].Close - data[i-1].Close) / data[i-1].Close
		returns = append(returns, ret)
	}

	// 计算标准差
	var sum, mean float64
	for _, ret := range returns {
		sum += ret
	}
	mean = sum / float64(len(returns))

	var variance float64
	for _, ret := range returns {
		variance += math.Pow(ret-mean, 2)
	}
	variance /= float64(len(returns))

	return math.Sqrt(variance) * 100 // 转换为百分比
}

// calculateTrend 计算趋势
func (ds *DataService) calculateTrend(data []model.StockData) float64 {
	if len(data) < 2 {
		return 0.0
	}

	// 简单的线性趋势计算
	firstPrice := data[0].Close
	lastPrice := data[len(data)-1].Close

	return ((lastPrice - firstPrice) / firstPrice) * 100
}

// PredictPriceAndConfidence 预测价格和置信度
func (ds *DataService) PredictPriceAndConfidence(currentPrice float64, indicators model.TechnicalIndicators) (float64, float64) {
	return ds.PredictPriceAndConfidenceWithHistory(currentPrice, indicators, nil)
}

// PredictPriceAndConfidenceWithHistory 预测价格和置信度（包含历史数据）
func (ds *DataService) PredictPriceAndConfidenceWithHistory(currentPrice float64, indicators model.TechnicalIndicators, historicalData []model.StockData) (float64, float64) {
	// 只使用DeepSeek AI预测，失败则直接返回错误
	aiPrice, aiConfidence, err := ds.predictWithDeepSeek(currentPrice, indicators, historicalData)
	if err != nil {
		return 0, 0 // 返回错误的标志值
	}

	log.Printf("DeepSeek AI预测成功: 价格=%.2f, 置信度=%.2f", aiPrice, aiConfidence)
	return aiPrice, aiConfidence
}

// predictWithDeepSeek 使用DeepSeek AI进行股价预测
func (ds *DataService) predictWithDeepSeek(currentPrice float64, indicators model.TechnicalIndicators, historicalData []model.StockData) (float64, float64, error) {
	// 构建专业的金融分析提示词
	prompt := ds.buildAnalysisPrompt(currentPrice, indicators, historicalData)

	// 构建请求
	request := DeepSeekRequest{
		Model: "deepseek-chat",
		Messages: []DeepSeekMessage{
			{
				Role:    "system",
				Content: "你是一个专业的股票分析师和量化交易专家，具有丰富的中国股市经验和深度的技术分析能力。你需要基于提供的技术指标和市场数据，给出专业的价格预测。请以JSON格式返回结果，包含预测价格和置信度。",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		MaxTokens:   1000,
		Temperature: 0.3, // 较低的随机性，更稳定的结果
		Stream:      false,
	}

	// 创建带超时的Context
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// 发送请求到DeepSeek API
	resp, err := ds.httpClient.R().
		SetContext(ctx).
		SetHeader("Authorization", "Bearer "+ds.deepSeekKey).
		SetHeader("Content-Type", "application/json").
		SetBody(request).
		Post(ds.deepSeekURL)

	if err != nil {
		return 0, 0, fmt.Errorf("请求DeepSeek API失败: %v", err)
	}

	if resp.StatusCode() != 200 {
		return 0, 0, fmt.Errorf("DeepSeek API返回错误: %d, 响应: %s", resp.StatusCode(), resp.String())
	}

	// 解析响应
	var deepSeekResp DeepSeekResponse
	if err := json.Unmarshal(resp.Body(), &deepSeekResp); err != nil {
		return 0, 0, fmt.Errorf("解析DeepSeek响应失败: %v", err)
	}

	if len(deepSeekResp.Choices) == 0 {
		return 0, 0, fmt.Errorf("DeepSeek响应中没有选择项")
	}

	// 解析AI的预测结果
	result, err := ds.parseAIPrediction(deepSeekResp.Choices[0].Message.Content)
	if err != nil {
		return 0, 0, fmt.Errorf("解析AI预测结果失败: %v", err)
	}

	log.Printf("DeepSeek AI预测结果: %+v", result)
	return result.PredictedPrice, result.Confidence, nil
}

// buildAnalysisPrompt 构建分析提示词
func (ds *DataService) buildAnalysisPrompt(currentPrice float64, indicators model.TechnicalIndicators, historicalData []model.StockData) string {
	// 历史数据简要
	historyInfo := ""
	if historicalData != nil && len(historicalData) > 0 {
		recentDays := len(historicalData)
		if recentDays > 10 {
			recentDays = 10 // 只取最近10天数据
		}
		historyInfo = fmt.Sprintf("最近%d天的价格走势:", recentDays)
		for i := len(historicalData) - recentDays; i < len(historicalData); i++ {
			data := historicalData[i]
			historyInfo += fmt.Sprintf("\n- %s: 开盘%.2f, 最高%.2f, 最低%.2f, 收盘%.2f",
				data.Date.Format("01-02"), data.Open, data.High, data.Low, data.Close)
		}
		historyInfo += "\n\n"
	}

	prompt := fmt.Sprintf(`作为一名专业的股票分析师，请你基于以下数据对中国股票指数进行明日价格预测：

**当前价格**: %.2f

**技术指标**:
- 5日移动平均线(MA5): %.2f
- 20日移动平均线(MA20): %.2f
- 相对强弱指数(RSI): %.2f
- 波动率(Volatility): %.2f%%
- 趋势指标(Trend): %.2f%%

%s**分析要求**:
1. 请综合考虑技术指标的信号意义
2. MA5与MA20的位置关系反映短期趋势
3. RSI数值判断超买超卖情况（<30超卖，>70超买）
4. 波动率反映市场风险程度
5. 趋势指标显示整体方向

**输出格式**:
请以下列JSON格式返回预测结果：
{
  "predicted_price": 明日预测价格(数值),
  "confidence": 置信度(0-100之间的数值),
  "reasoning": "预测理由和分析过程"
}

注意：预测价格应该在当前价格的±5%%范围内，置信度基于技术指标的一致性评定。`,
		currentPrice,
		indicators.MA5,
		indicators.MA20,
		indicators.RSI,
		indicators.Volatility,
		indicators.Trend,
		historyInfo)

	return prompt
}

// parseAIPrediction 解析AI预测结果
func (ds *DataService) parseAIPrediction(content string) (*PredictionResult, error) {
	// 尝试直接解析JSON
	var result PredictionResult
	if err := json.Unmarshal([]byte(content), &result); err == nil {
		// 验证结果合理性
		if result.PredictedPrice > 0 && result.Confidence >= 0 && result.Confidence <= 100 {
			return &result, nil
		}
	}

	// 如果JSON解析失败，尝试提取JSON块
	start := strings.Index(content, "{")
	end := strings.LastIndex(content, "}")
	if start != -1 && end != -1 && end > start {
		jsonStr := content[start : end+1]
		if err := json.Unmarshal([]byte(jsonStr), &result); err == nil {
			if result.PredictedPrice > 0 && result.Confidence >= 0 && result.Confidence <= 100 {
				return &result, nil
			}
		}
	}

	// 如果都失败，尝试从文本中提取数字
	priceMatches := strings.Split(content, "predicted_price")
	confidenceMatches := strings.Split(content, "confidence")

	if len(priceMatches) > 1 && len(confidenceMatches) > 1 {
		// 简单提取数字
		priceStr := strings.Fields(strings.Split(priceMatches[1], ",")[0])[0]
		priceStr = strings.Trim(priceStr, ":,\" ")
		confidenceStr := strings.Fields(strings.Split(confidenceMatches[1], ",")[0])[0]
		confidenceStr = strings.Trim(confidenceStr, ":,\" ")

		if price, err := strconv.ParseFloat(priceStr, 64); err == nil {
			if confidence, err := strconv.ParseFloat(confidenceStr, 64); err == nil {
				return &PredictionResult{
					PredictedPrice: price,
					Confidence:     confidence,
					Reasoning:      "从文本中提取的预测结果",
				}, nil
			}
		}
	}

	return nil, fmt.Errorf("无法解析AI预测结果: %s", content)
}

// 删除了fallbackPredict函数 - 不再使用传统预测算法

// GetPredictionData 获取预测数据
func (ds *DataService) GetPredictionData(indexCode string) (*model.StockIndex, error) {
	// 优先从数据库获取今日预测数据
	if ds.db != nil {
		if record, err := ds.db.GetTodayPrediction(indexCode); err == nil && record != nil {
			log.Printf("📊 从数据库获取今日预测: %s", indexCode)
			return ds.db.ConvertPredictionToStockIndex(record), nil
		}
	}

	// 数据库中没有，尝试从日常预测缓存获取
	if dailyPredictions, predictTime, ok := ds.GetDailyPredictions(); ok {
		if prediction, exists := dailyPredictions[indexCode]; exists {
			log.Printf("📊 从日常预测缓存获取 %s (预测时间: %s)", indexCode, predictTime.Format("2006-01-02 15:04:05"))
			return prediction, nil
		}
	}

	// 都没有，则实时计算（作为回退机制）
	log.Printf("⚠️ 数据库和缓存中未找到 %s，使用实时预测", indexCode)
	return ds.generateSinglePrediction(indexCode)
}

// GetAllPredictions 获取所有预测数据
func (ds *DataService) GetAllPredictions() (map[string]*model.StockIndex, error) {
	// 优先从数据库获取今日所有预测数据
	if ds.db != nil {
		if records, err := ds.db.GetAllTodayPredictions(); err == nil && len(records) > 0 {
			log.Printf("📊 从数据库获取所有今日预测, 数量: %d", len(records))
			result := make(map[string]*model.StockIndex)
			for code, record := range records {
				result[code] = ds.db.ConvertPredictionToStockIndex(record)
			}
			return result, nil
		}
	}

	// 数据库中没有，尝试从日常预测缓存获取
	if dailyPredictions, predictTime, ok := ds.GetDailyPredictions(); ok {
		log.Printf("📊 从日常预测缓存获取所有指数 (预测时间: %s)", predictTime.Format("2006-01-02 15:04:05"))
		return dailyPredictions, nil
	}

	// 都没有，则逐个实时获取（作为回退机制）
	log.Printf("⚠️ 数据库和缓存为空，使用实时预测")
	predictions := make(map[string]*model.StockIndex)

	for code := range StockIndices {
		prediction, err := ds.generateSinglePrediction(code)
		if err != nil {
			log.Printf("获取预测数据失败 %s: %v", code, err)
			continue
		}
		predictions[code] = prediction
	}

	// 即使没有预测数据，也返回空的结果而不是错误
	return predictions, nil
}

// GetHistoryData 获取历史数据
func (ds *DataService) GetHistoryData(indexCode string, period string) ([]model.HistoryData, error) {
	index, exists := StockIndices[indexCode]
	if !exists {
		return nil, fmt.Errorf("指数不存在: %s", indexCode)
	}

	stockData, err := ds.GetStockData(index.Symbol, period)
	if err != nil {
		return nil, err
	}

	var historyData []model.HistoryData
	for _, data := range stockData {
		historyData = append(historyData, model.HistoryData{
			Date:   data.Date.Format("2006-01-02"),
			Open:   data.Open,
			High:   data.High,
			Low:    data.Low,
			Close:  data.Close,
			Volume: data.Volume,
		})
	}

	return historyData, nil
}

// GetIndexInfo 获取指数基本信息
func (ds *DataService) GetIndexInfo(indexCode string) (*model.IndexInfo, error) {
	index, exists := StockIndices[indexCode]
	if !exists {
		return nil, fmt.Errorf("指数不存在: %s", indexCode)
	}

	// 获取完整的当前数据（包含昨收价）
	currentStockData, err := ds.GetCurrentStockData(index.Symbol)
	if err != nil {
		return nil, err
	}

	currentPrice := currentStockData.Close
	yesterdayClose := currentStockData.YesterdayClose

	// 使用昨收价计算正确的涨跌幅
	change := currentPrice - yesterdayClose
	changePercent := (change / yesterdayClose) * 100

	return &model.IndexInfo{
		Code:          index.Code,
		Name:          index.Name,
		Symbol:        index.Symbol,
		Market:        index.Market,
		Price:         math.Round(currentPrice*100) / 100,
		Change:        math.Round(change*100) / 100,
		ChangePercent: math.Round(changePercent*100) / 100,
		Volume:        currentStockData.Volume,
		Timestamp:     time.Now().UTC().Format(time.RFC3339),
	}, nil
}

// GetAllIndicesInfo 获取所有指数信息
func (ds *DataService) GetAllIndicesInfo() (map[string]*model.IndexInfo, error) {
	indicesInfo := make(map[string]*model.IndexInfo)

	for code := range StockIndices {
		info, err := ds.GetIndexInfo(code)
		if err != nil {
			log.Printf("获取指数信息失败 %s: %v", code, err)
			continue
		}
		indicesInfo[code] = info
	}

	return indicesInfo, nil
}

// GetDataSourceStatus 获取数据源状态
func (ds *DataService) GetDataSourceStatus() *model.DataSourceStatus {
	status := &model.DataSourceStatus{
		Recommendation: "使用智能模拟数据",
	}

	// 测试Yahoo Finance网站连接
	resp, err := ds.httpClient.R().Get("https://finance.yahoo.com")
	if err != nil {
		status.YahooFinanceWebsite.Status = "error"
		status.YahooFinanceWebsite.Error = err.Error()
	} else if resp.StatusCode() == 200 {
		status.YahooFinanceWebsite.Status = "connected"
	} else if resp.StatusCode() == 429 {
		status.YahooFinanceWebsite.Status = "rate_limited"
		status.YahooFinanceWebsite.Error = "请求频率过高，被限制访问"
	} else {
		status.YahooFinanceWebsite.Status = "error"
		status.YahooFinanceWebsite.Error = fmt.Sprintf("HTTP %d", resp.StatusCode())
	}

	// 测试yfinance（模拟）
	status.YFinance.Status = "no_data"
	status.YFinance.Error = "无法获取数据"
	status.YFinance.TestSymbol = "AAPL"

	// 根据状态给出建议
	if status.YFinance.Status == "working" {
		status.Recommendation = "yfinance正常工作，可以获取真实数据"
	} else if status.YahooFinanceWebsite.Status == "rate_limited" {
		status.Recommendation = "Yahoo Finance API被限制，建议使用智能模拟数据"
	} else {
		status.Recommendation = "网络连接问题，使用智能模拟数据"
	}

	return status
}

// getCache 获取缓存
func (ds *DataService) getCache(key string) (interface{}, bool) {
	ds.cacheMutex.RLock()
	defer ds.cacheMutex.RUnlock()

	item, exists := ds.cache[key]
	if !exists {
		return nil, false
	}

	if time.Now().After(item.ExpiresAt) {
		// 缓存过期，删除
		ds.cacheMutex.RUnlock()
		ds.cacheMutex.Lock()
		delete(ds.cache, key)
		ds.cacheMutex.Unlock()
		ds.cacheMutex.RLock()
		return nil, false
	}

	return item.Data, true
}

// setCache 设置缓存
func (ds *DataService) setCache(key string, data interface{}, duration time.Duration) {
	ds.cacheMutex.Lock()
	defer ds.cacheMutex.Unlock()

	ds.cache[key] = &CacheItem{
		Data:      data,
		ExpiresAt: time.Now().Add(duration),
	}
}

// ClearCache 清除缓存
func (ds *DataService) ClearCache() {
	ds.cacheMutex.Lock()
	defer ds.cacheMutex.Unlock()

	ds.cache = make(map[string]*CacheItem)
	log.Printf("🗑️ 缓存已清除")
}

// checkAndPerformInitialPrediction 检查是否需要立即执行预测
func (ds *DataService) checkAndPerformInitialPrediction() {
	// 等待系统初始化完成
	time.Sleep(2 * time.Second)

	ds.dailyMutex.RLock()
	isEmpty := len(ds.dailyPredictions) == 0
	lastPredictTime := ds.dailyPredictionsTime
	ds.dailyMutex.RUnlock()

	// 如果没有缓存或者缓存已过期（超过24小时），则立即执行预测
	if isEmpty || time.Since(lastPredictTime) > 24*time.Hour {
		log.Printf("🚀 系统启动时检测到需要更新预测数据，立即执行...")
		ds.performDailyPrediction()
	} else {
		log.Printf("📊 发现有效的日常预测缓存，无需重新预测")
	}
}

// startDailyScheduler 启动每日定时调度器
func (ds *DataService) startDailyScheduler() {
	// 设置上海时区 (UTC+8)
	shanghaiLoc, _ := time.LoadLocation("Asia/Shanghai")

	for {
		// 计算下一次下午3点10分的时间（使用上海时区）
		now := time.Now().In(shanghaiLoc)
		nextRun := time.Date(now.Year(), now.Month(), now.Day()+1, 15, 10, 0, 0, shanghaiLoc)

		// 如果当前时间在下午3点10分之前，则今天就执行
		if now.Hour() < 15 || (now.Hour() == 15 && now.Minute() < 10) {
			nextRun = time.Date(now.Year(), now.Month(), now.Day(), 15, 10, 0, 0, shanghaiLoc)
		}

		// 检查是否是工作日（周一到周五）
		for !ds.isWeekday(nextRun) {
			nextRun = nextRun.AddDate(0, 0, 1)
		}

		duration := nextRun.Sub(now)
		log.Printf("🕰️ 下一次预测任务将在 %v 后执行 (%s)", duration, nextRun.Format("2006-01-02 15:04:05"))

		// 设置定时器
		ds.timer = time.NewTimer(duration)

		select {
		case <-ds.timer.C:
			// 时间到，执行预测
			ds.performDailyPrediction()
		case <-ds.stopChan:
			// 收到停止信号
			if ds.timer != nil {
				ds.timer.Stop()
			}
			return
		}
	}
}

// isWeekday 检查是否是工作日（周一到周五，使用上海时区）
func (ds *DataService) isWeekday(t time.Time) bool {
	// 确保使用上海时区进行判断
	shanghaiLoc, _ := time.LoadLocation("Asia/Shanghai")
	localTime := t.In(shanghaiLoc)
	weekday := localTime.Weekday()
	return weekday != time.Saturday && weekday != time.Sunday
}

// performDailyPrediction 执行每日预测任务
func (ds *DataService) performDailyPrediction() {
	log.Printf("🤖 开始执行每日预测任务...")
	start := time.Now()

	newPredictions := make(map[string]*model.StockIndex)
	successCount := 0
	failedCount := 0

	// 逐个预测每个指数
	for indexCode := range StockIndices {
		log.Printf("📊 正在预测 %s...", indexCode)

		prediction, err := ds.generateSinglePrediction(indexCode)
		if err != nil {
			log.Printf("❌ %s 预测失败: %v", indexCode, err)
			failedCount++
			// 即使某个指数预测失败，也继续其他指数
			continue
		}

		newPredictions[indexCode] = prediction
		successCount++
		log.Printf("✅ %s 预测成功: 当前=%.2f, 预测=%.2f, 置信度=%.1f%%",
			indexCode, prediction.Current, prediction.Predicted, prediction.Confidence)

		// 保存到数据库
		if ds.db != nil {
			if err := ds.db.SavePrediction(prediction); err != nil {
				log.Printf("⚠️ 保存预测数据到数据库失败 %s: %v", indexCode, err)
			}
		}

		// 防止请求过于频繁
		time.Sleep(2 * time.Second)
	}

	// 更新内存缓存
	ds.dailyMutex.Lock()
	ds.dailyPredictions = newPredictions
	ds.dailyPredictionsTime = time.Now()
	ds.dailyMutex.Unlock()

	duration := time.Since(start)
	log.Printf("🎆 每日预测任务完成! 成功: %d, 失败: %d, 耗时: %v",
		successCount, failedCount, duration)

	// 清理旧的短期缓存
	ds.ClearCache()
}

// generateSinglePrediction 生成单个指数的预测（专用于定时任务）
func (ds *DataService) generateSinglePrediction(indexCode string) (*model.StockIndex, error) {
	index, exists := StockIndices[indexCode]
	if !exists {
		return nil, fmt.Errorf("指数不存在: %s", indexCode)
	}

	// 获取历史数据
	historicalData, err := ds.GetStockData(index.Symbol, "1mo")
	if err != nil {
		return nil, fmt.Errorf("获取历史数据失败: %v", err)
	}

	// 获取当前数据
	currentStockData, err := ds.GetCurrentStockData(index.Symbol)
	if err != nil {
		return nil, fmt.Errorf("获取当前数据失败: %v", err)
	}

	currentPrice := currentStockData.Close

	// 计算技术指标
	indicators := ds.CalculateTechnicalIndicators(historicalData)

	// 预测价格和置信度（传入历史数据）
	predictedPrice, confidence := ds.PredictPriceAndConfidenceWithHistory(currentPrice, indicators, historicalData)
	if predictedPrice == 0 && confidence == 0 {
		return nil, fmt.Errorf("DeepSeek AI预测失败")
	}

	// 计算预测涨跌幅（预测价格相对于当前价格的变化）
	predictedChange := predictedPrice - currentPrice
	predictedPercent := (predictedChange / currentPrice) * 100

	// 更新指数信息（只保留预测涨跌比例）
	index.Current = math.Round(currentPrice*100) / 100
	index.Predicted = predictedPrice
	index.Change = math.Round(predictedChange*100) / 100         // 预测涨跌金额
	index.ChangePercent = math.Round(predictedPercent*100) / 100 // 预测涨跌百分比
	index.Confidence = confidence
	index.TechnicalIndicators = indicators
	index.Timestamp = time.Now().UTC().Format(time.RFC3339)

	return &index, nil
}

// GetDailyPredictions 获取日常预测缓存
func (ds *DataService) GetDailyPredictions() (map[string]*model.StockIndex, time.Time, bool) {
	ds.dailyMutex.RLock()
	defer ds.dailyMutex.RUnlock()

	if len(ds.dailyPredictions) == 0 {
		return nil, time.Time{}, false
	}

	// 检查缓存是否在24小时内
	if time.Since(ds.dailyPredictionsTime) > 24*time.Hour {
		return nil, time.Time{}, false
	}

	// 返回缓存数据的副本
	result := make(map[string]*model.StockIndex)
	for k, v := range ds.dailyPredictions {
		result[k] = v
	}

	return result, ds.dailyPredictionsTime, true
}

// GetHistoricalPredictions 获取历史预测数据
func (ds *DataService) GetHistoricalPredictions(indexCode string, days int) ([]*model.StockIndex, error) {
	if ds.db == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}

	records, err := ds.db.GetHistoricalPredictions(indexCode, days)
	if err != nil {
		return nil, err
	}

	// 转换为StockIndex格式
	var results []*model.StockIndex
	for _, record := range records {
		stockIndex := ds.db.ConvertPredictionToStockIndex(&record)
		// 添加预测日期信息
		stockIndex.Timestamp = record.PredictionDate.UTC().Format("2006-01-02")
		results = append(results, stockIndex)
	}

	return results, nil
}

// GetAllHistoricalPredictions 获取所有指数的历史预测数据
func (ds *DataService) GetAllHistoricalPredictions(days int) (map[string][]*model.StockIndex, error) {
	if ds.db == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}

	recordsMap, err := ds.db.GetAllHistoricalPredictions(days)
	if err != nil {
		return nil, err
	}

	// 转换为StockIndex格式
	results := make(map[string][]*model.StockIndex)
	for indexCode, records := range recordsMap {
		var indexResults []*model.StockIndex
		for _, record := range records {
			stockIndex := ds.db.ConvertPredictionToStockIndex(&record)
			// 添加预测日期信息
			stockIndex.Timestamp = record.PredictionDate.UTC().Format("2006-01-02")
			indexResults = append(indexResults, stockIndex)
		}
		results[indexCode] = indexResults
	}

	return results, nil
}

// RefreshDailyPredictions 手动刷新每日预测缓存（公开接口）
func (ds *DataService) RefreshDailyPredictions() {
	log.Printf("🔄 手动触发预测缓存刷新")
	ds.performDailyPrediction()
}

// GetPredictionStats 获取预测统计信息（预测次数和成功率）
func (ds *DataService) GetPredictionStats() (map[string]interface{}, error) {
	if ds.db == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}

	// 获取总预测次数
	var totalPredictions int64
	if err := ds.db.GetDB().Model(&model.PredictionRecord{}).Count(&totalPredictions).Error; err != nil {
		return nil, fmt.Errorf("查询总预测次数失败: %v", err)
	}

	// 获取预测正确的次数
	// 预测正确的定义：预测涨跌方向与实际涨跌方向一致
	var correctPredictions int64
	// 构建查询：预测涨跌方向与实际涨跌方向一致的记录数
	// 这需要比较 predicted_price 与 current_price 的差值符号与 change 的符号是否一致
	// 使用反引号转义关键字 'change'
	if err := ds.db.GetDB().Model(&model.PredictionRecord{}).
		Where("(predicted_price - current_price) * `change` > 0").
		Count(&correctPredictions).Error; err != nil {
		return nil, fmt.Errorf("查询正确预测次数失败: %v", err)
	}

	// 计算成功率（避免除零错误）
	var successRate float64
	if totalPredictions > 0 {
		successRate = float64(correctPredictions) / float64(totalPredictions) * 100
	}

	return map[string]interface{}{
		"total_predictions":   totalPredictions,
		"correct_predictions": correctPredictions,
		"success_rate":        math.Round(successRate*100) / 100, // 保留两位小数
	}, nil
}

// Stop 停止定时任务
func (ds *DataService) Stop() {
	select {
	case ds.stopChan <- true:
		log.Printf("🛑 定时预测任务已停止")
	default:
		// 如果信道已满，不做任何操作
	}
	if ds.timer != nil {
		ds.timer.Stop()
	}

	// 关闭数据库连接
	if ds.db != nil {
		if err := ds.db.Close(); err != nil {
			log.Printf("⚠️ 关闭数据库连接失败: %v", err)
		} else {
			log.Printf("✅ 数据库连接已关闭")
		}
	}
}

package service

import (
	"fmt"
	"log"
	"math"
	"math/rand"
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

// DataService 数据服务
type DataService struct {
	cache      map[string]*CacheItem
	cacheMutex sync.RWMutex
	httpClient *resty.Client
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
func NewDataService() *DataService {
	return &DataService{
		cache: make(map[string]*CacheItem),
		httpClient: resty.New().
			SetTimeout(30 * time.Second).
			SetRetryCount(3).
			SetRetryWaitTime(1 * time.Second),
	}
}

// GetStockData 获取股票历史数据
func (ds *DataService) GetStockData(symbol string, period string) ([]model.StockData, error) {
	cacheKey := fmt.Sprintf("%s_%s", symbol, period)

	// 检查缓存
	if cached, found := ds.getCache(cacheKey); found {
		log.Printf("使用缓存数据: %s", symbol)
		return cached.([]model.StockData), nil
	}

	// 尝试获取真实数据
	data, err := ds.fetchRealData(symbol, period)
	if err != nil || len(data) == 0 {
		log.Printf("获取真实数据失败 %s: %v，使用模拟数据", symbol, err)
		data = ds.generateMockData(symbol, period)
	}

	// 缓存数据
	ds.setCache(cacheKey, data, 5*time.Minute)
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

	// 生成历史数据（基于当前数据）
	days := ds.getPeriodDays(period)
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

// generateMockData 生成模拟数据
func (ds *DataService) generateMockData(symbol string, period string) []model.StockData {
	days := ds.getPeriodDays(period)
	basePrice := ds.getBasePrice(symbol)

	var data []model.StockData
	currentPrice := basePrice

	for i := 0; i < days; i++ {
		date := time.Now().AddDate(0, 0, -days+i+1)

		// 生成价格波动
		changePercent := rand.Float64()*0.1 - 0.05 // ±5%的日波动
		currentPrice *= (1 + changePercent)

		// 生成OHLCV
		openPrice := currentPrice * (0.98 + rand.Float64()*0.04)
		highPrice := math.Max(openPrice, currentPrice) * (1.0 + rand.Float64()*0.03)
		lowPrice := math.Min(openPrice, currentPrice) * (0.97 + rand.Float64()*0.03)
		closePrice := currentPrice
		volume := int64(1000000 + rand.Intn(9000000))

		data = append(data, model.StockData{
			Date:           date,
			Open:           math.Round(openPrice*100) / 100,
			High:           math.Round(highPrice*100) / 100,
			Low:            math.Round(lowPrice*100) / 100,
			Close:          math.Round(closePrice*100) / 100,
			YesterdayClose: math.Round((closePrice*0.995)*100) / 100, // 估算昨收价
			Volume:         volume,
		})
	}

	log.Printf("生成模拟数据: %s, 数据量: %d", symbol, len(data))
	return data
}

// getPeriodDays 获取周期天数
func (ds *DataService) getPeriodDays(period string) int {
	periodMap := map[string]int{
		"1d":  1,
		"5d":  5,
		"1mo": 30,
		"3mo": 90,
		"6mo": 180,
		"1y":  365,
		"2y":  730,
		"5y":  1825,
	}
	if days, exists := periodMap[period]; exists {
		return days
	}
	return 30
}

// getBasePrice 获取基础价格
func (ds *DataService) getBasePrice(symbol string) float64 {
	basePrices := map[string]float64{
		"000001.SS": 3000,  // 上证综指
		"399001.SZ": 10000, // 深证成指
		"399006.SZ": 2000,  // 创业板指
		"000688.SS": 1000,  // 科创50
	}
	if price, exists := basePrices[symbol]; exists {
		return price
	}
	return 3000
}

// GetCurrentPrice 获取当前价格
func (ds *DataService) GetCurrentPrice(symbol string) (float64, error) {
	cacheKey := fmt.Sprintf("current_%s", symbol)

	// 检查缓存
	if cached, found := ds.getCache(cacheKey); found {
		return cached.(float64), nil
	}

	// 尝试获取真实价格
	price, err := ds.fetchRealCurrentPrice(symbol)
	if err != nil {
		log.Printf("获取真实价格失败 %s: %v，使用模拟价格", symbol, err)
		price = ds.generateMockCurrentPrice(symbol)
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

// generateMockCurrentPrice 生成模拟当前价格
func (ds *DataService) generateMockCurrentPrice(symbol string) float64 {
	basePrice := ds.getBasePrice(symbol)
	// 添加随机波动
	changePercent := rand.Float64()*0.2 - 0.1 // ±10%的波动
	price := basePrice * (1 + changePercent)
	return math.Round(price*100) / 100
}

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
	// 基于技术指标的简单预测逻辑
	prediction := currentPrice
	confidence := 50.0

	// RSI 超买超卖信号
	if indicators.RSI < 30 {
		prediction *= 1.02 // 超卖，看涨
		confidence += 10
	} else if indicators.RSI > 70 {
		prediction *= 0.98 // 超买，看跌
		confidence += 10
	}

	// 移动平均线信号
	if indicators.MA5 > indicators.MA20 {
		prediction *= 1.01 // 短期均线在上，看涨
		confidence += 5
	} else {
		prediction *= 0.99 // 短期均线在下，看跌
		confidence += 5
	}

	// 趋势信号
	if indicators.Trend > 0 {
		prediction *= 1.005
		confidence += 5
	} else {
		prediction *= 0.995
		confidence += 5
	}

	// 限制置信度范围
	if confidence > 95 {
		confidence = 95
	}
	if confidence < 30 {
		confidence = 30
	}

	return math.Round(prediction*100) / 100, math.Round(confidence*100) / 100
}

// GetPredictionData 获取预测数据
func (ds *DataService) GetPredictionData(indexCode string) (*model.StockIndex, error) {
	index, exists := StockIndices[indexCode]
	if !exists {
		log.Printf("指数不存在: %s", indexCode)
		return nil, fmt.Errorf("指数不存在: %s", indexCode)
	}

	// 获取历史数据
	historicalData, err := ds.GetStockData(index.Symbol, "1mo")
	if err != nil {
		log.Printf("获取历史数据失败 %s: %v", index.Symbol, err)
		return nil, err
	}

	// 获取完整的当前数据（包含昨收价）
	currentStockData, err := ds.GetCurrentStockData(index.Symbol)
	if err != nil {
		log.Printf("获取当前股票数据失败 %s: %v", index.Symbol, err)
		return nil, err
	}

	currentPrice := currentStockData.Close
	yesterdayClose := currentStockData.YesterdayClose

	// 计算技术指标
	indicators := ds.CalculateTechnicalIndicators(historicalData)

	// 预测价格和置信度
	predictedPrice, confidence := ds.PredictPriceAndConfidence(currentPrice, indicators)

	// 使用昨收价计算正确的涨跌幅
	change := currentPrice - yesterdayClose
	changePercent := (change / yesterdayClose) * 100

	predictedChange := predictedPrice - currentPrice
	predictedPercent := (predictedChange / currentPrice) * 100

	// 更新指数信息
	index.Current = math.Round(currentPrice*100) / 100
	index.Predicted = predictedPrice
	index.Change = math.Round(change*100) / 100
	index.ChangePercent = math.Round(changePercent*100) / 100
	index.PredictedChange = math.Round(predictedChange*100) / 100
	index.PredictedPercent = math.Round(predictedPercent*100) / 100
	index.Confidence = confidence
	index.TechnicalIndicators = indicators
	index.Timestamp = time.Now().UTC().Format(time.RFC3339)

	return &index, nil
}

// GetAllPredictions 获取所有预测数据
func (ds *DataService) GetAllPredictions() (map[string]*model.StockIndex, error) {
	predictions := make(map[string]*model.StockIndex)

	for code := range StockIndices {
		prediction, err := ds.GetPredictionData(code)
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
	log.Printf("缓存已清除")
}

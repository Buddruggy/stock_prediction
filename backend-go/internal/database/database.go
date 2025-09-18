package database

import (
	"fmt"
	"log"
	"math"
	"stock-prediction-backend/internal/config"
	"stock-prediction-backend/internal/model"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DatabaseService 数据库服务
type DatabaseService struct {
	db     *gorm.DB
	config *config.Config
}

// NewDatabaseService 创建数据库服务实例
func NewDatabaseService(cfg *config.Config) (*DatabaseService, error) {
	// 连接数据库
	db, err := gorm.Open(mysql.Open(cfg.Database.GetDSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %v", err)
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取数据库连接失败: %v", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	service := &DatabaseService{
		db:     db,
		config: cfg,
	}

	// 初始化数据库表
	if err := service.initTables(); err != nil {
		return nil, fmt.Errorf("初始化数据库表失败: %v", err)
	}

	log.Printf("✅ 数据库连接成功: %s", cfg.Database.DBName)
	return service, nil
}

// initTables 初始化数据库表
func (ds *DatabaseService) initTables() error {
	// 创建预测记录表
	if err := ds.db.AutoMigrate(&model.PredictionRecord{}); err != nil {
		return fmt.Errorf("创建预测记录表失败: %v", err)
	}

	// 创建统一的历史数据表
	if err := ds.db.AutoMigrate(&model.HistoricalData{}); err != nil {
		return fmt.Errorf("创建历史数据表失败: %v", err)
	}

	// 为历史数据表添加唯一约束索引
	if err := ds.db.Exec("ALTER TABLE historical_data ADD UNIQUE INDEX idx_unique_index_date (index_code, date)").Error; err != nil {
		// 如果索引已存在，忽略错误
		if !strings.Contains(err.Error(), "Duplicate key name") {
			log.Printf("⚠️ 创建历史数据表索引失败: %v", err)
		}
	}

	log.Printf("📊 数据库表初始化完成")
	return nil
}

// SavePrediction 保存预测记录
func (ds *DatabaseService) SavePrediction(prediction *model.StockIndex) error {
	record := &model.PredictionRecord{
		IndexCode:      prediction.Code,
		IndexName:      prediction.Name,
		PredictionDate: time.Now().UTC().Truncate(24 * time.Hour), // 使用UTC时区确保一致性
		CurrentPrice:   prediction.Current,
		PredictedPrice: prediction.Predicted,
		Change:         prediction.Change,
		ChangePercent:  prediction.ChangePercent,
		Confidence:     prediction.Confidence,
		MA5:            prediction.TechnicalIndicators.MA5,
		MA20:           prediction.TechnicalIndicators.MA20,
		RSI:            prediction.TechnicalIndicators.RSI,
		Volatility:     prediction.TechnicalIndicators.Volatility,
		Trend:          prediction.TechnicalIndicators.Trend,
	}

	// 使用 UPSERT 操作，如果记录存在则更新，不存在则插入
	result := ds.db.Where("index_code = ? AND prediction_date = ?",
		record.IndexCode, record.PredictionDate).
		Assign(record).
		FirstOrCreate(record)

	if result.Error != nil {
		return fmt.Errorf("保存预测记录失败 %s: %v", prediction.Code, result.Error)
	}

	log.Printf("💾 保存预测记录: %s (当前=%.2f, 预测=%.2f, 置信度=%.1f%%)",
		prediction.Code, prediction.Current, prediction.Predicted, prediction.Confidence)
	return nil
}

// GetLatestPrediction 获取最新预测记录
func (ds *DatabaseService) GetLatestPrediction(indexCode string) (*model.PredictionRecord, error) {
	var record model.PredictionRecord

	result := ds.db.Where("index_code = ?", indexCode).
		Order("prediction_date DESC").
		First(&record)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // 没有找到记录，返回 nil
		}
		return nil, fmt.Errorf("查询预测记录失败 %s: %v", indexCode, result.Error)
	}

	return &record, nil
}

// GetTodayPrediction 获取今日预测记录
func (ds *DatabaseService) GetTodayPrediction(indexCode string) (*model.PredictionRecord, error) {
	var record model.PredictionRecord
	// 使用UTC时区获取当前日期，确保与数据库时区一致
	today := time.Now().UTC().Truncate(24 * time.Hour)

	result := ds.db.Where("index_code = ? AND prediction_date = ?", indexCode, today).
		First(&record)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // 没有找到记录，返回 nil
		}
		return nil, fmt.Errorf("查询今日预测记录失败 %s: %v", indexCode, result.Error)
	}

	log.Printf("📊 从数据库成功获取今日预测: %s (日期: %s)", indexCode, today.Format("2006-01-02"))
	return &record, nil
}

// GetAllTodayPredictions 获取所有指数的今日预测记录
func (ds *DatabaseService) GetAllTodayPredictions() (map[string]*model.PredictionRecord, error) {
	var records []model.PredictionRecord
	// 使用UTC时区获取当前日期，确保与数据库时区一致
	today := time.Now().UTC().Truncate(24 * time.Hour)

	result := ds.db.Where("prediction_date = ?", today).Find(&records)
	if result.Error != nil {
		return nil, fmt.Errorf("查询今日预测记录失败: %v", result.Error)
	}

	predictionMap := make(map[string]*model.PredictionRecord)
	for i := range records {
		predictionMap[records[i].IndexCode] = &records[i]
	}

	if len(predictionMap) > 0 {
		log.Printf("📊 从数据库成功获取所有今日预测: %d 条记录 (日期: %s)", len(predictionMap), today.Format("2006-01-02"))
	}

	return predictionMap, nil
}

// SaveHistoricalData 保存历史数据
func (ds *DatabaseService) SaveHistoricalData(indexCode, indexName string, data []model.StockData) error {
	if len(data) == 0 {
		return nil
	}

	// 批量插入或更新历史数据到统一表
	var records []model.HistoricalData
	for _, stockData := range data {
		records = append(records, model.HistoricalData{
			IndexCode: indexCode,
			IndexName: indexName,
			Date:      stockData.Date.Truncate(24 * time.Hour), // 只保留日期部分
			Open:      stockData.Open,
			High:      stockData.High,
			Low:       stockData.Low,
			Close:     stockData.Close,
			Volume:    stockData.Volume,
		})
	}

	// 批量插入，遇到重复则忽略
	result := ds.db.Create(&records)
	if result.Error != nil {
		// 如果是重复键错误，尝试逐个更新
		for _, record := range records {
			ds.db.Where("index_code = ? AND date = ?",
				record.IndexCode, record.Date).
				Assign(record).
				FirstOrCreate(&record)
		}
	}

	log.Printf("💾 保存历史数据: %s, 数据量: %d", indexCode, len(records))
	return nil
}

// GetHistoricalData 获取历史数据
func (ds *DatabaseService) GetHistoricalData(indexCode string, days int) ([]model.StockData, error) {
	var records []model.HistoricalData
	result := ds.db.Where("index_code = ?", indexCode).
		Order("date DESC").
		Limit(days).
		Find(&records)

	if result.Error != nil {
		return nil, fmt.Errorf("查询历史数据失败 %s: %v", indexCode, result.Error)
	}

	// 转换为 StockData 格式
	var stockData []model.StockData
	for i := len(records) - 1; i >= 0; i-- { // 反转顺序，使其按日期升序
		record := records[i]
		stockData = append(stockData, model.StockData{
			Date:   record.Date,
			Open:   record.Open,
			High:   record.High,
			Low:    record.Low,
			Close:  record.Close,
			Volume: record.Volume,
		})
	}

	return stockData, nil
}

// GetHistoricalPredictions 获取历史预测记录
func (ds *DatabaseService) GetHistoricalPredictions(indexCode string, days int) ([]model.PredictionRecord, error) {
	var records []model.PredictionRecord

	// 计算起始日期
	startDate := time.Now().UTC().AddDate(0, 0, -days).Truncate(24 * time.Hour)

	result := ds.db.Where("index_code = ? AND prediction_date >= ?", indexCode, startDate).
		Order("prediction_date DESC").
		Find(&records)

	if result.Error != nil {
		return nil, fmt.Errorf("查询历史预测记录失败 %s: %v", indexCode, result.Error)
	}

	log.Printf("📊 从数据库获取历史预测记录: %s, 数量: %d, 天数: %d", indexCode, len(records), days)
	return records, nil
}

// GetAllHistoricalPredictions 获取所有指数的历史预测记录
func (ds *DatabaseService) GetAllHistoricalPredictions(days int) (map[string][]model.PredictionRecord, error) {
	var records []model.PredictionRecord

	// 计算起始日期
	startDate := time.Now().UTC().AddDate(0, 0, -days).Truncate(24 * time.Hour)

	result := ds.db.Where("prediction_date >= ?", startDate).
		Order("index_code, prediction_date DESC").
		Find(&records)

	if result.Error != nil {
		return nil, fmt.Errorf("查询所有历史预测记录失败: %v", result.Error)
	}

	// 按指数代码分组
	predictionMap := make(map[string][]model.PredictionRecord)
	for _, record := range records {
		predictionMap[record.IndexCode] = append(predictionMap[record.IndexCode], record)
	}

	log.Printf("📊 从数据库获取所有历史预测记录: %d 条记录, 天数: %d", len(records), days)
	return predictionMap, nil
}

// ConvertPredictionToStockIndex 将预测记录转换为StockIndex
func (ds *DatabaseService) ConvertPredictionToStockIndex(record *model.PredictionRecord) *model.StockIndex {
	return &model.StockIndex{
		Code:          record.IndexCode,
		Name:          record.IndexName,
		Current:       record.CurrentPrice,
		Predicted:     record.PredictedPrice,
		Change:        record.Change,
		ChangePercent: record.ChangePercent,
		Confidence:    record.Confidence,
		TechnicalIndicators: model.TechnicalIndicators{
			MA5:        record.MA5,
			MA20:       record.MA20,
			RSI:        record.RSI,
			Volatility: record.Volatility,
			Trend:      record.Trend,
		},
		Timestamp: record.CreatedAt.UTC().Format(time.RFC3339),
	}
}

// GetDB 获取底层gorm.DB对象
func (ds *DatabaseService) GetDB() *gorm.DB {
	return ds.db
}

// Close 关闭数据库连接
func (ds *DatabaseService) Close() error {
	sqlDB, err := ds.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// GetHistoricalPredictionsForDate 获取指定日期的历史预测记录
func (ds *DatabaseService) GetHistoricalPredictionsForDate(date time.Time) ([]model.PredictionRecord, error) {
	var records []model.PredictionRecord

	result := ds.db.Where("prediction_date = ? AND is_correct IS NULL", date).
		Find(&records)

	if result.Error != nil {
		return nil, fmt.Errorf("查询历史预测记录失败: %v", result.Error)
	}

	return records, nil
}

// UpdatePredictionAccuracy 更新预测准确性
func (ds *DatabaseService) UpdatePredictionAccuracy(recordID uint, isCorrect bool) error {
	result := ds.db.Model(&model.PredictionRecord{}).
		Where("id = ?", recordID).
		Updates(map[string]interface{}{
			"is_correct": isCorrect,
		})

	if result.Error != nil {
		return fmt.Errorf("更新预测准确性失败: %v", result.Error)
	}

	return nil
}

// GetPredictionStats 获取预测统计信息
func (ds *DatabaseService) GetPredictionStats() (map[string]interface{}, error) {
	// 获取总预测次数（已验证的）
	var totalPredictions int64
	if err := ds.db.Model(&model.PredictionRecord{}).
		Where("is_correct IS NOT NULL").
		Count(&totalPredictions).Error; err != nil {
		return nil, fmt.Errorf("查询总预测次数失败: %v", err)
	}

	// 获取预测正确的次数
	var correctPredictions int64
	if err := ds.db.Model(&model.PredictionRecord{}).
		Where("is_correct = ?", true).
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

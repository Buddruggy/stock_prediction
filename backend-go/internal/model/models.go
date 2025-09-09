package model

import (
	"time"
)

// StockIndex 股票指数信息
type StockIndex struct {
	Code                string              `json:"code"`
	Name                string              `json:"name"`
	Symbol              string              `json:"symbol"`
	Market              string              `json:"market"`
	Current             float64             `json:"current"`
	Predicted           float64             `json:"predicted"`
	Change              float64             `json:"change"`        // 预测涨跌金额
	ChangePercent       float64             `json:"changePercent"` // 预测涨跌百分比
	Confidence          float64             `json:"confidence"`
	TechnicalIndicators TechnicalIndicators `json:"technical_indicators"`
	Timestamp           string              `json:"timestamp"`
}

// TechnicalIndicators 技术指标
type TechnicalIndicators struct {
	MA5        float64 `json:"ma_5"`
	MA20       float64 `json:"ma_20"`
	RSI        float64 `json:"rsi"`
	Volatility float64 `json:"volatility"`
	Trend      float64 `json:"trend"`
}

// APIResponse API响应结构
type APIResponse struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	Timestamp string      `json:"timestamp"`
}

// HistoryData 历史数据
type HistoryData struct {
	Date   string  `json:"date"`
	Open   float64 `json:"open"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
	Volume int64   `json:"volume"`
}

// IndexInfo 指数基本信息
type IndexInfo struct {
	Code          string  `json:"code"`
	Name          string  `json:"name"`
	Symbol        string  `json:"symbol"`
	Market        string  `json:"market"`
	Price         float64 `json:"price"`
	Change        float64 `json:"change"`
	ChangePercent float64 `json:"changePercent"`
	Volume        int64   `json:"volume"`
	Timestamp     string  `json:"timestamp"`
}

// DataSourceStatus 数据源状态
type DataSourceStatus struct {
	YFinance struct {
		Status     string `json:"status"`
		Error      string `json:"error"`
		TestSymbol string `json:"test_symbol"`
	} `json:"yfinance"`
	YahooFinanceWebsite struct {
		Status string `json:"status"`
		Error  string `json:"error"`
	} `json:"yahoo_finance_website"`
	Recommendation string `json:"recommendation"`
}

// StockData 股票数据
type StockData struct {
	Date           time.Time `json:"date"`
	Open           float64   `json:"open"`
	High           float64   `json:"high"`
	Low            float64   `json:"low"`
	Close          float64   `json:"close"`
	YesterdayClose float64   `json:"yesterday_close"` // 昨收价
	Volume         int64     `json:"volume"`
}

// ===== 数据库模型 =====

// PredictionRecord 预测记录数据库模型
type PredictionRecord struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	IndexCode      string    `gorm:"type:varchar(20);not null;index" json:"index_code"`  // 指数代码
	IndexName      string    `gorm:"type:varchar(50);not null" json:"index_name"`        // 指数名称
	PredictionDate time.Time `gorm:"type:date;not null;index" json:"prediction_date"`    // 预测日期
	CurrentPrice   float64   `gorm:"type:decimal(10,2);not null" json:"current_price"`   // 当前价格
	PredictedPrice float64   `gorm:"type:decimal(10,2);not null" json:"predicted_price"` // 预测价格
	Change         float64   `gorm:"type:decimal(10,2);not null" json:"change"`          // 预测涨跌金额
	ChangePercent  float64   `gorm:"type:decimal(5,2);not null" json:"change_percent"`   // 预测涨跌百分比
	Confidence     float64   `gorm:"type:decimal(5,2);not null" json:"confidence"`       // 置信度
	MA5            float64   `gorm:"type:decimal(10,2)" json:"ma5"`                      // 5日移动平均线
	MA20           float64   `gorm:"type:decimal(10,2)" json:"ma20"`                     // 20日移动平均线
	RSI            float64   `gorm:"type:decimal(5,2)" json:"rsi"`                       // RSI指标
	Volatility     float64   `gorm:"type:decimal(5,2)" json:"volatility"`                // 波动率
	Trend          float64   `gorm:"type:decimal(5,2)" json:"trend"`                     // 趋势指标
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`                   // 创建时间
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`                   // 更新时间
}

// TableName 设置表名
func (PredictionRecord) TableName() string {
	return "predictions"
}

// HistoricalData 历史数据数据库模型
type HistoricalData struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	IndexCode string    `gorm:"type:varchar(20);not null;index" json:"index_code"` // 指数代码
	IndexName string    `gorm:"type:varchar(50);not null" json:"index_name"`       // 指数名称
	Date      time.Time `gorm:"type:date;not null;index" json:"date"`              // 日期
	Open      float64   `gorm:"type:decimal(10,2);not null" json:"open"`           // 开盘价
	High      float64   `gorm:"type:decimal(10,2);not null" json:"high"`           // 最高价
	Low       float64   `gorm:"type:decimal(10,2);not null" json:"low"`            // 最低价
	Close     float64   `gorm:"type:decimal(10,2);not null" json:"close"`          // 收盘价
	Volume    int64     `gorm:"type:bigint;not null" json:"volume"`                // 成交量
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`                  // 创建时间
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`                  // 更新时间
}

// TableName 设置表名为统一的历史数据表
func (HistoricalData) TableName() string {
	return "historical_data"
}

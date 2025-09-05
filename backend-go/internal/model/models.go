package model

import "time"

// StockIndex 股票指数信息
type StockIndex struct {
	Code                string              `json:"code"`
	Name                string              `json:"name"`
	Symbol              string              `json:"symbol"`
	Market              string              `json:"market"`
	Current             float64             `json:"current"`
	Predicted           float64             `json:"predicted"`
	Change              float64             `json:"change"`
	ChangePercent       float64             `json:"changePercent"`
	PredictedChange     float64             `json:"predictedChange"`
	PredictedPercent    float64             `json:"predictedPercent"`
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
	Date   time.Time `json:"date"`
	Open   float64   `json:"open"`
	High   float64   `json:"high"`
	Low    float64   `json:"low"`
	Close  float64   `json:"close"`
	Volume int64     `json:"volume"`
}

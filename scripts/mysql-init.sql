-- 智投预测数据库初始化脚本
-- 创建用户及赋权
CREATE USER IF NOT EXISTS 'stock_user'@'%' IDENTIFIED BY 'stock_pass';
GRANT ALL PRIVILEGES ON stock_prediction.* TO 'stock_user'@'%';
FLUSH PRIVILEGES;

-- 使用数据库
USE stock_prediction;

-- 创建预测记录表
CREATE TABLE IF NOT EXISTS predictions (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    index_code VARCHAR(20) NOT NULL COMMENT '指数代码',
    index_name VARCHAR(50) NOT NULL COMMENT '指数名称',
    prediction_date DATE NOT NULL COMMENT '预测日期',
    current_price DECIMAL(10,2) NOT NULL COMMENT '当前价格',
    predicted_price DECIMAL(10,2) NOT NULL COMMENT '预测价格',
    `change` DECIMAL(10,2) NOT NULL COMMENT '预测涨跌金额',
    change_percent DECIMAL(5,2) NOT NULL COMMENT '预测涨跌百分比',
    confidence DECIMAL(5,2) NOT NULL COMMENT '置信度',
    ma5 DECIMAL(10,2) DEFAULT NULL COMMENT '5日移动平均线',
    ma20 DECIMAL(10,2) DEFAULT NULL COMMENT '20日移动平均线',
    rsi DECIMAL(5,2) DEFAULT NULL COMMENT 'RSI指标',
    volatility DECIMAL(5,2) DEFAULT NULL COMMENT '波动率',
    trend DECIMAL(5,2) DEFAULT NULL COMMENT '趋势指标',
    is_correct BOOLEAN DEFAULT NULL COMMENT '预测是否正确',
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY unique_prediction (index_code, prediction_date),
    KEY idx_index_code (index_code),
    KEY idx_prediction_date (prediction_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='预测记录表';

-- 创建统一的历史数据表
CREATE TABLE IF NOT EXISTS historical_data (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    index_code VARCHAR(20) NOT NULL COMMENT '指数代码',
    index_name VARCHAR(50) NOT NULL COMMENT '指数名称',
    date DATE NOT NULL COMMENT '日期',
    open DECIMAL(10,2) NOT NULL COMMENT '开盘价',
    high DECIMAL(10,2) NOT NULL COMMENT '最高价',
    low DECIMAL(10,2) NOT NULL COMMENT '最低价',
    close DECIMAL(10,2) NOT NULL COMMENT '收盘价',
    volume BIGINT NOT NULL COMMENT '成交量',
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY unique_date_code (index_code, date),
    KEY idx_index_code (index_code),
    KEY idx_date (date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='统一历史数据表';

-- 插入一些示例数据以供测试
INSERT IGNORE INTO predictions (index_code, index_name, prediction_date, current_price, predicted_price, `change`, change_percent, confidence, ma5, ma20, rsi, volatility, trend) 
VALUES 
('sh000001', '上证综指', CURDATE(), 3000.00, 3050.00, 50.00, 1.67, 75.5, 2980.5, 2950.2, 55.8, 1.2, 0.8),
('sz399001', '深证成指', CURDATE(), 10000.00, 10150.00, 150.00, 1.50, 72.3, 9950.8, 9900.5, 52.3, 1.5, 1.2);

-- 插入一些示例历史数据
INSERT IGNORE INTO historical_data (index_code, index_name, date, open, high, low, close, volume)
VALUES 
('sh000001', '上证综指', CURDATE() - INTERVAL 1 DAY, 2990.00, 3010.00, 2980.00, 3000.00, 1000000000),
('sz399001', '深证成指', CURDATE() - INTERVAL 1 DAY, 9900.00, 10100.00, 9850.00, 10000.00, 500000000);

-- 数据库初始化完成
SELECT 'MySQL数据库初始化完成' AS status;
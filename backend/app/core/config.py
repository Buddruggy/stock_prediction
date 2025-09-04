"""
配置管理 - 使用 Pydantic Settings
"""

from pydantic_settings import BaseSettings
from typing import List
import os

class Settings(BaseSettings):
    """应用配置"""
    
    # 应用基本信息
    APP_NAME: str = "智投预测 API"
    APP_VERSION: str = "2.0.0"
    APP_DESCRIPTION: str = "AI股市指数预测平台 - RESTful API服务"
    
    # 服务器配置
    HOST: str = "0.0.0.0"
    PORT: int = 8000
    DEBUG: bool = False
    ENVIRONMENT: str = "production"
    
    # CORS配置
    ALLOWED_ORIGINS: List[str] = [
        "http://localhost:3000",  # Vue.js 开发服务器
        "http://localhost:5173",  # Vite 开发服务器
        "http://localhost:8080",  # 备用前端端口
        "https://zhitou-predict.com",  # 生产域名
    ]
    
    # 股票指数配置
    STOCK_INDICES: dict = {
        "sh000001": {
            "name": "上证综指",
            "symbol": "000001.SS",
            "code": "SH000001",
            "market": "上海证券交易所"
        },
        "sz399001": {
            "name": "深证成指", 
            "symbol": "399001.SZ",
            "code": "SZ399001",
            "market": "深圳证券交易所"
        },
        "sz399006": {
            "name": "创业板指",
            "symbol": "399006.SZ", 
            "code": "SZ399006",
            "market": "深圳证券交易所"
        },
        "sh000688": {
            "name": "科创50",
            "symbol": "000688.SS",
            "code": "SH000688", 
            "market": "上海证券交易所"
        }
    }
    
    # ML模型配置
    MODEL_CONFIG: dict = {
        "cache_duration": 300,      # 缓存时间（秒）
        "lookback_days": 30,        # 历史数据天数
        "confidence_threshold": 0.5, # 置信度阈值
        "feature_dimensions": 15,    # 特征维度
        "prediction_horizon": 1,     # 预测天数
    }
    
    # API配置
    API_V1_PREFIX: str = "/api/v1"
    REQUEST_TIMEOUT: int = 30
    MAX_REQUESTS_PER_MINUTE: int = 100
    
    # 数据库配置（可选）
    DATABASE_URL: str = ""
    
    # Redis配置（可选）
    REDIS_URL: str = ""
    
    # 日志配置
    LOG_LEVEL: str = "INFO"
    LOG_FILE: str = "logs/app.log"
    
    class Config:
        env_file = ".env"
        case_sensitive = True

# 创建全局配置实例
settings = Settings()

# 开发环境配置覆盖
if os.getenv("ENVIRONMENT") == "development":
    settings.DEBUG = True
    settings.ENVIRONMENT = "development"
    settings.LOG_LEVEL = "DEBUG"

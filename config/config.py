# 智投预测 - 配置文件
# AI股市指数预测平台配置管理

import os
from typing import Dict, Any

class Config:
    """基础配置类"""
    
    # 应用配置
    APP_NAME = "智投预测"
    APP_VERSION = "1.0.0"
    APP_DESCRIPTION = "AI股市指数预测平台"
    
    # 服务器配置
    HOST = os.getenv('HOST', '0.0.0.0')
    PORT = int(os.getenv('PORT', 9000))
    DEBUG = os.getenv('FLASK_ENV', 'production') == 'development'
    
    # 股票指数配置
    STOCK_INDICES = {
        'sh000001': {
            'name': '上证综合指数',
            'symbol': '000001.SS',
            'code': 'SH000001'
        },
        'sz399001': {
            'name': '深证成分指数',
            'symbol': '399001.SZ',
            'code': 'SZ399001'
        },
        'sz399006': {
            'name': '创业板综合指数',
            'symbol': '399006.SZ',
            'code': 'SZ399006'
        },
        'sh000688': {
            'name': '科创板50指数',
            'symbol': '000688.SS',
            'code': 'SH000688'
        }
    }
    
    # 预测模型配置
    MODEL_CONFIG = {
        'cache_duration': 300,  # 缓存时间（秒）
        'lookback_days': 30,    # 历史数据天数
        'confidence_threshold': 0.5,  # 置信度阈值
        'feature_dimensions': 15,      # 特征维度
    }
    
    # API配置
    API_CONFIG = {
        'rate_limit': '100/hour',
        'cors_origins': ['*'],
        'timeout': 30,
    }

class DevelopmentConfig(Config):
    """开发环境配置"""
    DEBUG = True
    FLASK_ENV = 'development'

class ProductionConfig(Config):
    """生产环境配置"""
    DEBUG = False
    FLASK_ENV = 'production'

class TestingConfig(Config):
    """测试环境配置"""
    DEBUG = True
    TESTING = True
    FLASK_ENV = 'testing'

# 配置字典
config_dict = {
    'development': DevelopmentConfig,
    'production': ProductionConfig,
    'testing': TestingConfig,
    'default': DevelopmentConfig
}

def get_config() -> Config:
    """获取当前环境配置"""
    env = os.getenv('FLASK_ENV', 'development')
    return config_dict.get(env, config_dict['default'])

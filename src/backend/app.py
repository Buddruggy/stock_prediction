#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
股票指数预测API服务器
提供股票指数数据获取和预测功能
"""

from flask import Flask, jsonify, request, send_from_directory
import sys
import os

# 添加项目根目录到Python路径
sys.path.insert(0, os.path.join(os.path.dirname(__file__), '../..'))

from config import get_config
from flask_cors import CORS
import numpy as np
import pandas as pd
from datetime import datetime, timedelta
import yfinance as yf
import warnings
warnings.filterwarnings('ignore')

# 尝试导入机器学习库
try:
    from sklearn.preprocessing import MinMaxScaler
    from sklearn.ensemble import RandomForestRegressor
    from sklearn.metrics import mean_absolute_error
    ML_AVAILABLE = True
except ImportError:
    ML_AVAILABLE = False
    print("警告: 机器学习库未安装，将使用简单预测模型")

# 设置静态文件和模板路径
static_folder_path = os.path.join(os.path.dirname(__file__), '../frontend/static')
template_folder_path = os.path.join(os.path.dirname(__file__), '../frontend/templates')

app = Flask(__name__, 
           static_folder=static_folder_path,
           static_url_path='/static',
           template_folder=template_folder_path)
CORS(app)

# 中国股票指数代码映射
STOCK_INDICES = {
    'sh000001': {
        'name': '上证综指',
        'symbol': '000001.SS',
        'code': 'SH000001'
    },
    'sz399001': {
        'name': '深证成指',
        'symbol': '399001.SZ',
        'code': 'SZ399001'
    },
    'sz399006': {
        'name': '创业板指',
        'symbol': '399006.SZ',
        'code': 'SZ399006'
    },
    'sh000688': {
        'name': '科创50',
        'symbol': '000688.SS',
        'code': 'SH000688'
    }
}

class StockPredictor:
    """股票指数预测器"""
    
    def __init__(self):
        self.scalers = {}
        self.models = {}
        self.data_cache = {}
        self.cache_time = {}
        self.cache_duration = 300  # 缓存5分钟
    
    def get_stock_data(self, symbol, period="1y"):
        """获取股票数据"""
        try:
            # 检查缓存
            if (symbol in self.data_cache and 
                symbol in self.cache_time and 
                datetime.now() - self.cache_time[symbol] < timedelta(seconds=self.cache_duration)):
                return self.data_cache[symbol]
            
            # 获取数据
            stock = yf.Ticker(symbol)
            data = stock.history(period=period)
            
            if data.empty:
                return None
            
            # 计算技术指标
            data = self.calculate_technical_indicators(data)
            
            # 缓存数据
            self.data_cache[symbol] = data
            self.cache_time[symbol] = datetime.now()
            
            return data
            
        except Exception as e:
            print(f"获取股票数据失败 {symbol}: {e}")
            return None
    
    def calculate_technical_indicators(self, data):
        """计算技术指标"""
        # 移动平均线
        data['MA5'] = data['Close'].rolling(window=5).mean()
        data['MA10'] = data['Close'].rolling(window=10).mean()
        data['MA20'] = data['Close'].rolling(window=20).mean()
        
        # RSI
        delta = data['Close'].diff()
        gain = (delta.where(delta > 0, 0)).rolling(window=14).mean()
        loss = (-delta.where(delta < 0, 0)).rolling(window=14).mean()
        rs = gain / loss
        data['RSI'] = 100 - (100 / (1 + rs))
        
        # MACD
        exp1 = data['Close'].ewm(span=12).mean()
        exp2 = data['Close'].ewm(span=26).mean()
        data['MACD'] = exp1 - exp2
        data['MACD_Signal'] = data['MACD'].ewm(span=9).mean()
        
        # 布林带
        data['BB_Middle'] = data['Close'].rolling(window=20).mean()
        bb_std = data['Close'].rolling(window=20).std()
        data['BB_Upper'] = data['BB_Middle'] + (bb_std * 2)
        data['BB_Lower'] = data['BB_Middle'] - (bb_std * 2)
        
        # 成交量相关
        data['Volume_MA'] = data['Volume'].rolling(window=20).mean()
        data['Volume_Ratio'] = data['Volume'] / data['Volume_MA']
        
        return data
    
    def prepare_features(self, data, lookback=30):
        """准备特征数据"""
        features = []
        targets = []
        
        # 选择特征列
        feature_columns = ['Open', 'High', 'Low', 'Close', 'Volume',
                          'MA5', 'MA10', 'MA20', 'RSI', 'MACD', 'MACD_Signal',
                          'BB_Upper', 'BB_Middle', 'BB_Lower', 'Volume_Ratio']
        
        # 删除包含NaN的行
        data = data.dropna()
        
        if len(data) < lookback + 1:
            return None, None
        
        for i in range(lookback, len(data)):
            # 取过去lookback天的数据作为特征
            feature_window = data[feature_columns].iloc[i-lookback:i].values
            features.append(feature_window.flatten())
            
            # 目标是下一天的收盘价
            targets.append(data['Close'].iloc[i])
        
        return np.array(features), np.array(targets)
    
    def simple_predict(self, data):
        """简单预测模型（当ML库不可用时）"""
        if data is None or len(data) < 30:
            return None, 0.5
        
        # 使用最近30天的数据
        recent_data = data.tail(30)
        
        # 计算趋势
        prices = recent_data['Close'].values
        
        # 线性回归趋势
        x = np.arange(len(prices))
        coeffs = np.polyfit(x, prices, 1)
        trend = coeffs[0]
        
        # 移动平均预测
        ma5 = recent_data['MA5'].iloc[-1]
        ma20 = recent_data['MA20'].iloc[-1]
        current_price = prices[-1]
        
        # 综合预测
        trend_prediction = current_price + trend
        ma_prediction = (ma5 + ma20) / 2
        
        # 加权平均
        predicted_price = 0.6 * trend_prediction + 0.4 * ma_prediction
        
        # 简单的置信度计算
        volatility = np.std(prices[-10:]) / current_price
        confidence = max(0.3, min(0.9, 1 - volatility * 10))
        
        return predicted_price, confidence
    
    def ml_predict(self, data):
        """机器学习预测模型"""
        if not ML_AVAILABLE or data is None:
            return self.simple_predict(data)
        
        try:
            # 准备数据
            X, y = self.prepare_features(data)
            if X is None or len(X) < 50:
                return self.simple_predict(data)
            
            # 数据归一化
            if 'scaler' not in self.scalers:
                self.scalers['scaler'] = MinMaxScaler()
                X_scaled = self.scalers['scaler'].fit_transform(X)
            else:
                X_scaled = self.scalers['scaler'].transform(X)
            
            # 训练模型
            if 'model' not in self.models:
                self.models['model'] = RandomForestRegressor(
                    n_estimators=100,
                    random_state=42,
                    n_jobs=-1
                )
            
            # 使用最近的数据训练
            train_size = int(len(X_scaled) * 0.8)
            X_train = X_scaled[:train_size]
            y_train = y[:train_size]
            X_test = X_scaled[train_size:]
            y_test = y[train_size:]
            
            self.models['model'].fit(X_train, y_train)
            
            # 预测
            if len(X_test) > 0:
                y_pred_test = self.models['model'].predict(X_test)
                mae = mean_absolute_error(y_test, y_pred_test)
                confidence = max(0.3, min(0.9, 1 - mae / np.mean(y_test)))
            else:
                confidence = 0.7
            
            # 预测明天的价格
            last_features = X_scaled[-1:] if len(X_scaled) > 0 else None
            if last_features is not None:
                predicted_price = self.models['model'].predict(last_features)[0]
            else:
                return self.simple_predict(data)
            
            return predicted_price, confidence
            
        except Exception as e:
            print(f"ML预测失败: {e}")
            return self.simple_predict(data)
    
    def predict(self, symbol):
        """预测股票价格"""
        data = self.get_stock_data(symbol)
        if data is None:
            return None
        
        # 使用ML预测或简单预测
        predicted_price, confidence = self.ml_predict(data)
        
        if predicted_price is None:
            return None
        
        # 获取当前价格和变化
        current_price = data['Close'].iloc[-1]
        previous_price = data['Close'].iloc[-2] if len(data) > 1 else current_price
        
        current_change = current_price - previous_price
        current_change_percent = (current_change / previous_price) * 100
        
        predicted_change = predicted_price - current_price
        predicted_change_percent = (predicted_change / current_price) * 100
        
        return {
            'current': float(current_price),
            'change': float(current_change),
            'changePercent': float(current_change_percent),
            'predicted': float(predicted_price),
            'predictedChange': float(predicted_change),
            'predictedPercent': float(predicted_change_percent),
            'confidence': float(confidence * 100),
            'timestamp': datetime.now().isoformat()
        }

# 创建预测器实例
predictor = StockPredictor()

@app.route('/')
def index():
    """提供前端页面"""
    from flask import render_template
    return render_template('index.html')

@app.route('/api/indices')
def get_indices():
    """获取所有支持的指数列表"""
    return jsonify(STOCK_INDICES)

@app.route('/api/predict/<index_code>')
def predict_index(index_code):
    """预测指定指数"""
    if index_code not in STOCK_INDICES:
        return jsonify({'error': '不支持的指数代码'}), 400
    
    try:
        symbol = STOCK_INDICES[index_code]['symbol']
        prediction = predictor.predict(symbol)
        
        if prediction is None:
            return jsonify({'error': '无法获取预测数据'}), 500
        
        result = {
            'code': index_code,
            'name': STOCK_INDICES[index_code]['name'],
            'symbol': STOCK_INDICES[index_code]['code'],
            **prediction
        }
        
        return jsonify(result)
        
    except Exception as e:
        return jsonify({'error': f'预测失败: {str(e)}'}), 500

@app.route('/api/predict/all')
def predict_all():
    """预测所有指数"""
    results = {}
    
    for index_code, index_info in STOCK_INDICES.items():
        try:
            symbol = index_info['symbol']
            prediction = predictor.predict(symbol)
            
            if prediction is not None:
                results[index_code] = {
                    'code': index_code,
                    'name': index_info['name'],
                    'symbol': index_info['code'],
                    **prediction
                }
            else:
                results[index_code] = {
                    'code': index_code,
                    'name': index_info['name'],
                    'symbol': index_info['code'],
                    'error': '无法获取数据'
                }
                
        except Exception as e:
            results[index_code] = {
                'code': index_code,
                'name': index_info['name'],
                'symbol': index_info['code'],
                'error': str(e)
            }
    
    return jsonify(results)

@app.route('/api/history/<index_code>')
def get_history(index_code):
    """获取历史数据"""
    if index_code not in STOCK_INDICES:
        return jsonify({'error': '不支持的指数代码'}), 400
    
    try:
        symbol = STOCK_INDICES[index_code]['symbol']
        period = request.args.get('period', '1mo')  # 默认1个月
        
        data = predictor.get_stock_data(symbol, period)
        if data is None:
            return jsonify({'error': '无法获取历史数据'}), 500
        
        # 转换数据格式
        history = []
        for date, row in data.iterrows():
            history.append({
                'date': date.strftime('%Y-%m-%d'),
                'open': float(row['Open']),
                'high': float(row['High']),
                'low': float(row['Low']),
                'close': float(row['Close']),
                'volume': int(row['Volume'])
            })
        
        return jsonify({
            'code': index_code,
            'name': STOCK_INDICES[index_code]['name'],
            'history': history
        })
        
    except Exception as e:
        return jsonify({'error': f'获取历史数据失败: {str(e)}'}), 500

@app.route('/api/status')
def get_status():
    """获取API状态"""
    return jsonify({
        'status': 'running',
        'timestamp': datetime.now().isoformat(),
        'ml_available': ML_AVAILABLE,
        'supported_indices': len(STOCK_INDICES),
        'cache_duration': predictor.cache_duration
    })

if __name__ == '__main__':
    import os
    
    # 获取环境变量
    debug_mode = os.getenv('FLASK_ENV', 'development') == 'development'
    port = int(os.getenv('PORT', 9000))
    host = os.getenv('HOST', '0.0.0.0')
    
    print("🚀 启动智投预测 - AI股市指数预测平台")
    print("=" * 50)
    print(f"环境: {os.getenv('FLASK_ENV', 'development')}")
    print(f"机器学习库: {'✅ 可用' if ML_AVAILABLE else '❌ 不可用'}")
    print(f"支持指数: {len(STOCK_INDICES)}个 {list(STOCK_INDICES.keys())}")
    print(f"服务地址: http://{host}:{port}")
    print("=" * 50)
    print("📚 API接口:")
    print(f"  🏠 主页: http://{host}:{port}")
    print(f"  📊 状态: http://{host}:{port}/api/status")
    print(f"  📈 预测: http://{host}:{port}/api/predict/all")
    print(f"  📋 指数: http://{host}:{port}/api/indices")
    print("=" * 50)
    
    try:
        app.run(debug=debug_mode, host=host, port=port)
    except Exception as e:
        print(f"❌ 服务启动失败: {e}")
        exit(1)

"""
API路由定义
"""

from fastapi import APIRouter, HTTPException
from datetime import datetime
from typing import Dict, Any
import random

from app.core.config import settings

router = APIRouter()

@router.get("/indices")
async def get_indices():
    """获取支持的股票指数列表"""
    return {
        "code": 200,
        "message": "success",
        "data": settings.STOCK_INDICES,
        "timestamp": datetime.utcnow().isoformat() + "Z"
    }

@router.get("/predict/{index_code}")
async def predict_index(index_code: str):
    """获取指定指数的预测数据"""
    
    if index_code not in settings.STOCK_INDICES:
        raise HTTPException(status_code=404, detail="Index not found")
    
    index_info = settings.STOCK_INDICES[index_code]
    
    # 模拟预测数据 - 实际项目中这里会调用ML模型
    current_price = 3000 + random.uniform(-100, 100)
    predicted_price = current_price + random.uniform(-50, 50)
    confidence = random.uniform(60, 90)
    
    return {
        "code": 200,
        "message": "success",
        "data": {
            "code": index_code,
            "name": index_info["name"],
            "symbol": index_info["symbol"],
            "market": index_info["market"],
            "current": round(current_price, 2),
            "change": round(predicted_price - current_price, 2),
            "changePercent": round((predicted_price - current_price) / current_price * 100, 2),
            "predicted": round(predicted_price, 2),
            "predictedChange": round(predicted_price - current_price, 2),
            "predictedPercent": round((predicted_price - current_price) / current_price * 100, 2),
            "confidence": round(confidence, 1),
            "timestamp": datetime.utcnow().isoformat() + "Z"
        },
        "timestamp": datetime.utcnow().isoformat() + "Z"
    }

@router.get("/predict/all")
async def predict_all():
    """获取所有指数的预测数据"""
    
    results = {}
    
    for index_code in settings.STOCK_INDICES.keys():
        index_info = settings.STOCK_INDICES[index_code]
        
        # 模拟预测数据
        current_price = 3000 + random.uniform(-100, 100)
        predicted_price = current_price + random.uniform(-50, 50)
        confidence = random.uniform(60, 90)
        
        results[index_code] = {
            "code": index_code,
            "name": index_info["name"],
            "symbol": index_info["symbol"],
            "market": index_info["market"],
            "current": round(current_price, 2),
            "predicted": round(predicted_price, 2),
            "change": round(predicted_price - current_price, 2),
            "changePercent": round((predicted_price - current_price) / current_price * 100, 2),
            "confidence": round(confidence, 1)
        }
    
    return {
        "code": 200,
        "message": "success",
        "data": results,
        "timestamp": datetime.utcnow().isoformat() + "Z"
    }

@router.get("/history/{index_code}")
async def get_history(index_code: str, period: str = "1mo"):
    """获取历史数据"""
    
    if index_code not in settings.STOCK_INDICES:
        raise HTTPException(status_code=404, detail="Index not found")
    
    index_info = settings.STOCK_INDICES[index_code]
    
    # 模拟历史数据
    history = []
    base_price = 3000
    
    for i in range(30):  # 30天历史数据
        date = datetime.utcnow().replace(day=i+1 if i < 28 else 28).date()
        price = base_price + random.uniform(-100, 100)
        
        history.append({
            "date": date.isoformat(),
            "open": round(price + random.uniform(-10, 10), 2),
            "high": round(price + random.uniform(0, 20), 2),
            "low": round(price - random.uniform(0, 20), 2),
            "close": round(price, 2),
            "volume": random.randint(1000000, 10000000)
        })
    
    return {
        "code": 200,
        "message": "success",
        "data": {
            "code": index_code,
            "name": index_info["name"],
            "period": period,
            "history": history
        },
        "timestamp": datetime.utcnow().isoformat() + "Z"
    }

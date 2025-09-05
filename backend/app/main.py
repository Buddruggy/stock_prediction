"""
智投预测 V2.0 - FastAPI 后端主应用
AI股市指数预测平台 - 纯API服务
"""

from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import JSONResponse
import uvicorn
import os
from datetime import datetime

from app.core.config import settings
from app.api import router as api_router
from app.core.logging import setup_logging

# 设置日志
setup_logging()

# 创建FastAPI应用
app = FastAPI(
    title=settings.APP_NAME,
    description=settings.APP_DESCRIPTION,
    version=settings.APP_VERSION,
    docs_url="/docs",
    redoc_url="/redoc",
    openapi_url="/openapi.json"
)

# CORS中间件配置
app.add_middleware(
    CORSMiddleware,
    allow_origins=settings.ALLOWED_ORIGINS,
    allow_credentials=True,
    allow_methods=["GET", "POST", "PUT", "DELETE", "OPTIONS"],
    allow_headers=["*"],
)

# 全局异常处理
@app.exception_handler(HTTPException)
async def http_exception_handler(request, exc):
    return JSONResponse(
        status_code=exc.status_code,
        content={
            "code": exc.status_code,
            "message": exc.detail,
            "timestamp": datetime.utcnow().isoformat() + "Z"
        }
    )

@app.exception_handler(Exception)
async def general_exception_handler(request, exc):
    return JSONResponse(
        status_code=500,
        content={
            "code": 500,
            "message": "Internal server error",
            "timestamp": datetime.utcnow().isoformat() + "Z"
        }
    )

# 根路径
@app.get("/")
async def root():
    """API根路径"""
    return {
        "code": 200,
        "message": "智投预测 API V2.0",
        "data": {
            "name": settings.APP_NAME,
            "version": settings.APP_VERSION,
            "description": settings.APP_DESCRIPTION,
            "docs": "/docs",
            "health": "/health"
        },
        "timestamp": datetime.utcnow().isoformat() + "Z"
    }

# 健康检查
@app.get("/health")
async def health_check():
    """健康检查端点"""
    return {
        "code": 200,
        "message": "healthy",
        "data": {
            "status": "running",
            "timestamp": datetime.utcnow().isoformat() + "Z",
            "version": settings.APP_VERSION
        }
    }

# 注册API路由
app.include_router(api_router, prefix="/api/v1")

# 启动事件
@app.on_event("startup")
async def startup_event():
    """应用启动事件"""
    print("🚀 智投预测 API V2.0 启动成功")
    print(f"📊 环境: {settings.ENVIRONMENT}")
    print(f"🌐 文档地址: http://{settings.HOST}:{settings.PORT}/docs")

# 关闭事件
@app.on_event("shutdown")
async def shutdown_event():
    """应用关闭事件"""
    print("👋 智投预测 API V2.0 已关闭")

if __name__ == "__main__":
    uvicorn.run(
        "main:app",
        host=settings.HOST,
        port=settings.PORT,
        reload=settings.DEBUG,
        log_level="info"
    )

"""
API路由模块
"""

from fastapi import APIRouter
from .routes import router as routes_router

# 创建主路由
router = APIRouter()

# 注册子路由
router.include_router(routes_router, tags=["prediction"])

__all__ = ["router"]

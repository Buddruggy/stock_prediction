"""Configuration loaded from environment variables."""
import os


class Settings:
    # DeepSeek LLM (used ONLY for natural-language interpretation, not price prediction)
    deepseek_api_key: str = os.getenv("DEEPSEEK_API_KEY", "")
    deepseek_url: str = os.getenv("DEEPSEEK_API_URL", "https://api.deepseek.com/chat/completions")
    deepseek_model: str = os.getenv("DEEPSEEK_MODEL", "deepseek-chat")

    # Storage
    db_path: str = os.getenv("DB_PATH", os.path.join(os.path.dirname(__file__), "..", "predictions.db"))

    # Server
    port: int = int(os.getenv("PORT", "8001"))
    host: str = os.getenv("HOST", "0.0.0.0")

    # Scheduler
    predict_hour: int = int(os.getenv("PREDICT_HOUR", "15"))
    predict_minute: int = int(os.getenv("PREDICT_MINUTE", "10"))
    timezone: str = "Asia/Shanghai"

    # Data / model
    history_days: int = int(os.getenv("HISTORY_DAYS", "400"))  # trading days pulled for training
    train_window: int = int(os.getenv("TRAIN_WINDOW", "250"))  # rolling window size
    cache_ttl_seconds: int = int(os.getenv("CACHE_TTL", "300"))

    # Tencent K-line endpoint (REAL historical OHLCV)
    tencent_kline_url: str = "https://web.ifzq.gtimg.cn/appstock/app/fqkline/get"


settings = Settings()

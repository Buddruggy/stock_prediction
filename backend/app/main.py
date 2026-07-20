"""FastAPI entrypoint. API contract kept compatible with the existing Vue frontend."""
import logging
import os
import sys

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from app import config  # noqa: E402
from app.service import PredictionService  # noqa: E402
from app.scheduler import start_scheduler  # noqa: E402
from app.indices import INDICES  # noqa: E402


logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s %(levelname)s %(name)s: %(message)s",
)
logger = logging.getLogger("main")

app = FastAPI(title="gogotou prediction API", version="3.0.0")

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_methods=["*"],
    allow_headers=["*"],
)

service = PredictionService()
start_scheduler(service)


def _resp(code: int, message: str, data=None):
    return {
        "code": code,
        "message": message,
        "data": data,
        "timestamp": __import__("time").strftime("%Y-%m-%dT%H:%M:%SZ", __import__("time").gmtime()),
    }


@app.get("/health")
def health():
    return {"status": "running", "version": "3.0.0"}


@app.get("/")
def root():
    from fastapi.responses import RedirectResponse
    return RedirectResponse(url="/docs")


@app.get("/api/v1/indices")
def list_indices():
    return _resp(200, "success", INDICES)


@app.get("/api/v1/predict/all")
def predict_all():
    try:
        return _resp(200, "success", service.get_all_predictions())
    except Exception as e:  # noqa: BLE001
        logger.exception("predict/all failed")
        return _resp(500, str(e))


@app.get("/api/v1/prediction-stats")
def prediction_stats():
    return _resp(200, "success", service.get_stats())


@app.get("/api/v1/predict/history/all")
def prediction_history_all(days: int = 30):
    data = service.get_all_historical_predictions(days)
    return _resp(200, "success", data)


@app.get("/api/v1/predict/{index_code}")
def predict_one(index_code: str):
    res = service.get_prediction(index_code)
    if res is None:
        return _resp(404, "index not found")
    return _resp(200, "success", res)


@app.get("/api/v1/history/{index_code}")
def history(index_code: str, period: str = "1mo"):
    res = service.get_history(index_code, period)
    if res is None:
        return _resp(404, "index not found")
    return _resp(200, "success", res)


@app.get("/api/v1/predict/history/{index_code}")
def prediction_history(index_code: str, days: int = 30):
    if service.db is None:
        return _resp(500, "db not available")
    rows = service.db.get_historical_predictions(index_code, days)
    return _resp(200, "success", rows)


@app.get("/api/v1/predict/stats")
def stats():
    return _resp(200, "success", service.get_stats())


@app.post("/api/v1/predict/refresh")
def refresh():
    import threading

    threading.Thread(target=service.run_daily, daemon=True).start()
    return _resp(200, "daily prediction triggered")


if __name__ == "__main__":
    import uvicorn

    uvicorn.run(app, host=config.settings.host, port=config.settings.port)

"""Prediction service: orchestrates data -> indicators -> model -> (optional) LLM.

Holds an in-memory cache of the latest daily prediction per index so the API
is fast; heavy model training happens once per day (scheduler) or on first hit.
"""
import logging
import math
import time
from datetime import date

from . import config
from . import data as data_layer
from . import indicators as ind
from . import model as model_layer
from . import llm as llm_layer
from . import db as db_layer
from .indices import INDICES

logger = logging.getLogger("service")


class PredictionService:
    def __init__(self):
        self.db = None
        try:
            self.db = db_layer.Database()
        except Exception as e:  # noqa: BLE001
            logger.warning("DB init failed, running cache-only: %s", e)

        self.cache = {}          # code -> prediction dict
        self.cache_time = {}     # code -> timestamp
        self.cache_ttl = config.settings.cache_ttl_seconds

    # ---------- public API ----------
    def get_prediction(self, code: str):
        if code not in INDICES:
            return None
        cached = self._cached(code)
        if cached is not None:
            return cached
        return self.compute_single(code)

    def get_all_predictions(self):
        out = {}
        for code in INDICES:
            out[code] = self.get_prediction(code)
        return out

    def get_history(self, code: str, period: str = "1mo"):
        if code not in INDICES:
            return None
        symbol = INDICES[code]["symbol"]
        days = self._period_days(period)
        bars = data_layer.fetch_kline(symbol, limit=days)
        return [
            {
                "date": b["date"],
                "open": b["open"],
                "high": b["high"],
                "low": b["low"],
                "close": b["close"],
                "volume": b["volume"],
            }
            for b in bars
        ]

    def get_stats(self):
        if self.db:
            s = self.db.get_prediction_stats()
            return {
                "total_predictions": s.get("total_predictions", 0),
                "success_rate": s.get("avg_confidence", 0),
                "avg_confidence": s.get("avg_confidence", 0),
            }
        return {"total_predictions": 0, "success_rate": 0, "avg_confidence": 0}

    def get_all_historical_predictions(self, days: int = 30):
        """Return {index_code: [prediction records]} for the History page.

        Maps DB column names (current_price/...) to the snake/camel fields the
        Vue frontend actually reads (current/predicted/changePercent/...).
        """
        if self.db is None:
            return {}
        out = {}
        for code in INDICES:
            rows = self.db.get_historical_predictions(code, days)
            out[code] = [
                {
                    "code": code,
                    "current": r.get("current_price"),
                    "predicted": r.get("predicted_price"),
                    "change": (r.get("predicted_price") or 0) - (r.get("current_price") or 0),
                    "changePercent": (r.get("predicted_return") or 0) * 100,
                    "direction": r.get("direction"),
                    "confidence": r.get("confidence"),
                    "prediction_date": r.get("prediction_date"),
                    "timestamp": r.get("prediction_date"),
                    "reasoning": r.get("reasoning"),
                }
                for r in rows
            ]
        return out

    def run_daily(self):
        logger.info("=== daily prediction run start ===")
        for code in INDICES:
            try:
                pred = self.compute_single(code, save=True)
                logger.info(
                    "predicted %s: cur=%.2f pred=%.2f ret=%+.2f%% conf=%.1f",
                    code, pred["current"], pred["predicted"],
                    pred["predicted_return"] * 100, pred["confidence"],
                )
            except Exception as e:  # noqa: BLE001
                logger.error("daily prediction failed for %s: %s", code, e)
        self.validate_previous()
        logger.info("=== daily prediction run done ===")

    # ---------- core computation ----------
    def compute_single(self, code: str, save: bool = False):
        symbol = INDICES[code]["symbol"]
        name = INDICES[code]["name"]

        bars = data_layer.fetch_kline(symbol)  # full ~3.5y real history
        df = ind.compute_indicators(bars)
        latest = df.iloc[-1]
        current_price = float(latest["close"])

        m = model_layer.run_model(df, train_window=config.settings.train_window)
        pred_return = m["predicted_return"]
        prob_up = m["prob_up"]
        predicted_price = current_price * (1 + pred_return)

        # Direction: primarily from the regression head; fall back to the
        # classifier when the regression signal is essentially flat.
        if abs(pred_return) < 0.0005:
            direction = "up" if prob_up >= 0.5 else "down"
        else:
            direction = "up" if pred_return > 0 else "down"
        # Probability that the predicted direction is correct (always aligned
        # with `direction`, to avoid contradictory UI signals).
        direction_prob = prob_up if direction == "up" else (1 - prob_up)

        bt = m["backtest"]
        confidence = round(bt["direction_accuracy"] * 100, 1)

        ind_latest = self._indicator_snapshot(latest)

        llm_out = llm_layer.generate_interpretation(
            index_name=name,
            current_price=current_price,
            predicted_price=predicted_price,
            predicted_return=pred_return,
            prob_up=prob_up,
            indicators=ind_latest,
            backtest=bt,
        )
        reasoning = (llm_out or {}).get("reasoning", "")
        llm_direction = (llm_out or {}).get("llm_direction")
        llm_conf = (llm_out or {}).get("llm_confidence")

        result = {
            "code": code,
            "name": name,
            "symbol": symbol,
            "market": INDICES[code]["market"],
            "current": round(current_price, 2),
            "predicted": round(predicted_price, 2),
            "change": round(predicted_price - current_price, 2),
            "changePercent": round(pred_return * 100, 2),
            "predicted_return": pred_return,
            "prob_up": round(prob_up, 4),
            "direction_prob": round(direction_prob, 4),
            "direction": direction,
            "confidence": confidence,
            "technical_indicators": ind_latest,
            "reasoning": reasoning,
            "llm_direction": llm_direction,
            "llm_confidence": llm_conf,
            "backtest": bt,
            "feature_importance": m["feature_importance"],
            "prediction_date": date.today().isoformat(),
            "timestamp": time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime()),
        }

        self.cache[code] = result
        self.cache_time[code] = time.time()

        if save and self.db:
            self.db.save_prediction({
                "index_code": code,
                "prediction_date": result["prediction_date"],
                "current_price": result["current"],
                "predicted_price": result["predicted"],
                "predicted_return": result["predicted_return"],
                "prob_up": result["prob_up"],
                "confidence": result["confidence"],
                "direction": result["direction"],
                "reasoning": result["reasoning"],
                "llm_direction": result["llm_direction"],
                "llm_confidence": result["llm_confidence"],
                "backtest_acc": bt["direction_accuracy"],
            })
            self.db.save_backtest(code, bt)

        return result

    # ---------- helpers ----------
    def validate_previous(self):
        """Use today's REAL price to check yesterday's directional call."""
        if not self.db:
            return
        yest = date.fromordinal(date.today().toordinal() - 1).isoformat()
        for code in INDICES:
            prev = self.db.get_today_prediction(code, day=yest)
            if not prev:
                continue
            try:
                bars = data_layer.fetch_kline(INDICES[code]["symbol"], limit=10)
                today_close = bars[-1]["close"]
                actual_ret = today_close / prev["current_price"] - 1
                predicted_up = prev["predicted_return"] > 0
                actual_up = actual_ret > 0
                ok = predicted_up == actual_up
                logger.info(
                    "verify %s (pred %s): predicted_ret=%+.2f%% actual_ret=%+.2f%% -> %s",
                    code, yest, prev["predicted_return"] * 100,
                    actual_ret * 100, "CORRECT" if ok else "wrong",
                )
            except Exception as e:  # noqa: BLE001
                logger.error("verify failed %s: %s", code, e)

    def _cached(self, code: str):
        t = self.cache_time.get(code)
        if t and (time.time() - t) < self.cache_ttl and code in self.cache:
            return self.cache[code]
        return None

    def _indicator_snapshot(self, row):
        return {
            "ma5": self._r(row.get("ma5")),
            "ma10": self._r(row.get("ma10")),
            "ma20": self._r(row.get("ma20")),
            "rsi": self._r(row.get("rsi")),
            "macd_hist": self._r(row.get("macd_hist")),
            "bb_position": self._r(row.get("bb_position")),
            "momentum20": self._r(row.get("momentum20")),
            "volatility20": self._r(row.get("volatility20")),
            "vol_ratio": self._r(row.get("vol_ratio")),
        }

    @staticmethod
    def _r(v, nd=4):
        try:
            return round(float(v), nd)
        except (TypeError, ValueError):
            return None

    @staticmethod
    def _period_days(period: str) -> int:
        return {
            "1d": 5, "5d": 10, "1mo": 30, "3mo": 90,
            "6mo": 180, "1y": 365,
        }.get(period, 30)

"""SQLite persistence for predictions and backtest accuracy."""
import logging
import os
import sqlite3
from datetime import datetime, date

from . import config

logger = logging.getLogger("db")


class Database:
    def __init__(self, path: str = None):
        self.path = os.path.abspath(path or config.settings.db_path)
        self._conn = sqlite3.connect(self.path, check_same_thread=False)
        self._conn.row_factory = sqlite3.Row
        self._init()

    def _init(self):
        self._conn.execute(
            """CREATE TABLE IF NOT EXISTS predictions (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                index_code TEXT NOT NULL,
                prediction_date TEXT NOT NULL,
                current_price REAL,
                predicted_price REAL,
                predicted_return REAL,
                prob_up REAL,
                confidence REAL,
                direction TEXT,
                reasoning TEXT,
                llm_direction TEXT,
                llm_confidence REAL,
                backtest_acc REAL,
                created_at TEXT DEFAULT CURRENT_TIMESTAMP
            )"""
        )
        self._conn.execute(
            """CREATE TABLE IF NOT EXISTS backtests (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                index_code TEXT NOT NULL,
                run_date TEXT NOT NULL,
                direction_accuracy REAL,
                recent_direction_accuracy REAL,
                mae REAL,
                rmse REAL,
                samples INTEGER,
                created_at TEXT DEFAULT CURRENT_TIMESTAMP
            )"""
        )
        self._conn.commit()

    # ---- predictions ----
    def save_prediction(self, rec: dict):
        self._conn.execute(
            """INSERT INTO predictions (
                index_code, prediction_date, current_price, predicted_price,
                predicted_return, prob_up, confidence, direction, reasoning,
                llm_direction, llm_confidence, backtest_acc
            ) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)""",
            (
                rec["index_code"], rec["prediction_date"], rec["current_price"],
                rec["predicted_price"], rec["predicted_return"], rec["prob_up"],
                rec["confidence"], rec["direction"], rec.get("reasoning"),
                rec.get("llm_direction"), rec.get("llm_confidence"),
                rec.get("backtest_acc"),
            ),
        )
        self._conn.commit()

    def get_today_prediction(self, index_code: str, day: str = None):
        day = day or date.today().isoformat()
        row = self._conn.execute(
            "SELECT * FROM predictions WHERE index_code=? AND prediction_date=? ORDER BY id DESC LIMIT 1",
            (index_code, day),
        ).fetchone()
        return dict(row) if row else None

    def get_all_today_predictions(self, day: str = None):
        day = day or date.today().isoformat()
        rows = self._conn.execute(
            "SELECT * FROM predictions WHERE prediction_date=?", (day,)
        ).fetchall()
        return {r["index_code"]: dict(r) for r in rows}

    def get_historical_predictions(self, index_code: str, days: int = 30):
        rows = self._conn.execute(
            """SELECT * FROM predictions WHERE index_code=? 
               AND prediction_date >= date('now', ?) ORDER BY prediction_date""",
            (index_code, f"-{days} days"),
        ).fetchall()
        return [dict(r) for r in rows]

    # ---- backtests ----
    def save_backtest(self, index_code: str, bt: dict):
        self._conn.execute(
            """INSERT INTO backtests (
                index_code, run_date, direction_accuracy,
                recent_direction_accuracy, mae, rmse, samples
            ) VALUES (?,?,?,?,?,?,?)""",
            (
                index_code, date.today().isoformat(),
                bt.get("direction_accuracy"), bt.get("recent_direction_accuracy"),
                bt.get("mae"), bt.get("rmse"), bt.get("samples"),
            ),
        )
        self._conn.commit()

    def get_prediction_stats(self):
        row = self._conn.execute(
            "SELECT COUNT(*) AS n, AVG(confidence) AS avg_conf FROM predictions"
        ).fetchone()
        return {"total_predictions": row["n"], "avg_confidence": row["avg_conf"]}

    def close(self):
        self._conn.close()

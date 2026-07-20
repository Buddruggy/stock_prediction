"""Real historical K-line data from Tencent finance API.

Replaces the old `generateHistoryFromCurrent` random generator.
Endpoint returns REAL daily OHLCV: [date, open, close, high, low, volume].
A single request can return ~850 rows (~3.5y of trading days).
"""
import time
import logging
from datetime import datetime, timedelta

import requests

from . import config

logger = logging.getLogger("data")

_CACHE = {}
_CACHE_TTL = config.settings.cache_ttl_seconds


def _cache_get(key):
    item = _CACHE.get(key)
    if item and time.time() < item["expires"]:
        return item["data"]
    return None


def _cache_set(key, data):
    _CACHE[key] = {"data": data, "expires": time.time() + _CACHE_TTL}


def fetch_kline(symbol: str, limit: int = None) -> list[dict]:
    """Fetch REAL daily K-line (qfq / forward-adjusted) for a Tencent symbol.

    Pulls ~3.5y of history (Tencent returns up to ~850 rows per request),
    then optionally returns only the most recent `limit` bars.
    Returns list of {date, open, close, high, low, volume} sorted ascending.
    """
    end = datetime.now()
    start = end - timedelta(days=1300)  # ~3.5y, guarantees >=2y of trading days
    start_s = start.strftime("%Y-%m-%d")
    end_s = end.strftime("%Y-%m-%d")

    cache_key = f"{symbol}_full"
    data = _cache_get(cache_key)
    if data is None:
        url = f"{config.settings.tencent_kline_url}?param={symbol},day,{start_s},{end_s},1200,qfq"
        try:
            resp = requests.get(
                url,
                timeout=20,
                headers={
                    "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
                    "Referer": "http://gu.qq.com",
                },
            )
            resp.raise_for_status()
            payload = resp.json()
        except Exception as e:  # noqa: BLE001
            logger.error("tencent kline request failed for %s: %s", symbol, e)
            raise

        node = (payload.get("data") or {}).get(symbol)
        if not node:
            raise ValueError(f"no data returned for {symbol}")

        rows = node.get("qfqday") or node.get("day") or []
        data = []
        for r in rows:
            if len(r) < 6:
                continue
            try:
                data.append(
                    {
                        "date": r[0],
                        "open": float(r[1]),
                        "close": float(r[2]),
                        "high": float(r[3]),
                        "low": float(r[4]),
                        "volume": float(r[5]),
                    }
                )
            except (ValueError, TypeError):
                continue

        if not data:
            raise ValueError(f"empty parsed kline for {symbol}")

        data.sort(key=lambda x: x["date"])
        _cache_set(cache_key, data)
        logger.info("fetched %d real kline rows for %s", len(data), symbol)
    else:
        logger.info("kline cache hit: %s (%d rows)", symbol, len(data))

    if limit:
        data = data[-limit:]
    return data


def get_latest_quote(symbol: str) -> dict:
    data = fetch_kline(symbol, limit=10)
    if not data:
        raise ValueError(f"no quote for {symbol}")
    return data[-1]

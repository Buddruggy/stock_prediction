"""Technical indicators computed from REAL K-line data (pandas)."""
import pandas as pd
import numpy as np


def compute_indicators(bars: list[dict]) -> pd.DataFrame:
    """Compute a rich feature set from real OHLCV bars.

    Returns a DataFrame with one row per bar (sorted ascending by date).
    The last row keeps its `next_return` as NaN (no future yet) — callers
    must drop rows with a missing label before training, but may still use
    the final row's features for prediction.
    """
    df = pd.DataFrame(bars)
    df["date"] = pd.to_datetime(df["date"])
    df = df.sort_values("date").reset_index(drop=True)

    close = df["close"]
    high = df["high"]
    low = df["low"]
    vol = df["volume"]

    # Moving averages
    df["ma5"] = close.rolling(5).mean()
    df["ma10"] = close.rolling(10).mean()
    df["ma20"] = close.rolling(20).mean()

    # RSI(14)
    delta = close.diff()
    gain = delta.clip(lower=0)
    loss = -delta.clip(upper=0)
    avg_gain = gain.rolling(14).mean()
    avg_loss = loss.rolling(14).mean()
    rs = avg_gain / avg_loss.replace(0, np.nan)
    df["rsi"] = (100 - 100 / (1 + rs)).fillna(50)

    # MACD(12,26,9)
    ema12 = close.ewm(span=12, adjust=False).mean()
    ema26 = close.ewm(span=26, adjust=False).mean()
    df["macd"] = ema12 - ema26
    df["macd_signal"] = df["macd"].ewm(span=9, adjust=False).mean()
    df["macd_hist"] = df["macd"] - df["macd_signal"]

    # Bollinger Bands(20, 2)
    bb_mid = close.rolling(20).mean()
    bb_std = close.rolling(20).std()
    df["bb_upper"] = bb_mid + 2 * bb_std
    df["bb_lower"] = bb_mid - 2 * bb_std
    df["bb_position"] = (close - df["bb_lower"]) / (df["bb_upper"] - df["bb_lower"]).replace(0, np.nan)
    df["bb_position"] = df["bb_position"].fillna(0.5)

    # Momentum (price relative to N days ago)
    df["momentum5"] = close / close.shift(5) - 1
    df["momentum10"] = close / close.shift(10) - 1
    df["momentum20"] = close / close.shift(20) - 1

    # Volatility (rolling std of daily returns)
    ret = close.pct_change()
    df["volatility20"] = ret.rolling(20).std()
    df["return"] = ret

    # Volume ratio (vs 5-day average)
    df["vol_ratio"] = vol / vol.rolling(5).mean()

    # Position of price vs MAs (normalized)
    df["price_vs_ma20"] = close / df["ma20"] - 1

    # Label: NEXT day return & direction (what we want to predict)
    df["next_return"] = close.shift(-1) / close - 1
    df["direction"] = (df["next_return"] > 0).astype(int)

    # fill trivial infinities / NaN from divisions where appropriate
    df = df.replace([np.inf, -np.inf], np.nan)
    return df


# Feature columns used by the model
FEATURES = [
    "ma5", "ma10", "ma20",
    "rsi", "macd", "macd_hist",
    "bb_position",
    "momentum5", "momentum10", "momentum20",
    "volatility20", "vol_ratio",
    "return", "price_vs_ma20",
]

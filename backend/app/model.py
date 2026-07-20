"""Machine-learning prediction + honest walk-forward backtest.

The model predicts next-day return / direction from REAL technical features.
DeepSeek is NOT used for the price — only for interpretation (see llm.py).
"""
import logging

import numpy as np
import pandas as pd
from sklearn.ensemble import GradientBoostingClassifier, GradientBoostingRegressor
from sklearn.preprocessing import StandardScaler

from .indicators import FEATURES

logger = logging.getLogger("model")


def run_model(
    df_with_indicators: pd.DataFrame,
    features=None,
    train_window: int = 250,
    step: int = 10,
) -> dict:
    """Train on real history, run a walk-forward backtest, and predict tomorrow.

    Args:
        df_with_indicators: ascending-sorted DataFrame from compute_indicators();
                            the LAST row has a NaN `next_return` (no future yet).
    Returns dict with predicted_return, prob_up, backtest metrics, feature_importance.
    """
    features = features or FEATURES
    train_df = (
        df_with_indicators.dropna(subset=features + ["next_return"])
        .reset_index(drop=True)
    )
    if len(train_df) < 60:
        raise ValueError(f"not enough training samples: {len(train_df)}")

    n = len(train_df)
    tw = min(train_window, n - 10)

    preds, actuals, dp, da = [], [], [], []
    for t in range(tw, n, step):
        if t >= n:
            break
        Xtr = train_df[features].iloc[t - tw : t]
        ytr = train_df["next_return"].iloc[t - tw : t]
        sc = StandardScaler().fit(Xtr)
        reg = GradientBoostingRegressor(
            n_estimators=120, max_depth=3, learning_rate=0.05,
            subsample=0.9, random_state=42,
        )
        reg.fit(sc.transform(Xtr), ytr)
        Xt = train_df[features].iloc[t : t + 1]
        pr = float(reg.predict(sc.transform(Xt))[0])
        ar = float(train_df["next_return"].iloc[t])
        preds.append(pr)
        actuals.append(ar)
        dp.append(1 if pr > 0 else 0)
        da.append(int(ar > 0))

    # Final model trained on ALL labelled history -> predict the latest bar
    sc = StandardScaler().fit(train_df[features])
    Xs = sc.transform(train_df[features])
    reg = GradientBoostingRegressor(
        n_estimators=120, max_depth=3, learning_rate=0.05,
        subsample=0.9, random_state=42,
    )
    reg.fit(Xs, train_df["next_return"])
    clf = GradientBoostingClassifier(
        n_estimators=120, max_depth=3, learning_rate=0.05, random_state=42
    )
    clf.fit(Xs, train_df["direction"])

    last_feat = df_with_indicators[features].iloc[[-1]]
    if last_feat.isna().any().any():
        raise ValueError("latest bar has missing features (need more history)")
    pred_return = float(reg.predict(sc.transform(last_feat))[0])
    prob_up = float(clf.predict_proba(sc.transform(last_feat))[0][1])

    arr_p = np.array(preds)
    arr_a = np.array(actuals)
    dir_acc = float(np.mean([p == a for p, a in zip(dp, da)]))
    mae = float(np.mean(np.abs(arr_p - arr_a)))
    rmse = float(np.sqrt(np.mean((arr_p - arr_a) ** 2)))
    k = max(1, len(dp) // 3)
    recent_dir_acc = float(np.mean([p == a for p, a in zip(dp[-k:], da[-k:])]))

    importances = {
        f: float(x) for f, x in zip(features, clf.feature_importances_)
    }

    return {
        "predicted_return": pred_return,
        "prob_up": prob_up,
        "backtest": {
            "direction_accuracy": round(dir_acc, 4),
            "recent_direction_accuracy": round(recent_dir_acc, 4),
            "mae": mae,
            "rmse": rmse,
            "samples": len(dp),
        },
        "feature_importance": importances,
    }

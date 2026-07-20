"""DeepSeek integration — used ONLY for natural-language interpretation.

Price prediction is done by the ML model (model.py). Here DeepSeek receives
the REAL indicators + the model's prediction and writes a human-readable
rationale + an independent directional sanity check. The API key is read from
DEEPSEEK_API_KEY and is NEVER hardcoded.
"""
import json
import logging

import requests

from . import config

logger = logging.getLogger("llm")


def generate_interpretation(
    index_name: str,
    current_price: float,
    predicted_price: float,
    predicted_return: float,
    prob_up: float,
    indicators: dict,
    backtest: dict,
    api_key: str = None,
) -> dict:
    """Return {reasoning, llm_agrees, llm_confidence} or None if no key configured."""
    api_key = api_key or config.settings.deepseek_api_key
    if not api_key:
        logger.info("no DeepSeek key configured — skipping interpretation")
        return None

    prompt = _build_prompt(
        index_name, current_price, predicted_price, predicted_return,
        prob_up, indicators, backtest,
    )
    try:
        resp = requests.post(
            config.settings.deepseek_url,
            timeout=60,
            headers={
                "Authorization": f"Bearer {api_key}",
                "Content-Type": "application/json",
            },
            json={
                "model": config.settings.deepseek_model,
                "messages": [
                    {
                        "role": "system",
                        "content": (
                            "你是一名严谨的中国股市量化分析师。基于提供的真实技术面数据"
                            "和模型的统计预测，用中文给出简洁、克制、可解释的解读，"
                            "并独立判断明日方向。不要夸大，不要承诺收益。"
                            "只输出 JSON："
                            '{"reasoning":"...","direction":"up/down/flat",'
                            '"confidence":0到100的整数}。'
                        ),
                    },
                    {"role": "user", "content": prompt},
                ],
                "max_tokens": 600,
                "temperature": 0.3,
                "stream": False,
            },
        )
        resp.raise_for_status()
        content = resp.json()["choices"][0]["message"]["content"]
        return _parse(content)
    except Exception as e:  # noqa: BLE001
        logger.error("DeepSeek call failed: %s", e)
        return None


def _build_prompt(name, cur, pred, ret, prob_up, ind, bt):
    return f"""指数：{name}
当前价：{cur:.2f}
模型预测明日价：{pred:.2f}（预期收益率 {ret*100:+.2f}%，模型看多概率 {prob_up*100:.1f}%）

真实技术指标：
- MA5/MA10/MA20：{ind.get('ma5'):.2f} / {ind.get('ma10'):.2f} / {ind.get('ma20'):.2f}
- RSI(14)：{ind.get('rsi'):.1f}
- MACD柱：{ind.get('macd_hist'):.2f}
- 布林带位置(0-1)：{ind.get('bb_position'):.2f}
- 20日动量：{ind.get('momentum20')*100:+.2f}%
- 20日波动率：{ind.get('volatility20')*100:.2f}%
- 量比：{ind.get('vol_ratio'):.2f}

历史回测（walk-forward，样本 {bt.get('samples')}）：
- 方向准确率：{bt.get('direction_accuracy')*100:.1f}%
- 近1/3样本方向准确率：{bt.get('recent_direction_accuracy')*100:.1f}%

请结合以上真实数据给出明日走势的中文解读，并独立判断方向（up/down/flat）。"""


def _parse(content: str):
    try:
        obj = json.loads(content)
    except json.JSONDecodeError:
        s, e = content.find("{"), content.rfind("}")
        if s != -1 and e != -1:
            try:
                obj = json.loads(content[s : e + 1])
            except json.JSONDecodeError:
                return None
        else:
            return None
    direction = str(obj.get("direction", "")).lower()
    llm_agrees = None
    if direction in ("up", "涨", "看涨"):
        llm_agrees = "up"
    elif direction in ("down", "跌", "看跌"):
        llm_agrees = "down"
    else:
        llm_agrees = "flat"
    conf = obj.get("confidence", 50)
    try:
        conf = max(0, min(100, int(conf)))
    except (TypeError, ValueError):
        conf = 50
    return {
        "reasoning": obj.get("reasoning", ""),
        "llm_direction": llm_agrees,
        "llm_confidence": conf,
    }

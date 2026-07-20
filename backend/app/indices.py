"""Supported stock indices. Tencent K-line symbol == index code (e.g. sh000001)."""

INDICES = {
    "sh000001": {
        "code": "sh000001",
        "name": "上证综指",
        "symbol": "sh000001",
        "market": "上海证券交易所",
    },
    "sz399001": {
        "code": "sz399001",
        "name": "深证成指",
        "symbol": "sz399001",
        "market": "深圳证券交易所",
    },
    "sz399006": {
        "code": "sz399006",
        "name": "创业板指",
        "symbol": "sz399006",
        "market": "深圳证券交易所",
    },
    "sh000688": {
        "code": "sh000688",
        "name": "科创50",
        "symbol": "sh000688",
        "market": "上海证券交易所",
    },
}


def get_index(code: str):
    return INDICES.get(code)

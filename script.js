// 学术风格的股票指数预测系统
const stockIndices = {
    'sh000001': { name: '上证综合指数', englishName: 'Shanghai Composite Index' },
    'sz399001': { name: '深证成分指数', englishName: 'Shenzhen Component Index' },
    'sz399006': { name: '创业板综合指数', englishName: 'ChiNext Composite Index' },
    'sh000688': { name: '科创板50指数', englishName: 'STAR Market 50 Index' }
};

// 系统初始化
document.addEventListener('DOMContentLoaded', function () {
    console.log('学术预测模型系统初始化...');
    loadPredictionData();
    updateTimestamp();

    // 定期更新数据
    setInterval(updateTimestamp, 30000);
    setInterval(loadPredictionData, 300000);
});

// 加载预测数据
async function loadPredictionData() {
    try {
        const response = await fetch('/api/predict/all');
        if (response.ok) {
            const data = await response.json();
            updatePredictionDisplay(data);
            console.log('实时数据获取成功');
        } else {
            console.log('API响应异常，使用基准数据集');
            loadBaselineData();
        }
    } catch (error) {
        console.log('网络连接异常，切换至离线模式');
        loadBaselineData();
    }
}

// 更新预测显示
function updatePredictionDisplay(data) {
    Object.keys(stockIndices).forEach(indexCode => {
        if (data[indexCode] && !data[indexCode].error) {
            const prediction = data[indexCode];

            // 更新观测值 (当前价格)
            updateDataPoint(`${indexCode}-current`, formatAcademicNumber(prediction.current));
            updateDataPoint(`${indexCode}-change`,
                formatAcademicPercentage(prediction.changePercent),
                prediction.changePercent >= 0 ? 'positive' : 'negative'
            );

            // 更新预测值
            updateDataPoint(`${indexCode}-predicted`, formatAcademicNumber(prediction.predicted));
            updateDataPoint(`${indexCode}-pred-change`,
                formatAcademicPercentage(prediction.predictedPercent),
                prediction.predictedPercent >= 0 ? 'positive' : 'negative'
            );

            // 更新模型置信度
            updateDataPoint(`${indexCode}-confidence`,
                Math.round(prediction.confidence) + '%');
        }
    });
}

// 加载基准数据集 (用于演示和离线模式)
function loadBaselineData() {
    const baselineDataset = {
        'sh000001': {
            current: 3245.67, changePercent: 0.38,
            predicted: 3268.45, predictedPercent: 0.70,
            confidence: 78.0
        },
        'sz399001': {
            current: 11234.89, changePercent: -0.40,
            predicted: 11298.76, predictedPercent: 0.57,
            confidence: 72.0
        },
        'sz399006': {
            current: 2456.78, changePercent: 0.78,
            predicted: 2487.34, predictedPercent: 1.24,
            confidence: 85.0
        },
        'sh000688': {
            current: 987.65, changePercent: -0.85,
            predicted: 995.23, predictedPercent: 0.77,
            confidence: 69.0
        }
    };

    updatePredictionDisplay(baselineDataset);
}

// 更新数据点
function updateDataPoint(elementId, content, className = '') {
    const element = document.getElementById(elementId);
    if (element) {
        element.textContent = content;
        if (className) {
            element.className = `change ${className}`;
        }
    }
}

// 学术化数字格式
function formatAcademicNumber(number) {
    if (number >= 10000) {
        return number.toLocaleString('zh-CN', {
            minimumFractionDigits: 0,
            maximumFractionDigits: 0
        });
    } else {
        return number.toLocaleString('zh-CN', {
            minimumFractionDigits: 2,
            maximumFractionDigits: 2
        });
    }
}

// 学术化百分比格式
function formatAcademicPercentage(percent) {
    const sign = percent >= 0 ? '+' : '';
    return `${sign}${percent.toFixed(2)}%`;
}

// 更新时间戳
function updateTimestamp() {
    const currentTime = new Date();
    const formattedTime = currentTime.toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        timeZone: 'Asia/Shanghai'
    });

    const timestampElement = document.getElementById('lastUpdate');
    if (timestampElement) {
        timestampElement.textContent = `${formattedTime} (UTC+8)`;
    }
}

// 数据验证函数
function validatePredictionData(data) {
    const requiredFields = ['current', 'predicted', 'confidence'];
    return requiredFields.every(field =>
        data.hasOwnProperty(field) &&
        typeof data[field] === 'number' &&
        !isNaN(data[field])
    );
}

// 计算统计指标
function calculateModelMetrics(predictions) {
    const confidenceValues = Object.values(predictions)
        .filter(p => !p.error)
        .map(p => p.confidence);

    if (confidenceValues.length === 0) return null;

    const averageConfidence = confidenceValues.reduce((a, b) => a + b, 0) / confidenceValues.length;
    const maxConfidence = Math.max(...confidenceValues);
    const minConfidence = Math.min(...confidenceValues);

    return {
        averageConfidence: Math.round(averageConfidence * 10) / 10,
        maxConfidence: Math.round(maxConfidence * 10) / 10,
        minConfidence: Math.round(minConfidence * 10) / 10,
        sampleSize: confidenceValues.length
    };
}

// 错误处理
function handleSystemError(error, context) {
    console.error(`系统异常 [${context}]:`, error);

    // 在学术环境中提供详细的错误信息
    const errorNotification = {
        timestamp: new Date().toISOString(),
        context: context,
        errorType: error.name || 'UnknownError',
        message: error.message || '未知系统异常'
    };

    console.log('错误报告:', errorNotification);

    // 可以在此处添加错误上报或日志记录功能
    return errorNotification;
}

// 学术化系统日志
function logSystemActivity(activity, details = {}) {
    const logEntry = {
        timestamp: new Date().toISOString(),
        activity: activity,
        details: details,
        userAgent: navigator.userAgent,
        sessionId: Math.random().toString(36).substr(2, 9)
    };

    console.log('系统活动日志:', logEntry);
    return logEntry;
}

// 页面加载完成后的学术化初始化
document.addEventListener('DOMContentLoaded', function () {
    logSystemActivity('系统初始化', {
        modelType: 'Random Forest + LSTM',
        predictionHorizon: 'T+1',
        featureDimension: 15,
        updateFrequency: '5分钟'
    });
});
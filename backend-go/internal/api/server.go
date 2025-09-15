package api

import (
	"log"
	"net/http"
	"stock-prediction-backend/internal/config"
	"stock-prediction-backend/internal/model"
	"stock-prediction-backend/internal/service"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Server API服务器
type Server struct {
	config      *config.Config
	dataService *service.DataService
	router      *gin.Engine
}

// NewServer 创建新的API服务器
func NewServer(cfg *config.Config) *Server {
	server := &Server{
		config:      cfg,
		dataService: service.NewDataService(cfg),
	}

	server.setupRouter()
	return server
}

// setupRouter 设置路由
func (s *Server) setupRouter() {
	// 设置Gin模式
	if s.config.LogLevel == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	s.router = gin.New()

	// 中间件
	s.router.Use(gin.Logger())
	s.router.Use(gin.Recovery())
	s.router.Use(s.corsMiddleware())

	// 健康检查
	s.router.GET("/health", s.healthCheck)

	// API路由组
	v1 := s.router.Group("/api/v1")
	{
		// 预测相关
		v1.GET("/predict/all", s.getAllPredictions)
		v1.GET("/predict/:index_code", s.getPrediction)

		// 历史预测数据
		v1.GET("/predict/history/all", s.getAllHistoricalPredictions)
		v1.GET("/predict/history/:index_code", s.getHistoricalPredictions)

		// 历史数据
		v1.GET("/history/:index_code", s.getHistoryData)

		// 指数信息
		v1.GET("/indices/all", s.getAllIndicesInfo)
		v1.GET("/indices/:index_code", s.getIndexInfo)

		// 数据源状态
		v1.GET("/data-source/status", s.getDataSourceStatus)

		// 预测缓存管理
		v1.GET("/prediction-cache/status", s.getPredictionCacheStatus)
		v1.POST("/prediction-cache/refresh", s.refreshPredictionCache)
		
		// 预测统计信息
		v1.GET("/prediction-stats", s.getPredictionStats)
	}
}

// Run 启动服务器
func (s *Server) Run() error {
	return s.router.Run(":" + s.config.Port)
}

// corsMiddleware CORS中间件
func (s *Server) corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// healthCheck 健康检查
func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, model.APIResponse{
		Code:    200,
		Message: "healthy",
		Data: map[string]interface{}{
			"status":    "running",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"version":   "2.0.0",
		},
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

// getAllPredictions 获取所有预测数据
func (s *Server) getAllPredictions(c *gin.Context) {
	predictions, err := s.dataService.GetAllPredictions()
	if err != nil {
		log.Printf("获取所有预测数据失败: %v", err)
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Code:      500,
			Message:   "Internal server error",
			Data:      nil,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse{
		Code:      200,
		Message:   "success",
		Data:      predictions,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

// getPrediction 获取指定指数的预测数据
func (s *Server) getPrediction(c *gin.Context) {
	indexCode := c.Param("index_code")

	prediction, err := s.dataService.GetPredictionData(indexCode)
	if err != nil {
		log.Printf("获取预测数据失败 %s: %v", indexCode, err)
		c.JSON(http.StatusNotFound, model.APIResponse{
			Code:      404,
			Message:   "Index not found",
			Data:      nil,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse{
		Code:      200,
		Message:   "success",
		Data:      prediction,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

// getHistoryData 获取历史数据
func (s *Server) getHistoryData(c *gin.Context) {
	indexCode := c.Param("index_code")
	period := c.DefaultQuery("period", "1mo")

	historyData, err := s.dataService.GetHistoryData(indexCode, period)
	if err != nil {
		log.Printf("获取历史数据失败 %s: %v", indexCode, err)
		c.JSON(http.StatusNotFound, model.APIResponse{
			Code:      404,
			Message:   "Index not found",
			Data:      nil,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse{
		Code:      200,
		Message:   "success",
		Data:      historyData,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

// getAllIndicesInfo 获取所有指数信息
func (s *Server) getAllIndicesInfo(c *gin.Context) {
	indicesInfo, err := s.dataService.GetAllIndicesInfo()
	if err != nil {
		log.Printf("获取所有指数信息失败: %v", err)
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Code:      500,
			Message:   "Internal server error",
			Data:      nil,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse{
		Code:      200,
		Message:   "success",
		Data:      indicesInfo,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

// getIndexInfo 获取指定指数的信息
func (s *Server) getIndexInfo(c *gin.Context) {
	indexCode := c.Param("index_code")

	indexInfo, err := s.dataService.GetIndexInfo(indexCode)
	if err != nil {
		log.Printf("获取指数信息失败 %s: %v", indexCode, err)
		c.JSON(http.StatusNotFound, model.APIResponse{
			Code:      404,
			Message:   "Index not found",
			Data:      nil,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse{
		Code:      200,
		Message:   "success",
		Data:      indexInfo,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

// getDataSourceStatus 获取数据源状态
func (s *Server) getDataSourceStatus(c *gin.Context) {
	status := s.dataService.GetDataSourceStatus()

	c.JSON(http.StatusOK, model.APIResponse{
		Code:      200,
		Message:   "success",
		Data:      status,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

// getPredictionCacheStatus 获取预测缓存状态
func (s *Server) getPredictionCacheStatus(c *gin.Context) {
	dailyPredictions, predictTime, hasPredictions := s.dataService.GetDailyPredictions()

	status := map[string]interface{}{
		"has_cache":       hasPredictions,
		"cache_count":     len(dailyPredictions),
		"predict_time":    "",
		"cache_age_hours": 0,
		"is_valid":        false,
	}

	if hasPredictions {
		status["predict_time"] = predictTime.Format("2006-01-02 15:04:05")
		cacheAge := time.Since(predictTime)
		status["cache_age_hours"] = int(cacheAge.Hours())
		status["is_valid"] = cacheAge < 24*time.Hour
	}

	c.JSON(http.StatusOK, model.APIResponse{
		Code:      200,
		Message:   "success",
		Data:      status,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

// refreshPredictionCache 手动刷新预测缓存
func (s *Server) refreshPredictionCache(c *gin.Context) {
	// 在后台异步执行预测任务
	go s.dataService.RefreshDailyPredictions()

	c.JSON(http.StatusOK, model.APIResponse{
		Code:      200,
		Message:   "预测缓存刷新任务已启动，请稍后检查状态",
		Data:      nil,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

// getPredictionStats 获取预测统计信息
func (s *Server) getPredictionStats(c *gin.Context) {
	stats, err := s.dataService.GetPredictionStats()
	if err != nil {
		log.Printf("获取预测统计信息失败: %v", err)
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Code:      500,
			Message:   "Internal server error",
			Data:      nil,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse{
		Code:      200,
		Message:   "success",
		Data:      stats,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

// getAllHistoricalPredictions 获取所有指数的历史预测数据
func (s *Server) getAllHistoricalPredictions(c *gin.Context) {
	days := 30 // 默认获取30天的历史数据
	if daysParam := c.Query("days"); daysParam != "" {
		if parsedDays, err := strconv.Atoi(daysParam); err == nil && parsedDays > 0 && parsedDays <= 365 {
			days = parsedDays
		}
	}

	historyData, err := s.dataService.GetAllHistoricalPredictions(days)
	if err != nil {
		log.Printf("获取所有历史预测数据失败: %v", err)
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Code:      500,
			Message:   "Internal server error",
			Data:      nil,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse{
		Code:      200,
		Message:   "success",
		Data:      historyData,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

// getHistoricalPredictions 获取指定指数的历史预测数据
func (s *Server) getHistoricalPredictions(c *gin.Context) {
	indexCode := c.Param("index_code")
	days := 30 // 默认获取30天的历史数据
	if daysParam := c.Query("days"); daysParam != "" {
		if parsedDays, err := strconv.Atoi(daysParam); err == nil && parsedDays > 0 && parsedDays <= 365 {
			days = parsedDays
		}
	}

	historyData, err := s.dataService.GetHistoricalPredictions(indexCode, days)
	if err != nil {
		log.Printf("获取历史预测数据失败 %s: %v", indexCode, err)
		c.JSON(http.StatusNotFound, model.APIResponse{
			Code:      404,
			Message:   "Historical predictions not found",
			Data:      nil,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse{
		Code:      200,
		Message:   "success",
		Data:      historyData,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

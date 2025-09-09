package main

import (
	"log"
	"stock-prediction-backend/internal/api"
	"stock-prediction-backend/internal/config"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 创建API服务器
	server := api.NewServer(cfg)

	// 启动服务器
	log.Printf("🚀 启动股票预测后端服务，端口: %s", cfg.Port)
	if err := server.Run(); err != nil {
		log.Fatalf("❌ 服务器启动失败: %v", err)
	}
}

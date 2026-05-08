package main

import (
	"fmt"
	"log"

	grpcclients "foodRatingSystem/api-gateway/grpc-clients"
	"foodRatingSystem/api-gateway/router"
	"foodRatingSystem/shared/config"
)

func main() {
	config.LoadConfig()

	err := grpcclients.InitClients()
	if err != nil {
		log.Fatalf("gRPC 客户端初始化失败: %v", err)
	}

	r := router.SetupRouter()
	fmt.Printf("API 网关已启动，监听端口 :%s\n", config.AppConfig.ServerPort)
	r.Run(":" + config.AppConfig.ServerPort)
}

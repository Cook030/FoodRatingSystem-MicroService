package grpc_clients

import (
	"context"
	"fmt"
	"log"
	"time"

	"foodRatingSystem/proto/rating"
	"foodRatingSystem/proto/restaurant"
	"foodRatingSystem/proto/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type ServiceConfig struct {
	Name    string
	Address string
	Service string
}

var (
	UserClient       user.UserServiceClient
	RestaurantClient restaurant.RestaurantServiceClient
	RatingClient     rating.RatingServiceClient

	services []ServiceConfig
	conns    []*grpc.ClientConn
)

func InitClients() error {
	services = []ServiceConfig{
		{Name: "user-service", Address: "localhost:50051", Service: "user.UserService"},
		{Name: "restaurant-service", Address: "localhost:50052", Service: "restaurant.RestaurantService"},
		{Name: "rating-service", Address: "localhost:50053", Service: "rating.RatingService"},
	}

	for _, svc := range services {
		conn, err := grpc.NewClient(svc.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return fmt.Errorf("连接 %s 失败: %v", svc.Name, err)
		}
		conns = append(conns, conn)
	}

	UserClient = user.NewUserServiceClient(conns[0])
	RestaurantClient = restaurant.NewRestaurantServiceClient(conns[1])
	RatingClient = rating.NewRatingServiceClient(conns[2])

	go startHealthCheckLoop()

	return nil
}

type ServiceHealth struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

func CheckAllServices() []ServiceHealth {
	var results []ServiceHealth

	for i, svc := range services {
		status := checkService(conns[i], svc.Service)
		results = append(results, ServiceHealth{
			Name:   svc.Name,
			Status: status,
		})
	}

	return results
}

func checkService(conn *grpc.ClientConn, service string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	client := grpc_health_v1.NewHealthClient(conn)
	resp, err := client.Check(ctx, &grpc_health_v1.HealthCheckRequest{
		Service: service,
	})
	if err != nil {
		return "NOT_SERVING"
	}

	switch resp.GetStatus() {
	case grpc_health_v1.HealthCheckResponse_SERVING:
		return "SERVING"
	case grpc_health_v1.HealthCheckResponse_NOT_SERVING:
		return "NOT_SERVING"
	case grpc_health_v1.HealthCheckResponse_SERVICE_UNKNOWN:
		return "SERVICE_UNKNOWN"
	default:
		return "UNKNOWN"
	}
}

func startHealthCheckLoop() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		results := CheckAllServices()
		for _, r := range results {
			if r.Status != "SERVING" {
				log.Printf("[健康检查] %s 状态异常: %s", r.Name, r.Status)
			}
		}
	}
}

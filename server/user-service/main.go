package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "foodRatingSystem/proto/user"
	"foodRatingSystem/shared/config"
	"foodRatingSystem/shared/database"
	"foodRatingSystem/shared/model"
	"foodRatingSystem/shared/registry"
	"foodRatingSystem/shared/utils"
	"foodRatingSystem/user-service/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type UserServer struct {
	pb.UnimplementedUserServiceServer
}

func (s *UserServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	user := &model.User{
		UserName:     req.Username,
		PasswordHash: req.Password,
	}

	registeredUser, err := service.Register(user)
	if err != nil {
		return nil, err
	}

	token, err := utils.GenerateToken(fmt.Sprintf("%d", registeredUser.ID), registeredUser.UserName)
	if err != nil {
		return nil, fmt.Errorf("生成token失败")
	}

	return &pb.RegisterResponse{
		Id:       uint32(registeredUser.ID),
		Username: registeredUser.UserName,
		Token:    token,
	}, nil
}

func (s *UserServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := service.Login(req.Username, req.Password)
	if err != nil {
		return nil, err
	}

	token, err := utils.GenerateToken(fmt.Sprintf("%d", user.ID), user.UserName)
	if err != nil {
		return nil, fmt.Errorf("生成token失败")
	}

	return &pb.LoginResponse{
		Id:       uint32(user.ID),
		Username: user.UserName,
		Token:    token,
	}, nil
}

func (s *UserServer) GetUserByID(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.GetUserByIDResponse, error) {
	user, err := service.GetUserByID(uint(req.Id))
	if err != nil {
		return nil, err
	}

	return &pb.GetUserByIDResponse{
		Id:       uint32(user.ID),
		Username: user.UserName,
	}, nil
}

func (s *UserServer) VerifyUser(ctx context.Context, req *pb.VerifyUserRequest) (*pb.VerifyUserResponse, error) {
	user, err := service.GetUserByID(uint(req.UserId))
	if err != nil {
		return &pb.VerifyUserResponse{Valid: false}, nil
	}

	if user.UserName != req.Username {
		return &pb.VerifyUserResponse{Valid: false}, nil
	}

	return &pb.VerifyUserResponse{
		Valid:    true,
		UserId:   uint32(user.ID),
		Username: user.UserName,
	}, nil
}

func main() {
	config.LoadConfig()
	database.Connectdb()

	serviceName := "user-service"
	serviceAddr := getEnvDefault("SERVICE_ADDR", "localhost:50051")
	port := getEnvDefault("SERVICE_PORT", "50051")

	etcdRegistry, err := registry.NewEtcdRegistry(config.AppConfig.EtcdEndpoints)
	if err != nil {
		log.Fatalf("连接 etcd 失败: %v", err)
	}

	ctx := context.Background()
	if err := etcdRegistry.Register(ctx, serviceName, serviceAddr); err != nil {
		log.Fatalf("注册服务到 etcd 失败: %v", err)
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("用户服务监听失败: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, &UserServer{})

	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("user.UserService", grpc_health_v1.HealthCheckResponse_SERVING)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("正在关闭用户服务...")
		grpcServer.GracefulStop()
		_ = etcdRegistry.Deregister(ctx, serviceName, serviceAddr)
		_ = etcdRegistry.Close()
	}()

	fmt.Printf("用户 gRPC 服务已启动，监听端口 :%s\n", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("用户服务启动失败: %v", err)
	}
}

func getEnvDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

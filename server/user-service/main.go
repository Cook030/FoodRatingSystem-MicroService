package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "foodRatingSystem/proto/user"
	"foodRatingSystem/shared/config"
	"foodRatingSystem/shared/database"
	"foodRatingSystem/shared/model"
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

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("用户服务监听失败: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, &UserServer{})

	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("user.UserService", grpc_health_v1.HealthCheckResponse_SERVING)

	fmt.Println("用户 gRPC 服务已启动，监听端口 :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("用户服务启动失败: %v", err)
	}
}

package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "foodRatingSystem/proto/rating"
	"foodRatingSystem/proto/restaurant"
	grpcclient "foodRatingSystem/rating-service/grpc-client"
	ratingservice "foodRatingSystem/rating-service/service"
	"foodRatingSystem/shared/config"
	"foodRatingSystem/shared/database"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type RatingServer struct {
	pb.UnimplementedRatingServiceServer
}

func (s *RatingServer) SubmitRating(ctx context.Context, req *pb.SubmitRatingRequest) (*pb.SubmitRatingResponse, error) {
	var targetrest interface{}
	var resID int
	if req.RestaurantId > 0 {
		targetrest = int(req.RestaurantId)
		resID = int(req.RestaurantId)
	} else if req.RestaurantName != "" {
		targetrest = req.RestaurantName
	} else {
		return nil, fmt.Errorf("必须提供餐厅ID或餐厅名称")
	}

	err := ratingservice.SubmitReview(targetrest, uint(req.UserId), req.Stars, req.Comment)
	if err != nil {
		return nil, err
	}

	if resID > 0 {
		go func() {
			_, _ = grpcclient.RestaurantClient.InvalidateCache(context.Background(), &restaurant.InvalidateCacheRequest{
				RestaurantId: int32(resID),
			})
		}()
	}

	return &pb.SubmitRatingResponse{Message: "评价成功！"}, nil
}

func (s *RatingServer) GetRatingsByRestaurantID(ctx context.Context, req *pb.RestaurantIDRequest) (*pb.RatingsResponse, error) {
	ratings, err := ratingservice.GetRatingsByRestaurantID(int(req.Id))
	if err != nil {
		return nil, err
	}

	var ratingMessages []*pb.RatingMessage
	for _, r := range ratings {
		ratingMessages = append(ratingMessages, &pb.RatingMessage{
			Id:           uint32(r.ID),
			UserId:       uint32(r.UserID),
			RestaurantId: uint32(r.RestaurantID),
			Stars:        r.Stars,
			Comment:      r.Comment,
			CreatedAt:    r.CreatedAt.Format("2006-01-02T15:04:05Z"),
			Username:     r.User.UserName,
		})
	}

	return &pb.RatingsResponse{Ratings: ratingMessages}, nil
}

func main() {
	config.LoadConfig()
	database.Connectdb()
	database.ConnectRedis()

	err := grpcclient.InitRestaurantClient()
	if err != nil {
		log.Fatalf("gRPC 客户端初始化失败: %v", err)
	}

	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("评分服务监听失败: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterRatingServiceServer(grpcServer, &RatingServer{})

	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("rating.RatingService", grpc_health_v1.HealthCheckResponse_SERVING)

	fmt.Println("评分 gRPC 服务已启动，监听端口 :50053")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("评分服务启动失败: %v", err)
	}
}

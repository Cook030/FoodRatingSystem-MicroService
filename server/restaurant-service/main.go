package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "foodRatingSystem/proto/restaurant"
	rservice "foodRatingSystem/restaurant-service/service"
	"foodRatingSystem/shared/config"
	"foodRatingSystem/shared/database"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type RestaurantServer struct {
	pb.UnimplementedRestaurantServiceServer
}

func (s *RestaurantServer) GetNearbyRestaurants(ctx context.Context, req *pb.NearbyRequest) (*pb.NearbyResponse, error) {
	data, err := rservice.GetNearbyRestaurants(req.Lat, req.Lon)
	if err != nil {
		return nil, err
	}

	var restaurants []*pb.RestaurantWithDistance
	for _, r := range data {
		restaurants = append(restaurants, &pb.RestaurantWithDistance{
			Id:          uint32(r.ID),
			Name:        r.Name,
			Latitude:    r.Latitude,
			Longitude:   r.Longitude,
			AvgScore:    r.AverageScore,
			Category:    r.Category,
			ReviewCount: int32(r.ReviewCount),
			Distance:    r.Distance,
		})
	}

	return &pb.NearbyResponse{Restaurants: restaurants}, nil
}

func (s *RestaurantServer) GetRecommendedRestaurants(ctx context.Context, req *pb.RecommendRequest) (*pb.RecommendResponse, error) {
	data, err := rservice.GetRecommendedRestaurants(req.Lat, req.Lon)
	if err != nil {
		return nil, err
	}

	var restaurants []*pb.RestaurantWithScore
	for _, r := range data {
		restaurants = append(restaurants, &pb.RestaurantWithScore{
			Id:          uint32(r.ID),
			Name:        r.Name,
			Latitude:    r.Latitude,
			Longitude:   r.Longitude,
			AvgScore:    r.AverageScore,
			Category:    r.Category,
			ReviewCount: int32(r.ReviewCount),
			Distance:    r.Distance,
			FinalScore:  r.FinalScore,
		})
	}

	return &pb.RecommendResponse{Restaurants: restaurants}, nil
}

func (s *RestaurantServer) GetRestaurants(ctx context.Context, req *pb.SearchRequest) (*pb.SearchResponse, error) {
	data, err := rservice.GetRestaurants(req.Lat, req.Lon, req.Search, req.Sort, req.HasLocation)
	if err != nil {
		return nil, err
	}

	var restaurants []*pb.RestaurantWithScore
	for _, r := range data {
		restaurants = append(restaurants, &pb.RestaurantWithScore{
			Id:          uint32(r.ID),
			Name:        r.Name,
			Latitude:    r.Latitude,
			Longitude:   r.Longitude,
			AvgScore:    r.AverageScore,
			Category:    r.Category,
			ReviewCount: int32(r.ReviewCount),
			Distance:    r.Distance,
			FinalScore:  r.FinalScore,
		})
	}

	return &pb.SearchResponse{Restaurants: restaurants}, nil
}

func (s *RestaurantServer) GetRestaurantDetail(ctx context.Context, req *pb.RestaurantIDRequest) (*pb.RestaurantDetailResponse, error) {
	restaurant, err := rservice.GetRestaurantByID(int(req.Id))
	if err != nil {
		return nil, err
	}

	return &pb.RestaurantDetailResponse{
		Id:          uint32(restaurant.ID),
		Name:        restaurant.Name,
		Latitude:    restaurant.Latitude,
		Longitude:   restaurant.Longitude,
		AvgScore:    restaurant.AverageScore,
		Category:    restaurant.Category,
		ReviewCount: int32(restaurant.ReviewCount),
	}, nil
}

func (s *RestaurantServer) GetRestaurantRatings(ctx context.Context, req *pb.RestaurantIDRequest) (*pb.RestaurantRatingsResponse, error) {
	ratings, err := rservice.GetRatingsByRestaurantID(int(req.Id))
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

	return &pb.RestaurantRatingsResponse{Ratings: ratingMessages}, nil
}

func (s *RestaurantServer) CreateRestaurant(ctx context.Context, req *pb.CreateRestaurantRequest) (*pb.CreateRestaurantResponse, error) {
	restaurant, err := rservice.CreateRestaurant(req.Name, req.Latitude, req.Longitude, req.Category)
	if err != nil {
		return nil, err
	}

	return &pb.CreateRestaurantResponse{
		Id:          uint32(restaurant.ID),
		Name:        restaurant.Name,
		Latitude:    restaurant.Latitude,
		Longitude:   restaurant.Longitude,
		AvgScore:    restaurant.AverageScore,
		Category:    restaurant.Category,
		ReviewCount: int32(restaurant.ReviewCount),
	}, nil
}

func (s *RestaurantServer) InvalidateCache(ctx context.Context, req *pb.InvalidateCacheRequest) (*pb.InvalidateCacheResponse, error) {
	rservice.ClearRatingCache(uint(req.RestaurantId))
	return &pb.InvalidateCacheResponse{Success: true}, nil
}

func main() {
	config.LoadConfig()
	database.Connectdb()
	database.ConnectRedis()

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("餐厅服务监听失败: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterRestaurantServiceServer(grpcServer, &RestaurantServer{})

	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("restaurant.RestaurantService", grpc_health_v1.HealthCheckResponse_SERVING)

	fmt.Println("餐厅 gRPC 服务已启动，监听端口 :50052")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("餐厅服务启动失败: %v", err)
	}
}

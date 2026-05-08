package grpc_client

import (
	"foodRatingSystem/proto/restaurant"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var RestaurantClient restaurant.RestaurantServiceClient

func InitRestaurantClient() error {
	conn, err := grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	RestaurantClient = restaurant.NewRestaurantServiceClient(conn)
	return nil
}

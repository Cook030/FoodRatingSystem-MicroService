package grpc_client

import (
	"foodRatingSystem/proto/restaurant"
	"foodRatingSystem/shared/registry"
)

var RestaurantClient restaurant.RestaurantServiceClient

func InitRestaurantClient(etcdEndpoints []string) error {
	conn, err := registry.NewEtcdGrpcConn(etcdEndpoints, "restaurant-service")
	if err != nil {
		return err
	}
	RestaurantClient = restaurant.NewRestaurantServiceClient(conn)
	return nil
}

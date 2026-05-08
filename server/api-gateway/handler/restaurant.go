package handler

import (
	"context"
	"net/http"
	"strconv"

	grpcclients "foodRatingSystem/api-gateway/grpc-clients"
	"foodRatingSystem/proto/restaurant"

	"github.com/gin-gonic/gin"
)

func GetNearbyRestaurants(c *gin.Context) {
	latStr := c.Query("lat")
	lonStr := c.Query("lon")

	lat, _ := strconv.ParseFloat(latStr, 64)
	lon, _ := strconv.ParseFloat(lonStr, 64)

	resp, err := grpcclients.RestaurantClient.GetNearbyRestaurants(context.Background(), &restaurant.NearbyRequest{
		Lat: lat,
		Lon: lon,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp.Restaurants)
}

func GetRestaurants(c *gin.Context) {
	latStr := c.Query("lat")
	lonStr := c.Query("lon")
	search := c.Query("search")
	sortBy := c.DefaultQuery("sort", "distance")

	var lat, lon float64
	var hasLocation bool
	if latStr != "" && lonStr != "" {
		lat, _ = strconv.ParseFloat(latStr, 64)
		lon, _ = strconv.ParseFloat(lonStr, 64)
		if lat != 0 && lon != 0 {
			hasLocation = true
		}
	}

	resp, err := grpcclients.RestaurantClient.GetRestaurants(context.Background(), &restaurant.SearchRequest{
		Lat:         lat,
		Lon:         lon,
		Search:      search,
		Sort:        sortBy,
		HasLocation: hasLocation,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp.Restaurants)
}

func GetRestaurantDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的餐厅ID"})
		return
	}

	resp, err := grpcclients.RestaurantClient.GetRestaurantDetail(context.Background(), &restaurant.RestaurantIDRequest{
		Id: int32(id),
	})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "餐厅不存在"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func GetRestaurantRatings(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的餐厅ID"})
		return
	}

	resp, err := grpcclients.RestaurantClient.GetRestaurantRatings(context.Background(), &restaurant.RestaurantIDRequest{
		Id: int32(id),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp.Ratings)
}

func GetRecommendedRestaurants(c *gin.Context) {
	latStr := c.Query("lat")
	lonStr := c.Query("lon")

	lat, _ := strconv.ParseFloat(latStr, 64)
	lon, _ := strconv.ParseFloat(lonStr, 64)

	resp, err := grpcclients.RestaurantClient.GetRecommendedRestaurants(context.Background(), &restaurant.RecommendRequest{
		Lat: lat,
		Lon: lon,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp.Restaurants)
}

func CreateRestaurant(c *gin.Context) {
	var req struct {
		Name      string  `json:"name"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Category  string  `json:"category"`
	}

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := grpcclients.RestaurantClient.CreateRestaurant(context.Background(), &restaurant.CreateRestaurantRequest{
		Name:      req.Name,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		Category:  req.Category,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, resp)
}

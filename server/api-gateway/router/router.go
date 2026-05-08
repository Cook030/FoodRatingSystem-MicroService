package router

import (
	"foodRatingSystem/api-gateway/handler"
	"foodRatingSystem/api-gateway/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(middleware.Cors())

	api := r.Group("/api")
	{
		api.GET("/health", handler.HealthCheck)

		api.GET("/restaurants/nearby", handler.GetNearbyRestaurants)
		api.GET("/restaurants/recommend", handler.GetRecommendedRestaurants)
		api.GET("/restaurants", handler.GetRestaurants)
		api.GET("/restaurants/:id", handler.GetRestaurantDetail)
		api.GET("/restaurants/:id/ratings", handler.GetRestaurantRatings)

		api.POST("/user/register", handler.Register)
		api.POST("/user/login", handler.Login)

		auth := api.Group("", middleware.JWTAuth())
		{
			auth.POST("/rating", handler.SubmitRating)
			auth.POST("/restaurants", handler.CreateRestaurant)
		}
	}

	return r
}

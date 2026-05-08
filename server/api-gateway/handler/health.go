package handler

import (
	"net/http"
	"time"

	grpcclients "foodRatingSystem/api-gateway/grpc-clients"

	"github.com/gin-gonic/gin"
)

func HealthCheck(c *gin.Context) {
	services := grpcclients.CheckAllServices()

	allServing := true
	for _, s := range services {
		if s.Status != "SERVING" {
			allServing = false
			break
		}
	}

	status := "healthy"
	httpCode := http.StatusOK
	if !allServing {
		status = "degraded"
		httpCode = http.StatusServiceUnavailable
	}

	c.JSON(httpCode, gin.H{
		"status":    status,
		"timestamp": time.Now().Format(time.RFC3339),
		"services":  services,
	})
}

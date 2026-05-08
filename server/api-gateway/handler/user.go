package handler

import (
	"context"
	"net/http"

	grpcclients "foodRatingSystem/api-gateway/grpc-clients"
	"foodRatingSystem/proto/user"

	"github.com/gin-gonic/gin"
)

type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供用户名、密码"})
		return
	}

	resp, err := grpcclients.UserClient.Register(context.Background(), &user.RegisterRequest{
		Username: input.Username,
		Password: input.Password,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "注册成功",
		"user": gin.H{
			"id":       resp.Id,
			"username": resp.Username,
		},
		"token": resp.Token,
	})
}

func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供用户名和密码"})
		return
	}

	resp, err := grpcclients.UserClient.Login(context.Background(), &user.LoginRequest{
		Username: input.Username,
		Password: input.Password,
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"user": gin.H{
			"id":       resp.Id,
			"username": resp.Username,
		},
		"token": resp.Token,
	})
}

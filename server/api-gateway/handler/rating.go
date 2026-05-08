package handler

import (
	"context"
	"net/http"

	grpcclients "foodRatingSystem/api-gateway/grpc-clients"
	"foodRatingSystem/proto/rating"
	pbuser "foodRatingSystem/proto/user"

	"github.com/gin-gonic/gin"
)

type RatingRequest struct {
	RestaurantID   int     `json:"restaurant_id"`
	RestaurantName string  `json:"restaurant_name"`
	Username       string  `json:"username"`
	Stars          float64 `json:"stars"`
	Comment        string  `json:"comment"`
}

func SubmitRating(c *gin.Context) {
	var req RatingRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数格式错误"})
		return
	}

	userID := c.GetString("user_id")
	userName := c.GetString("user_name")

	verifyResp, err := grpcclients.UserClient.VerifyUser(context.Background(), &pbuser.VerifyUserRequest{
		UserId:   parseUint32(userID),
		Username: userName,
	})
	if err != nil || !verifyResp.Valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户不存在"})
		return
	}

	_, err = grpcclients.RatingClient.SubmitRating(context.Background(), &rating.SubmitRatingRequest{
		UserId:         verifyResp.UserId,
		RestaurantId:   int32(req.RestaurantID),
		RestaurantName: req.RestaurantName,
		Stars:          req.Stars,
		Comment:        req.Comment,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "评价成功！"})
}

func parseUint32(s string) uint32 {
	var n uint32
	for _, c := range s {
		if c >= '0' && c <= '9' {
			n = n*10 + uint32(c-'0')
		}
	}
	return n
}

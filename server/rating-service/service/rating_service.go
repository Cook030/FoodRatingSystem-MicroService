package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"foodRatingSystem/rating-service/repository"
	"foodRatingSystem/shared/database"
	"foodRatingSystem/shared/model"
	"time"
)

func SubmitReview(targetrest interface{}, userid uint, stars float64, comment string) error {
	rating := model.Rating{
		UserID:    userid,
		Stars:     stars,
		Comment:   comment,
		CreatedAt: time.Now(),
	}
	var resID uint
	if v, ok := targetrest.(int); ok {
		resID = uint(v)
	} else if v, ok := targetrest.(string); ok {
		r, err := repository.FindRestaurantByName(v)
		if err != nil {
			return errors.New("找不到餐厅[" + v + "]")
		}
		resID = r.ID
	} else {
		return errors.New("第一个参数必须是餐厅ID(int)或者餐厅名(string)")
	}

	rating.RestaurantID = resID

	err := repository.AddRatingAndUpdateScore(rating)
	if err != nil {
		return err
	}

	ClearRatingCache(resID)
	return nil
}

func GetRatingsByRestaurantID(restaurantID int) ([]model.Rating, error) {
	cacheKey := fmt.Sprintf("ratings:%d", restaurantID)
	cacheData, err := database.RedisClient.Get(database.Ctx, cacheKey).Result()
	if err == nil {
		var ratings []model.Rating
		json.Unmarshal([]byte(cacheData), &ratings)
		return ratings, nil
	}

	ratings, err := repository.GetRatingsByRestaurantID(restaurantID)
	if err != nil {
		return nil, err
	}

	data, _ := json.Marshal(ratings)
	database.RedisClient.Set(database.Ctx, cacheKey, data, 10*time.Minute)

	return ratings, nil
}

func ClearRatingCache(restaurantID uint) {
	if database.RedisClient == nil {
		return
	}
	database.RedisClient.Del(database.Ctx, fmt.Sprintf("restaurant:%d", restaurantID))
	database.RedisClient.Del(database.Ctx, fmt.Sprintf("ratings:%d", restaurantID))

	patterns := []string{"recommend:*", "nearby:*", "search:*"}
	for _, pattern := range patterns {
		keys, _ := database.RedisClient.Keys(database.Ctx, pattern).Result()
		if len(keys) > 0 {
			database.RedisClient.Del(database.Ctx, keys...)
		}
	}
}

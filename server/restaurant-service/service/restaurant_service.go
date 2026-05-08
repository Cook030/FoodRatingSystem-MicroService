package service

import (
	"encoding/json"
	"fmt"
	"foodRatingSystem/restaurant-service/repository"
	"foodRatingSystem/shared/database"
	"foodRatingSystem/shared/model"
	"foodRatingSystem/shared/utils"
	"math"
	"sort"
	"strings"
	"time"
)

type RestaurantWithDistance struct {
	model.Restaurant
	Distance float64 `json:"distance"`
}

type RestaurantWithScore struct {
	model.Restaurant
	Distance   float64 `json:"distance"`
	FinalScore float64 `json:"final_score"`
}

func GetNearbyRestaurants(userLat, userLon float64) ([]RestaurantWithDistance, error) {
	cacheKey := fmt.Sprintf("nearby:%.2f:%.2f", userLat, userLon)
	cacheData, err := database.RedisClient.Get(database.Ctx, cacheKey).Result()
	if err == nil {
		var results []RestaurantWithDistance
		json.Unmarshal([]byte(cacheData), &results)
		return results, nil
	}

	rests, err := repository.GetAllRestaurants()
	if err != nil {
		return nil, err
	}
	var rwd []RestaurantWithDistance
	for _, rest := range rests {
		dist := utils.Distance(userLat, userLon, rest.Latitude, rest.Longitude)
		rwd = append(rwd, RestaurantWithDistance{
			Restaurant: rest,
			Distance:   dist,
		})
	}
	sort.Slice(rwd, func(i, j int) bool {
		return rwd[i].Distance < rwd[j].Distance
	})

	data, _ := json.Marshal(rwd)
	database.RedisClient.Set(database.Ctx, cacheKey, data, 2*time.Hour)

	return rwd, nil
}

func GetRecommendedRestaurants(userLat, userLon float64) ([]RestaurantWithScore, error) {
	cacheKey := fmt.Sprintf("recommend:%.2f:%.2f", userLat, userLon)
	cacheData, err := database.RedisClient.Get(database.Ctx, cacheKey).Result()
	if err == nil {
		var results []RestaurantWithScore
		json.Unmarshal([]byte(cacheData), &results)
		return results, nil
	}

	rests, err := repository.GetAllRestaurants()
	if err != nil {
		return nil, err
	}

	var results []RestaurantWithScore
	for _, rest := range rests {
		dist := utils.Distance(userLat, userLon, rest.Latitude, rest.Longitude)

		scorePart := rest.AverageScore * 0.6
		distPart := (1.0 / (dist + 1.0)) * 0.3
		reviewPart := math.Log10(float64(rest.ReviewCount)+1.0) * 0.1

		finalScore := scorePart + distPart + reviewPart

		results = append(results, RestaurantWithScore{
			Restaurant: rest,
			Distance:   dist,
			FinalScore: finalScore,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].FinalScore > results[j].FinalScore
	})

	data, _ := json.Marshal(results)
	database.RedisClient.Set(database.Ctx, cacheKey, data, 2*time.Hour)

	return results, nil
}

func GetRestaurants(userLat, userLon float64, search, sortBy string, hasLocation bool) ([]RestaurantWithScore, error) {
	cacheKey := fmt.Sprintf("search:%.2f:%.2f:%s:%s:%v", userLat, userLon, search, sortBy, hasLocation)
	cacheData, err := database.RedisClient.Get(database.Ctx, cacheKey).Result()
	if err == nil {
		var results []RestaurantWithScore
		json.Unmarshal([]byte(cacheData), &results)
		return results, nil
	}

	rests, err := repository.GetAllRestaurants()
	if err != nil {
		return nil, err
	}

	var results []RestaurantWithScore
	for _, rest := range rests {
		if search != "" && !isContainsKeywords(rest.Name, search) {
			continue
		}

		var dist float64
		var finalScore float64

		if hasLocation {
			dist = utils.Distance(userLat, userLon, rest.Latitude, rest.Longitude)
			scorePart := rest.AverageScore * 0.6
			distPart := (1.0 / (dist + 1.0)) * 0.3
			reviewPart := math.Log10(float64(rest.ReviewCount)+1.0) * 0.1
			finalScore = scorePart + distPart + reviewPart
		} else {
			scorePart := rest.AverageScore * 0.7
			reviewPart := math.Log10(float64(rest.ReviewCount)+1.0) * 0.3
			finalScore = scorePart + reviewPart
			dist = -1
		}

		results = append(results, RestaurantWithScore{
			Restaurant: rest,
			Distance:   dist,
			FinalScore: finalScore,
		})
	}

	switch sortBy {
	case "score":
		sort.Slice(results, func(i, j int) bool {
			return results[i].AverageScore > results[j].AverageScore
		})
	case "reviews":
		sort.Slice(results, func(i, j int) bool {
			return results[i].ReviewCount > results[j].ReviewCount
		})
	case "recommended":
		sort.Slice(results, func(i, j int) bool {
			return results[i].FinalScore > results[j].FinalScore
		})
	default:
		if hasLocation {
			sort.Slice(results, func(i, j int) bool {
				return results[i].Distance < results[j].Distance
			})
		} else {
			sort.Slice(results, func(i, j int) bool {
				return results[i].AverageScore > results[j].AverageScore
			})
		}
	}

	data, _ := json.Marshal(results)
	database.RedisClient.Set(database.Ctx, cacheKey, data, 2*time.Hour)

	return results, nil
}

func isContainsKeywords(s, substr string) bool {
	s = strings.ToLower(s)
	substr = strings.ToLower(substr)
	return strings.Contains(s, substr)
}

func GetRestaurantByID(id int) (*model.Restaurant, error) {
	cacheKey := fmt.Sprintf("restaurant:%d", id)
	cacheData, err := database.RedisClient.Get(database.Ctx, cacheKey).Result()
	if err == nil {
		var rest model.Restaurant
		json.Unmarshal([]byte(cacheData), &rest)
		return &rest, nil
	}

	rest, err := repository.GetRestaurantByID(id)
	if err != nil {
		return nil, err
	}
	data, _ := json.Marshal(rest)
	database.RedisClient.Set(database.Ctx, cacheKey, data, 10*time.Minute)
	return rest, nil
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

func CreateRestaurant(name string, lat, lon float64, category string) (*model.Restaurant, error) {
	rest := model.Restaurant{
		Name:      name,
		Latitude:  lat,
		Longitude: lon,
		Category:  category,
	}

	result, err := repository.CreateRestaurant(rest)
	if err == nil {
		ClearListCache()
		restaurantCacheKey := fmt.Sprintf("restaurant:%d", result.ID)
		database.RedisClient.Del(database.Ctx, restaurantCacheKey)
	}
	return result, err
}

func ClearListCache() {
	if database.RedisClient == nil {
		return
	}
	patterns := []string{"recommend:*", "nearby:*", "search:*"}
	for _, pattern := range patterns {
		keys, _ := database.RedisClient.Keys(database.Ctx, pattern).Result()
		if len(keys) > 0 {
			database.RedisClient.Del(database.Ctx, keys...)
		}
	}
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

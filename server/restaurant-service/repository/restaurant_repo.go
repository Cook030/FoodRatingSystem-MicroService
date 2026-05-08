package repository

import (
	"foodRatingSystem/shared/database"
	"foodRatingSystem/shared/model"
)

func GetAllRestaurants() ([]model.Restaurant, error) {
	var restaurants []model.Restaurant
	err := database.DB.Find(&restaurants).Error
	if err != nil {
		return nil, err
	}
	for i := range restaurants {
		var count int64
		database.DB.Model(&model.Rating{}).
			Where("restaurant_id = ?", restaurants[i].ID).
			Count(&count)
		restaurants[i].ReviewCount = int(count)
	}
	return restaurants, nil
}

func GetRestaurantByID(id int) (*model.Restaurant, error) {
	var r model.Restaurant
	err := database.DB.First(&r, id).Error
	if err != nil {
		return nil, err
	}
	var count int64
	database.DB.Model(&model.Rating{}).Where("restaurant_id = ?", r.ID).Count(&count)
	r.ReviewCount = int(count)
	return &r, nil
}

func CreateRestaurant(rest model.Restaurant) (*model.Restaurant, error) {
	err := database.DB.Create(&rest).Error
	if err != nil {
		return nil, err
	}
	return &rest, nil
}

func GetRatingsByRestaurantID(restaurantID int) ([]model.Rating, error) {
	var ratings []model.Rating
	err := database.DB.Preload("User").Where("restaurant_id = ?", restaurantID).Order("created_at DESC").Find(&ratings).Error
	return ratings, err
}

package repository

import (
	"foodRatingSystem/shared/database"
	"foodRatingSystem/shared/model"

	"gorm.io/gorm"
)

func AddRatingAndUpdateScore(r model.Rating) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&r).Error
		if err != nil {
			return err
		}
		var avg float64
		var count int64
		tx.Model(&model.Rating{}).Where("restaurant_id = ?", r.RestaurantID).Count(&count)
		tx.Model(&model.Rating{}).Where("restaurant_id = ?", r.RestaurantID).Select("COALESCE(AVG(stars), 0)").Scan(&avg)
		result := tx.Model(&model.Restaurant{}).Where("id = ?", r.RestaurantID).Updates(map[string]interface{}{
			"avg_score":    avg,
			"review_count": count,
		})
		return result.Error
	})
}

func GetRatingsByRestaurantID(restaurantID int) ([]model.Rating, error) {
	var ratings []model.Rating
	err := database.DB.Preload("User").Where("restaurant_id = ?", restaurantID).Order("created_at DESC").Find(&ratings).Error
	return ratings, err
}

func FindRestaurantByName(name string) (*model.Restaurant, error) {
	var r model.Restaurant
	err := database.DB.Where("name = ?", name).First(&r).Error
	if err != nil {
		return nil, err
	}
	return &r, nil
}

package repository

import (
	"errors"
	"foodRatingSystem/shared/database"
	"foodRatingSystem/shared/model"

	"gorm.io/gorm"
)

func CreateUser(user *model.User) (*model.User, error) {
	err := database.DB.Create(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetUserByUsername(userName string) (*model.User, error) {
	var u model.User
	err := database.DB.Where("user_name = ?", userName).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &u, nil
}

func GetUserByID(userID uint) (*model.User, error) {
	var u model.User
	err := database.DB.First(&u, userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &u, nil
}

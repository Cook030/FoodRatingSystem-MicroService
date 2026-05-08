package model

import "time"

type Rating struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       uint      `gorm:"index;not null" json:"user_id"`
	RestaurantID uint      `gorm:"index;not null" json:"restaurant_id"`
	Stars        float64   `gorm:"type:decimal(2,1)" json:"stars"`
	Comment      string    `gorm:"type:text" json:"comment"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`

	User       User       `gorm:"foreignKey:UserID" json:"user"`
	Restaurant Restaurant `gorm:"foreignKey:RestaurantID" json:"-"`
}

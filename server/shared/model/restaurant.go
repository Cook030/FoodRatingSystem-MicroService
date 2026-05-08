package model

import "time"

type Restaurant struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string    `gorm:"type:varchar(255);index;not null" json:"name"`
	Latitude     float64   `gorm:"column:latitude;type:decimal(10,7)" json:"latitude"`
	Longitude    float64   `gorm:"column:longitude;type:decimal(10,7)" json:"longitude"`
	AverageScore float64   `gorm:"column:avg_score;type:decimal(3,2);default:0" json:"avg_score"`
	Category     string    `gorm:"column:category;type:varchar(100)" json:"category"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	ReviewCount  int       `gorm:"default:0" json:"review_count"`
	Ratings      []Rating  `gorm:"foreignKey:RestaurantID" json:"-"`
}

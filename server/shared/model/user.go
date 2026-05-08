package model

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserName     string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"user_name"`
	PasswordHash string    `gorm:"type:varchar(255);not null" json:"-"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	Ratings      []Rating  `gorm:"foreignKey:UserID" json:"-"`
}

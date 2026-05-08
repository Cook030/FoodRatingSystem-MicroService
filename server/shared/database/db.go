package database

import (
	"fmt"

	"foodRatingSystem/shared/config"
	"foodRatingSystem/shared/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connectdb() {
	cfg := config.AppConfig
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode)
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("无法连接数据库：", err)
		return
	}
	fmt.Println("成功连接到数据库")

	DB.AutoMigrate(&model.Restaurant{}, &model.Rating{}, &model.User{})
	fmt.Println("数据库迁移完成")
}

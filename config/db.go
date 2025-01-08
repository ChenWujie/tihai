package config

import (
	"log"
	"tihai/global"
	"tihai/internal/model"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func initDB() {
	dsn := AppConfig.Database.Dsn
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to initialize database, got error: %v", err)
	}

	sqlDB, err := db.DB()

	sqlDB.SetMaxIdleConns(AppConfig.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(AppConfig.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err != nil {
		log.Fatalf("Failed to configure database, got error: %v", err)
	}
	err = db.AutoMigrate(&model.User{}, &model.Question{}, &model.Paper{}, &model.Score{},
		&model.StudentAnswer{}, &model.Comment{}, &model.Class{})
	if err != nil {
		log.Fatalf("Failed to migrate database, got error: %v", err)
	}
	global.Db = db
}

package db

import (
	"fmt"
	"log"
	"time"

	"github.com/iragsraghu/user-service/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewMySQL connects to MySQL, ensures DB exists, and runs migrations
func NewMySQL(user, pass, host, port, name string) (*gorm.DB, error) {
	if port == "" {
		port = "3306"
	}
	if host == "" {
		host = "127.0.0.1"
	}
	// 1. Connect without DB
	dsnRoot := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port)

	rootDB, err := gorm.Open(mysql.Open(dsnRoot), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("failed root connect: %w", err)
	}

	// 2. Create DB if not exists
	if err := rootDB.Exec("CREATE DATABASE IF NOT EXISTS " + name).Error; err != nil {
		return nil, fmt.Errorf("failed create db: %w", err)
	}
	log.Printf("✅ Database %s is ready", name)

	// 3. Now connect to actual DB
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port, name)
	fmt.Println(dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("failed final connect: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	log.Println("✅ Connected to MySQL (GORM)")

	// 4. Run AutoMigrate for all models
	if err := db.AutoMigrate(&models.User{}); err != nil {
		return nil, fmt.Errorf("failed to migrate: %w", err)
	}
	log.Println("✅ User table migrated")

	return db, nil
}

package database

import (
	"fmt"
	"log"
	"time"

	"order-service/internal/config"
	"order-service/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	DB *gorm.DB
}

func NewDatabase(cfg *config.Config) (*Database, error) {
	dsn := cfg.GetDatabaseDSN()
	
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Order service database connection established successfully")

	return &Database{DB: db}, nil
}

func (d *Database) AutoMigrate() error {
	if err := d.DB.AutoMigrate(
		&models.Order{},
		&models.OrderItem{},
	); err != nil {
		return fmt.Errorf("failed to auto-migrate database: %w", err)
	}

	log.Println("Order service database migration completed successfully")
	return nil
}

func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}
	return sqlDB.Close()
}

func (d *Database) GetDB() *gorm.DB {
	return d.DB
}

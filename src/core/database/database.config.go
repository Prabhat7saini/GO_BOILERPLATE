package database

import (
	"boiler-platecode/src/core/config"
	"log"
	"os"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db   *gorm.DB
	once sync.Once
)

// InitDB initializes the database singleton
func InitDB() {
	once.Do(func() {
		dsn := config.AppConfig.DbUrl

		// Configure GORM logger
		newLogger := logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  logger.Info, // Use logger.Warn to avoid SQL printing
				IgnoreRecordNotFoundError: true,
				Colorful:                  true,
			},
		)

		var err error
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: newLogger,
		})
		if err != nil {
			log.Fatal("❌ Failed to connect DB: ", err)
		}

		log.Printf("✅ Connected to DB: %s\n", dsn)
	})
}

// GetDB provides access to the singleton DB instance
func GetDB() *gorm.DB {
	if db == nil {
		InitDB()
	}
	return db
}

// CloseDB closes the database connection pool
func CloseDB() error {
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

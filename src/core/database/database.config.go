package database

import (
	"boiler-platecode/src/core/config"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	dsn := config.AppConfig.DbUrl

	// ✅ Set GORM to log only warnings and avoid 'record not found' logs
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,     // Slow SQL threshold
			LogLevel:                  logger.Info,     // Only log Warn and above// info print sql and warn not print the sql
			IgnoreRecordNotFoundError: true,            // suppress 'record not found'
			Colorful:                  true,            // Pretty output
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatal("Failed to connect DB................:", err)
	}

	log.Printf("Connected to DB...............: %s\n", dsn)
	DB = db
}

func GetDB() *gorm.DB {
	return DB
}

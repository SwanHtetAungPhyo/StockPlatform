package database

import (
	"fmt"
	"log"
	"os"
	"time"

	logging "github.com/SwanHtetAungPhyo/closure/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DB_INIT() {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", dbHost, dbPort, dbUser, dbName, dbPassword)

	var err error
	maxAttempts := 10
	retryDelay := 5 * time.Second


	for attempts := 1; attempts <= maxAttempts; attempts++ {
		DB, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
		if err == nil {
			fmt.Println("Successfully connected to the database")
			return
		}
		log.Printf("Failed to connect to the database (Attempt %d/%d): %v", attempts, maxAttempts, err)
		time.Sleep(retryDelay)
	}
	
	sqlDB , err := DB.DB()
	if err != nil {
		logging.Logger.Fatalf("%s", err.Error())
		return 
	}
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxIdleTime(10)
	sqlDB.SetConnMaxLifetime(1 * time.Hour)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)
	logging.Logger.Println("Connection pooling is set.......")
}

func Migration(models ...interface{}) {
	for _, model := range models {
		if model == nil {
			log.Fatalf("Migration received a nil model")
		}

		if err := DB.AutoMigrate(model); err != nil {
			log.Printf("Migration failed for model %T: %v", model, err)
		} else {
			fmt.Printf("Migration succeeded for model %T\n", model)
		}
	}
}

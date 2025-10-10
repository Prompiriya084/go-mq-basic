package database

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	models "github.com/Prompiriya084/go-mq/Models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDb() *gorm.DB {
	// Load .env from specific path
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	host := os.Getenv("DB_Host")
	port, _ := strconv.Atoi(os.Getenv("DB_Port"))
	username := os.Getenv("DB_Username")
	password := os.Getenv("DB_Password")
	database := os.Getenv("DB_Name")

	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, username, password, database)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			// IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			// ParameterizedQueries:      true,          // Don't include params in the SQL log
			Colorful: true, // Disable color
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic("Failed to connect database.")
	}
	fmt.Printf("Connect successful.")

	db.AutoMigrate(models.Order{}, models.Inventory{})

	return db
}

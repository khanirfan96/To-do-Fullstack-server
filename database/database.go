package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBCollections struct {
	TodoCollection    *mongo.Collection
	CalorieCollection *mongo.Collection
	UserCollection    *mongo.Collection
}

var (
	client *mongo.Client
	DB     DBCollections
)

// Initialize sets up the database connection and collections
func Initialize() error {
	if err := loadEnv(); err != nil {
		return fmt.Errorf("failed to load environment variables: %v", err)
	}

	if err := connectDB(); err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	setupCollections()
	return nil
}

// loadEnv loads the environment variables from .env file
func loadEnv() error {
	if err := godotenv.Load(".env"); err != nil {
		return fmt.Errorf("error loading .env file: %v", err)
	}
	return nil
}

// connectDB establishes connection to MongoDB
func connectDB() error {
	connectionString := os.Getenv("DB_URI")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(connectionString)
	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// Ping the database
	if err = client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	fmt.Println("Successfully connected to MongoDB!")
	return nil
}

// setupCollections initializes the database collections
func setupCollections() {
	dbName := os.Getenv("DB_NAME")
	database := client.Database(dbName)

	DB = DBCollections{
		TodoCollection:    database.Collection("todolist"),
		CalorieCollection: database.Collection("calorietracker"),
		UserCollection:    database.Collection("user"),
	}

	fmt.Printf("Collections initialized:\n")
	fmt.Printf("- Todo Collection: %v\n", DB.TodoCollection.Name())
	fmt.Printf("- Calorie Collection: %v\n", DB.CalorieCollection.Name())
}

// GetContext returns a context with timeout
func GetContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 100*time.Second)
}

// Close closes the database connection
func Close() {
	if client != nil {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}
}

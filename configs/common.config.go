package configs

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

type Config struct {
	Port             int
	MongoDB          *mongo.Client
	MongoDB_URI      string
	MongoDB_Database string
	SECRET_KEY       string
}

var Configs *Config

func LoadConfigs() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	port, err := strconv.Atoi(getEnv("PORT", "8080"))
	if err != nil {
		log.Printf("Invalid PORT value, using default 8080")
		port = 8080
	}

	secretKey := getEnv("SECRET_JWT", "MY_SECRET")
	mongoURI := getEnv("MONGO_URI", "mongodb://admin:admin1234@127.0.0.1:27017/go_fiber_jwt_v2?authSource=admin")
	mongoDB := getEnv("MONGO_DATABASE", "go_fiber_jwt_v2")

	config := &Config{
		Port:             port,
		MongoDB_URI:      mongoURI,
		MongoDB_Database: mongoDB,
		SECRET_KEY:       secretKey,
	}

	Configs = config
	return config
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

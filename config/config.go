package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Token                 string
	GuildID               string
	MongoURI              string
	DBName                string
	TwitterConsumerKey    string
	TwitterConsumerSecret string
	TwitterAccessToken    string
	TwitterAccessSecret   string
}

var cfg *Config

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	cfg = &Config{
		Token:                 os.Getenv("DISCORD_BOT_TOKEN"),
		GuildID:               os.Getenv("GUILD_ID"),
		MongoURI:              os.Getenv("MONGO_URI"),
		DBName:                os.Getenv("DB_NAME"),
		TwitterConsumerKey:    os.Getenv("TWITTER_CONSUMER_KEY"),
		TwitterConsumerSecret: os.Getenv("TWITTER_CONSUMER_SECRET"),
		TwitterAccessToken:    os.Getenv("TWITTER_ACCESS_TOKEN"),
		TwitterAccessSecret:   os.Getenv("TWITTER_ACCESS_SECRET"),
	}
	return cfg
}

func GetConfig() *Config {
	return cfg
}

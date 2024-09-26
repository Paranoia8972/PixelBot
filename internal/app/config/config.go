package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Token         string
	GuildID       string
	MongoURI      string
	DBName        string
	Port          string
	TranscriptUrl string
}

var cfg *Config

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	cfg = &Config{
		Token:         os.Getenv("DISCORD_BOT_TOKEN"),
		GuildID:       os.Getenv("GUILD_ID"),
		MongoURI:      os.Getenv("MONGO_URI"),
		DBName:        os.Getenv("DB_NAME"),
		Port:          os.Getenv("SERVER_PORT"),
		TranscriptUrl: os.Getenv("TRANSCRIPT_URL"),
	}
	return cfg
}

func GetConfig() *Config {
	return cfg
}

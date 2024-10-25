package config

import (
	"log"
	"os"
	"path/filepath"

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
	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("Error getting executable path: %v", err)
	}
	exeDir := filepath.Dir(exePath)
	envPath := filepath.Join(exeDir, ".env")

	err = godotenv.Load(envPath)
	if err != nil {
		log.Fatalf("Error loading .env file from %s: %v", envPath, err)
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

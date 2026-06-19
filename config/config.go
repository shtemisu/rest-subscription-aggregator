package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseHost       string
	DatabasePort       string
	DatabaseUser       string
	DatabasePassword   string
	DatabaseName       string
	JWT_Access_Secret  string
	JWT_Refresh_Secret string
}

func NewConfig() *Config {
	_ = godotenv.Load()
	log.Printf("Loading config from environment variables: DB_HOST=%s, DB_PORT=%s, DB_USER=%s\n", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"))
	return &Config{
		DatabaseHost:       os.Getenv("DB_HOST"),
		DatabasePort:       os.Getenv("DB_PORT"),
		DatabaseUser:       os.Getenv("DB_USER"),
		DatabasePassword:   os.Getenv("DB_PASSWORD"),
		DatabaseName:       os.Getenv("DB_NAME"),
		JWT_Access_Secret:  os.Getenv("JWT_ACCESS"),
		JWT_Refresh_Secret: os.Getenv("JWT_REFRESH"),
	}
}

package config

import (
	"os"

	"github.com/joho/godotenv"
)

type APP_ENV string

func (a APP_ENV) ToString() string {
	return string(a)
}

const (
	DEV  APP_ENV = "development"
	PROD APP_ENV = "production"
	TEST APP_ENV = "test"
)

type Config struct {
	APP_ENV APP_ENV
	PORT    string
	DB_URL  string
}

var Env = getEnvs()

func (c Config) IsProduction() bool {
	return c.APP_ENV == PROD
}

func getEnvs() Config {
	return Config{
		APP_ENV: APP_ENV((getDefaults("APP_ENV", DEV.ToString()))),
		PORT:    getDefaults("PORT", ":5001"),
		DB_URL:  getDefaults("DATABASE_URL", "postgresql://localhost:5432/poll-app"),
	}
}

func getDefaults(key, fallback string) string {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

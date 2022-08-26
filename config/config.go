package config

import (
	"go-dynamodb-crud/utils/env"
	"strconv"
)

type Config struct {
	Port        int
	TimeOut     int
	Dialect     string
	DatabaseURI string
}

func GetConfig() Config {
	return Config{
		Port:        parseEnvToInt("PORT", "8090"),
		TimeOut:     parseEnvToInt("TIMEOUT", "30"),
		Dialect:     env.GetEnv("DIALECT", "sqlite3"),
		DatabaseURI: env.GetEnv("DATABASE_URI", ":memory:"),
	}
}

func parseEnvToInt(envName string, defaultValue string) int {
	num, err := strconv.Atoi(env.GetEnv(envName, defaultValue))
	if err != nil {
		return 0
	}
	return num
}

package config

import (
	"fmt"
	"os"
)

type Config struct {
	DBUser string
	DBPass string
	DBHost string
	DBName string
	Port   string
}

func Load() *Config {
	return &Config{
		DBUser: os.Getenv("DB_USER"),
		DBPass: os.Getenv("DB_PASS"),
		DBHost: os.Getenv("DB_HOST"),
		DBName: os.Getenv("DB_NAME"),
		Port:   os.Getenv("DB_PORT"),
	}
}

func (c *Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true",
		c.DBUser, c.DBPass, c.DBHost, c.DBName)
}

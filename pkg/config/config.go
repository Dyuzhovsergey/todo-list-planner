// Package config for configuration parametrs
package config

import (
	"log"
	"os"
)

type Config struct {
	DBFile   string
	HTTPPort string
	Password string
	JWTKey   []byte
}

var Cfg Config

func Init() {
	Cfg = Config{
		DBFile:   getEnv("TODO_DBFILE", "data/scheduler.db"),
		HTTPPort: getEnv("TODO_PORT", "7540"),
		Password: os.Getenv("TODO_PASSWORD"),
		JWTKey:   []byte(os.Getenv("TODO_SECRET")),
	}
	if len(Cfg.JWTKey) == 0 {
		log.Println("TODO_SECRET not set! Used default key")
		Cfg.JWTKey = []byte("default_secret")
	}
}

func getEnv(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}

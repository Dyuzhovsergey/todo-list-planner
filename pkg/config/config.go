package config

import (
	"log"
	"os"
)

type Config struct {
	Password string
	JWTKey   []byte
}

var Cfg Config

func Init() {
	Cfg = Config{
		Password: os.Getenv("TODO_PASSWORD"),
		JWTKey:   []byte(os.Getenv("TODO_PASSWORD")),
	}
	if len(Cfg.JWTKey) == 0 {
		log.Println("TODO_SECRET not set! Used default key")
		Cfg.JWTKey = []byte("default_secret")
	}
}

package config

import (
	"log"
	"os"
)

type Config struct {
	Password string
	JWRKey   []byte
}

var Cfg Config

func Init() {
	Cfg = Config{
		Password: os.Getenv("TODO_PASSWORD"),
		JWRKey:   []byte(os.Getenv("TODO_PASSWORD")),
	}
	if len(Cfg.JWRKey) == 0 {
		log.Println("TODO_SECRET not set! Used default key")
		Cfg.JWRKey = []byte("default_secret")
	}
}

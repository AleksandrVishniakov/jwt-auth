package configs

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env  string `env:"ENV" env-default:"production"`
	JWTSignature string `env:"JWT_SIGNATURE"`
	HTTP HTTP
	DB DB
	Admin Admin
}

type HTTP struct {
	Port int `env:"HTTP_PORT" env-default:"8080"`
}

type Admin struct {
	Login string `env:"ADMIN_LOGIN"`
	Password string `env:"ADMIN_PASSWORD"`
}

type DB struct {
	Host string `env:"DB_HOST"`
	Port string `env:"DB_PORT"`
	DBName string `env:"DB_NAME"`
	User string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
}

func MustConfig() Config {
	cfg := Config{}

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Fatalf("Failed to read configs: %s", err.Error())
	}

	return cfg
}

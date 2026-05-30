package config

import (
	"log"

	"github.com/joho/godotenv"
)

type Config struct {
	Server struct {
		Port         string
		ReadTimeout  string
		WriteTimeout string
	}
	Database struct {
		Path string
	}
	JWT struct {
		Secret      string
		ExpireHours int
	}
	Logging struct {
		Level string
	}
}

func LoadConfig(path string) (*Config, error) {
	config := &Config{}
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using defaults")
	}
	
	config.Server.Port = "8080"
	config.Server.ReadTimeout = "30s"
	config.Server.WriteTimeout = "30s"
	config.Database.Path = "./data/k8s_ui_admin.db"
	config.JWT.Secret = "k8s-ui-admin-secret-key-change-in-production"
	config.JWT.ExpireHours = 24
	config.Logging.Level = "debug"
	
	return config, nil
}
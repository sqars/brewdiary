package config

import (
	"encoding/json"
	"os"
)

// Config - struct for application configuration data
type Config struct {
	DbHost     string   `json:"database-host"`
	DbPort     string   `json:"database-port"`
	DbName     string   `json:"database-name"`
	DbUser     string   `json:"database-user"`
	DbPassword string   `json:"database-password"`
	Host       string   `json:"application-host"`
	Cors       []string `json:"cors"`
}

// GetConfig returs pointer to Config
func GetConfig() (Config, error) {
	var conf Config
	file, err := os.Open("config/config.json")
	if err != nil {
		return conf, err
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&conf)
	if err != nil {
		return conf, err
	}
	return conf, nil
}

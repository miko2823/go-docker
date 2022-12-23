package main

import (
	"encoding/json"
	"os"
)

type Environment struct {
	DB_HOST     string `json:"DB_HOST"`
	PORT        int64  `json:"PORT"`
	DB_USER     string `json:"DB_USER"`
	DB_PASSWORD string `json:"DB_PASSWORD"`
	DB_NAME     string `json:"DB_NAME"`
}

func getEnvironment() (Environment, error) {
	f, err := os.Open("./config/dev.json")
	os_env := os.Getenv("env")

	if os_env == "prod" {
		f, err = os.Open("./config/prod.json")
		if err != nil {
			return Environment{}, err
		}
	}
	defer f.Close()

	var cfg Environment

	err = json.NewDecoder(f).Decode(&cfg)
	if err != nil {
		return Environment{}, err
	}

	return cfg, nil
}

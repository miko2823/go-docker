package config

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"
)

var dbConnCount int64

type Environment struct {
	DB_HOST     string `json:"DB_HOST"`
	DB_PORT     int64  `json:"PORT"`
	DB_USER     string `json:"DB_USER"`
	DB_PASSWORD string `json:"DB_PASSWORD"`
	DB_NAME     string `json:"DB_NAME"`
}

func GetEnvironment() (Environment, error) {
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

func ConnectToDB() (*sql.DB, error) {

	environment, err := GetEnvironment()
	if err != nil {
		return nil, err
	}

	fmt.Println("ENV", environment)
	dns := fmt.Sprintf(
		"host=%v port=%v user=%v password=%v dbname=%v",
		environment.DB_HOST,
		environment.DB_PORT,
		environment.DB_USER,
		environment.DB_PASSWORD,
		environment.DB_NAME)

	fmt.Println("dns", dns)

	for {
		connection, err := openDB(dns)
		if err != nil {
			log.Println("Postgres not yet ready...", err)
			dbConnCount++
		} else {
			log.Println("Connected Successfully!")
			return connection, nil
		}

		if dbConnCount > 10 {
			return nil, errors.New("Failed to connect postgres...")
		}

		log.Println("Backing off for two seconds...")
		time.Sleep(2 * time.Second)
		continue
	}
}

func openDB(dns string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dns)

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

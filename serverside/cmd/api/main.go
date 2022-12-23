package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const webPort = "9000"

var flagConfig = flag.String("config", "./config/local.yml", "path to the config file")

type Config struct {
	Env     Environment
	Rabbit  *amqp.Connection
	Routing Routing
}

func main() {

	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()
	log.Println("Connected to RabbitMQ")

	env, err := getEnvironment()
	log.Println("get Config", env)

	// app := Config{
	// 	Env:    env,
	// 	Rabbit: rabbitConn,
	// }

	var routing = Routing{}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: routing.buildHandler(),
	}
	log.Printf("Starting broker service on port %s\n", webPort)

	err = server.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")

		if err != nil {
			fmt.Println("RabbitMQ is not ready")
			counts++
		} else {
			connection = c
			break
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backing off...")
		time.Sleep(backOff)
		continue
	}
	return connection, nil
}

package main

import (
	"authentication/data"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var count int64

const (
	webPort  = "8088"
	gRpcPort = "50001"
	rpcPort  = "5001"
)

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("Starting Authentication service")

	// connect DB
	conn := connectToDB()
	if conn == nil {
		log.Panic("Cant't connect to DB")
	}

	//set up config
	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	// Register the RPC Server
	err := rpc.Register(new(RPCServer))
	go app.rpcListen()

	go app.gRPCListen()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
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

func connectToDB() *sql.DB {
	dns := os.Getenv("DNS")
	fmt.Println("DNS", dns)

	for {
		connection, err := openDB(dns)
		if err != nil {
			log.Println("Postgres not yet ready...")
			count++
		} else {
			log.Println("Connected Successfully!")
			return connection
		}
		if count > 10 {
			log.Println(err)
			return nil
		}
		log.Println("Backing off for two seconds...")
		time.Sleep(2 * time.Second)
		continue
	}
}

func (app *Config) rpcListen() error {
	log.Println("Starting RPC server on port ", rpcPort)
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))
	if err != nil {
		return err
	}
	defer listen.Close()

	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(rpcConn)
	}

}

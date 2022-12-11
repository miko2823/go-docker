package main

import (
	"authentication/data"
	"log"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type RPCServer struct{}

type RPCPayload struct {
	Email    string
	Password string
}

func (r *RPCServer) Auth(payload RPCPayload, response *string) error {
	conn := connectToDB()
	if conn == nil {
		log.Panic("Cant't connect to DB")
	}
	model := data.New(conn)
	user, err := model.User.GetByEmail(payload.Email)
	if err != nil {
		return err
	}

	valid, err := user.PasswordMatches(payload.Password)
	if err != nil || !valid {
		return err
	}

	*response = "Processed payload via RPC: " + payload.Email
	return nil
}

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/rpc"
	"time"

	"github.com/miko2823/go-docker/auth"
	"github.com/miko2823/go-docker/event"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPyaload  `json:"log,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPyaload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = app.writeJson(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := app.readJSON(w, r, &requestPayload)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	switch requestPayload.Action {

	case "auth":
		app.authenticate(w, requestPayload.Auth)
	case "authViaRPC":
		app.authViaRPC(w, requestPayload.Auth)
	case "log":
		app.logEventViaRabbit(w, requestPayload.Log)

	default:
		app.errorJSON(w, errors.New("unknown action"))
	}
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	// call auth service
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	request, err := http.NewRequest(
		"POST",
		"http://authentication:8088/authenticate",
		bytes.NewBuffer(jsonData))

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}

	response, err := client.Do(request)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	defer response.Body.Close()

	// return status code
	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling auth service"))
		return
	}

	var jsonFromService jsonResponse

	err = json.NewDecoder(response.Body).Decode(&jsonFromService)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated"
	payload.Data = jsonFromService.Data

	app.writeJson(w, http.StatusAccepted, payload)
}

func (app *Config) logEventViaRabbit(w http.ResponseWriter, l LogPyaload) {
	err := app.pushToQueue(l.Name, l.Data)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged via RabbitMQ"

	app.writeJson(w, http.StatusAccepted, payload)
}

func (app *Config) pushToQueue(name, msg string) error {
	emitter, err := event.NewEventEmitter(app.Rabbit)

	if err != nil {
		return err
	}

	payload := LogPyaload{
		Name: name,
		Data: msg,
	}

	j, _ := json.MarshalIndent(&payload, "", "\t")
	err = emitter.Push(string(j), "log.INFO")
	if err != nil {
		return err
	}
	return nil
}

type RPCPayload struct {
	Email    string
	Password string
}

func (app *Config) authViaRPC(w http.ResponseWriter, a AuthPayload) {
	client, err := rpc.Dial("tcp", "authentication:5001")

	if err != nil {
		app.errorJSON(w, err)
		return
	}
	rpcPayload := RPCPayload{
		Email:    a.Email,
		Password: a.Password,
	}
	var result string
	err = client.Call("RPCServer.Auth", rpcPayload, &result)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	payload := jsonResponse{
		Error:   false,
		Message: result,
	}
	app.writeJson(w, http.StatusAccepted, payload)
}

func (app *Config) AuthViaGRPC(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	conn, err := grpc.Dial(
		"authentication:50001",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock())

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	defer conn.Close()
	c := auth.NewAuthServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	_, err = c.CheckAuth(ctx, &auth.AuthRequest{
		AuthEntry: &auth.Auth{
			Email:    requestPayload.Auth.Email,
			Password: requestPayload.Auth.Password,
		},
	})

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Auth Success!"

	app.writeJson(w, http.StatusAccepted, payload)
}

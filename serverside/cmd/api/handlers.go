package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/miko2823/go-docker/pkg"
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

func HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := pkg.ReadJSON(w, r, &requestPayload)

	if err != nil {
		pkg.ErrorJSON(w, err)
		return
	}

	switch requestPayload.Action {

	case "auth":
		authenticate(w, requestPayload.Auth)

	default:
		pkg.ErrorJSON(w, errors.New("unknown action"))
	}
}

func authenticate(w http.ResponseWriter, a AuthPayload) {
	// call auth service
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	request, err := http.NewRequest(
		"POST",
		"http://authentication:8088/authenticate",
		bytes.NewBuffer(jsonData))

	if err != nil {
		pkg.ErrorJSON(w, err)
		return
	}

	client := &http.Client{}

	response, err := client.Do(request)

	if err != nil {
		pkg.ErrorJSON(w, err)
		return
	}

	defer response.Body.Close()

	// return status code
	if response.StatusCode == http.StatusUnauthorized {
		pkg.ErrorJSON(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		pkg.ErrorJSON(w, errors.New("error calling auth service"))
		return
	}

	var jsonFromService pkg.JsonResponse

	err = json.NewDecoder(response.Body).Decode(&jsonFromService)

	if err != nil {
		pkg.ErrorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		pkg.ErrorJSON(w, err, http.StatusUnauthorized)
		return
	}

	var payload pkg.JsonResponse
	payload.Error = false
	payload.Message = "Authenticated"
	payload.Data = jsonFromService.Data

	pkg.WriteJson(w, http.StatusAccepted, payload)
}

// func logEventViaRabbit(w http.ResponseWriter, l LogPyaload) {
// 	err := app.pushToQueue(l.Name, l.Data)

// 	if err != nil {
// 		pkg.ErrorJSON(w, err)
// 		return
// 	}

// 	var payload pkg.JsonResponse
// 	payload.Error = false
// 	payload.Message = "logged via RabbitMQ"

// 	pkg.WriteJson(w, http.StatusAccepted, payload)
// }

// func pushToQueue(name, msg string) error {
// 	emitter, err := event.NewEventEmitter(app.Rabbit)

// 	if err != nil {
// 		return err
// 	}

// 	payload := LogPyaload{
// 		Name: name,
// 		Data: msg,
// 	}

// 	j, _ := json.MarshalIndent(&payload, "", "\t")
// 	err = emitter.Push(string(j), "log.INFO")
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// type RPCPayload struct {
// 	Email    string
// 	Password string
// }

// func authViaRPC(w http.ResponseWriter, a AuthPayload) {
// 	client, err := rpc.Dial("tcp", "authentication:5001")

// 	if err != nil {
// 		pkg.ErrorJSON(w, err)
// 		return
// 	}
// 	rpcPayload := RPCPayload{
// 		Email:    a.Email,
// 		Password: a.Password,
// 	}
// 	var result string
// 	err = client.Call("RPCServer.Auth", rpcPayload, &result)
// 	if err != nil {
// 		pkg.ErrorJSON(w, err)
// 		return
// 	}
// 	payload := pkg.JsonResponse{
// 		Error:   false,
// 		Message: result,
// 	}
// 	pkg.WriteJson(w, http.StatusAccepted, payload)
// }

// func AuthViaGRPC(w http.ResponseWriter, r *http.Request) {
// 	var requestPayload RequestPayload
// 	err := pkg.ReadJSON(w, r, &requestPayload)
// 	if err != nil {
// 		pkg.ErrorJSON(w, err)
// 		return
// 	}
// 	conn, err := grpc.Dial(
// 		"authentication:50001",
// 		grpc.WithTransportCredentials(insecure.NewCredentials()),
// 		grpc.WithBlock())

// 	if err != nil {
// 		pkg.ErrorJSON(w, err)
// 		return
// 	}

// 	defer conn.Close()
// 	c := auth.NewAuthServiceClient(conn)
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

// 	defer cancel()

// 	_, err = c.CheckAuth(ctx, &auth.AuthRequest{
// 		AuthEntry: &auth.Auth{
// 			Email:    requestPayload.Auth.Email,
// 			Password: requestPayload.Auth.Password,
// 		},
// 	})

// 	if err != nil {
// 		pkg.ErrorJSON(w, err)
// 		return
// 	}

// 	var payload pkg.JsonResponse
// 	payload.Error = false
// 	payload.Message = "Auth Success!"

// 	pkg.WriteJson(w, http.StatusAccepted, payload)
// }

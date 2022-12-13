package main

import (
	"authentication/auth"
	"authentication/data"
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

type AuthServer struct {
	auth.UnimplementedAuthServiceServer
	Models data.Models
}

func (a *AuthServer) CheckAuth(ctx context.Context, req *auth.AuthRequest) (*auth.AuthResponse, error) {
	input := req.GetAuthEntry()

	user, err := a.Models.User.GetByEmail(input.Email)
	if err != nil {
		return &auth.AuthResponse{Result: "failed"}, err
	}

	valid, err := user.PasswordMatches(input.Password)
	if err != nil || !valid {
		return &auth.AuthResponse{Result: "failed"}, err
	}

	return &auth.AuthResponse{Result: "success??"}, nil

}

func (app *Config) gRPCListen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", gRpcPort))
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}

	s := grpc.NewServer()

	auth.RegisterAuthServiceServer(s, &AuthServer{Models: app.Models})

	log.Printf("gRPC Server started on port %s", gRpcPort)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}
}

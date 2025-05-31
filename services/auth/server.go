package main

import (
	"context"
	"errors"
	"log"
	"net"

	pb "github.com/KatyaPark11/Sudoku-Golang/generated/auth"
	"google.golang.org/grpc"
)

// Временное хранилище пользователей (в реальности — база данных)
var users = map[string]string{} // username -> password

type authServer struct {
	pb.UnimplementedAuthServiceServer
}

func (a *authServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	if _, exists := users[req.Username]; exists {
		return &pb.RegisterResponse{
			Success: false,
			Message: "User already exists",
		}, nil
	}
	users[req.Username] = req.Password
	return &pb.RegisterResponse{
		Success: true,
		Message: "Registration successful",
	}, nil
}

func (a *authServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	pass, exists := users[req.Username]
	if !exists || pass != req.Password {
		return &pb.LoginResponse{
			Success: false,
			Token:   "",
		}, errors.New("invalid credentials")
	}

	token := "dummy-token-for-" + req.Username // В реальности — JWT или другой токен.

	return &pb.LoginResponse{
		Success: true,
		Token:   token,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, &authServer{})

	log.Println("Auth service listening on :50052")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

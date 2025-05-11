package main

import (
	"log"
	"net"
	"net/http"

	"github.com/zubrodin/calc-service/internal/app"
	"github.com/zubrodin/calc-service/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	application := app.New(cfg)

	// Запуск HTTP сервера
	go func() {
		router := application.SetupRouter()
		log.Printf("Starting HTTP server on %s", cfg.ServerAddress)
		if err := http.ListenAndServe(cfg.ServerAddress, router); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Запуск gRPC сервера
	lis, err := net.Listen("tcp", cfg.GrpcAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := application.GRPCHandler()

	log.Printf("Starting gRPC server on %s", cfg.GrpcAddress)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve gRPC: %v", err)
	}
}

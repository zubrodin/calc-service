package config

import "os"

type Config struct {
	ServerAddress string
	GrpcAddress   string
	DatabasePath  string
}

func Load() (*Config, error) {
	databasePath := os.Getenv("DB_PATH")
	if databasePath == "" {
		databasePath = "./calc.db"
	}

	serverAddress := os.Getenv("SERVER_ADDRESS")
	if serverAddress == "" {
		serverAddress = ":8080"
	}

	grpcAddress := os.Getenv("GRPC_ADDRESS")
	if grpcAddress == "" {
		grpcAddress = ":50051"
	}

	return &Config{
		ServerAddress: serverAddress,
		GrpcAddress:   grpcAddress,
		DatabasePath:  databasePath,
	}, nil
}

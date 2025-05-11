package app

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/zubrodin/calc-service/internal/config"
	pb "github.com/zubrodin/calc-service/internal/grpc"
	"github.com/zubrodin/calc-service/internal/handler"
	"github.com/zubrodin/calc-service/internal/repository"
	"github.com/zubrodin/calc-service/internal/service"
	"github.com/zubrodin/calc-service/pkg/calculator"
	"github.com/zubrodin/calc-service/pkg/validator"
	"google.golang.org/grpc"
)

type calculatorServer struct {
	pb.UnimplementedCalculatorServer
	service *service.Service
	repo    repository.Repository
}

func (s *calculatorServer) GetTask(ctx context.Context, req *pb.TaskRequest) (*pb.TaskResponse, error) {
	task, err := s.repo.GetPendingTask()
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	if task == nil {
		return &pb.TaskResponse{}, nil
	}

	return &pb.TaskResponse{
		Id:        task.ID,
		Arg1:      task.Arg1,
		Arg2:      task.Arg2,
		Operation: task.Operation,
	}, nil
}

func (s *calculatorServer) SubmitResult(ctx context.Context, req *pb.ResultRequest) (*pb.ResultResponse, error) {
	if err := s.repo.SaveResult(req.Id, req.Result); err != nil {
		return &pb.ResultResponse{Success: false}, fmt.Errorf("failed to save result: %w", err)
	}
	return &pb.ResultResponse{Success: true}, nil
}

type App struct {
	config  *config.Config
	handler *handler.Handler
	service *service.Service
	repo    repository.Repository
}

func New(cfg *config.Config) *App {
	validator := validator.New()
	calculator := calculator.New()

	repo, err := repository.NewSQLiteRepository(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}

	service := service.New(calculator, validator, repo)
	handler := handler.New(service, repo)

	return &App{
		config:  cfg,
		handler: handler,
		service: service,
		repo:    repo,
	}
}

func (a *App) GRPCHandler() *grpc.Server {
	s := grpc.NewServer()
	pb.RegisterCalculatorServer(s, &calculatorServer{
		service: a.service,
		repo:    a.repo,
	})
	return s
}

func (a *App) SetupRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/register", a.handler.Register)
	mux.HandleFunc("/api/v1/login", a.handler.Login)
	mux.HandleFunc("/api/v1/calculate", a.handler.Authenticate(a.handler.Calculate))
	return mux
}

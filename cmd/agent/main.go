package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	pb "github.com/zubrodin/calc-service/internal/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	orchestratorAddr := os.Getenv("ORCHESTRATOR_ADDRESS")
	if orchestratorAddr == "" {
		orchestratorAddr = "localhost:50051"
	}

	conn, err := grpc.Dial(orchestratorAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithTimeout(5*time.Second))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewCalculatorClient(conn)

	for {
		task, err := client.GetTask(context.Background(), &pb.TaskRequest{})
		if err != nil {
			log.Printf("Error getting task: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		if task.Id == "" {
			time.Sleep(1 * time.Second)
			continue
		}

		result, err := calculate(task)
		if err != nil {
			log.Printf("Calculation error: %v", err)
			continue
		}

		_, err = client.SubmitResult(context.Background(), &pb.ResultRequest{
			Id:     task.Id,
			Result: result,
		})
		if err != nil {
			log.Printf("Error submitting result: %v", err)
		}

		time.Sleep(500 * time.Millisecond)
	}
}

func calculate(task *pb.TaskResponse) (float64, error) {
	arg1, err := strconv.ParseFloat(task.Arg1, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid arg1: %w", err)
	}

	arg2, err := strconv.ParseFloat(task.Arg2, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid arg2: %w", err)
	}

	switch task.Operation {
	case "+":
		return arg1 + arg2, nil
	case "-":
		return arg1 - arg2, nil
	case "*":
		return arg1 * arg2, nil
	case "/":
		if arg2 == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return arg1 / arg2, nil
	default:
		return 0, fmt.Errorf("unknown operation: %s", task.Operation)
	}
}

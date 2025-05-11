package repository

import (
	"errors"
	"time"
)

type User struct {
	ID       int
	Login    string
	Password string
}

type Task struct {
	ID          string
	UserID      int
	Expression  string
	Arg1        string
	Arg2        string
	Operation   string
	Result      float64
	Status      string
	CreatedAt   time.Time
	StartedAt   time.Time
	CompletedAt time.Time
}

type Repository interface {
	CreateUser(login, password string) (int64, error)
	Authenticate(login, password string) (*User, error)
	CreateTask(userID int, expr string) (string, error)
	GetPendingTask() (*Task, error)
	SaveResult(id string, result float64) error
	GetUserTasks(userID int) ([]Task, error)
	GetTaskByID(id string) (*Task, error)
}

var (
	ErrUserExists      = errors.New("user already exists")
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidPassword = errors.New("invalid password")
	ErrTaskNotFound    = errors.New("task not found")
)

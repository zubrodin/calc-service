package service

import (
	"github.com/zubrodin/calc-service/internal/storage"
	"github.com/zubrodin/calc-service/pkg/calculator"
	"github.com/zubrodin/calc-service/pkg/validator"
)

var ErrInvalidExpression = validator.ErrInvalidExpression

type Service struct {
	calculator *calculator.Calculator
	validator  *validator.Validator
	storage    storage.Storage
}

func New(calc *calculator.Calculator, valid *validator.Validator, storage storage.Storage) *Service {
	return &Service{
		calculator: calc,
		validator:  valid,
		storage:    storage,
	}
}

func (s *Service) Calculate(expr string) (float64, error) {
	if err := s.validator.Validate(expr); err != nil {
		return 0, ErrInvalidExpression
	}

	result, err := s.calculator.Calculate(expr)
	if err != nil {
		return 0, err
	}

	return result, nil
}

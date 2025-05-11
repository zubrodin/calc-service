package validator

import (
	"errors"
	"regexp"
)

var ErrInvalidExpression = errors.New("Expression is not valid")

type Validator struct {
	validExpr *regexp.Regexp
}

func New() *Validator {
	// Разрешаем цифры, пробелы, скобки и основные арифметические операции
	pattern := `^[\d\s\(\)\+\-\*\/\.]+$`
	return &Validator{
		validExpr: regexp.MustCompile(pattern),
	}
}

func (v *Validator) Validate(expr string) error {
	if !v.validExpr.MatchString(expr) {
		return ErrInvalidExpression
	}
	return nil
}

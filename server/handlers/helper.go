package handlers

import (
	"github.com/go-playground/validator"
	"go.uber.org/zap"
)

type Helper struct {
	l *zap.Logger
	v *validator.Validate
}

func NewHelper(l *zap.Logger) *Helper {
	v := validator.New()
	return &Helper{l, v}
}

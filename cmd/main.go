package main

import (
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"justcheck/internal/server"
)

func main() {
	logger, _ := zap.NewProduction()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {

		}
	}(logger)

	validate := validator.New()


	s := server.NewServer(logger,validate)
	err := s.Run()
	if err != nil {
		logger.Sugar().Error(err)
	}

}

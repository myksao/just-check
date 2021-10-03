package server

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"io"
	justCheckHandler "justcheck/internal/justcheck/delivery/http"
	"justcheck/internal/justcheck/usecase"
	"os"
)

type server struct {
	engine *gin.Engine
	logger *zap.Logger
	validator *validator.Validate
}

//NewServer constructor
func NewServer(logger *zap.Logger,validator *validator.Validate) *server{
	return &server{
		engine: gin.Default(),
		logger: logger,
		validator: validator,
	}
}

func (server *server) Run()  error{
	gin.SetMode(gin.ReleaseMode)
	// Disable Console Color, you don't need console color when writing the logs to file.
	gin.DisableConsoleColor()

	// Logging to a file.
	f, _ := os.Create("server.log")
	gin.DefaultWriter = io.MultiWriter(f)

	server.MapRoutes()

	justChkUC := usecase.NewJustCheckUseCase(server.logger,server.validator)
	justCheckH := justCheckHandler.NewJustCheckHandler(justChkUC,server.validator,server.engine,server.logger)
	justCheckH.MapJustCheckRoute()

	ginRunErr := make(chan error)
	go func() {
		err := server.engine.Run(":8080")
		if err != nil {
			ginRunErr <- errors.Wrap(err,"")
		}
	}()


	if ginRunErr!=nil{
		msg := <- ginRunErr
		server.logger.Error(msg.Error())
		return msg
	}

	return nil
}


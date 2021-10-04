package http

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"justcheck/internal/justcheck"
	"mime/multipart"
	"net/http"
)

type justCheckHandler struct {
	justCheckUC justcheck.UseCase
	validator *validator.Validate
	router *gin.Engine
	logger *zap.Logger
}

func NewJustCheckHandler(justCheckUC justcheck.UseCase,validator *validator.Validate,router *gin.Engine,logger *zap.Logger) *justCheckHandler{
	return &justCheckHandler{
		router: router,
		justCheckUC: justCheckUC,
		logger: logger,
		validator: validator,
	}
}

func (justChk *justCheckHandler) Check()  gin.HandlerFunc{
	return func(context *gin.Context) {
		type request struct{
			File *multipart.FileHeader `form:"file" binding:"required"`
			ContentType string `form:"type"`
		}

		var requestData request

		err := context.ShouldBind(&requestData)
		if err != nil {
			justChk.logger.Sugar().Error(err)
			context.JSON(http.StatusNotAcceptable,gin.H{
				"message":err.Error(),
			})
			return
		}

		justChk.logger.Sugar().Info(requestData)

		open, openFilErr := requestData.File.Open()
		if openFilErr != nil {
			justChk.logger.Sugar().Error(openFilErr)
			context.JSON(http.StatusNotAcceptable,gin.H{
				"message":openFilErr.Error(),
			})
			return
		}

		defer func(open multipart.File) {
			closeErr := open.Close()
			if closeErr != nil {
				justChk.logger.Sugar().Error(closeErr)
				//context.JSON(http.StatusNotAcceptable,gin.H{
				//	"message":closeErr.Error(),
				//})
				return
			}
		}(open)

		check, checkErr := justChk.justCheckUC.Check(context.Request.Context(),open,requestData.ContentType)
		if checkErr != nil {
			justChk.logger.Sugar().Error(checkErr.Error())
			context.JSON(http.StatusNotAcceptable,gin.H{
				"message":checkErr.Error(),
			})
			return
		}


		context.JSON(http.StatusOK,gin.H{
			"message":check,
		})



	}
}

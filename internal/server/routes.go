package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func (server *server) MapRoutes()  {

	server.engine.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {

		// your custom format
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	server.engine.HandleMethodNotAllowed = true

	server.engine.NoMethod(func(context *gin.Context) {
		context.JSON(
			http.StatusMethodNotAllowed,
			gin.H{
				"code": "PAGE_NOT_FOUND",
				"message": "Page not found 🤡",
			},
		)
	})

	server.engine.NoRoute(func(context *gin.Context) {
		context.JSON(
			http.StatusNotFound,
			gin.H{
				"code": "PAGE_NOT_FOUND",
				"message": "Page not found 🤡",
			},
		)
	})
}

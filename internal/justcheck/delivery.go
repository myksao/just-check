package justcheck

import "github.com/gin-gonic/gin"

type Delivery interface {
	Check() gin.HandlerFunc
}

package httpapi

import "github.com/gin-gonic/gin"

// RegisterRoutes 注册所有HTTP接口
func RegisterRoutes(r *gin.Engine) {

	v1 := r.Group("/v1")

	{
		//v1.POST("/chat", ChatHandler)
		v1.Any("/chat", ChatHandler)
		//v1.GET("/sessions/:id", GetSessionHandler)
	}
}

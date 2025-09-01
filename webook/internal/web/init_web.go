package web

import "github.com/gin-gonic/gin"

func RegisterRoutes() *gin.Engine {
	server := gin.Default()
	//REST 风格
	//server.PUT("/users",func (context *gin.Context) {
	//
	//})
	registerUserRoutes(server)
	return server
}

func registerUserRoutes(server *gin.Engine) {

}

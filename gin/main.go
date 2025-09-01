package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default() //Default方法返回Engine，逻辑上的服务器，用来监听端口

	router.GET("/hello", func(c *gin.Context) {
		//当一个HTTP请求用GET方法访问的时候，如果访问的的路径是/hello，就执行下面的代码
		c.String(http.StatusOK, "hello , go")
	})

	router.POST("post", func(c *gin.Context) {
		c.String(http.StatusOK, "this is post")
	})

	router.GET("/users/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "这里是参数路由"+name)
	})

	router.GET("/views/*.html", func(c *gin.Context) {
		page := c.Param(".html")
		c.String(http.StatusOK, "这里是通配符路由"+page)
	})
	//查询参数需要用到Query方法
	router.GET("/order", func(c *gin.Context) {
		oid := c.Query("id")
		c.String(http.StatusOK, "这里是查询参数 "+oid)
	})

	router.Run(":8080") // 监听并在 0.0.0.0:8080 上启动服务
}

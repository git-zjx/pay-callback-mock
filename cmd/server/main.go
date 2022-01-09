package main

import (
	"net/http"
	"payment-mocker/handlers"

	"github.com/gin-gonic/gin"
)

var resources = "../../resources"

func main() {
	r := gin.Default()

	r.LoadHTMLGlob(resources + "/html/*")
	r.Static("/assets", resources+"/")
	r.StaticFile("/favicon.ico", resources+"/favicon.ico")

	r.GET("/alipay", func(c *gin.Context) {
		c.HTML(http.StatusOK, "alipay.html", gin.H{})
	})

	r.POST("/alipay", handlers.AlipayHandler)

	r.GET("/wechat", func(c *gin.Context) {
		c.HTML(http.StatusOK, "wechat.html", gin.H{})
	})

	r.POST("/wechat", handlers.WechatHandler)

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	r.Run(":8081")
}

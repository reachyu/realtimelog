package httpserver

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	_const "reallog/common/const"
	"reallog/ws"
	"time"
)

/*
开启http server
封装对外提供的接口服务
*/

func InitHttpServer() {
	var router *gin.Engine
	// Engin gin.Default() 默认是加载了一些框架内置的中间件
	router = gin.Default()
	log.Println(">>>> HTTP Server starting <<<<")
	// gin.New() 没有默认加载框架内置的中间件，根据需要自己手动加载中间件
	//router := gin.New()
	// go语言panic() 的时候，造成崩溃退出。而gin.Recovery这个中间是处理这个异常然后返回http code 500
	router.Use(gin.Recovery())

	// 设置跨域
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,  // 这是允许访问所有域
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE", "OPTIONS"},   //服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
		AllowHeaders:     []string{"x-xq5-jwt", "Content-Type", "Origin", "Content-Length"},  // 允许跨域设置
		ExposeHeaders:    []string{"x-xq5-jwt"},  // 跨域关键设置 让浏览器可以解析
		AllowCredentials: true,  //  跨域请求是否需要带cookie信息 默认设置为true
		MaxAge:           12 * time.Hour,
	}))

	router.Static("/html", "./public")

	hub := ws.NewHub()
	go hub.Run()
	router.GET("/ws/logs/:trainId", func(c *gin.Context) { ws.ServeWs(hub, c) })

	router.Run(":" + _const.HTTP_SERVER_PORT)
}
package controller

import "github.com/gin-gonic/gin"

func StartServer() {
	//Set release mode
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	//person route
	Customer(router)

	router.Run()
}

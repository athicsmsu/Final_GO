package controller

import "github.com/gin-gonic/gin"

func StartServer() {
	//Set release mode
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	//customer route
	Customer(router)
	//product route
	Product(router)

	router.Run()
}

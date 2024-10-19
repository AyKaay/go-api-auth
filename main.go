package main

import (
	"user-auth-api/controllers"
	"user-auth-api/initializers"
	"user-auth-api/middlewares"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvs()
}

func main() {
	router := gin.Default()

	router.POST("/auth/signup", controllers.CreateUser)
	router.POST("/auth/login", controllers.Login)
	router.POST("/auth/logout", middlewares.CheckAuth, controllers.Logout)

	router.POST("/transaction/create", middlewares.CheckAuth, controllers.CreateTransaction)
	router.Run()
}

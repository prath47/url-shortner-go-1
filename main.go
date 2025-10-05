package main

import (
	"os"
	"url-shortner-1/controllers"
	"url-shortner-1/initializers"
	"url-shortner-1/middleware"

	"github.com/gin-gonic/gin"
)

func Initialize() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.MigrateDB()
}

func main() {
	Initialize()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default
	}

	router := gin.Default()

	// user Routes
	// login
	router.POST("/api/user/login", controllers.Login)

	// signup
	router.POST("/api/user/signup", controllers.Signup)

	router.POST("/api/user/check-login", middleware.IsAuthenticated, checkAuth)

	router.GET("/api/url/get-all-url", middleware.IsAuthenticated, controllers.GetAllUserUrl)

	router.POST("/api/url/generate-url", middleware.IsAuthenticated, controllers.EncryptTheUrl)
	
	router.GET("/api/url/redirectTo/:Url", middleware.IsAuthenticated, controllers.RedirectToPage)
	
	router.DELETE("/api/url/deleteUrl", middleware.IsAuthenticated, controllers.DeleteUrl)

	router.Run(":" + port)
}

func checkAuth(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hi man you are logged in",
	})
}

/*
url shortner

models
User
UrlMap

User {
	id int
	name string
	email string
	password string
}

UrlMap {
	id int
	FromUrl string unique
	ToUrl string
	userId id
}

User Routes
/api/user/login
/api/user/signup

Url Routes
/api/url/create -> post -> protected
/api/url/delete -> delete -> protected
/api/url/get/:urlUniqueId -> get -> not protected
/api/url/getUserUrls/:userId -> get -> protected
/api/url/redirect/:uniqueUrlId -> get -> not protected

*/

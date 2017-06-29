package router

import (
	"gopkg.in/gin-gonic/gin.v1"
	handler "github.com/NirmalVatsyayan/go-url-shortner/Handlers"
	middleware "github.com/NirmalVatsyayan/go-url-shortner/Middlewares"
)

func Routes(port string) *gin.Engine{

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	apis := router.Group("/apis/")

	apis.POST("/login", handler.LoginHandler)
	apis.POST("/register", handler.RegisterHandler)
	apis.GET("/redirect-url/:encodedUrl", handler.EncodedUrlRetriveHandler)


	auth := apis.Group("/auth")
	auth.Use(middleware.Auth())
	{
		auth.GET("/hello", handler.HelloHandler)
		auth.GET("/refresh_token", handler.RefreshHandler)
	}

	authorized := apis.Group("/user/")
	authorized.Use(middleware.Auth())
	{
		authorized.GET("/profile", handler.ProfileHandler)
	}

	shortner := apis.Group("/shortner/")
	shortner.Use(middleware.Auth())
	{
		shortner.POST("submit", handler.UrlPostHandler)
		shortner.GET("urls", handler.GetUserUrlsHandler)
		shortner.GET("urls/:encodedUrl", handler.GetUrlInfoHandler)
	}
	return router
}

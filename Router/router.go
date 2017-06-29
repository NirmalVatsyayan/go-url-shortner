package router

import (
	"gopkg.in/gin-gonic/gin.v1"
	handler "github.com/NirmalVatsyayan/UrlShortnerRepo/Handlers"
	middleware "github.com/NirmalVatsyayan/UrlShortnerRepo/Middlewares"
)

func Routes(port string) *gin.Engine{

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.POST("/login", handler.LoginHandler)
	router.POST("/register", handler.RegisterHandler)
	router.GET("/url/:encodedUrl", handler.EncodedUrlRetriveHandler)

	auth := router.Group("/auth")
	auth.Use(middleware.Auth())
	{
		auth.GET("/hello", handler.HelloHandler)
		auth.GET("/refresh_token", handler.RefreshHandler)
	}

	urlShortner := router.Group("/url-shortner")
	urlShortner.Use(middleware.Auth())
	{
		urlShortner.POST("/submit", handler.UrlPostHandler)
	}

	authorized := router.Group("/user/")
	authorized.Use(middleware.Auth())
	{
		authorized.GET("/profile", handler.ProfileHandler)
		authorized.GET("urls", handler.GetUserUrlsHandler)
		authorized.GET("urlinfo/:encodedUrl", handler.GetUrlInfoHandler)
	}

	return router
}

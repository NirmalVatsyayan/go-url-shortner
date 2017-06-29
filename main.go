package main

import (
	db "github.com/NirmalVatsyayan/UrlShortnerRepo/Database"
	router "github.com/NirmalVatsyayan/UrlShortnerRepo/Router"
)


func StartServer(){
	//port := os.Getenv("PORT")
	//if port == "" {
	//	port = "8000"
	//}
	port := ":8000"
	routes := router.Routes(port)
	routes.Run(port)
}

func main(){
	db.InitDB()
	StartServer()
}

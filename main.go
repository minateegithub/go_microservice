package main

import (
	"log"

	"github.com/joho/godotenv"

	"github.com/minateegithub/go_microservice/controllers"

	"github.com/gin-gonic/gin"
	"github.com/minateegithub/go_microservice/config"
	"github.com/minateegithub/go_microservice/routes"
)

func main() {
	//Load env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Database
	config.Connect()

	//Init validators
	controllers.InitValidator()

	// Init Router
	router := gin.Default()

	// Route Handlers / Endpoints
	routes.Routes(router)

	log.Fatal(router.Run(":4747"))
}

package config

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/minateegithub/go_microservice/controllers"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func Connect() {

	dbURL := os.Getenv("DB_URL")
	schemaName := os.Getenv("DB_DATA_BASE_NAME")

	// Database Config
	clientOptions := options.Client().ApplyURI(dbURL)
	client, err := mongo.NewClient(clientOptions)

	//Set up a context required by mongo.Connect
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)

	//To close the connection at the end
	defer cancel()

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal("Couldn't connect to the database", err)
	} else {
		log.Println("Connected!")
	}
	db := client.Database(schemaName)
	controllers.EnrolleeCollection(db)
	return
}

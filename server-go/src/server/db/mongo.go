package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Mongo *mongo.Database
var Ctx context.Context

func Conn() {
	clientOpts := options.Client().ApplyURI(fmt.Sprintf("mongodb+srv://%s:%s@cluster0.umd8v.mongodb.net/%s?retryWrites=true&w=majority", MONGO_USER, MONGO_PASSWORD, MONGO_DB_NAME))
	client, err := mongo.NewClient(clientOpts)
	if err != nil {
		log.Fatalln(err)
	}

	Ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = client.Connect(Ctx); err != nil {
		log.Fatalln(err)
	}

	log.Println("===== mongodb connected =====")
	Mongo = client.Database(MONGO_DB_NAME)
}

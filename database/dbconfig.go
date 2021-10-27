package database

import (
	"context"
	"log"
	"time"

	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func Connect() *mongo.Client  {

	MONGOURI, _ := os.LookupEnv("MONGO_URI")
    clientOptions := options.Client().ApplyURI(MONGOURI)
    // clientOptions := options.Client().ApplyURI("mongodb+srv://Admin:linktorydatabase@cluster0.1oyqm.mongodb.net/link-dir?retryWrites=true&w=majority")
    client, err := mongo.NewClient(clientOptions)
    if err != nil {
        log.Fatal(err)
    }
    ctx, _ := context.WithTimeout(context.Background(),10 * time.Second)
    err = client.Connect(ctx)
    if err != nil {
        log.Fatal(err)
    }
    ctxerr := client.Ping(ctx, readpref.Primary())
    if err != nil {
        log.Fatal("Couldn't connect to the database", ctxerr)
    } else {
        log.Println("Connected!")
    }
    return client
}

func GetCollection(collectionName string) *mongo.Collection{
    
    client := Connect()
    collection := client.Database("link-dir").Collection(collectionName)
    
    return collection
}


package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var DBclient = Connect()
var DB = DBclient.Database("link-dir")

func CreateDoc(collection string, filter bson.M ) (*mongo.InsertOneResult){
	coll := DB.Collection(collection)
	res, err := coll.InsertOne(context.TODO(), filter)
	if err != nil {
		return nil
	}
	return res
}


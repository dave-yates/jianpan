package db

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Config holds the mongo configuration details
var Config configuration
var chineseDB *mongo.Database
var collection *mongo.Collection

type configuration struct {
	User     string
	Password string
	URI      string
}

//SetupConfig sets the mongo config settings
func SetupConfig() {
	Config.User = "dave"
	Config.Password = "pass1"
	Config.URI = fmt.Sprintf("mongodb://%v:%v@localhost:27017/admin", Config.User, Config.Password)
}

//InitDB sets up the mongo database ready for the data import
func InitDB(ctx context.Context, client mongo.Client) error {

	chineseDB = client.Database("chinesedb")
	collection = chineseDB.Collection("items")

	return nil
}

//Insert inserts the translation into the mongoDB
func Insert(ctx context.Context, data bson.D) error {

	opts := options.Update().SetUpsert(true)

	update := bson.D{
		{Key: "$set", Value: data},
	}
	_, err := collection.UpdateOne(ctx, data, update, opts)

	if err != nil {
		log.Fatal(err)
	}
	return nil
}

//GetTranslations takes a context and a search term and returns the results from the database
func GetTranslations(ctx context.Context, search string) []bson.D {

	pattern := primitive.Regex{Pattern: "^" + search}
	filter := bson.D{
		{Key: "pinyin", Value: bson.D{
			{Key: "$regex", Value: pattern}},
		},
	}

	//options := bson.M

	cursor, err := collection.Find(ctx, filter) //, options*/)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer cursor.Close(ctx)

	var results []bson.D
	if err = cursor.All(ctx, &results); err != nil {
		fmt.Println(err)

		return nil
	}

	return results

}
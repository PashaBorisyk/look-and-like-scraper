package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"look-and-like-web-scrapper/models"
	"time"
)

var database *mongo.Database
var collectionsMap map[string]*Collection

type Collection struct {
	collection *mongo.Collection
}

func init() {

	mongoUri := "mongodb://localhost:27017"
	productDatabaseName := "look-and-like"

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoUri))
	if err != nil {
		log.Fatal("Unable to create Mongo client with address: ", mongoUri, "; ", err)
		return
	}

	ctx := createContext()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal("Unable to connect to Mongo database with address: ", mongoUri, "; ", err)
		return
	}

	database = client.Database(productDatabaseName)
	collectionsMap = make(map[string]*Collection)
}

func GetCollection(name string) *Collection {
	if collectionsMap[name] == nil {
		collectionsMap[name] = &Collection{database.Collection(name)}
	}
	return collectionsMap[name]
}

func (holder *Collection) Insert(docs ...interface{}) error {
	ctx := createContext()
	_, err := holder.collection.InsertMany(ctx, docs)
	if err != nil {
		log.Println("Error inserting document: ", err)
	}
	return err
}

func decodeMultipleResult(cursor *mongo.Cursor, foreach func(product models.Product, err error) error) error {
	ctx := createContext()
	var product models.Product
	for cursor.Next(ctx) {
		err := cursor.Decode(&product)
		if err != nil {
			log.Println("Unable to decode document: ", err)
		}
		//product.ID = cursor.Current.Lookup("_id").ObjectID()
		if foreach(product, err) != nil {
			break
		}
	}
	_ = cursor.Close(ctx)
	err := cursor.Err()
	if err != nil {
		log.Println(err)
	}
	return err
}

func createContext() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	return ctx
}

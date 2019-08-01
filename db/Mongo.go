package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"look-and-like-scraper/models"
	"time"
)

var database *mongo.Database
var collectionsMap map[string]*Collection

type Collection struct {
	collection *mongo.Collection
}

func init() {

	mongoUri := "mongodb://look-and-like-test:unE8DZr3T7yA6SLDPjknaT8Bj0MzLD4O4604EDq0OE44Lv9BxAslwWXTqLvJFzvqLCBoCDshGgUUJuKoahpT6w==@look-and-like-test.documents.azure.com:10255/?ssl=true&replicaSet=globaldb"
	productDatabaseName := "look-and-like-test"

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

func (holder *Collection) Insert(doc interface{}) (interface{}, error) {
	switch object := doc.(type) {
	case models.Product:
		object.ID = primitive.NewObjectID()
		doc = object
	}
	ctx := createContext()
	result, err := holder.collection.InsertOne(ctx,doc)
	if err != nil {
		log.Println("Error inserting document: ", err)
		return nil, err
	}
	return result.InsertedID, err
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

package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Disc struct {
	List     string `bson:"list"`
	Title    string `bson:"title"`
	Artist   string `bson:"artist"`
	Cover    string `bson:"cover,omitempty"`
	Released int    `bson:"released"`
}

var client *mongo.Client
var dbName string

func init() {
	dbName = os.Getenv("DB_NAME")
	mongoHostname := os.Getenv("MONGO_HOSTNAME")
	mongoPassword := os.Getenv("MONGO_PASSWORD")
	mongoUsername := os.Getenv("MONGO_USERNAME")

	fmtUri := "mongodb+srv://%s:%s@%s/%s?retryWrites=true&w=majority"
	url := fmt.Sprintf(fmtUri, mongoUsername, mongoPassword, mongoHostname, dbName)

	clientOptions := options.Client().ApplyURI(url)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	client = c
}

func GetDb() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	databases, err := client.ListDatabaseNames(ctx, bson.M{})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(databases)
}

func FindByList(list string) []Disc {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	coll := client.Database(dbName).Collection("disc-list")

	filter := bson.M{"list": list}
	opts := options.Find().SetProjection(bson.M{"_id": 0})
	cursor, err := coll.Find(ctx, filter, opts)

	if err != nil {
		log.Fatal(err)
	}

	var results []Disc

	if err = cursor.All(ctx, &results); err != nil {
		log.Fatal(err)
	}

	return results
}

func FindBy(list, title, artist string, released int) Disc {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	coll := client.Database(dbName).Collection("disc-list")

	filter := bson.M{
		"list":     list,
		"title":    title,
		"artist":   artist,
		"released": released,
	}

	opts := options.FindOne().SetProjection(bson.M{"_id": 0})

	var d Disc
	r := coll.FindOne(ctx, filter, opts)
	r.Decode(&d)

	return d
}

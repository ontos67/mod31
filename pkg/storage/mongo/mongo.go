package mongo

import (
	"GoNews/pkg/storage"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type Mong struct {
	client *mongo.Client
}

func New() (*Mong, error) {
	mongoOpts := options.Client().ApplyURI("mongodb://server:27017/")
	client, err := mongo.Connect(context.Background(), mongoOpts)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	db := Mong{
		client: client,
	}
	return &db, err
}
func (db *Mong) AddPost(p storage.Post) error {
	collection := db.client.Database("mongodb").Collection("posts")
	_, err := collection.InsertOne(context.Background(), p)
	if err != nil {
		return err
	}
	return nil
}
func (db *Mong) Posts() ([]storage.Post, error) {
	collection := db.client.Database("mongodb").Collection("posts")
	filter := bson.D{}
	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())
	var data []storage.Post
	for cur.Next(context.Background()) {
		var l storage.Post
		err := cur.Decode(&l)
		if err != nil {
			return nil, err
		}
		data = append(data, l)
	}
	return data, cur.Err()
}
func (db *Mong) DeletePost(p storage.Post) error {
	collection := db.client.Database("mongodb").Collection("posts")
	filter := bson.M{"_id": p.ID}
	_, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	return err
}
func (db *Mong) UpdatePost(p storage.Post) error {

	collection := db.client.Database("mongodb").Collection("posts")
	filter := bson.M{"_id": p.ID}
	res := collection.FindOneAndUpdate(context.Background(), filter, p)
	if res.Err() != nil {
		log.Fatal("Update failed", res.Err())
	}
	return res.Err()
}

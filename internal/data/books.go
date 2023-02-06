package data

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

type Book struct {
	Id       primitive.ObjectID `bson:"_id" json:"_id"`
	Title    string             `bson:"title" json:"title"`
	Created  time.Time          `json:"created" json:"created"`
	Author   string             `json:"author" bson:"author"`
	TextArea string             `bson:"textArea" json:"textArea"`
}

type BookModel struct {
	Collection *mongo.Collection
}

func (b BookModel) Get(id primitive.ObjectID) ([]Book, error) {

	if id.String() == "" {
		return nil, errors.New("id can't be empty")
	}

	filter := bson.D{{"_id", id}}
	cursor, err := b.Collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	var books []Book

	for cursor.Next(context.TODO()) {
		var result Book
		if err := cursor.Decode(&result); err != nil {
			log.Fatal(err)
		}
		books = append(books, result)
	}
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	return books, nil
}

func (b BookModel) Insert(book Book) (*mongo.InsertOneResult, error) {

	js, err := bson.Marshal(book)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := b.Collection.InsertOne(ctx, js)
	if err != nil {
		return nil, err
	}
	return result, nil

}

func (b BookModel) Update(id primitive.ObjectID, newBook Book) (*mongo.UpdateResult, error) {
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"title", newBook.Title},
		{"textArea", newBook.TextArea}}}}
	result, err := b.Collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (b BookModel) Delete(id primitive.ObjectID) (*mongo.DeleteResult, error) {
	filter := bson.D{{"_id", id}}
	result, err := b.Collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (b BookModel) GetFilteredData(title string) []Book {
	filter := bson.D{{"title", title}}
	cursor, err := b.Collection.Find(context.TODO(), filter)
	if err != nil {
		return nil
	}
	var books []Book

	for cursor.Next(context.TODO()) {
		var result Book
		if err := cursor.Decode(&result); err != nil {
			log.Fatal(err)
		}
		books = append(books, result)
	}
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	return books

}

package data

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	options2 "go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Book struct {
	Id       primitive.ObjectID `bson:"_id" json:"_id"`
	Title    string             `bson:"title" json:"title"`
	Created  time.Time          `json:"created" json:"created"`
	Author   string             `json:"author" bson:"author"`
	TextArea string             `bson:"textArea" json:"textArea"`
}
type BookWithoutText struct {
	Id       primitive.ObjectID `bson:"_id" json:"_id"`
	Title    string             `bson:"title" json:"title"`
	Created  time.Time          `json:"created" json:"created"`
	Author   string             `json:"author" bson:"author"`
	TextArea string             `bson:"textArea" json:"textArea"`
}

type BookModel struct {
	Collection *mongo.Collection
}

func (b BookModel) Get(id primitive.ObjectID) (any, error) {

	if id.String() == "" {
		return nil, errors.New("id can't be empty")
	}

	filter := bson.D{{"_id", id}}
	result := b.Collection.FindOne(context.Background(), filter)

	var book Book

	err := result.Decode(&book)
	if err != nil {
		return nil, err
	}

	return &book, nil
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

func (b BookModel) GetFilteredData(urlFilter string) ([]BookWithoutText, error) {
	var books []Book
	var test []BookWithoutText
	options := options2.Find().SetSort(bson.D{{"created", -1}})
	filter := bson.D{{"title", urlFilter}}
	cursor, err := b.Collection.Find(context.Background(), filter, options)
	if err != nil {
		return nil, err
	}

	for cursor.Next(context.Background()) {
		var result BookWithoutText
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		test = append(test, result)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	filter = bson.D{{"author", urlFilter}}
	cursor, err = b.Collection.Find(context.Background(), filter, options)

	if err != nil {
		return nil, err
	}
	for cursor.Next(context.Background()) {
		var result BookWithoutText
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		//books = append(books, result)
		for i, book := range books {
			if result.Id == book.Id {
				break
			} else {
				if i == len(books)-1 {
					test = append(test, result)
				}
			}
		}
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return test, nil
}

func (b BookModel) GetLatestBooks() ([]BookWithoutText, error) {
	var books []BookWithoutText

	sort := options2.Find().SetSort(bson.D{{"created", -1}})
	limit := options2.Find().SetLimit(20)
	cursor, err := b.Collection.Find(context.Background(), bson.D{}, limit, sort)
	if err != nil {
		return nil, err
	}

	for cursor.Next(context.Background()) {
		var result BookWithoutText
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		books = append(books, result)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return books, nil
}

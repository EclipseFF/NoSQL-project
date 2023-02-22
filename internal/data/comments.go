package data

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Comment struct {
	CommentId primitive.ObjectID `bson:"_id"`
	UserId    primitive.ObjectID `bson:"userId"`
	BookId    primitive.ObjectID `bson:"bookId"`
	Text      string             `bson:"text"`
}

type CommentModel struct {
	Collection *mongo.Collection
}

func (m CommentModel) Get(bookId primitive.ObjectID) ([]Comment, error) {
	filter := bson.D{{"_id", bookId.Hex()}}
	cursor, err := m.Collection.Find(context.Background(), filter)

	if err != nil {
		return nil, err

	}
	var result []Comment

	for cursor.Next(context.Background()) {
		var comm Comment
		if err := cursor.Decode(&comm); err != nil {
			return nil, err
		}
		result = append(result, comm)
	}

	return result, nil
}

func (m CommentModel) Insert(comment Comment) (*mongo.InsertOneResult, error) {
	js, err := bson.Marshal(comment)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := m.Collection.InsertOne(ctx, js)
	if err != nil {
		return nil, err
	}
	return result, nil
}

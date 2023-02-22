package data

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Comment struct {
	CommentId primitive.ObjectID
	UserId    primitive.ObjectID
	BookId    primitive.ObjectID
	Text      string
}

type CommentModel struct {
	Collection *mongo.Collection
}

/*func (m CommentModel) Get(bookId primitive.ObjectID) {
	booksStage := bson.D{{"$lookup", bson.D{{"from", "Books"}, {"localField", "bookId"}, {"foreignField", "_id"}, {"as", "book"}}}}
	usersStage := bson.D{{"$lookup", bson.D{{"from", "Users"}, {"localField", "userID"}, {"foreignField", "_id"}, {"as", "commentAuthor"}}}}
	result, err := m.Collection.Aggregate(context.Background(), mongo.Pipeline{booksStage, usersStage})
	if err != nil {
		return
	}
}*/

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

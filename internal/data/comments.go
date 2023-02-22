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
	lookup := bson.D{[
	bson.D{{"$lookup":{from:"users",
		localField:"author",
		foreignField:"_id",
		as:"PostAuthor"
	}}},
	bson.D{{"$lookup":{from:"users",
		localField:"comments.coment_creator",
		foreignField:"_id",
		as:"CommentAuthor"
	}}}
	]}
	result, err := m.Collection.Aggregate(context.Background())
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

package data

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Models struct {
	Books    BookModel
	Users    UserModel
	Comments CommentModel
}

func NewModels(client mongo.Client) Models {
	booksCollection := client.Database("ELibrary").Collection("Books")
	usersCollection := client.Database("ELibrary").Collection("Users")
	commentsCollection := client.Database("ELibrary").Collection("Comments")
	return Models{
		Books:    BookModel{Collection: booksCollection},
		Users:    UserModel{Collection: usersCollection},
		Comments: CommentModel{Collection: commentsCollection},
	}
}

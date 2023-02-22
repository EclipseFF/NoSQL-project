package main

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang/internal/data"
	"net/http"
)

func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserId string `json:"userId"`
		BookId string `json:"bookId"`
		Text   string `json:"text"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if input.Text == "" {
		app.emptyComment(w, r)
		return
	}

	var comment data.Comment
	comment.CommentId = primitive.NewObjectID()
	comment.BookId, err = primitive.ObjectIDFromHex(input.BookId)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	comment.UserId, err = primitive.ObjectIDFromHex(input.UserId)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	comment.Text = input.Text

	result, err := app.models.Comments.Insert(comment)
	if result.InsertedID != comment.CommentId {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusAccepted, envelope{"comment": comment}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}

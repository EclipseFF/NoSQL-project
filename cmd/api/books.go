package main

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang/internal/data"
	"net/http"
	"time"
)

func (app *application) createBookHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title    string `json:"title"`
		Author   string `json:"author"`
		TextArea string `json:"textArea"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	book := &data.Book{
		Id:       primitive.NewObjectID(),
		Title:    input.Title,
		Created:  time.Now(),
		Author:   input.Author,
		TextArea: input.TextArea,
	}

	result, err := app.models.Books.Insert(*book)

	if err != nil || result.InsertedID != book.Id {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusAccepted, envelope{"book": book}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getFilteredData(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title string `json:"title"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	books := app.models.Books.GetFilteredData(input.Title)

	err = app.writeJSON(w, http.StatusAccepted, envelope{"books": books}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

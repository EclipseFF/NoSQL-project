package main

import (
	"errors"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
	params := httprouter.ParamsFromContext(r.Context())
	filter := params.ByName("filter")
	if filter == "" {
		app.badRequestResponse(w, r, errors.New("invalid filter parameter"))
	}

	books, err := app.models.Books.GetFilteredData(filter)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	err = app.writeJSON(w, http.StatusAccepted, envelope{"books": books}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) getLatest(w http.ResponseWriter, r *http.Request) {
	books, err := app.models.Books.GetLatestBooks()
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	err = app.writeJSON(w, http.StatusAccepted, envelope{"books": books}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) getById(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	idFromParam := params.ByName("id")
	if idFromParam == "" {
		app.errorResponse(w, r, http.StatusBadRequest, "wrong id")
	}

	id, err := primitive.ObjectIDFromHex(idFromParam)
	book, err := app.models.Books.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}

	}

	err = app.writeJSON(w, http.StatusOK, envelope{"book": book}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

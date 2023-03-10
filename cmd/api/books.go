package main

import (
	"errors"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang/internal/data"
	"net/http"
	"strings"
	"time"
)

func (app *application) createBookHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Title       string `json:"title"`
		Author      string `json:"author"`
		Description string `json:"description"`
		TextArea    string `json:"textArea"`
		Url         string `json:"url"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	book := &data.Book{
		Id:          primitive.NewObjectID(),
		Title:       input.Title,
		Created:     time.Now(),
		Author:      input.Author,
		Description: input.Description,
		TextArea:    input.TextArea,
		Url:         input.Url,
	}

	result, err := app.models.Books.Insert(*book)

	if err != nil || result.InsertedID != book.Id {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusAccepted, envelope{"book": book}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) getFilteredData(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	filter := params.ByName("filter")
	filter = strings.ReplaceAll(filter, "_", " ")

	if filter == "" {
		app.badRequestResponse(w, r, errors.New("invalid filter parameter"))
		return
	}

	books, err := app.models.Books.GetFilteredData(filter)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusAccepted, envelope{"books": books}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}

func (app *application) getLatest(w http.ResponseWriter, r *http.Request) {
	books, err := app.models.Books.GetLatestBooks()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusAccepted, envelope{"books": books}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}

func (app *application) getById(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	idFromParam := params.ByName("id")
	if idFromParam == "" {
		app.errorResponse(w, r, http.StatusBadRequest, "wrong id")
		return
	}

	id, err := primitive.ObjectIDFromHex(idFromParam)
	book, err := app.models.Books.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			app.notFoundResponse(w, r)
			return
		}

		app.serverErrorResponse(w, r, err)
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"book": book}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}
func (app *application) getFavList(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	userId := params.ByName("userId")
	if userId == "" {
		app.badRequestResponse(w, r, errors.New("wrong id"))
		return
	}
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	user, err := app.models.Users.GetUserById(id)
	if errors.Is(err, mongo.ErrNoDocuments) {
		app.notFoundResponse(w, r)
		return
	}

	var books []data.BookWithoutText

	for _, id := range user.FavoriteBooks {
		result, err := app.models.Books.Get(id)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		var book data.BookWithoutText
		book.Id = result.Id
		book.Title = result.Title
		book.Created = result.Created
		book.Author = result.Author
		book.Description = result.Description
		book.Url = result.Url
		books = append(books, book)
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"favorite_books": books}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

package main

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"golang/internal/data"
	"net/http"
	"time"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		ID:            primitive.NewObjectID(),
		CreatedAt:     time.Now(),
		Name:          input.Name,
		Email:         input.Email,
		IsAuthor:      false,
		FavoriteBooks: nil,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	result, err := app.models.Users.Insert(user)

	if err != nil && mongo.IsDuplicateKeyError(err) {
		app.resourceAlreadyUsed(w, r, "email")
		return
	}

	if err != nil || result.InsertedID != user.ID {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusAccepted, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `bson:"email" json:"email"`
		Password string `bson:"password" json:"password"`
	}

	err := app.readJSON(w, r, &input)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user, err := app.models.Users.GetByEmail(input.Email)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = bcrypt.CompareHashAndPassword(user.Password.Get(), []byte(input.Password))
	if err != nil {
		app.userNotRegistered(w, r)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"success": true}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) addToFavoriteHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserId string `json:"userId"`
		BookId string `json:"bookId"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	userId, err := primitive.ObjectIDFromHex(input.UserId)
	bookId, err := primitive.ObjectIDFromHex(input.BookId)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user, err := app.models.Users.GetUserById(userId)

	if err != nil {
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
	}

	for _, book := range user.FavoriteBooks {
		if book == bookId {
			app.resourceAlreadyUsed(w, r, "book")
			return
		}
	}

	user.FavoriteBooks = append(user.FavoriteBooks, bookId)
	result, err := app.models.Users.Update(user)
	if result.ModifiedCount == 0 {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) removeFromFavoriteHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserId string `json:"userId"`
		BookId string `json:"bookId"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	userId, err := primitive.ObjectIDFromHex(input.UserId)
	bookId, err := primitive.ObjectIDFromHex(input.BookId)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user, err := app.models.Users.GetUserById(userId)

	if err != nil {
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
	}

	for i, book := range user.FavoriteBooks {
		if book == bookId {
			user.FavoriteBooks = append(user.FavoriteBooks[:i], user.FavoriteBooks[i+1:]...)
			result, err := app.models.Users.Update(user)
			if result.ModifiedCount == 0 {
				app.serverErrorResponse(w, r, err)
				return
			}
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}
			return
		}
	}

	app.badRequestResponse(w, r, errors.New("no books were deleted"))
}

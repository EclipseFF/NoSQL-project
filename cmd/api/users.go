package main

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		ID:        primitive.NewObjectID(),
		CreatedAt: time.Now(),
		Name:      input.Name,
		Email:     input.Email,
		IsAuthor:  false,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	result, err := app.models.Users.Insert(user)

	if err != nil || result.InsertedID != user.ID {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusAccepted, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
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
	fmt.Println(user)
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
	}
}

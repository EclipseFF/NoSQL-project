package main

import (
	"fmt"
	"net/http"
)

func (app *application) logError(r *http.Request, err error) {
	app.logger.Print(err)
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := envelope{"error": message}

	err := app.writeJSON(w, status, env, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	app.errorResponse(w, r, http.StatusNotFound, message)
}

func (app *application) userNotRegistered(w http.ResponseWriter, r *http.Request) {
	message := "the requested user could not be found"
	app.errorResponse(w, r, http.StatusNotFound, message)
}

func (app *application) resourceAlreadyUsed(w http.ResponseWriter, r *http.Request, resource string) {
	message := fmt.Sprintf("the requested %s is already used", resource)
	app.errorResponse(w, r, http.StatusConflict, message)
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)
	message := "the server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (app *application) emptyComment(w http.ResponseWriter, r *http.Request) {
	app.errorResponse(w, r, http.StatusBadRequest, "comment can't be empty")
}

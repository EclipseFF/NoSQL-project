package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.HandlerFunc(http.MethodPost, "/usersRegister", app.registerUserHandler)
	router.HandlerFunc(http.MethodPost, "/usersLogin", app.loginHandler)
	return router
}

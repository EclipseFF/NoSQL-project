package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.HandlerFunc(http.MethodPost, "/usersRegister", app.registerUserHandler)
	router.HandlerFunc(http.MethodPost, "/usersLogin", app.loginHandler)
	router.HandlerFunc(http.MethodPost, "/bookInsert", app.createBookHandler)
	router.HandlerFunc(http.MethodGet, "/getBook/:filter", app.getFilteredData)
	router.HandlerFunc(http.MethodGet, "/getBook", app.getLatest)
	router.HandlerFunc(http.MethodGet, "/getById/:id", app.getById)

	return app.recoverPanic(router)
}

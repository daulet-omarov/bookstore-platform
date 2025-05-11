package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodPost, "/users", app.createUserHandler)
	router.HandlerFunc(http.MethodGet, "/user/:id", app.showUserHandler)
	router.HandlerFunc(http.MethodPut, "/user/:id", app.updateUserHandler)
	router.HandlerFunc(http.MethodPost, "/user/authenticate", app.authenticateUserHandler)

	return router
}

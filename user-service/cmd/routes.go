package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	dynamic := alice.New(app.sessionManager.LoadAndSave, app.authenticate)

	router.Handler(http.MethodPost, "/users", dynamic.ThenFunc(app.createUserHandler))
	router.Handler(http.MethodGet, "/users/:id", dynamic.ThenFunc(app.showUserHandler))
	router.Handler(http.MethodPost, "/users/authenticate", dynamic.ThenFunc(app.authenticateUserHandler))

	protected := dynamic.Append(app.requireAuthentication)
	router.Handler(http.MethodPut, "/users/update", protected.ThenFunc(app.updateUserHandler))

	return router
}

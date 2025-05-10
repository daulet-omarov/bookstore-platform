package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	//router.NotFound = http.HandlerFunc(app.notFoundResponse)
	//router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	//
	//router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	//
	router.HandlerFunc(http.MethodGet, "/books", app.listBookHandler)
	router.HandlerFunc(http.MethodPost, "/books", app.createBookHandler)
	router.HandlerFunc(http.MethodGet, "/books/:id", app.showBookHandler)
	router.HandlerFunc(http.MethodPut, "/v1/books/:id", app.updateBookHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/books/:id", app.deleteBookHandler)

	return router
}

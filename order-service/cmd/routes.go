package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodPost, "/orders", app.createOrderHandler)
	router.HandlerFunc(http.MethodGet, "/orders/:id", app.showOrderHandler)
	router.HandlerFunc(http.MethodPut, "/orders/:id", app.updateOrderHandler)
	router.HandlerFunc(http.MethodDelete, "/orders/:id", app.deleteOrderHandler)
	router.HandlerFunc(http.MethodGet, "/users/:id/orders", app.showUserOrderHandler)

	return router
}

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
	protected := dynamic.Append(app.requireAuthentication)

	router.Handler(http.MethodPost, "/orders", protected.ThenFunc(app.createOrderHandler))
	router.Handler(http.MethodGet, "/orders/:id", protected.ThenFunc(app.showOrderHandler))
	router.Handler(http.MethodPut, "/orders/:id", protected.ThenFunc(app.updateOrderHandler))
	router.Handler(http.MethodDelete, "/orders/:id", protected.ThenFunc(app.deleteOrderHandler))
	router.Handler(http.MethodGet, "/users/me/orders", protected.ThenFunc(app.showUserOrderHandler))

	return router
}

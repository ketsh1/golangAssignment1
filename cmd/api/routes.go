package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/modules", app.createModuleInfoHandler)
	router.HandlerFunc(http.MethodGet, "/v1/modules/:id", app.getModuleInfoHandler)
	router.HandlerFunc(http.MethodPut, "/v1/modules/:id", app.editModuleInfoHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/modules/:id", app.deleteModuleInfoHandler)
	return router
}

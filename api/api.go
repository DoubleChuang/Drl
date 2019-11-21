package api

import (
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()
	router.GET("/", helloHandler)
	return router
}

func SendErrorResponse(w http.ResponseWriter, sc int, errMsg string) {
	w.WriteHeader(sc)
	io.WriteString(w, errMsg)
}

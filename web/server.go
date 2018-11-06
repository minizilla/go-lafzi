package web

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var r = mux.NewRouter()

// Server server instance.
var Server = &http.Server{
	Handler:      r,
	Addr:         ":8080",
	WriteTimeout: 15 * time.Second,
	ReadTimeout:  15 * time.Second,
}

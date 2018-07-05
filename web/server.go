package web

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var router = mux.NewRouter()

// Server server instance.
var Server = &http.Server{
	Handler:      router,
	Addr:         "127.0.0.1:8080",
	WriteTimeout: 15 * time.Second,
	ReadTimeout:  15 * time.Second,
}

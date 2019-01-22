package http

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Service ...
type Service func(r *mux.Router)

// NewServer ...
func NewServer(addr string, services ...Service) *http.Server {
	r := mux.NewRouter()

	services = append(services, asset, index, about)
	for _, service := range services {
		service(r)
	}

	return &http.Server{
		Handler:      r,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}

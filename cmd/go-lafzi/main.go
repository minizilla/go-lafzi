package main

import (
	"log"

	"github.com/billyzaelani/go-lafzi/web"
)

func main() {
	server := web.Server
	log.Printf("Listening on %s", server.Addr)
	log.Fatal(server.ListenAndServe())
}

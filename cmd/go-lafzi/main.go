package main

import (
	"log"

	"github.com/billyzaelani/go-lafzi/web"
)

func main() {
	log.Fatal(web.Server.ListenAndServe())
}

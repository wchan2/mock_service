package main

import (
	"log"
	"net/http"

	"github.com/wchan2/mockservice"
)

func main() {
	mockService, err := mockservice.New("/mocks")
	if err != nil {
		log.Fatalf("Failed to created mock service %s", err)
	}
	http.ListenAndServe(":8080", mockService)
}

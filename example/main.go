package main

import (
	"net/http"

	"github.com/wchan2/mock_service"
)

func main() {
	mockService := mock_service.New("/mocks")
	http.ListenAndServe(":8080", mockService)
}

package mockservice_test

import (
	"net/http"
	"testing"

	"github.com/wchan2/mock_service"
)

func TestLookup(t *testing.T) {
	mockEndpoint := &mockservice.MockEndpoint{
		Method:          http.MethodPost,
		Endpoint:        "/test/endpoint",
		StatusCode:      http.StatusOK,
		ResponseBody:    "sample response",
		ResponseHeaders: map[string]string{},
	}
	endpoints := mockservice.NewEndpoints()
	if err := endpoints.Create(mockEndpoint); err != nil {
		t.Errorf("Expected creating endpoint to return a nil error but got %s", err)
	}

	storedEp, err := endpoints.Lookup(http.MethodPost, "/test/endpoint")
	if err != nil {
		t.Errorf("Expected endpoint lookup error to be nil but got %s", err)
	}

	if storedEp != mockEndpoint {
		t.Errorf("Expected %+v endpoint but got %+v", mockEndpoint, storedEp)
	}
}

func TestLookupNonExistentEndpoint(t *testing.T) {
	endpoints := mockservice.NewEndpoints()

	ep, err := endpoints.Lookup("POST", "/does/not/exist")
	if ep != nil {
		t.Errorf("Expected nil endpoint but got %+v", ep)
	}

	if err != mockservice.ErrEndpointDoesNotExist {
		t.Errorf("Expected %s error but got %s", mockservice.ErrEndpointDoesNotExist, err)
	}
}

func TestCreateErrEmptyHTTPMethod(t *testing.T) {
	mockEndpoint := &mockservice.MockEndpoint{
		Method:          " ",
		Endpoint:        "/test/endpoint",
		StatusCode:      http.StatusOK,
		ResponseBody:    "sample response",
		ResponseHeaders: map[string]string{},
	}
	endpoints := mockservice.NewEndpoints()

	if err := endpoints.Create(mockEndpoint); err != mockservice.ErrEmptyHTTPMethod {
		t.Errorf("Expected %s error but got %s", mockservice.ErrEmptyHTTPMethod, err)
	}
}

func TestCreateErrEmptyEndpoint(t *testing.T) {
	mockEndpoint := &mockservice.MockEndpoint{
		Method:          "POST",
		Endpoint:        " ",
		StatusCode:      http.StatusOK,
		ResponseBody:    "sample response",
		ResponseHeaders: map[string]string{},
	}
	endpoints := mockservice.NewEndpoints()

	if err := endpoints.Create(mockEndpoint); err != mockservice.ErrEmptyEndpoint {
		t.Errorf("Expected %s error but got %s", mockservice.ErrEmptyEndpoint, err)
	}
}

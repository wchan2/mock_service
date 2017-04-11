package mock_service_test

import (
	"net/http"
	"testing"

	"github.com/wchan2/mock_service"
)

func TestLookup(t *testing.T) {
	mockEndpoint := &mock_service.MockEndpoint{
		Method:          http.MethodPost,
		Endpoint:        "/test/endpoint",
		StatusCode:      http.StatusOK,
		ResponseBody:    "sample response",
		ResponseHeaders: map[string]string{},
	}
	endpoints := mock_service.NewEndpoints()
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
	endpoints := mock_service.NewEndpoints()

	ep, err := endpoints.Lookup("POST", "/does/not/exist")
	if ep != nil {
		t.Errorf("Expected nil endpoint but got %s", ep)
	}

	if err != mock_service.ErrEndpointDoesNotExist {
		t.Errorf("Expected %s error but got %s", mock_service.ErrEndpointDoesNotExist, err)
	}
}

func TestCreateErrEmptyHTTPMethod(t *testing.T) {
	mockEndpoint := &mock_service.MockEndpoint{
		Method:          " ",
		Endpoint:        "/test/endpoint",
		StatusCode:      http.StatusOK,
		ResponseBody:    "sample response",
		ResponseHeaders: map[string]string{},
	}
	endpoints := mock_service.NewEndpoints()

	if err := endpoints.Create(mockEndpoint); err != mock_service.ErrEmptyHTTPMethod {
		t.Errorf("Expected %s error but got %s", mock_service.ErrEmptyHTTPMethod, err)
	}
}

func TestCreateErrEmptyEndpoint(t *testing.T) {
	mockEndpoint := &mock_service.MockEndpoint{
		Method:          "POST",
		Endpoint:        " ",
		StatusCode:      http.StatusOK,
		ResponseBody:    "sample response",
		ResponseHeaders: map[string]string{},
	}
	endpoints := mock_service.NewEndpoints()

	if err := endpoints.Create(mockEndpoint); err != mock_service.ErrEmptyEndpoint {
		t.Errorf("Expected %s error but got %s", mock_service.ErrEmptyEndpoint, err)
	}
}

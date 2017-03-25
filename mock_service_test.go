package mock_service_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/wchan2/mock_service"
)

func TestServeEndpointRegistration_NilReqBody(t *testing.T) {
	service := mock_service.New("/mocks")
	req, err := http.NewRequest(http.MethodPost, "/mocks", nil)
	if err != nil {
		t.Fatalf("Expected error to create new request to be nil but got %s", err)
	}

	recorder := httptest.NewRecorder()
	service.ServeHTTP(recorder, req)
	if recorder.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("Expected %d but received %d", http.StatusBadRequest, recorder.Result().StatusCode)
	}

	if recorder.Body.String() != "Registering an endpoint requires a payload" {
		t.Errorf(
			`Expected message: "%s" but received "%s"`,
			"Registering an endpoint requires a payload",
			recorder.Body.String(),
		)
	}
}

func TestServeEndpointRegistration_InvalidJSONReqBody(t *testing.T) {
	service := mock_service.New("/mocks")
	req, err := http.NewRequest(http.MethodPost, "/mocks", strings.NewReader(""))
	if err != nil {
		t.Fatalf("Expected error to create new request to be nil but got %s", err)
	}

	recorder := httptest.NewRecorder()
	service.ServeHTTP(recorder, req)
	if recorder.Result().StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected %d but received %d", http.StatusBadRequest, recorder.Result().StatusCode)
	}

	if recorder.Body.String() != "Unable to Unmarshal request body : unexpected end of JSON input" {
		t.Errorf(
			`Expected message: "%s" but received "%s"`,
			"Unable to Unmarshal request body : unexpected end of JSON input",
			recorder.Body.String(),
		)
	}
}

func TestServeEndpointRegistration_Success(t *testing.T) {
	service := mock_service.New("/mocks")
	const jsonRequestBody = `{"method": "GET", "endpoint": "/mock/test", "httpStatusCode": 201, "responseBody": "hello world", "responseHeaders": {"Foo": "Bar"}}`
	req, err := http.NewRequest(http.MethodPost, "/mocks", strings.NewReader(jsonRequestBody))
	if err != nil {
		t.Fatalf("Expected error to create new request to be nil but got %s", err)
	}

	recorder := httptest.NewRecorder()
	service.ServeHTTP(recorder, req)
	if recorder.Result().StatusCode != http.StatusCreated {
		t.Errorf("Expected %i status but got %i", http.StatusCreated, recorder.Result().StatusCode)
	}
}

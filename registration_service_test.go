package mockservice_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/wchan2/mockservice"

	"testing"
)

type errReader struct{}

func (e *errReader) Read(b []byte) (n int, err error) {
	err = errors.New("read error")
	return
}

func TestRegistrationService_ServeHTTP(t *testing.T) {
	t.Run("Endpoint_registration_with_no_request_body", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/test", nil)
		if err != nil {
			t.Errorf("Expected new request error to be nil but got %s", err)
		}

		endpoints := mockservice.NewEndpoints()
		svc := mockservice.NewRegistrationService(endpoints)
		recorder := httptest.NewRecorder()
		svc.ServeHTTP(recorder, req)
		if recorder.Code != http.StatusBadRequest {
			t.Errorf("Expected response status code to be %d but got %d", http.StatusBadRequest, recorder.Code)
		}
	})

	t.Run("Endpoint_registration_read_body_error", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/test", &errReader{})
		if err != nil {
			t.Errorf("Expected new request error to be nil but got %s", err)
		}

		endpoints := mockservice.NewEndpoints()
		svc := mockservice.NewRegistrationService(endpoints)
		recorder := httptest.NewRecorder()
		svc.ServeHTTP(recorder, req)
		if recorder.Code != http.StatusBadRequest {
			t.Errorf("Expected response status code to be %d but got %d", http.StatusBadRequest, recorder.Code)
		}
	})

	t.Run("Endpoint_registration_with_bad_request_body", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/test", strings.NewReader("bad request body"))
		if err != nil {
			t.Errorf("Expected new request error to be nil but got %s", err)
		}

		endpoints := mockservice.NewEndpoints()
		svc := mockservice.NewRegistrationService(endpoints)
		recorder := httptest.NewRecorder()
		svc.ServeHTTP(recorder, req)
		if recorder.Code != http.StatusInternalServerError {
			t.Errorf("Expected response status code to be %d but got %d", http.StatusInternalServerError, recorder.Code)
		}
	})

	t.Run("Endpoint_registration_with_empty_HTTP_method", func(t *testing.T) {
		emptyMethodJSONBody := `{
			"method": " ",
			"endpoint": "/test",
			"httpStatusCode": 204,
			"responseBody": "",
			"responseHeaders": { "foo": "bar" }
		}`

		req, err := http.NewRequest(http.MethodPost, "/test", strings.NewReader(emptyMethodJSONBody))
		if err != nil {
			t.Errorf("Expected new request error to be nil but got %s", err)
		}

		endpoints := mockservice.NewEndpoints()
		svc := mockservice.NewRegistrationService(endpoints)
		recorder := httptest.NewRecorder()
		svc.ServeHTTP(recorder, req)
		if recorder.Code != http.StatusBadRequest {
			t.Errorf("Expected response status code to be %d but got %d", http.StatusBadRequest, recorder.Code)
		}
	})

	t.Run("Endpoint_registration_with_empty_endpoint", func(t *testing.T) {
		emptyEndpointJSONBody := `{
			"method": "GET",
			"endpoint": " ",
			"httpStatusCode": 204,
			"responseBody": "",
			"responseHeaders": { "foo": "bar" }
		}`

		req, err := http.NewRequest(http.MethodPost, "/test", strings.NewReader(emptyEndpointJSONBody))
		if err != nil {
			t.Errorf("Expected new request error to be nil but got %s", err)
		}

		endpoints := mockservice.NewEndpoints()
		svc := mockservice.NewRegistrationService(endpoints)
		recorder := httptest.NewRecorder()
		svc.ServeHTTP(recorder, req)
		if recorder.Code != http.StatusBadRequest {
			t.Errorf("Expected response status code to be %d but got %d", http.StatusBadRequest, recorder.Code)
		}
	})
}

package mockservice_test

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/wchan2/mock_service"
)

func TestEndpointService_ServeHTTP(t *testing.T) {
	mockEndpoint := mockservice.MockEndpoint{
		Method:          http.MethodGet,
		Endpoint:        "/test",
		StatusCode:      http.StatusOK,
		ResponseBody:    "test",
		ResponseHeaders: map[string]string{"foo": "bar"},
	}
	endpoints := mockservice.NewEndpoints()
	endpoints.Create(&mockEndpoint)
	svc := mockservice.NewEndpointService(endpoints)

	t.Run("Endpoint_found", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/test", nil)
		if err != nil {
			t.Errorf("Expected new request error to be nil but got %s", err)
		}

		svc.ServeHTTP(recorder, req)
		if recorder.Code != http.StatusOK {
			t.Errorf("Expected response status code to be %d but got %d", http.StatusOK, recorder.Code)
		}

		if reflect.DeepEqual(recorder.Header(), mockEndpoint.ResponseHeaders) {
			t.Errorf("Expected %+v but got %+v", mockEndpoint.ResponseHeaders, recorder.Header())
		}

		if recorder.Body.String() != mockEndpoint.ResponseBody {
			t.Errorf("Expected %s but got %s", mockEndpoint.ResponseBody, recorder.Body.String())
		}
	})

	t.Run("Endpoint_not_found", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/not-found", nil)
		if err != nil {
			t.Errorf("Expected new request error to be nil but got %s", err)
		}

		svc.ServeHTTP(recorder, req)
		if recorder.Code != http.StatusNotFound {
			t.Errorf("Expected response status code to be %d but got %d", http.StatusNotFound, recorder.Code)
		}

		if reflect.DeepEqual(recorder.Header(), mockEndpoint.ResponseHeaders) {
			t.Errorf("Expected %+v but got %+v", mockEndpoint.ResponseHeaders, recorder.Header())
		}

		if recorder.Body.String() != "Endpoint does not exist\n" {
			t.Errorf(`Expected "Endpoint does not exist\n" but got %s`, recorder.Body.String())
		}
	})
}

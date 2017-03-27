package mock_service_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/wchan2/mock_service"
)

const successfulRegistrationRequest = `{"method": "GET", "endpoint": "/mock/test", "httpStatusCode": 203, "responseBody": "hello world", "responseHeaders": {"Foo": "Bar"}}`

func TestServeEndpointRegistration_NilReqBody(t *testing.T) {
    conf := mock_service.NewConf("","")
	service := mock_service.New("/mocks", conf)
	req, err := http.NewRequest(http.MethodPost, "/mocks", nil)
	if err != nil {
		t.Fatalf("Expected error to create new request to be nil but got %s", err)
	}

	recorder := httptest.NewRecorder()
	service.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusBadRequest {
		t.Errorf("Expected %d but received %d", http.StatusBadRequest, recorder.Code)
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
    conf := mock_service.NewConf("","")
	service := mock_service.New("/mocks", conf)
	req, err := http.NewRequest(http.MethodPost, "/mocks", strings.NewReader(""))
	if err != nil {
		t.Fatalf("Expected error to create new request to be nil but got %s", err)
	}

	recorder := httptest.NewRecorder()
	service.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusInternalServerError {
		t.Errorf("Expected %d but received %d", http.StatusBadRequest, recorder.Code)
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
    conf := mock_service.NewConf("","")
	service := mock_service.New("/mocks", conf)
	req, err := http.NewRequest(http.MethodPost, "/mocks", strings.NewReader(successfulRegistrationRequest))
	if err != nil {
		t.Fatalf("Expected error to create new request to be nil but got %s", err)
	}

	recorder := httptest.NewRecorder()
	service.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusCreated {
		t.Errorf("Expected %i status but got %i", http.StatusCreated, recorder.Code)
	}
}
func TestServeMockHTTP(t *testing.T) {
    conf := mock_service.NewConf("","")
	service := mock_service.New("/mocks", conf)
	req, err := http.NewRequest(http.MethodPost, "/mocks", strings.NewReader(successfulRegistrationRequest))
	if err != nil {
		t.Fatalf("Expected error to create new request to be nil but got %s", err)
	}

	recorder := httptest.NewRecorder()
	service.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusCreated {
		t.Errorf("Expected %d status but got %d when registering the mock endpoint", http.StatusCreated, recorder.Code)
	}

	// test the mock endpoint
	testReq, err := http.NewRequest("GET", "/mock/test", nil)
	if err != nil {
		t.Fatalf("Expected error to create new request ot be nil but got %s", err)
	}
	testRecorder := httptest.NewRecorder()
	service.ServeHTTP(testRecorder, testReq)
	if testRecorder.Code != http.StatusNonAuthoritativeInfo {
		t.Errorf("Expected %d status but got %d when sending a mock request", http.StatusNonAuthoritativeInfo, testRecorder.Code)
	}

	if testRecorder.Body.String() != "hello world" {
		t.Errorf(`Expected "%s" response body but got "%s"`, "hello world", testRecorder.Body.String())
	}
}

func TestPreloadEndpointsFromConf(t *testing.T){
    conf := mock_service.NewConf("test/test_mocker.json","JSON")
    service := mock_service.New("/mocks", conf)
	testReq, err := http.NewRequest("GET", "/mock/conf_test1", nil)
	if err != nil {
		t.Fatalf("Expected error to create new request ot be nil but got %s", err)
	}
	testRecorder := httptest.NewRecorder()
	service.ServeHTTP(testRecorder, testReq)
    if testRecorder.Code != http.StatusOK {
		t.Errorf("Expected %d status but got %d when sending a mock request", http.StatusOK, testRecorder.Code)
    }
    if testRecorder.Body.String() != "Configurable mocker" {
		t.Errorf(`Expected "%s" response body but got "%s"`, "Configurable mocker", testRecorder.Body.String())
    }
	testRecorder = httptest.NewRecorder()
	testReq, err = http.NewRequest("GET", "/mock/conf_test2", nil)
	service.ServeHTTP(testRecorder, testReq)
    if testRecorder.Code != http.StatusOK {
		t.Errorf("Expected %d status but got %d when sending a mock request", http.StatusOK, testRecorder.Code)
    }
    if testRecorder.Body.String() != "Configurable mocker" {
		t.Errorf(`Expected "%s" response body but got "%s"`, "Configurable mocker", testRecorder.Body.String())
    }

}

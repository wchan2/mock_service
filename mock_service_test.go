package mock_service_test

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/wchan2/mock_service"
)

const (
	successfulRegistrationRequest    = `{"method": "GET", "endpoint": "/mock/test", "httpStatusCode": 203, "responseBody": "hello world", "responseHeaders": {"Foo": "Bar"}}`
	emptyMethodRegistrationRequest   = `{"method": " ", "endpoint": "/mock/test", "httpStatusCode": 203, "responseBody": "hello world", "responseHeaders": {"Foo": "Bar"}}`
	emptyEndpointRegistrationRequest = `{"method": "GET", "endpoint": " ", "httpStatusCode": 203, "responseBody": "hello world", "responseHeaders": {"Foo": "Bar"}}`
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

func TestServeEndpointRegistration_Success(t *testing.T) {
	service, err := mock_service.New("/mocks")
	if err != nil {
		t.Errorf("Expected err in creating new mock service to be nil but got %s", err)
	}
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

func TestServeMockHTTP_MockCreated(t *testing.T) {
	service, err := mock_service.New("/mocks")
	if err != nil {
		t.Errorf("Expected err in creating new mock service to be nil but got %s", err)
	}
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

func TestServeMockHTTP_MockCreatedWithConf(t *testing.T) {
	conf := mock_service.MockServiceConf{
		RegistrationEndpoint: "/mocks",
		Endpoints: []mock_service.MockEndpoint{
			{
				Method:          http.MethodPost,
				Endpoint:        "/mock/test",
				StatusCode:      http.StatusCreated,
				ResponseBody:    `hello world`,
				ResponseHeaders: map[string]string{"Foo": "Bar"},
			},
		},
	}
	service, err := mock_service.NewWithConf(&conf)
	if err != nil {
		t.Errorf("Expected err to be nil when creating the mock service but got %s", err)
	}
	recorder := httptest.NewRecorder()

	// test that a mocked endpoint is mocked
	req, err := http.NewRequest(http.MethodPost, "/mock/test", nil)
	if err != nil {
		t.Fatalf("Expected error to create new request to be nil but got %s", err)
	}
	service.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusCreated {
		t.Errorf("Expected %d status but got %d", http.StatusCreated, recorder.Code)
	}

	if recorder.Body.String() != "hello world" {
		t.Errorf(`Expected "%s" response body but got "%s"`, "hello world", recorder.Body.String())
	}
}

func TestServeMockHTTP_MockEndpointDoesNotExist(t *testing.T) {
	service, err := mock_service.New("/mocks")
	if err != nil {
		t.Errorf("Expected err in creating new mock service to be nil but got %s", err)
	}

	recorder := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/mock/test", nil)
	service.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusNotFound {
		t.Errorf("Expected %d status but got %d", http.StatusNotFound, recorder.Code)
	}

	if recorder.Body.String() != mock_service.ErrEndpointDoesNotExist.Error() {
		t.Errorf(`Expected "%s" response body but got "%s"`, mock_service.ErrEndpointDoesNotExist, recorder.Body.String())
	}
}

func TestMockService_EmptyRegistrationEndpoint(t *testing.T) {
	service, err := mock_service.New(" ")
	if err != mock_service.ErrEmptyRegistrationEndpoint {
		t.Errorf(`Expected err to be "%s" but received "%s"`, mock_service.ErrEmptyRegistrationEndpoint, err)
	}

	if service != nil {
		t.Errorf("Expected service to be nil when receiving an error but got %+v", service)
	}
}

func TestMockService_RegisterEmptyEndpoint(t *testing.T) {
	service, err := mock_service.New("/mocks")
	if err != nil {
		t.Errorf("Expected err in creating new mock service to be nil but got %s", err)
	}

	recorder := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/mocks", strings.NewReader(emptyEndpointRegistrationRequest))
	service.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusBadRequest {
		t.Errorf("Expected %d status but got %d", http.StatusBadRequest, recorder.Code)
	}

	if recorder.Body.String() != mock_service.ErrEmptyEndpoint.Error() {
		t.Errorf(`Expected "%s" response body but got "%s"`, mock_service.ErrEmptyEndpoint, recorder.Body.String())
	}
}

func TestMockService_RegisterEmptyMethod(t *testing.T) {
	service, err := mock_service.New("/mocks")
	if err != nil {
		t.Errorf("Expected err in creating new mock service to be nil but got %s", err)
	}

	recorder := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/mocks", strings.NewReader(emptyMethodRegistrationRequest))
	service.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusBadRequest {
		t.Errorf("Expected %d status but got %d", http.StatusBadRequest, recorder.Code)
	}

	if recorder.Body.String() != mock_service.ErrEmptyHTTPMethod.Error() {
		t.Errorf(`Expected "%s" response body but got "%s"`, mock_service.ErrEmptyHTTPMethod, recorder.Body.String())
	}
}

func TestNewWithConf_EmptyRegisterEndpointInConf(t *testing.T) {
	conf := mock_service.MockServiceConf{
		RegistrationEndpoint: "",
		Endpoints: []mock_service.MockEndpoint{
			{
				Method:          http.MethodPost,
				Endpoint:        "/mock/test",
				StatusCode:      http.StatusCreated,
				ResponseBody:    `hello world`,
				ResponseHeaders: map[string]string{"Foo": "Bar"},
			},
		},
	}
	service, err := mock_service.NewWithConf(&conf)
	if err != mock_service.ErrEmptyRegistrationEndpoint {
		t.Errorf(`Expected err to be "%s" but received "%s"`, mock_service.ErrEmptyRegistrationEndpoint, err)
	}

	if service != nil {
		t.Errorf("Expected service to be nil when receiving an error but got %+v", service)
	}
}

func TestNewWithConf_EmptyHTTPMethod(t *testing.T) {
	conf := mock_service.MockServiceConf{
		RegistrationEndpoint: "/mocks",
		Endpoints: []mock_service.MockEndpoint{
			{
				Method:          " ",
				Endpoint:        "/mock/test",
				StatusCode:      http.StatusCreated,
				ResponseBody:    `hello world`,
				ResponseHeaders: map[string]string{"Foo": "Bar"},
			},
		},
	}
	service, err := mock_service.NewWithConf(&conf)
	if err != mock_service.ErrEmptyHTTPMethod {
		t.Errorf(`Expected "%s" error when creating a mock service with empty http method in the configuration but got %v`, mock_service.ErrEmptyHTTPMethod, err)
	}

	if service != nil {
		t.Errorf("Expected service to be nil when receiving an error but got %+v", service)
	}
}

func TestNewWithConf_EmptyEndpoint(t *testing.T) {
	conf := mock_service.MockServiceConf{
		RegistrationEndpoint: "/mocks",
		Endpoints: []mock_service.MockEndpoint{
			{
				Method:          http.MethodPost,
				Endpoint:        " ",
				StatusCode:      http.StatusCreated,
				ResponseBody:    `hello world`,
				ResponseHeaders: map[string]string{"Foo": "Bar"},
			},
		},
	}
	service, err := mock_service.NewWithConf(&conf)
	if err != mock_service.ErrEmptyEndpoint {
		t.Errorf(`Expected "%s" error when creating a mock service with empty endpoint in the configuration but got %v`, mock_service.ErrEmptyEndpoint, err)
	}

	if service != nil {
		t.Errorf("Expected service to be nil when receiving an error but got %+v", service)
	}
}

func TestServeEndpointRegistration_NilReqBody(t *testing.T) {
	service, err := mock_service.New("/mocks")
	if err != nil {
		t.Errorf("Expected err in creating new mock service to be nil but got %s", err)
	}
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
	service, err := mock_service.New("/mocks")
	if err != nil {
		t.Errorf("Expected err in creating new mock service to be nil but got %s", err)
	}
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

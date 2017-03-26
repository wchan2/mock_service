package mock_service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

var ErrEndpointDoesNotExist = errors.New("Endpoint does not exist")

type MockService struct {
	createMockEndpoint string
	mockEndpoints      map[string]map[string]MockEndpoint

}

type MockEndpoint struct {
	Method          string            `json:"method"`
	Endpoint        string            `json:"endpoint"`
	StatusCode      int               `json:"httpStatusCode"`
	ResponseBody    string            `json:"responseBody"`
	ResponseHeaders map[string]string `json:"responseHeaders"`
	sync.Mutex

	// TODO:
	// AcceptHeaders   map[string]string `json:"responseHeaders"`
	// RequestBody     string            `json:"requestBody"`
}

func New(createMockEndpoint string) *MockService {
	return &MockService{
		createMockEndpoint: createMockEndpoint,
		mockEndpoints:      map[string]map[string]MockEndpoint{},
	}
}

func (m *MockService) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost && req.URL.Path == m.createMockEndpoint {
		m.serveRegistrationHTTP(w, req)
		return
	}

	m.serveMockHTTP(w, req)
}

func (m *MockService) CreateMockEndpoint(endpoint MockEndpoint) {
	if _, ok := m.mockEndpoints[endpoint.Method]; !ok {
		m.mockEndpoints[endpoint.Method] = map[string]MockEndpoint{}
	}
	m.mockEndpoints[endpoint.Method][endpoint.Endpoint].Lock()
	m.mockEndpoints[endpoint.Method][endpoint.Endpoint] = endpoint
	m.mockEndpoints[endpoint.Method][endpoint.Endpoint].Unlock()
}

func (m *MockService) LookupEndpoint(method, path string) (*MockEndpoint, error) {
	mockEndpoint, ok := m.mockEndpoints[method][path]
	if !ok {
		return nil, ErrEndpointDoesNotExist
	}
	return &mockEndpoint, nil
}

func (m *MockService) serveRegistrationHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Registering an endpoint requires a payload")
		log.Printf("Registration HTTP received empty payload")
		return
	}
	reqPayload, err := ioutil.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unable to read from payload due to: %s", err)
		log.Printf("Unable to read from payload due to: %s", err)
		return
	}
	endpointRequest := MockEndpoint{}
	if err := json.Unmarshal(reqPayload, &endpointRequest); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to Unmarshal request body %s: %s", reqPayload, err)
		log.Printf("Unable to Unmarshal request body %s: %s", reqPayload, err)
		return
	}

	m.CreateMockEndpoint(endpointRequest)
	w.WriteHeader(http.StatusCreated)
}

func (m *MockService) serveMockHTTP(w http.ResponseWriter, req *http.Request) {
	endpoint, err := m.LookupEndpoint(req.Method, req.URL.Path)
	if err == ErrEndpointDoesNotExist {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to lookup endpoint: %s", err)
		return
	}

	w.WriteHeader(endpoint.StatusCode)
	fmt.Fprint(w, endpoint.ResponseBody)
	for headerKey, headerVal := range endpoint.ResponseHeaders {
		w.Header().Add(headerKey, headerVal)
	}
}

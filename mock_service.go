package mock_service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
)

var (
	ErrEndpointDoesNotExist         = errors.New("Endpoint does not exist")
	ErrRegistrationEndpointConflict = errors.New("Endpoint conflicts with registration endpoint")
	ErrEmptyRegistrationEndpoint    = errors.New("Empty registration endpoint provided")
	ErrEmptyHTTPMethod              = errors.New("Empty HTTP method provided")
	ErrEmptyEndpoint                = errors.New("Empty endpoint provided")
)

type MockService struct {
	mockRegistrationEndpoint string
	mockedEndpoints          map[string]map[string]MockEndpoint
	sync.Mutex
}

type MockEndpoint struct {
	Method          string            `json:"method" xml:"method"`
	Endpoint        string            `json:"endpoint" xml:"endpoint"`
	StatusCode      int               `json:"httpStatusCode" xml:"httpStatusCode"`
	ResponseBody    string            `json:"responseBody" xml:"responseBody"`
	ResponseHeaders map[string]string `json:"responseHeaders" xml:"responseHeaders"`

	// TODO:
	// AcceptHeaders   map[string]string `json:"responseHeaders"`
	// RequestBody     string            `json:"requestBody"`
}

type MockServiceConf struct {
	RegistrationEndpoint string         `json:"regisgtrationEndpoint" xml:"registrationEndpoint"`
	Endpoints            []MockEndpoint `json:"endpoints" xml:"endpoints"`
}

func New(mockRegistrationEndpoint string) (*MockService, error) {
	if strings.Trim(mockRegistrationEndpoint, " ") == "" {
		return nil, ErrEmptyRegistrationEndpoint
	}
	return &MockService{
		mockRegistrationEndpoint: mockRegistrationEndpoint,
		mockedEndpoints:          map[string]map[string]MockEndpoint{},
	}, nil
}

func NewWithConf(conf *MockServiceConf) (*MockService, error) {
	if strings.Trim(conf.RegistrationEndpoint, " ") == "" {
		return nil, ErrEmptyRegistrationEndpoint
	}
	m := &MockService{
		mockRegistrationEndpoint: conf.RegistrationEndpoint,
		mockedEndpoints:          map[string]map[string]MockEndpoint{},
	}
	if err := m.loadEndpoints(conf.Endpoints); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *MockService) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// TODO: validate against requests that attempt to re-register the mock registration endpoint
	if req.Method == http.MethodPost && req.URL.Path == m.mockRegistrationEndpoint {
		m.serveRegistrationHTTP(w, req)
		return
	}

	m.serveMockHTTP(w, req)
}

func (m *MockService) loadEndpoints(endpoints []MockEndpoint) error {
	for i := range endpoints {
		if err := m.CreateMockEndpoint(endpoints[i]); err != nil {
			return err
		}
	}
	return nil
}

func (m *MockService) CreateMockEndpoint(endpoint MockEndpoint) error {
	if strings.Trim(endpoint.Method, " ") == "" {
		return ErrEmptyHTTPMethod
	}

	if strings.Trim(endpoint.Endpoint, " ") == "" {
		return ErrEmptyEndpoint
	}
	m.Lock()
	if _, ok := m.mockedEndpoints[endpoint.Method]; !ok {
		m.mockedEndpoints[endpoint.Method] = map[string]MockEndpoint{}
	}
	m.mockedEndpoints[endpoint.Method][endpoint.Endpoint] = endpoint
	m.Unlock()
	return nil
}

func (m *MockService) LookupEndpoint(method, path string) (*MockEndpoint, error) {
	mockEndpoint, ok := m.mockedEndpoints[method][path]
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
	if err := m.CreateMockEndpoint(endpointRequest); err != nil {
		switch err {
		case ErrEmptyEndpoint:
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, err)
			return
		case ErrEmptyHTTPMethod:
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, err)
		}
	}
	w.WriteHeader(http.StatusCreated)
}

func (m *MockService) serveMockHTTP(w http.ResponseWriter, req *http.Request) {
	endpoint, err := m.LookupEndpoint(req.Method, req.URL.Path)
	if err == ErrEndpointDoesNotExist {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, err)
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

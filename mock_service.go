package mock_service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var (
	ErrRegistrationEndpointConflict = errors.New("Endpoint conflicts with registration endpoint")
	ErrEmptyRegistrationEndpoint    = errors.New("Empty registration endpoint provided")
)

type MockService struct {
	mockRegistrationEndpoint string
	mockedEndpoints          *Endpoints
}

type MockServiceConf struct {
	RegistrationEndpoint string          `json:"regisgtrationEndpoint" xml:"registrationEndpoint"`
	Endpoints            []*MockEndpoint `json:"endpoints" xml:"endpoints"`
}

func New(mockRegistrationEndpoint string) (*MockService, error) {
	if strings.Trim(mockRegistrationEndpoint, " ") == "" {
		return nil, ErrEmptyRegistrationEndpoint
	}
	return &MockService{
		mockRegistrationEndpoint: mockRegistrationEndpoint,
		mockedEndpoints:          NewEndpoints(),
	}, nil
}

func NewWithConf(conf *MockServiceConf) (*MockService, error) {
	if strings.Trim(conf.RegistrationEndpoint, " ") == "" {
		return nil, ErrEmptyRegistrationEndpoint
	}
	m := &MockService{
		mockRegistrationEndpoint: conf.RegistrationEndpoint,
		mockedEndpoints:          NewEndpoints(),
	}
	if err := m.loadEndpoints(conf.Endpoints); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *MockService) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// TODO: validate against requests that attempt to re-register the mock registration endpoint
	if req.Method == http.MethodPost && req.URL.Path == m.mockRegistrationEndpoint {
		m.RegisterMockEndpoint(w, req)
		return
	}

	m.ServeMockEndpoint(w, req)
}

func (m *MockService) loadEndpoints(endpoints []*MockEndpoint) error {
	for i := range endpoints {
		if err := m.mockedEndpoints.Create(endpoints[i]); err != nil {
			return err
		}
	}
	return nil
}

func (m *MockService) RegisterMockEndpoint(w http.ResponseWriter, req *http.Request) {
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

	if err := m.mockedEndpoints.Create(&endpointRequest); err != nil {
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

func (m *MockService) ServeMockEndpoint(w http.ResponseWriter, req *http.Request) {
	endpoint, err := m.mockedEndpoints.Lookup(req.Method, req.URL.Path)
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

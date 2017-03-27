package mock_service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
)

var ErrEndpointDoesNotExist = errors.New("Endpoint does not exist")
var ErrCreatorNotExist = errors.New("Mocker Creator Endpoint does not exist.")

type MockService struct {
	mockCreatorEndpoint string
	mockedEndpoints     map[string]map[string]MockEndpoint
	sync.Mutex
}

type MockEndpoint struct {
	Method          string            `json:"method"`
	Endpoint        string            `json:"endpoint"`
	StatusCode      int               `json:"httpStatusCode"`
	ResponseBody    string            `json:"responseBody"`
	ResponseHeaders map[string]string `json:"responseHeaders"`

	// TODO:
	// AcceptHeaders   map[string]string `json:"responseHeaders"`
	// RequestBody     string            `json:"requestBody"`
}

type MockEndpointSeries []MockEndpoint

func New(mockCreatorEndpoint string) *MockService {
	return &MockService{
		mockCreatorEndpoint: mockCreatorEndpoint,
		mockedEndpoints:     map[string]map[string]MockEndpoint{},
	}
}

func (m *MockService) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost && req.URL.Path == m.mockCreatorEndpoint {
		m.serveRegistrationHTTP(w, req)
		return
	}

	m.serveMockHTTP(w, req)
}

func (m *MockService) PreloadEndpointsFromConf(confFilePath string) (error, bool) {
	if "" == strings.Trim(m.mockRegistrationEndpoint) {
		return ErrCreatorNotExist, false
	}
	fi, err := os.Open(confFilePath)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	fd, _ := ioutil.ReadAll(fi)
	var series MockEndpointSeries
	if err := json.Unmarshal(fd, &series); err != nil {
		log.Fatal("Unable to Unmarshal request body %s: %s", fd, err)
		return err, false
	}
	for _, endpoint := range series {
		m.CreateMockEndpoint(endpoint)
	}
	return nil, true
}

func (m *MockService) CreateMockEndpoint(endpoint MockEndpoint) {
	if "POST" == endpoint.Method && m.mockCreatorEndpoint == endpoint.Endpoint {
		log.Fatal("Ignore the overwriting for the mock creator endpoint.")
		return
	}
	m.Lock()
	if _, ok := m.mockedEndpoints[endpoint.Method]; !ok {
		m.mockedEndpoints[endpoint.Method] = map[string]MockEndpoint{}
	}
	m.mockedEndpoints[endpoint.Method][endpoint.Endpoint] = endpoint
	m.Unlock()
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

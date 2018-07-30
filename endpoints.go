package mockservice

import (
	"errors"
	"strings"
	"sync"
)

var (
	// ErrEndpointDoesNotExist is returned when the endpoint is not found
	ErrEndpointDoesNotExist = errors.New("Endpoint does not exist")

	// ErrEmptyHTTPMethod is returned when attempting to create add a mock endpoint with an empty HTTP method
	ErrEmptyHTTPMethod = errors.New("Empty HTTP method provided")

	// ErrEmptyEndpoint is returned when attempting to add a mock endpoint with an empty endpoint
	ErrEmptyEndpoint = errors.New("Empty endpoint provided")
)

// Endpoints includes all the registered endpoints broken down by http method and the path
type Endpoints struct {
	endpoints map[string]map[string]*MockEndpoint
	sync.Mutex
}

// MockEndpoint is used to match a request to a given response
type MockEndpoint struct {
	Method          string            `json:"method" xml:"method"`
	Endpoint        string            `json:"endpoint" xml:"endpoint"`
	StatusCode      int               `json:"httpStatusCode" xml:"httpStatusCode"`
	ResponseBody    string            `json:"responseBody" xml:"responseBody"`
	ResponseHeaders map[string]string `json:"responseHeaders" xml:"responseHeaders"`
	RequestHeaders  map[string]string `json:"responseHeaders" xml:"requestHeaders"`
	RequestBody     string            `json:"requestBody" xml:"requestBody"`
}

// NewEndpoints creates a parent struct that adding endpoints for lookup
func NewEndpoints() *Endpoints {
	return &Endpoints{endpoints: make(map[string]map[string]*MockEndpoint)}
}

// Lookup searches an endpoint by HTTP method and URL path
func (e *Endpoints) Lookup(method, path string) (*MockEndpoint, error) {
	mockEndpoint, ok := e.endpoints[method][path]
	if !ok {
		return nil, ErrEndpointDoesNotExist
	}
	return mockEndpoint, nil
}

// Create adds the endpoint to enable Lookup
func (e *Endpoints) Create(endpoint *MockEndpoint) error {
	if strings.Trim(endpoint.Method, " ") == "" {
		return ErrEmptyHTTPMethod
	}

	if strings.Trim(endpoint.Endpoint, " ") == "" {
		return ErrEmptyEndpoint
	}
	e.Lock()
	if _, ok := e.endpoints[endpoint.Method]; !ok {
		e.endpoints[endpoint.Method] = map[string]*MockEndpoint{}
	}
	e.endpoints[endpoint.Method][endpoint.Endpoint] = endpoint
	e.Unlock()

	return nil
}

// Load allows the bulk creation of a list of endpoints
func (e *Endpoints) Load(endpoints []*MockEndpoint) error {
	for i := range endpoints {
		if err := e.Create(endpoints[i]); err != nil {
			return err
		}
	}
	return nil
}

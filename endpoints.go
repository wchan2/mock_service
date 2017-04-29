package mock_service

import (
	"errors"
	"strings"
	"sync"
)

var (
	ErrEndpointDoesNotExist = errors.New("Endpoint does not exist")
	ErrEmptyHTTPMethod      = errors.New("Empty HTTP method provided")
	ErrEmptyEndpoint        = errors.New("Empty endpoint provided")
)

type Endpoints struct {
	endpoints map[string]map[string]*MockEndpoint
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

func NewEndpoints() *Endpoints {
	return &Endpoints{endpoints: make(map[string]map[string]*MockEndpoint)}
}

func (e *Endpoints) Lookup(method, path string) (*MockEndpoint, error) {
	mockEndpoint, ok := e.endpoints[method][path]
	if !ok {
		return nil, ErrEndpointDoesNotExist
	}
	return mockEndpoint, nil
}

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

func (m *Endpoints) Load(endpoints []*MockEndpoint) error {
	for i := range endpoints {
		if err := m.Create(endpoints[i]); err != nil {
			return err
		}
	}
	return nil
}

package mockservice

import (
	"errors"
	"net/http"
	"strings"
)

var (
	// ErrRegistrationEndpointConflict happens when attempting to register a mock endpoint that is the same as the registration endpoint
	ErrRegistrationEndpointConflict = errors.New("Endpoint conflicts with registration endpoint")

	// ErrEmptyRegistrationEndpoint happens when an empty registration endpoint is used to create a new mock service
	ErrEmptyRegistrationEndpoint = errors.New("Empty registration endpoint provided")
)

// MockService is a service that allows endpoints to be mocked
type MockService struct {
	mockRegistrationEndpoint string
	registrationService      *RegistrationService
	endpointService          *EndpointService
}

// Conf is a quick and easy way to configure the mock service with the registration endpoint and pre-determined mock endpoints
type Conf struct {
	RegistrationEndpoint string          `json:"regisgtrationEndpoint" xml:"registrationEndpoint"`
	Endpoints            []*MockEndpoint `json:"endpoints" xml:"endpoints"`
}

// New creates a mock service
func New(mockRegistrationEndpoint string) (*MockService, error) {
	if strings.Trim(mockRegistrationEndpoint, " ") == "" {
		return nil, ErrEmptyRegistrationEndpoint
	}

	mockEndpoints := NewEndpoints()
	return &MockService{
		mockRegistrationEndpoint: mockRegistrationEndpoint,
		registrationService:      NewRegistrationService(mockEndpoints),
		endpointService:          NewEndpointService(mockEndpoints),
	}, nil
}

// NewWithConf creates a mock service with a pre-determined configuration
func NewWithConf(conf *Conf) (*MockService, error) {
	if strings.Trim(conf.RegistrationEndpoint, " ") == "" {
		return nil, ErrEmptyRegistrationEndpoint
	}
	mockEndpoints := NewEndpoints()
	if err := mockEndpoints.Load(conf.Endpoints); err != nil {
		return nil, err
	}
	registrationService := NewRegistrationService(mockEndpoints)
	endpointService := NewEndpointService(mockEndpoints)

	m := &MockService{
		mockRegistrationEndpoint: conf.RegistrationEndpoint,
		endpointService:          endpointService,
		registrationService:      registrationService,
	}

	return m, nil
}

// ServeHTTP serves HTTP requests to the registration and mock endpoints
func (m *MockService) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost && req.URL.Path == m.mockRegistrationEndpoint {
		m.registrationService.ServeHTTP(w, req)
		return
	}

	m.endpointService.ServeHTTP(w, req)
}

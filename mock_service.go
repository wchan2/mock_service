package mock_service

import (
	"errors"
	"net/http"
	"strings"
)

var (
	ErrRegistrationEndpointConflict = errors.New("Endpoint conflicts with registration endpoint")
	ErrEmptyRegistrationEndpoint    = errors.New("Empty registration endpoint provided")
)

type MockService struct {
	mockRegistrationEndpoint string
	registrationService      *RegistrationService
	endpointService          *EndpointService
}

type MockServiceConf struct {
	RegistrationEndpoint string          `json:"regisgtrationEndpoint" xml:"registrationEndpoint"`
	Endpoints            []*MockEndpoint `json:"endpoints" xml:"endpoints"`
}

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

func NewWithConf(conf *MockServiceConf) (*MockService, error) {
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

func (m *MockService) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// TODO: validate against requests that attempt to re-register the mock registration endpoint
	if req.Method == http.MethodPost && req.URL.Path == m.mockRegistrationEndpoint {
		m.registrationService.ServeHTTP(w, req)
		return
	}

	m.endpointService.ServeHTTP(w, req)
}

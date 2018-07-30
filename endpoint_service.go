package mockservice

import (
	"fmt"
	"net/http"
)

// EndpointService matches HTTP requests to HTTP mock responses
type EndpointService struct {
	mockedEndpoints *Endpoints
}

// NewEndpointService creates an EndpointsService with endpoints to be used for matching
func NewEndpointService(endpoints *Endpoints) *EndpointService {
	return &EndpointService{mockedEndpoints: endpoints}
}

// ServeHTTP serves HTTP responses when a matched HTTP request is found
func (m *EndpointService) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	endpoint, err := m.mockedEndpoints.Lookup(req.Method, req.URL.Path)
	if err == ErrEndpointDoesNotExist {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, fmt.Sprintf("Unable to lookup endpoint: %s", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(endpoint.StatusCode)
	fmt.Fprint(w, endpoint.ResponseBody)
	for headerKey, headerVal := range endpoint.ResponseHeaders {
		w.Header().Add(headerKey, headerVal)
	}
}

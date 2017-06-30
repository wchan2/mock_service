package mock_service

import (
	"fmt"
	"net/http"
)

type EndpointService struct {
	mockedEndpoints *Endpoints
}

func NewEndpointService(endpoints *Endpoints) *EndpointService {
	return &EndpointService{mockedEndpoints: endpoints}
}

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

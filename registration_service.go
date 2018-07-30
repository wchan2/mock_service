package mockservice

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// RegistrationService allows endpoints to be registered
type RegistrationService struct {
	mockedEndpoints *Endpoints
}

// NewRegistrationService creates a registration service to support the registering of mock endpoints through HTTP
func NewRegistrationService(endpoints *Endpoints) *RegistrationService {
	return &RegistrationService{mockedEndpoints: endpoints}
}

// ServeHTTP creates mock endpoints via HTTP requests
func (m *RegistrationService) ServeHTTP(w http.ResponseWriter, req *http.Request) {
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
		case ErrEmptyEndpoint, ErrEmptyHTTPMethod:
			log.Printf("Validation error %s: %s", reqPayload, err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	w.WriteHeader(http.StatusCreated)
}

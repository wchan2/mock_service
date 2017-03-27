package mock_service

import (
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

var ErrEndpointDoesNotExist = errors.New("Endpoint does not exist")
var ErrCreatorNotExist = errors.New("Mocker Creator Endpoint does not exist.")
var ErrUnknownConfType = errors.New("Could not recongize the Configure file.")

type MockService struct {
	mockCreatorEndpoint string
	mockedEndpoints     map[string]map[string]MockEndpoint
	sync.Mutex
}

type MockEndpoint struct {
	Method          string            `json:"method" yaml:"method"`
	Endpoint        string            `json:"endpoint" yaml:"endpoint"`
	StatusCode      int               `json:"httpStatusCode" yaml:"httpStatusCode"`
	ResponseBody    string            `json:"responseBody" yaml:"responseBody"`
	ResponseHeaders map[string]string `json:"responseHeaders" yaml:"responseHeaders"`

	// TODO:
	// AcceptHeaders   map[string]string `json:"responseHeaders"`
	// RequestBody     string            `json:"requestBody"`
}

type ConfType int

const (
	UNKNOWN ConfType = iota
	JSON
	YAML
)

type MockServiceConf struct {
	Path string
	Type ConfType
}

func NewConf(confPath string, confType string) *MockServiceConf {
	c := &MockServiceConf{
		Path: confPath,
		Type: UNKNOWN,
	}
	switch strings.ToUpper(confType) {
	case "JSON":
		c.Type = JSON
	case "YAML":
		c.Type = YAML
	default:
		c.Type = UNKNOWN
	}
	return c
}

type MockEndpointSeries []MockEndpoint

func New(mockCreatorEndpoint string, conf *MockServiceConf) *MockService {
	m := &MockService{
		mockCreatorEndpoint: mockCreatorEndpoint,
		mockedEndpoints:     map[string]map[string]MockEndpoint{},
	}
	if "" != strings.TrimSpace(conf.Path) {
		m.PreloadEndpointsFromConf(conf)
	}
	return m
}

func (m *MockService) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost && req.URL.Path == m.mockCreatorEndpoint {
		m.serveRegistrationHTTP(w, req)
		return
	}

	m.serveMockHTTP(w, req)
}

func (m *MockService) PreloadEndpointsFromConf(conf *MockServiceConf) (error, bool) {
	if "" == strings.TrimSpace(m.mockCreatorEndpoint) {
		return ErrCreatorNotExist, false
	}
	fi, err := os.Open(conf.Path)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	fd, _ := ioutil.ReadAll(fi)
	var series MockEndpointSeries
	switch conf.Type {
	case JSON, UNKNOWN:
		if err := json.Unmarshal(fd, &series); nil != err {
			log.Fatal("Unable to Unmarshal JSON string %s: %s", string(fd), err)
			return err, false
		}
	case YAML:
		if err := yaml.Unmarshal(fd, &series); nil != err {
			log.Fatal("Unable to Unmarshal YAML string %s: %s", string(fd), err)
			return err, false
		}
	default:
		log.Fatal("Unable to handle the undefined conf file type: %s", conf.Type)
		return ErrUnknownConfType, false
	}
	fmt.Println("%v", series)
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

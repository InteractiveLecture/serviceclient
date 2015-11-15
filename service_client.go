package serviceclient

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	s "strings"
)

type ServiceClient struct {
	HttpClient     *http.Client
	Service        string
	backendAdapter BackendAdapter
}

type ServiceCaller interface {
	Get(path string, headers ...string) (*http.Response, error)
	GetSecure(path string, headers ...string) (*http.Response, error)
	Delete(path string, headers ...string) (*http.Response, error)
	DeleteSecure(path string, headers ...string) (*http.Response, error)
	Post(path string, bodyType string, body io.Reader, headers ...string) (*http.Response, error)
	PostSecure(path string, bodyType string, body io.Reader, headers ...string) (*http.Response, error)
	Put(path string, bodyType string, body io.Reader, headers ...string) (*http.Response, error)
	PutSecure(path string, bodyType string, body io.Reader, headers ...string) (*http.Response, error)
	Patch(path string, bodyType string, body io.Reader, headers ...string) (*http.Response, error)
	PatchSecure(path string, bodyType string, body io.Reader, headers ...string) (*http.Response, error)
}

type BackendAdapter interface {
	Resolve(name string) (string, error)
	Configure(services ...string) error
	Refresh() error
}

var instances = make(map[string]*ServiceClient)

func Configure(backendAdapter BackendAdapter, services ...string) error {
	for _, service := range services {
		instances[service] = New(service, backendAdapter)
	}
	return backendAdapter.Configure(services...)
}

func New(service string, backendAdapter BackendAdapter) *ServiceClient {
	return &ServiceClient{http.DefaultClient, service, backendAdapter}
}

func GetInstance(service string) *ServiceClient {
	return instances[service]
}

func (client *ServiceClient) Get(path string, headers ...string) (*http.Response, error) {
	req, err := prepareRequest("GET", client.resolvePath(path, "http"), nil, headers)
	if err != nil {
		return nil, err
	}
	return client.sendRequest(req)
}

func (client *ServiceClient) GetSecure(path string, headers ...string) (*http.Response, error) {
	req, err := prepareRequest("GET", client.resolvePath(path, "https"), nil, headers)
	if err != nil {
		return nil, err
	}
	return client.sendRequest(req)
}

func (client *ServiceClient) Post(path string, bodyType string, body io.Reader, headers ...string) (*http.Response, error) {
	headers = append(headers, fmt.Sprintf("Content-Type:%s", bodyType))
	req, err := prepareRequest("POST", client.resolvePath(path, "http"), body, headers)
	if err != nil {
		return nil, err
	}
	return client.sendRequest(req)
}

func (client *ServiceClient) PostSecure(path string, bodyType string, body io.Reader, headers ...string) (*http.Response, error) {
	headers = append(headers, fmt.Sprintf("Content-Type:%s", bodyType))
	req, err := prepareRequest("POST", client.resolvePath(path, "https"), body, headers)
	if err != nil {
		return nil, err
	}
	return client.sendRequest(req)
}

func (client *ServiceClient) Put(path string, bodyType string, body io.Reader, headers ...string) (*http.Response, error) {
	headers = append(headers, fmt.Sprintf("Content-Type:%s", bodyType))
	req, err := prepareRequest("PUT", client.resolvePath(path, "http"), body, headers)
	if err != nil {
		return nil, err
	}
	return client.sendRequest(req)
}

func (client *ServiceClient) PutSecure(path string, bodyType string, body io.Reader, headers ...string) (*http.Response, error) {
	headers = append(headers, fmt.Sprintf("Content-Type:%s", bodyType))
	req, err := prepareRequest("PUT", client.resolvePath(path, "https"), body, headers)
	if err != nil {
		return nil, err
	}
	return client.sendRequest(req)
}

func (client *ServiceClient) Patch(path string, bodyType string, body io.Reader, headers ...string) (*http.Response, error) {
	headers = append(headers, fmt.Sprintf("Content-Type:%s", bodyType))
	req, err := prepareRequest("PATCH", client.resolvePath(path, "http"), body, headers)
	if err != nil {
		return nil, err
	}
	return client.sendRequest(req)
}

func (client *ServiceClient) PatchSecure(path string, bodyType string, body io.Reader, headers ...string) (*http.Response, error) {
	headers = append(headers, fmt.Sprintf("Content-Type:%s", bodyType))
	req, err := prepareRequest("PATCH", client.resolvePath(path, "https"), body, headers)
	if err != nil {
		return nil, err
	}
	return client.sendRequest(req)
}

func (client *ServiceClient) Delete(path string, headers ...string) (*http.Response, error) {
	req, err := prepareRequest("DELETE", client.resolvePath(path, "http"), nil, headers)
	if err != nil {
		return nil, err
	}
	return client.sendRequest(req)
}

func (client *ServiceClient) DeleteSecure(path string, headers ...string) (*http.Response, error) {
	req, err := prepareRequest("DELETE", client.resolvePath(path, "https"), nil, headers)
	if err != nil {
		return nil, err
	}
	return client.sendRequest(req)
}

func (client *ServiceClient) sendRequest(req *http.Request) (resp *http.Response, err error) {
	resp, err = client.HttpClient.Do(req)
	if err != nil {
		client.backendAdapter.Refresh()
		resp, err = client.HttpClient.Do(req)
	}
	return
}

func addHeaders(req *http.Request, headers []string) error {
	for _, header := range headers {
		parts := s.Split(header, ":")
		if len(parts) < 2 {
			return errors.New("not a valid header")
		}
		req.Header.Set(parts[0], parts[1])
	}
	return nil
}

func prepareRequest(requestType string, url string, body io.Reader, headers []string) (req *http.Request, err error) {
	req, err = http.NewRequest(requestType, url, body)
	if err != nil {
		return
	}
	err = addHeaders(req, headers)
	if err != nil {
		req = nil
		return
	}
	return
}

func (client *ServiceClient) resolvePath(path string, schema string) string {
	address, _ := client.backendAdapter.Resolve(client.Service)
	if s.HasPrefix(address, "http") || s.HasPrefix(address, "https") {
		return fmt.Sprintf("%s/%s", address, path)
	}
	return fmt.Sprintf("%s://%s/%s", schema, address, path)
}

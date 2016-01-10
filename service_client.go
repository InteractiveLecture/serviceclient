package serviceclient

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	s "strings"

	"github.com/miekg/dns"
)

type ServiceClient struct {
	HttpClient *http.Client
	Service    string
	Resolver   AddressResolver
}

type AddressResolverFactory func() AddressResolver

var ResolverFactory = func() AddressResolver {
	return ConsulDnsAddressResolver{"discovery:53"}
}

func New(service string) *ServiceClient {
	return &ServiceClient{http.DefaultClient, service, ResolverFactory()}
}

func (client *ServiceClient) Get(path string, headers ...string) (*http.Response, error) {
	address, err := client.resolvePath(path, "http")
	if err != nil {
		return nil, err
	}
	req, err := prepareRequest("GET", address, nil, headers)
	if err != nil {
		return nil, err
	}
	return client.sendRequest(req)
}

func (client *ServiceClient) GetSecure(path string, headers ...string) (*http.Response, error) {
	address, err := client.resolvePath(path, "https")
	if err != nil {
		return nil, err
	}
	req, err := prepareRequest("GET", address, nil, headers)
	if err != nil {
		return nil, err
	}
	return client.sendRequest(req)
}

func (client *ServiceClient) Post(path string, bodyType string, body io.Reader, headers ...string) (*http.Response, error) {
	address, err := client.resolvePath(path, "http")
	if err != nil {
		return nil, err
	}
	headers = append(headers, fmt.Sprintf("Content-Type:%s", bodyType))
	req, err := prepareRequest("POST", address, body, headers)
	if err != nil {
		return nil, err
	}
	return client.sendRequest(req)
}

func (client *ServiceClient) PostSecure(path string, bodyType string, body io.Reader, headers ...string) (*http.Response, error) {
	address, err := client.resolvePath(path, "https")
	if err != nil {
		return nil, err
	}
	headers = append(headers, fmt.Sprintf("Content-Type:%s", bodyType))
	req, err := prepareRequest("POST", address, body, headers)
	if err != nil {
		return nil, err
	}
	return client.sendRequest(req)
}

func (client *ServiceClient) Put(path string, bodyType string, body io.Reader, headers ...string) (*http.Response, error) {
	address, err := client.resolvePath(path, "http")
	if err != nil {
		return nil, err
	}
	headers = append(headers, fmt.Sprintf("Content-Type:%s", bodyType))
	req, err := prepareRequest("PUT", address, body, headers)
	if err != nil {
		return nil, err
	}
	return client.sendRequest(req)
}

func (client *ServiceClient) PutSecure(path string, bodyType string, body io.Reader, headers ...string) (*http.Response, error) {
	address, err := client.resolvePath(path, "https")
	if err != nil {
		return nil, err
	}
	headers = append(headers, fmt.Sprintf("Content-Type:%s", bodyType))
	req, err := prepareRequest("PUT", address, body, headers)
	if err != nil {
		return nil, err
	}
	return client.sendRequest(req)
}

func (client *ServiceClient) Patch(path string, bodyType string, body io.Reader, headers ...string) (*http.Response, error) {
	address, err := client.resolvePath(path, "http")
	if err != nil {
		return nil, err
	}
	headers = append(headers, fmt.Sprintf("Content-Type:%s", bodyType))
	req, err := prepareRequest("PATCH", address, body, headers)
	if err != nil {
		return nil, err
	}
	return client.sendRequest(req)
}

func (client *ServiceClient) PatchSecure(path string, bodyType string, body io.Reader, headers ...string) (*http.Response, error) {
	address, err := client.resolvePath(path, "https")
	if err != nil {
		return nil, err
	}
	headers = append(headers, fmt.Sprintf("Content-Type:%s", bodyType))
	req, err := prepareRequest("PATCH", address, body, headers)
	if err != nil {
		return nil, err
	}
	return client.sendRequest(req)
}

func (client *ServiceClient) Delete(path string, headers ...string) (*http.Response, error) {
	address, err := client.resolvePath(path, "http")
	if err != nil {
		return nil, err
	}
	req, err := prepareRequest("DELETE", address, nil, headers)
	if err != nil {
		return nil, err
	}
	return client.sendRequest(req)
}

func (client *ServiceClient) DeleteSecure(path string, headers ...string) (*http.Response, error) {
	address, err := client.resolvePath(path, "https")
	if err != nil {
		return nil, err
	}
	req, err := prepareRequest("DELETE", address, nil, headers)
	if err != nil {
		return nil, err
	}
	return client.sendRequest(req)
}

func (client *ServiceClient) sendRequest(req *http.Request) (*http.Response, error) {
	return client.HttpClient.Do(req)
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
	log.Printf("sending %s request to: %s", requestType, url)
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

type AddressResolver interface {
	Resolve(serviceName string) (string, error)
}

type ConsulDnsAddressResolver struct {
	ServerAddress string
}

func (client *ServiceClient) resolvePath(path string, schema string) (string, error) {
	address, err := client.Resolver.Resolve(client.Service)
	if err != nil {
		return "", err
	}
	address = address + ":8080"
	if s.HasPrefix(address, "http") || s.HasPrefix(address, "https") {
		return fmt.Sprintf("%s/%s", address, path), nil
	}
	return fmt.Sprintf("%s://%s%s", schema, address, path), nil
}

func (resolver ConsulDnsAddressResolver) Resolve(service string) (string, error) {
	m1 := new(dns.Msg)
	m1.Id = dns.Id()
	m1.RecursionDesired = true
	m1.SetQuestion(service+".service.consul.", dns.TypeA)
	c := new(dns.Client)
	in, _, err := c.Exchange(m1, resolver.ServerAddress)
	if err != nil {
		log.Fatal(err)
	}
	if len(in.Answer) > 0 {
		return in.Answer[0].(*dns.A).A.String(), nil
	}
	return "", errors.New("Could not resolve service address")
}

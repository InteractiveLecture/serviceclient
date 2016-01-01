package servicetest

import (
	"net/http"
	"net/http/httptest"

	"github.com/InteractiveLecture/serviceclient"
)

type MockDnsResolver struct {
	ServiceAddress string
}

func (r MockDnsResolver) Resolve(service string) (string, error) {
	return r.ServiceAddress, nil
}

func Service(serviceName string, handler http.Handler) (*httptest.Server, *serviceclient.ServiceClient) {
	server := httptest.NewServer(handler)
	serviceclient.ResolverFactory = func() serviceclient.AddressResolver {
		return MockDnsResolver{server.URL}
	}
	client := serviceclient.New(serviceName)
	return server, client
}

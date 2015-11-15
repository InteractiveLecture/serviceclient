package servicetest

import (
	"net/http"
	"net/http/httptest"

	"github.com/InteractiveLecture/serviceclient"
	"github.com/InteractiveLecture/serviceclient/mocks"
	"github.com/golang/mock/gomock"
)

func Service(controller *gomock.Controller, serviceName string, handler http.Handler) *httptest.Server {
	mockBackend := mocks.NewMockBackendAdapter(controller)
	server := httptest.NewServer(handler)
	mockBackend.EXPECT().Resolve(serviceName).Return(server.URL, nil)
	mockBackend.EXPECT().Configure(serviceName).Return(nil)
	mockBackend.EXPECT().Refresh()
	serviceclient.Configure(mockBackend, serviceName)
	return server
}
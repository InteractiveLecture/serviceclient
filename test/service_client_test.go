package servicetest

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/InteractiveLecture/serviceclient"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func prepareMock(t *testing.T) (*servicemocks.MockBackendAdapter, *gomock.Controller) {
	mockCtrl := gomock.NewController(t)
	mock := servicemocks.NewMockBackendAdapter(mockCtrl)
	mock.EXPECT().Configure("acl-service", "authentication-service").Return(nil)
	return mock, mockCtrl
}

func prepareService(t *testing.T) (*servicemocks.MockBackendAdapter, *gomock.Controller) {
	mock, mockCtrl := prepareMock(t)
	serviceclient.Configure(mock, "acl-service", "authentication-service")
	return mock, mockCtrl
}

func prepareDefaultServer(controller *gomock.Controller, serviceName string) *httptest.Server {
	return Service(controller, serviceName, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"fake twitter json string"}`)
	}))
}

func TestConfigure(t *testing.T) {
	mock, mockCtrl := prepareMock(t)
	defer mockCtrl.Finish()
	assert.Nil(t, serviceclient.Configure(mock, "acl-service", "authentication-service"))
}

func TestGetInstance(t *testing.T) {
	_, mockCtrl := prepareService(t)
	defer mockCtrl.Finish()
	assert.NotNil(t, serviceclient.GetInstance("acl-service"))
}

type Context func(instance *serviceclient.ServiceClient) (*http.Response, error)

func TestGet(t *testing.T) {
	inContext(t, get)
}

func TestGetSecure(t *testing.T) {
	inContext(t, getSecure)
}

func TestDelete(t *testing.T) {
	inContext(t, deleteNormal)
}
func TestDeleteSecure(t *testing.T) {
	inContext(t, deleteSecure)
}
func TestePostSecure(t *testing.T) {
	inContext(t, postSecure)
}
func TestPost(t *testing.T) {
	inContext(t, post)
}
func TestPut(t *testing.T) {
	inContext(t, put)
}
func TestPutSecure(t *testing.T) {
	inContext(t, putSecure)
}
func TestPatchSecure(t *testing.T) {
	inContext(t, patchSecure)
}
func TestPatch(t *testing.T) {
	inContext(t, patch)
}
func inContext(t *testing.T, context Context) {
	controller := gomock.NewController(t)
	server := prepareDefaultServer(controller, "acl-service")
	defer server.Close()
	defer controller.Finish()
	instance := serviceclient.GetInstance("acl-service")
	resp, err := context(instance)
	checkResponse(t, resp, err)
}

func get(instance *serviceclient.ServiceClient) (*http.Response, error) {
	return instance.Get("/bla")
}

func getSecure(instance *serviceclient.ServiceClient) (*http.Response, error) {
	return instance.GetSecure("/bla")
}

func post(instance *serviceclient.ServiceClient) (*http.Response, error) {
	jsonString := `{"text":"fake twitter json string"}`
	reader := strings.NewReader(jsonString)
	return instance.Post("/bla/blubb", "application/json", reader)
}
func postSecure(instance *serviceclient.ServiceClient) (*http.Response, error) {
	jsonString := `{"text":"fake twitter json string"}`
	reader := strings.NewReader(jsonString)
	return instance.PostSecure("/bla/blubb", "application/json", reader)
}

func putSecure(instance *serviceclient.ServiceClient) (*http.Response, error) {
	jsonString := `{"text":"fake twitter json string"}`
	reader := strings.NewReader(jsonString)
	return instance.PutSecure("/bla/blubb", "application/json", reader)
}

func put(instance *serviceclient.ServiceClient) (*http.Response, error) {
	jsonString := `{"text":"fake twitter json string"}`
	reader := strings.NewReader(jsonString)
	return instance.Put("/bla/blubb", "application/json", reader)
}

func patch(instance *serviceclient.ServiceClient) (*http.Response, error) {
	jsonString := `{"text":"fake twitter json string"}`
	reader := strings.NewReader(jsonString)
	return instance.Patch("/bla/blubb", "application/json", reader)
}

func patchSecure(instance *serviceclient.ServiceClient) (*http.Response, error) {
	jsonString := `{"text":"fake twitter json string"}`
	reader := strings.NewReader(jsonString)
	return instance.PatchSecure("/bla/blubb", "application/json", reader)
}

func deleteNormal(instance *serviceclient.ServiceClient) (*http.Response, error) {
	return instance.Delete("/bla")
}

func deleteSecure(instance *serviceclient.ServiceClient) (*http.Response, error) {
	return instance.DeleteSecure("/bla")
}

func checkResponse(t *testing.T, resp *http.Response, err error) {
	assert.Nil(t, err)
	assert.NotNil(t, resp)
}

package serviceclient

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func prepareMock(t *testing.T) (*MockBackendAdapter, *gomock.Controller) {
	mockCtrl := gomock.NewController(t)
	mock := NewMockBackendAdapter(mockCtrl)
	defer mockCtrl.Finish()
	mock.EXPECT().Configure("acl-service", "authentication-service").Return(nil)
	return mock, mockCtrl
}

func prepareService(t *testing.T) (*MockBackendAdapter, *gomock.Controller) {
	mock, mockCtrl := prepareMock(t)
	Configure(mock, "acl-service", "authentication-service")
	return mock, mockCtrl
}

func prepareDefaultServer() *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"fake twitter json string"}`)
	}))
	return ts
}

func ConfigureTest(t *testing.T) {
	mock, mockCtrl := prepareMock(t)
	defer mockCtrl.Finish()
	assert.Nil(t, Configure(mock, "acl-service", "authentication-service"))
}

func GetInstanceTest(t *testing.T) {
	_, mockCtrl := prepareService(t)
	defer mockCtrl.Finish()
	assert.NotNil(t, GetInstance("acl-service"))
}

type Context func(t *testing.T, mock *MockBackendAdapter, server *httptest.Server)

func WithServerTest(t *testing.T) {
	ts := prepareDefaultServer()
	defer ts.Close()
	mock, mockCtrl := prepareService(t)
	defer mockCtrl.Finish()
	mock.EXPECT().Resolve("acl-service").Return(ts.URL)
	instance := GetInstance("acl-service")
	testGet(t, mock, instance)
	testGetSecure(t, mock, instance)
	testDelete(t, mock, instance)
	testDeleteSecure(t, mock, instance)
	jsonString := `{"text":"fake twitter json string"}`
	reader := strings.NewReader(jsonString)
	testPost(t, mock, instance, reader)
	testPostSecure(t, mock, instance, reader)
	testPutSecure(t, mock, instance, reader)
	testPut(t, mock, instance, reader)
	testPatch(t, mock, instance, reader)
	testPatchSecure(t, mock, instance, reader)
}

func checkResponse(t *testing.T, resp *http.Response, err error) {
	assert.Nil(t, err)
	assert.NotNil(t, resp)
}

func testGet(t *testing.T, mock *MockBackendAdapter, instance *ServiceClient) {
	resp, err := instance.Get("/bla/blubb")
	checkResponse(t, resp, err)
}

func testGetSecure(t *testing.T, mock *MockBackendAdapter, instance *ServiceClient) {
	resp, err := instance.GetSecure("/bla/blubb")
	checkResponse(t, resp, err)

}

func testDelete(t *testing.T, mock *MockBackendAdapter, instance *ServiceClient) {
	resp, err := instance.DeleteSecure("/bla/blubb")
	checkResponse(t, resp, err)

}
func testDeleteSecure(t *testing.T, mock *MockBackendAdapter, instance *ServiceClient) {
	resp, err := instance.DeleteSecure("/bla/blubb")
	checkResponse(t, resp, err)
}

func testPost(t *testing.T, mock *MockBackendAdapter, instance *ServiceClient, reader io.Reader) {
	resp, err := instance.Post("/bla/blubb", "application/json", reader)
	checkResponse(t, resp, err)
}

func testPostSecure(t *testing.T, mock *MockBackendAdapter, instance *ServiceClient, reader io.Reader) {
	resp, err := instance.PostSecure("/bla/blubb", "application/json", reader)
	checkResponse(t, resp, err)
}

func testPatch(t *testing.T, mock *MockBackendAdapter, instance *ServiceClient, reader io.Reader) {
	resp, err := instance.Patch("/bla/blubb", "application/json", reader)
	checkResponse(t, resp, err)
}
func testPatchSecure(t *testing.T, mock *MockBackendAdapter, instance *ServiceClient, reader io.Reader) {
	resp, err := instance.PatchSecure("/bla/blubb", "application/json", reader)
	checkResponse(t, resp, err)
}
func testPut(t *testing.T, mock *MockBackendAdapter, instance *ServiceClient, reader io.Reader) {
	resp, err := instance.Put("/bla/blubb", "application/json", reader)
	checkResponse(t, resp, err)
}
func testPutSecure(t *testing.T, mock *MockBackendAdapter, instance *ServiceClient, reader io.Reader) {
	resp, err := instance.PutSecure("/bla/blubb", "application/json", reader)
	checkResponse(t, resp, err)
}

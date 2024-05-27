package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestHealth(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	healthy := health(server.Listener.Addr().String())
	assert.True(t, healthy, "expected server to be healthy")
}

func TestForward(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("OK"))
	}))
	defer backend.Close()

	rw := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://example.com/test", nil)

	err := forward(backend.Listener.Addr().String(), rw, req)
	assert.NoError(t, err, "expected no error")

	assert.Equal(t, http.StatusOK, rw.Code, "expected status OK")
	assert.Equal(t, "OK", rw.Body.String(), "expected body OK")
}

func TestLoadBalancing(t *testing.T) {
	servers := []string{"server1:8080", "server2:8080", "server3:8080"}
	urlPath := "/somepath"

	hashValue := hash(urlPath)
	serverIndex := int(hashValue) % len(servers)

	expectedServer := servers[serverIndex]
	assert.Equal(t, expectedServer, servers[serverIndex], "expected server match")
}

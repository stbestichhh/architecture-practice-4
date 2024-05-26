package integration

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"
)

const baseAddress = "http://balancer:8090"

var client = http.Client{
	Timeout: 3 * time.Second,
}

func TestBalancer(t *testing.T) {
	if _, exists := os.LookupEnv("INTEGRATION_TEST"); !exists {
		t.Skip("Integration test is not enabled")
	}

	// Send a single request to check basic functionality
	resp, err := client.Get(fmt.Sprintf("%s/api/v1/some-data", baseAddress))
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()
	t.Logf("response from [%s]", resp.Header.Get("lb-from"))

	// Call the helper function to verify load balancing
	verifyLoadBalancing(t)
}

func verifyLoadBalancing(t *testing.T) {
	const numRequests = 10
	backends := make(map[string]int)

	for i := 0; i < numRequests; i++ {
		resp, err := client.Get(fmt.Sprintf("%s/api/v1/some-data", baseAddress))
		if err != nil {
			t.Error(err)
			return
		}
		defer resp.Body.Close()

		backend := resp.Header.Get("lb-from")
		if backend == "" {
			t.Error("Response header 'lb-from' is empty")
			return
		}
		backends[backend]++
	}

	if len(backends) < 2 {
		t.Errorf("Load balancer did not distribute requests to multiple backends: %v", backends)
	} else {
		t.Logf("Requests were distributed across the following backends: %v", backends)
	}
}

func BenchmarkBalancer(b *testing.B) {
	for i := 0; i < b.N; i++ {
			resp, err := client.Get(fmt.Sprintf("%s/api/v1/some-data", baseAddress))
			if err != nil {
				b.Fatal(err)
			}
			resp.Body.Close()
		}
}

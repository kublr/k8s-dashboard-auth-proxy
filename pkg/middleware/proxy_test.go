package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kublr/k8s-dashboard-auth-proxy/pkg/config"
)

func TestProxyMiddleware_Handler(t *testing.T) {
	proxy, err := NewProxyMiddleware(config.Config{
		DashboardAuthorizationHeadersPrefix: "Test-",
	})
	require.NoError(t, err)

	req := httptest.NewRequest("GET", "https://localhost:443", nil)
	req.Header.Set("Test-Authorization", "Bearer 1234567890")
	req.Header.Set("Test-Impersonate-User", "user")
	req.Header.Add("Test-Impersonate-Group", "group1")
	req.Header.Add("Test-Impersonate-Group", "group2")

	proxy.Handler().ServeHTTP(httptest.NewRecorder(), req)

	assert.EqualValues(t, map[string][]string{
		"Authorization":     {"Bearer 1234567890"},
		"Impersonate-User":  {"user"},
		"Impersonate-Group": {"group1", "group2"},
	}, req.Header)
}

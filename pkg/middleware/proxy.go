package middleware

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/textproto"
	"net/url"
	"strings"
	"time"

	"github.com/kublr/k8s-dashboard-auth-proxy/pkg/config"
)

// ProxyMiddleware is a HTTP middleware implementation that handles Proxy authentication
type ProxyMiddleware struct {
	headersPrefix string

	proxy *httputil.ReverseProxy
}

// NewProxyMiddleware creates new ProxyMiddleware instance
func NewProxyMiddleware(config config.Config) (*ProxyMiddleware, error) {
	upstreamURL, err := url.Parse(config.DashboardEndpoint)
	if err != nil {
		return nil, fmt.Errorf("cannot parse upstream URL: %s", config.DashboardEndpoint)
	}

	proxy := httputil.NewSingleHostReverseProxy(upstreamURL)
	proxy.Transport = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	return &ProxyMiddleware{
		headersPrefix: textproto.CanonicalMIMEHeaderKey(config.DashboardAuthorizationHeadersPrefix),
		proxy:         proxy,
	}, nil
}

// Handler returns http.Handler implementation
func (p *ProxyMiddleware) Handler() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		for headerName, headerValues := range req.Header {
			if strings.HasPrefix(headerName, p.headersPrefix) {
				newHeaderName := strings.TrimPrefix(headerName, p.headersPrefix)
				req.Header.Del(headerName)
				req.Header.Del(newHeaderName)

				for _, headerValue := range headerValues {
					req.Header.Add(newHeaderName, headerValue)
				}
			}
		}

		p.proxy.ServeHTTP(rw, req)
	})
}

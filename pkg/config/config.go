package config

import (
	"net"
	"time"
)

// Config describes application configuration
type Config struct {
	SecurePort          int
	SecureBindAddress   net.IP
	InsecurePort        int
	InsecureBindAddress net.IP

	AutoGenerateCertificates bool
	DefaultCertDir           string
	CertFile                 string
	KeyFile                  string

	DashboardEndpoint                   string
	DashboardAuthorizationHeadersPrefix string
}

// AuthConfig describes auth configuration
type AuthConfig struct {
	AuthorizationHeadersPrefix string `mapstructure:"authorization-headers-prefix"`
}

// UpstreamConfig describes proxy configuration
type UpstreamConfig struct {
	URL              string        `mapstructure:"url"`
	Timeout          time.Duration `mapstructure:"timeout"`
	KeepAliveTimeout time.Duration `mapstructure:"keep-alive-timeout"`
	SkipTLSVerify    bool          `mapstructure:"skip-tls-verify"`
	MaxIdleConns     int           `mapstructure:"max-idle-conns"`
}

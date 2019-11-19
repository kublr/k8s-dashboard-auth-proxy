package pkg

import (
	"crypto/elliptic"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/kubernetes/dashboard/src/app/backend/cert"
	"github.com/kubernetes/dashboard/src/app/backend/cert/ecdsa"

	"github.com/kublr/k8s-dashboard-auth-proxy/pkg/config"
	"github.com/kublr/k8s-dashboard-auth-proxy/pkg/middleware"
)

// Server is auth-proxy server
type Server struct {
	config config.Config
}

// NewServer returns new Server instance
func NewServer(config config.Config) *Server {
	return &Server{
		config: config,
	}
}

// Run registers API and starts server of incoming requests.
func (s *Server) Run() {
	proxyMiddleware, err := middleware.NewProxyMiddleware(s.config)
	if err != nil {
		log.Fatalf("Error creating proxy handler: %s", err)
	}

	http.Handle("/", handlers.LoggingHandler(os.Stdout, proxyMiddleware.Handler()))

	var servingCerts []tls.Certificate
	if s.config.AutoGenerateCertificates {
		log.Println("Auto-generating certificates")
		certCreator := ecdsa.NewECDSACreator(s.config.KeyFile, s.config.CertFile, elliptic.P256())
		certManager := cert.NewCertManager(certCreator, s.config.DefaultCertDir)
		servingCert, err := certManager.GetCertificates()
		if err != nil {
			log.Fatalf("Error while generating server certificates: %s", err)
		}
		servingCerts = []tls.Certificate{servingCert}
	} else if s.config.CertFile != "" && s.config.KeyFile != "" {
		certFilePath := s.config.DefaultCertDir + string(os.PathSeparator) + s.config.CertFile
		keyFilePath := s.config.DefaultCertDir + string(os.PathSeparator) + s.config.KeyFile
		servingCert, err := tls.LoadX509KeyPair(certFilePath, keyFilePath)
		if err != nil {
			log.Fatalf("Error while loading server certificates: %s", err)
		}
		servingCerts = []tls.Certificate{servingCert}
	}

	// Listen for http or https
	if servingCerts != nil {
		log.Printf("Serving securely on HTTPS port: %d", s.config.SecurePort)
		secureAddr := fmt.Sprintf("%s:%d", s.config.SecureBindAddress, s.config.SecurePort)
		server := &http.Server{
			Addr:      secureAddr,
			Handler:   http.DefaultServeMux,
			TLSConfig: &tls.Config{Certificates: servingCerts},
		}
		go func() {
			log.Fatal(server.ListenAndServeTLS("", ""))
		}()
	} else {
		log.Printf("Serving insecurely on HTTP port: %d", s.config.InsecurePort)
		addr := fmt.Sprintf("%s:%d", s.config.InsecureBindAddress, s.config.InsecurePort)
		go func() {
			log.Fatal(http.ListenAndServe(addr, nil))
		}()
	}
	select {}
}

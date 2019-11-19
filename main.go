package main

import (
	"flag"
	"fmt"
	"github.com/kublr/k8s-dashboard-auth-proxy/pkg"
	"log"
	"math/rand"
	"net"
	"os"
	"runtime"
	"time"

	"github.com/spf13/pflag"

	"github.com/kublr/k8s-dashboard-auth-proxy/pkg/config"
)

var (
	argSecurePort          = pflag.Int("port", 9443, "The secure port to listen to for incoming HTTPS requests.")
	argSecureBindAddress   = pflag.IP("bind-address", net.IPv4(0, 0, 0, 0), "The IP address on which to serve the --port.")
	argInsecurePort        = pflag.Int("insecure-port", 9999, "The port to listen to for incoming HTTP requests.")
	argInsecureBindAddress = pflag.IP("insecure-bind-address", net.IPv4(127, 0, 0, 1), "The IP address on which to serve the --insecure-port.")

	argAutoGenerateCertificates = pflag.Bool("auto-generate-certificates", false, "When set to true, proxy will automatically generate certificates used to serve HTTPS.")
	argDefaultCertDir           = pflag.String("default-cert-dir", "/certs", "Directory path containing '--tls-cert-file' and '--tls-key-file' files. Used also when auto-generating certificates flag is set.")
	argCertFile                 = pflag.String("tls-cert-file", "", "File containing the default x509 Certificate for HTTPS.")
	argKeyFile                  = pflag.String("tls-key-file", "", "File containing the default x509 private key matching --tls-cert-file.")

	argDashboardEndpoint          = pflag.String("dashboard-endpoint", "https://localhost:8443", "The address of the Kubernetes Dashboard to connect to in the format of protocol://address:port, e.g., http://localhost:8080.")
	argDashboardAuthHeadersPrefix = pflag.String("dashboard-auth-headers-prefix", "Dashboard-", "This prefix will be removed from HTTP headers before passing them to Kubernetes Dashboard")
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	rand.Seed(time.Now().Unix())

	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	_ = flag.CommandLine.Parse(make([]string, 0)) // Init for glog calls in kubernetes packages

	c := getConfig()
	log.Println("Using config:", fmt.Sprintf("%#v", c))

	pkg.NewServer(c).Run()
}

func getConfig() config.Config {
	return config.Config{
		SecurePort:          *argSecurePort,
		SecureBindAddress:   *argSecureBindAddress,
		InsecurePort:        *argInsecurePort,
		InsecureBindAddress: *argInsecureBindAddress,

		AutoGenerateCertificates: *argAutoGenerateCertificates,
		DefaultCertDir:           *argDefaultCertDir,
		CertFile:                 *argCertFile,
		KeyFile:                  *argKeyFile,

		DashboardEndpoint:                   *argDashboardEndpoint,
		DashboardAuthorizationHeadersPrefix: *argDashboardAuthHeadersPrefix,
	}
}

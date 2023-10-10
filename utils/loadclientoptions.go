package utils

import (
  "crypto/tls"
  "crypto/x509"
  "fmt"
  "os"
  "strconv"

  "go.temporal.io/sdk/client"
)

/* LoadClientOptions - Return client options for Temporal Cloud */
func LoadClientOptions() (client.Options, error) {

  // Read env variables
  targetHost := os.Getenv("TEMPORAL_HOST_URL")
  namespace := os.Getenv("TEMPORAL_NAMESPACE")
  clientCert := os.Getenv("TEMPORAL_TLS_CERT")
  clientKey := os.Getenv("TEMPORAL_TLS_KEY")

  // Optional:
  serverRootCACert := os.Getenv("TEMPORAL_SERVER_ROOT_CA_CERT")
  serverName := os.Getenv("TEMPORAL_SERVER_NAME")

  insecureSkipVerify, _ := strconv.ParseBool(os.Getenv("TEMPORAL_INSECURE_SKIP_VERIFY"))

  // Load client cert
  cert, err := tls.LoadX509KeyPair(clientCert, clientKey)
  if err != nil {
    return client.Options{}, fmt.Errorf("failed loading client cert and key: %w", err)
  }

  // Load server CA if given
  var serverCAPool *x509.CertPool
  if serverRootCACert != "" {
    serverCAPool = x509.NewCertPool()
    b, err := os.ReadFile(serverRootCACert)
    if err != nil {
      return client.Options{}, fmt.Errorf("failed reading server CA: %w", err)
    } else if !serverCAPool.AppendCertsFromPEM(b) {
      return client.Options{}, fmt.Errorf("server CA PEM file invalid")
    }
  }

  // Return client options
  return client.Options{
    HostPort:  targetHost,
    Namespace: namespace,
    ConnectionOptions: client.ConnectionOptions{
      TLS: &tls.Config{
        Certificates:       []tls.Certificate{cert},
        RootCAs:            serverCAPool,
        ServerName:         serverName,
        InsecureSkipVerify: insecureSkipVerify,
      },
    },
    Logger: NewTClientLogger(),
  }, nil

}

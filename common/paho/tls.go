package paho

// Configure TLS authentication for paho/mqtt

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	"github.com/sirupsen/logrus"
)

// get a new config for ssl
func configureTLSConfig(caCertFile, clientCertFile, clientKeyFile string) (*tls.Config, error) {
	// import CA from file
	certpool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(caCertFile)
	if err != nil {
		logrus.Errorf("problem reading CA file at %s: %s", caCertFile, err)
		return nil, err
	}
	certpool.AppendCertsFromPEM(ca)

	// match client cert and key
	cert, err := tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
	if err != nil {
		logrus.Errorf("problem with cert/key pair: %s", err)
		return nil, err
	}

	return &tls.Config{
		RootCAs:            certpool,
		ClientAuth:         tls.RequireAndVerifyClientCert,
		ClientCAs:          nil,
		InsecureSkipVerify: true, //nolint:gosec
		Certificates:       []tls.Certificate{cert},
	}, nil
}

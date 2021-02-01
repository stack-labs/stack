package helper

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/stack-labs/stack-rpc/pkg/metadata"
)

type TLS struct {
	CertFile     string `json:"cert_file"`
	KeyFile      string `json:"key_file"`
	ClientCAFile string `json:"client_ca_file"`
}

// RequestToContext of helper...
func RequestToContext(r *http.Request) context.Context {
	ctx := context.Background()
	md := make(metadata.Metadata)
	for k, v := range r.Header {
		md[k] = strings.Join(v, ",")
	}
	return metadata.NewContext(ctx, md)
}

// TLSConfig ...
func TLSConfig(conf *TLS) (*tls.Config, error) {
	cert := conf.CertFile
	key := conf.KeyFile
	ca := conf.ClientCAFile

	if len(cert) > 0 && len(key) > 0 {
		certs, err := tls.LoadX509KeyPair(cert, key)
		if err != nil {
			return nil, err
		}

		if len(ca) > 0 {
			caCert, err := ioutil.ReadFile(ca)
			if err != nil {
				return nil, err
			}

			caCertPool := x509.NewCertPool()
			caCertPool.AppendCertsFromPEM(caCert)

			return &tls.Config{
				Certificates: []tls.Certificate{certs},
				ClientCAs:    caCertPool,
				ClientAuth:   tls.RequireAndVerifyClientCert,
				NextProtos:   []string{"h2", "http/1.1"},
			}, nil
		}

		return &tls.Config{
			Certificates: []tls.Certificate{certs}, NextProtos: []string{"h2", "http/1.1"},
		}, nil
	}

	return nil, errors.New("TLS certificate and key files not specified")
}

// ServeCORS ...
func ServeCORS(w http.ResponseWriter, r *http.Request) {
	set := func(w http.ResponseWriter, k, v string) {
		if v := w.Header().Get(k); len(v) > 0 {
			return
		}
		w.Header().Set(k, v)
	}

	if origin := r.Header.Get("Origin"); len(origin) > 0 {
		set(w, "Access-Control-Allow-Origin", origin)
	} else {
		set(w, "Access-Control-Allow-Origin", "*")
	}

	set(w, "Access-Control-Allow-Methods", "POST, PATCH, GET, OPTIONS, PUT, DELETE")
	set(w, "Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

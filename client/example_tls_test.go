// +build TLS

/*
Copyright 2019 Stefan Miller

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package client_test

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/stfnmllr/go-resp3/client"
)

/*
This example is based on the TLS Support getting started of Redis 6.
Please see https://github.com/antirez/redis/blob/unstable/TLS.md for details.

The go-resp3 client TLS support is solely based on Go standard library capabilities, especially the
- TLS configuration (https://golang.org/pkg/crypto/tls/#Config) and
- client connection (https://golang.org/pkg/crypto/tls/#Client)

Please check the Go standard library documentation for further information.
*/

const (
	certFile   = "tls/redis.crt"
	keyFile    = "tls/redis.key"
	cacertFile = "tls/ca.crt"
)

func Example_TLS() {
	// Root certificate.
	cacert, err := ioutil.ReadFile(cacertFile)
	if err != nil {
		log.Fatal(err)
	}
	roots := x509.NewCertPool()
	if !roots.AppendCertsFromPEM(cacert) {
		log.Fatal("failed to parse root certificate")
	}

	// Key pair.
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatal(err)
	}

	// Create TLS configuration.
	config := &tls.Config{
		// whether the ServerName needs to be provided or InsecureSkipVerify needs to be set
		// please see the comments in https://golang.org/pkg/crypto/tls/#Client
		InsecureSkipVerify: true,
		RootCAs:            roots,
		Certificates:       []tls.Certificate{cert},
	}

	// Create dialer with TLS configuration.
	dialer := &client.Dialer{TLSConfig: config}

	// Open connection with standard host and port.
	conn, err := dialer.Dial("")
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	if err := conn.Set("mykey", "Hello Redis").Err(); err != nil {
		log.Fatal(err)
	}
	value, err := conn.Get("mykey").ToString()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(value)
	// Output:
	// Hello Redis
}

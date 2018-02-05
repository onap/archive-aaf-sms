/*
 * Copyright 2018 Intel Corporation, Inc
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package auth

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
)

var tlsConfig *tls.Config

// GetTLSConfig initializes a tlsConfig using the CA's certificate
// This config is then used to enable the server for mutual TLS
func GetTLSConfig(caCertFile string) *tls.Config {
	// Initialize tlsConfig once
	if tlsConfig == nil {
		caCert, err := ioutil.ReadFile(caCertFile)

		if err != nil {
			log.Fatal("Error reading CA Certificate")
			log.Fatal(err)
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		tlsConfig = &tls.Config{
			ClientAuth: tls.RequireAndVerifyClientCert,
			ClientCAs:  caCertPool,
			MinVersion: tls.VersionTLS12,
		}
		tlsConfig.BuildNameToCertificate()
	}
	return tlsConfig
}

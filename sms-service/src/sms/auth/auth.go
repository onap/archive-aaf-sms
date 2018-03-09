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
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"golang.org/x/crypto/openpgp"
	"io/ioutil"

	smslogger "sms/log"
)

var tlsConfig *tls.Config

// GetTLSConfig initializes a tlsConfig using the CA's certificate
// This config is then used to enable the server for mutual TLS
func GetTLSConfig(caCertFile string) (*tls.Config, error) {
	// Initialize tlsConfig once
	if tlsConfig == nil {
		caCert, err := ioutil.ReadFile(caCertFile)

		if err != nil {
			return nil, err
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
	return tlsConfig, nil
}

// GeneratePGPKeyPair produces a PGP key pair and returns
// two things:
// A base64 encoded form of the public part of the entity
// A base64 encoded form of the private key
func GeneratePGPKeyPair() (string, string, error) {
	var entity *openpgp.Entity
	entity, err := openpgp.NewEntity("aaf.sms.init", "PGP Key for unsealing", "", nil)
	if err != nil {
		smslogger.WriteError(err.Error())
		return "", "", err
	}

	// Sign the identity in the entity
	for _, id := range entity.Identities {
		err = id.SelfSignature.SignUserId(id.UserId.Id, entity.PrimaryKey, entity.PrivateKey, nil)
		if err != nil {
			smslogger.WriteError(err.Error())
			return "", "", err
		}
	}

	// Sign the subkey in the entity
	for _, subkey := range entity.Subkeys {
		err := subkey.Sig.SignKey(subkey.PublicKey, entity.PrivateKey, nil)
		if err != nil {
			smslogger.WriteError(err.Error())
			return "", "", err
		}
	}

	buffer := new(bytes.Buffer)
	entity.Serialize(buffer)
	pbkey := base64.StdEncoding.EncodeToString(buffer.Bytes())

	buffer.Reset()
	entity.SerializePrivate(buffer, nil)
	prkey := base64.StdEncoding.EncodeToString(buffer.Bytes())

	return pbkey, prkey, nil
}

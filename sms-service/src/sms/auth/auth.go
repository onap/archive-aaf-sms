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
	"golang.org/x/crypto/openpgp/packet"
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
			// Change to RequireAndVerify once we have mandatory certs
			ClientAuth: tls.VerifyClientCertIfGiven,
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

// DecryptPGPBytes decrypts a PGP encoded input string and returns
// a base64 representation of the decoded string
func DecryptPGPBytes(data string, prKey string) (string, error) {
	// Convert private key to bytes from base64
	prKeyBytes, err := base64.StdEncoding.DecodeString(prKey)
	if err != nil {
		smslogger.WriteError("Error Decoding base64 private key: " + err.Error())
		return "", err
	}

	dataBytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		smslogger.WriteError("Error Decoding base64 data: " + err.Error())
		return "", err
	}

	prEntity, err := openpgp.ReadEntity(packet.NewReader(bytes.NewBuffer(prKeyBytes)))
	if err != nil {
		smslogger.WriteError("Error reading entity from PGP key: " + err.Error())
		return "", err
	}

	prEntityList := &openpgp.EntityList{prEntity}
	message, err := openpgp.ReadMessage(bytes.NewBuffer(dataBytes), prEntityList, nil, nil)
	if err != nil {
		smslogger.WriteError("Error Decrypting message: " + err.Error())
		return "", err
	}

	var retBuf bytes.Buffer
	retBuf.ReadFrom(message.UnverifiedBody)

	return retBuf.String(), nil
}

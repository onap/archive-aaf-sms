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
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/packet"
	"io/ioutil"

	smslogger "sms/log"
)

var tlsConfig *tls.Config

func checkError(err error, topic string) error {
	if err != nil {
		smslogger.WriteError(topic + ": " + err.Error())
		return err
	}

	return nil
}

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
	config := &packet.Config{
		DefaultHash: crypto.SHA256,
	}

	entity, err := openpgp.NewEntity("aaf.sms.init", "PGP Key for unsealing", "", config)
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

// EncryptPGPString takes data and a public key and encrypts using that
// public key
func EncryptPGPString(data string, pbKey string) (string, error) {
	pbKeyBytes, err := base64.StdEncoding.DecodeString(pbKey)
	if checkError(err, "Decoding Base64 Public Key") != nil {
		return "", err
	}

	dataBytes := []byte(data)

	pbEntity, err := openpgp.ReadEntity(packet.NewReader(bytes.NewBuffer(pbKeyBytes)))
	if checkError(err, "Reading entity from PGP key") != nil {
		return "", err
	}

	// encrypt string
	buf := new(bytes.Buffer)
	out, err := openpgp.Encrypt(buf, []*openpgp.Entity{pbEntity}, nil, nil, nil)
	if checkError(err, "Creating Encryption Pipe") != nil {
		return "", err
	}

	_, err = out.Write(dataBytes)
	if checkError(err, "Writing to Encryption Pipe") != nil {
		return "", err
	}

	err = out.Close()
	if checkError(err, "Closing Encryption Pipe") != nil {
		return "", err
	}

	crp := base64.StdEncoding.EncodeToString(buf.Bytes())
	return crp, nil
}

// DecryptPGPString decrypts a PGP encoded input string and returns
// a base64 representation of the decoded string
func DecryptPGPString(data string, prKey string) (string, error) {
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

// ReadFromFile reads a file and loads the PGP key into
// a string
func ReadFromFile(fileName string) (string, error) {

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		smslogger.WriteError(err.Error())
		smslogger.WriteError("Cannot read file: " + fileName)
		return "", err
	}
	return string(data), nil

}

// WriteToFile writes a PGP key into a file.
// It will truncate the file if it exists
func WriteToFile(data string, fileName string) error {

	err := ioutil.WriteFile(fileName, []byte(data), 0600)
	if err != nil {
		smslogger.WriteError(err.Error())
		smslogger.WriteError("Cannot write to file: " + fileName)
		return err
	}
	return nil

}

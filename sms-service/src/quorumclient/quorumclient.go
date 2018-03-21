/*
* Copyright 2018 TechMahindra
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

package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	uuid "github.com/hashicorp/go-uuid"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	smsauth "sms/auth"
	smslogger "sms/log"
	"strings"
	"time"
)

func readMyID() (string, error) {
	data, err := ioutil.ReadFile("myid.cfg")
	if err != nil {
		smslogger.WriteWarn("Unable to read myid.cfg")
		return "", err
	}
	return string(data), nil
}

func writeMyID(data string) error {
	err := ioutil.WriteFile("myid.cfg", []byte(data), 0600)
	if err != nil {
		smslogger.WriteWarn("Unable to write myid.cfg")
		return err
	}
	return nil
}

//This application checks the backend status and
//calls necessary initialization endpoints on the
//SMS webservice
func main() {
	smslogger.Init("quorumclient.log")
	smslogger.WriteInfo("Starting Log for Quorum Client")

	//myID is used to uniquely identify the quorum client
	//Using any other information such as hostname is not
	//guaranteed to be unique.
	//In Kubernetes, pod restarts will also change the hostname
	myID, err := readMyID()
	if err != nil {
		smslogger.WriteWarn("Unable to find an ID for this client. Generating...")
		myID, _ = uuid.GenerateUUID()
		writeMyID(myID)
	}

	//Struct to read json configuration file
	type config struct {
		BackEndURL string `json:"url"`
		CAFile     string `json:"cafile"`
		ClientCert string `json:"clientcert"`
		ClientKey  string `json:"clientkey"`
		B64Key     string `json:"key"`
		TimeOut    string `json:"timeout"`
		DisableTLS bool   `json:"disable_tls"`
	}

	//Load the config File for reading
	vcf, err := os.Open("config.json")
	if err != nil {
		log.Fatalf("Error reading config file %v", err)
	}

	cfg := config{}
	decoder := json.NewDecoder(vcf)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatalf("Error while parsing config file %v", err)
	}

	transport := http.Transport{}

	if cfg.DisableTLS {
		// Read the CA cert. This can be the self-signed CA
		// or CA cert provided by AAF
		caCert, err := ioutil.ReadFile(cfg.CAFile)
		if err != nil {
			log.Fatalf("Error while reading CA file %v ", err)
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		// Load the client certificate files
		cert, err := tls.LoadX509KeyPair(cfg.ClientCert, cfg.ClientKey)
		if err != nil {
			log.Fatalf("Error while loading key pair %v ", err)
		}

		transport.TLSClientConfig = &tls.Config{
			RootCAs:      caCertPool,
			Certificates: []tls.Certificate{cert},
		}
	}

	client := &http.Client{
		Transport: &transport,
	}

	pbkey, prkey, err := smsauth.GeneratePGPKeyPair()
	if err != nil {
		smslogger.WriteError("Unable to create PGP Key Pair, using defaults: " + err.Error())
		// Following commands should always work
		pbkey, _ = smsauth.ReadKeysFromFile("def.pb")
		prkey, _ = smsauth.ReadKeysFromFile("def.pr")
	}

	smsauth.WriteKeysToFile(pbkey, "gen.pb")
	smsauth.WriteKeysToFile(prkey, "gen.pr")

	duration, _ := time.ParseDuration(cfg.TimeOut)
	ticker := time.NewTicker(duration)
	registrationDone := false

	for _ = range ticker.C {

		//URL and Port is configured in config file
		response, err := client.Get(cfg.BackEndURL + "/v1/sms/quorum/status")
		if err != nil {
			smslogger.WriteError("Unable to connect to SMS. Retrying...")
			continue
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatalf("Error while reading response %v", err)
		}

		var data map[string]interface{}
		json.Unmarshal(responseData, &data)
		sealed := data["sealed"].(bool)

		// Unseal the vault if sealed
		if sealed {
			//Register with SMS if not already done so
			if !registrationDone {
				body := strings.NewReader(
					`{"pgpkey":" + pbkey + ",
					"quorumid":" + myID  "}`)
				res, err := client.Post(cfg.BackEndURL+"/v1/sms/quorum/register", "application/json", body)
				registrationDone = true
			}

			decdB64Key, _ := base64.StdEncoding.DecodeString(cfg.B64Key)
			body := strings.NewReader(`{"key":"` + string(decdB64Key) + `"}`)
			//URL and PORT is configured via config file
			response, err = client.Post(cfg.BackEndURL+"/v1/sms/quorum/unseal", "application/json", body)
			if err != nil {
				log.Fatalf("Error while unsealing %v", err)
			}
		}
	}
}

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
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

//This application checks the backend status and
//calls necessary initialization endpoints on the
//SMS webservice
func main() {
	//Struct to read json configuration file
	type config struct {
		BackEndURL string `json:"url"`
		CAFile     string `json:"cafile"`
		ClientCert string `json:"clientcert"`
		ClientKey  string `json:"clientkey"`
		B64Key     string `json:"key"`
		TimeOut    string `json:"timeout"`
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

	duration, _ := time.ParseDuration(cfg.TimeOut)
	ticker := time.NewTicker(duration)

	for _ = range ticker.C {

		caCert, err := ioutil.ReadFile(cfg.CAFile)
		if err != nil {
			log.Fatalf("Error while reading CA file %v ", err)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		cert, err := tls.LoadX509KeyPair(cfg.ClientCert, cfg.ClientKey)
		if err != nil {
			log.Fatalf("Error while loading key pair %v ", err)
		}

		client := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs:      caCertPool,
					Certificates: []tls.Certificate{cert},
				},
			},
		}
		//URL and Port is configured in config file
		response, err := client.Get(cfg.BackEndURL + "v1/sms/status")
		if err != nil {
			log.Fatalf("Error while connecting to SMS webservice %v", err)
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
			decdB64Key, _ := base64.StdEncoding.DecodeString(cfg.B64Key)
			body := strings.NewReader(`{"key":"` + string(decdB64Key) + `"}`)
			//URL and PORT is configured via config file
			response, err = client.Post(cfg.BackEndURL+"v1/sms/unseal", "application/json", body)
			if err != nil {
				log.Fatalf("Error while unsealing %v", err)
			}
		}
	}
}

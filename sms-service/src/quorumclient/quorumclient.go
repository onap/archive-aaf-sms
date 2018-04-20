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
	"encoding/json"
	uuid "github.com/hashicorp/go-uuid"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	smsauth "sms/auth"
	smslogger "sms/log"
	"strings"
	"time"
)

func loadPGPKeys(prKeyPath string, pbKeyPath string) (string, string, error) {

	var pbkey, prkey string
	generated := false
	prkey, err := smsauth.ReadFromFile(prKeyPath)
	if smslogger.CheckError(err, "LoadPGP Private Key") != nil {
		smslogger.WriteInfo("No Private Key found. Generating...")
		pbkey, prkey, _ = smsauth.GeneratePGPKeyPair()
		generated = true
	} else {
		pbkey, err = smsauth.ReadFromFile(pbKeyPath)
		if smslogger.CheckError(err, "LoadPGP Public Key") != nil {
			smslogger.WriteWarn("No Public Key found. Generating...")
			pbkey, prkey, _ = smsauth.GeneratePGPKeyPair()
			generated = true
		}
	}

	// Storing the keys to file to allow for recovery during restarts
	if generated {
		smsauth.WriteToFile(prkey, prKeyPath)
		smsauth.WriteToFile(pbkey, pbKeyPath)
	}

	return pbkey, prkey, nil

}

//This application checks the backend status and
//calls necessary initialization endpoints on the
//SMS webservice
func main() {
	podName := os.Getenv("HOSTNAME")
	idFilePath := filepath.Join("auth", podName, "id")
	pbKeyPath := filepath.Join("auth", podName, "pbkey")
	prKeyPath := filepath.Join("auth", podName, "prkey")
	shardPath := filepath.Join("auth", podName, "shard")

	smslogger.Init("quorum.log")
	smslogger.WriteInfo("Starting Log for Quorum Client")

	/*
		myID is used to uniquely identify the quorum client
		Using any other information such as hostname is not
		guaranteed to be unique.
		In Kubernetes, pod restarts will also change the hostname
	*/
	myID, err := smsauth.ReadFromFile(idFilePath)
	if smslogger.CheckError(err, "Read ID") != nil {
		smslogger.WriteWarn("Unable to find an ID for this client. Generating...")
		myID, _ = uuid.GenerateUUID()
		smsauth.WriteToFile(myID, idFilePath)
	}

	/*
		readMyShard will read the shard from disk when this client
		instance restarts. It will return err when a shard is not found.
		This is the case for first startup
	*/
	registrationDone := true
	myShard, err := smsauth.ReadFromFile(shardPath)
	if smslogger.CheckError(err, "Read Shard") != nil {
		smslogger.WriteWarn("Unable to find a shard file. Registering with SMS...")
		registrationDone = false
	}

	pbkey, prkey, _ := loadPGPKeys(prKeyPath, pbKeyPath)

	//Struct to read json configuration file
	type config struct {
		BackEndURL string `json:"url"`
		CAFile     string `json:"cafile"`
		ClientCert string `json:"clientcert"`
		ClientKey  string `json:"clientkey"`
		TimeOut    string `json:"timeout"`
		DisableTLS bool   `json:"disable_tls"`
	}

	//Load the config File for reading
	vcf, err := os.Open("config.json")
	if err != nil {
		log.Fatalf("Error reading config file %v", err)
	}

	cfg := config{}
	err = json.NewDecoder(vcf).Decode(&cfg)
	if err != nil {
		log.Fatalf("Error while parsing config file %v", err)
	}

	transport := http.Transport{}

	if cfg.DisableTLS == false {
		// Read the CA cert. This can be the self-signed CA
		// or CA cert provided by AAF
		caCert, err := ioutil.ReadFile(cfg.CAFile)
		if err != nil {
			log.Fatalf("Error while reading CA file %v ", err)
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		// Load the client certificate files
		//cert, err := tls.LoadX509KeyPair(cfg.ClientCert, cfg.ClientKey)
		//if err != nil {
		//	log.Fatalf("Error while loading key pair %v ", err)
		//}

		transport.TLSClientConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
			RootCAs:    caCertPool,
			//Enable once we have proper client certificates
			//Certificates: []tls.Certificate{cert},
		}
	}

	client := &http.Client{
		Transport: &transport,
	}

	duration, _ := time.ParseDuration(cfg.TimeOut)
	ticker := time.NewTicker(duration)

	for _ = range ticker.C {

		//URL and Port is configured in config file
		response, err := client.Get(cfg.BackEndURL + "/v1/sms/quorum/status")
		if smslogger.CheckError(err, "Connect to SMS") != nil {
			continue
		}

		var data struct {
			Seal bool `json:"sealstatus"`
		}
		err = json.NewDecoder(response.Body).Decode(&data)

		sealed := data.Seal

		// Unseal the vault if sealed
		if sealed {
			//Register with SMS if not already done so
			if !registrationDone {
				body := strings.NewReader(`{"pgpkey":"` + pbkey + `","quorumid":"` + myID + `"}`)
				res, err := client.Post(cfg.BackEndURL+"/v1/sms/quorum/register", "application/json", body)
				if smslogger.CheckError(err, "Register with SMS") != nil {
					continue
				}
				registrationDone = true
				var data struct {
					Shard string `json:"shard"`
				}
				json.NewDecoder(res.Body).Decode(&data)
				myShard = data.Shard
				smsauth.WriteToFile(myShard, shardPath)
			}

			decShard, err := smsauth.DecryptPGPString(myShard, prkey)
			body := strings.NewReader(`{"unsealshard":"` + decShard + `"}`)
			//URL and PORT is configured via config file
			response, err = client.Post(cfg.BackEndURL+"/v1/sms/quorum/unseal", "application/json", body)
			if smslogger.CheckError(err, "Unsealing Vault") != nil {
				continue
			}
		}
	}
}

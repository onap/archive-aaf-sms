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
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

//This application will checks vault status,
//if found sealed it will call unseal vault webservice
func main() {
	//Struct to read vault json configuration file
	type vaultConfig struct {
		B64Key  string `json:"key"`
		TimeOut string `json:"timeout"`
	}
	//Load the VaultConfig File for reading
	vcf, err := os.Open("vaultConfig.json")
	if err != nil {
		log.Fatalf("Error reading vault config file %v", err)
	}

	var VaultConfig *vaultConfig
	decoder := json.NewDecoder(vcf)
	err = decoder.Decode(&VaultConfig)
	if err != nil {
		log.Fatalf("Error while parsing vault config file %v", err)
	}

	duration, _ := time.ParseDuration(VaultConfig.TimeOut)

	for _ = range time.NewTicker(duration).C {
		//Currently using a localhost host, later will be replaced with
		//exact url
		response, err := http.Get("http://localhost:8200/v1/sys/seal-status")
		if err != nil {
			log.Fatalf("Error while connecting to vault webservice %v", err)
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatalf("Error while reading response %v", err)
		}
		var vaultdata map[string]interface{}
		json.Unmarshal(responseData, &vaultdata)
		vaultsealed := vaultdata["sealed"].(bool)
		// Unseal the vault if sealed
		if vaultsealed {
			decdB64Key, _ := base64.StdEncoding.DecodeString(VaultConfig.B64Key)
			body := strings.NewReader(`{"key":"` + string(decdB64Key) + `"}`)
			//below url will be replaced with exact webservice
			response, err = http.Post("http://127.0.0.1:8200/v1/sys/unseal", "application/x-www-form-urlencoded", body)
			if err != nil {
				log.Fatalf("Error while unsealing %v", err)
			}
		}
	}

}

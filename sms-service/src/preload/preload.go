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

package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	pkgerrors "github.com/pkg/errors"
)

//DataYaml stores a list of domains from JSON file
type DataJSON struct {
	Domain SecretDomainJSON `json:"domain,omitempty"`
	Domains SecretDomainJSON `json:"domains,omitempty"`
}

//SecretDomainJSON stores a name for the Domain and a list of Secrets
type SecretDomainJSON struct {
	Name    string       `json:"name"`
	Secrets []SecretJSON `json:"secrets"`
}

//SecretJSON stores a name for the Secret and a list of Values
type SecretJSON struct {
	Name   string                 `json:"name"`
	Values map[string]interface{} `json:"values"`
}

func processJSONFile(name string) (DataJSON, error) {

	data, err := ioutil.ReadFile(name)
	if err != nil {
		return DataJSON{}, pkgerrors.Cause(err)
	}

	d := DataJSON{}
	err = json.Unmarshal(data, &d)
	if err != nil {
		return DataJSON{}, pkgerrors.Cause(err)
	}

	return d, nil
}

type smsClient struct {
	BaseURL *url.URL
	//In seconds
	Timeout    int
	CaCertPath string

	httpClient *http.Client
}

func (c *smsClient) init() error {

	skipVerify := false
	caCert, err := ioutil.ReadFile(c.CaCertPath)
	if err != nil {
		fmt.Println(pkgerrors.Cause(err))
		fmt.Println("Using Insecure Server Verification")
		skipVerify = true
	}

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	tlsConfig.InsecureSkipVerify = skipVerify

	// Add cert information when skipVerify is false
	if skipVerify == false {
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		tlsConfig.RootCAs = caCertPool
	}

	tr := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	c.httpClient = &http.Client{
		Transport: tr,
		Timeout:   time.Duration(c.Timeout) * time.Second,
	}

	return nil
}

func (c *smsClient) sendPostRequest(relURL string, message map[string]interface{}) error {

    fmt.Println(message)
	rel, err := url.Parse(relURL)
	if err != nil {
		return pkgerrors.Cause(err)
	}
	u := c.BaseURL.ResolveReference(rel)

	body, err := json.Marshal(message)
	if err != nil {
		fmt.Println(body)
		return pkgerrors.Cause(err)
	}

	resp, err := c.httpClient.Post(u.String(), "application/json", bytes.NewBuffer(body))
	if err != nil {
		return pkgerrors.Cause(err)
	}

	if resp.StatusCode >= 400 && resp.StatusCode < 600 {
		// Request Failed
		errText, _ := ioutil.ReadAll(resp.Body)
		return pkgerrors.Errorf("Request Failed with: %s and Error: %s",
			resp.Status, string(errText))
	}

	return nil
}

func (c *smsClient) createDomain(domain string) error {

	message := map[string]interface{}{
		"name": domain,
	}
	url := "/v1/sms/domain"
	err := c.sendPostRequest(url, message)
	if err != nil {
		return pkgerrors.Cause(err)
	}
	return nil
}

func (c *smsClient) createSecret(domain string, secret string,
	values map[string]interface{}) error {

	//Create a map out of the valueYaml array
	//mapValues := make(map[string]string)
	//for _, val := range values {
	//	mapValues[val.Key] = val.Value
	//}

	message := map[string]interface{}{
		"name":   secret,
		"values": values,
	}

	fmt.Println(message)

	url := "/v1/sms/domain/" + strings.TrimSpace(domain) + "/secret"
	err := c.sendPostRequest(url, message)
	if err != nil {
		return pkgerrors.Cause(err)
	}

	return nil
}

func (c *smsClient) uploadToSMS(data DataJSON) error {

	err := c.createDomain(data.Domain.Name)
	if err != nil {
		return pkgerrors.Cause(err)
	}

    for _, s := range data.Domain.Secrets {
		err = c.createSecret(data.Domain.Name, s.Name, s.Values)
		if err != nil {
			return pkgerrors.Cause(err)
		}
	}

	return nil
}

func main() {

	cacert := flag.String("cacert", "/sms/certs/aaf_root_ca.cer",
		"Path to the CA Certificate file")
	serviceurl := flag.String("serviceurl", "https://aaf-sms.onap",
		"Url for the SMS Service")
	serviceport := flag.Int("serviceport", 10443,
		"Service port if its different than the default")
	jsondir := flag.String("jsondir", ".",
		"Folder containing json files to upload")

	flag.Parse()

	files, err := ioutil.ReadDir(*jsondir)
	if err != nil {
		log.Fatal(pkgerrors.Cause(err))
	}

	serviceURL, err := url.Parse(*serviceurl + ":" + strconv.Itoa(*serviceport))
	if err != nil {
		log.Fatal(pkgerrors.Cause(err))
	}

	client := &smsClient{
		Timeout:    30,
		BaseURL:    serviceURL,
		CaCertPath: *cacert,
	}
	client.init()

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			fmt.Println("Processing   ", file.Name())
			d, err := processJSONFile(file.Name())
			if err != nil {
				log.Printf("Error Reading %s : %s", file.Name(), pkgerrors.Cause(err))
				continue
			}

			err = client.uploadToSMS(d)
			if err != nil {
				log.Printf("Error Uploading %s : %s", file.Name(), pkgerrors.Cause(err))
				continue
			}
		}
	}
}


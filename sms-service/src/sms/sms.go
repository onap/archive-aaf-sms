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
	"log"
	"net/http"

	smsauth "sms/auth"
	smsbackend "sms/backend"
	smsconfig "sms/config"
	smshandler "sms/handler"
)

func main() {
	// Read Configuration File
	smsConf, err := smsconfig.ReadConfigFile("smsconfig.json")
	if err != nil {
		log.Fatal(err)
	}

	backendImpl, err := smsbackend.InitSecretBackend()
	if err != nil {
		log.Fatal(err)
	}

	httpRouter := smshandler.CreateRouter(backendImpl)

	// TODO: Use CA certificate from AAF
	tlsConfig, err := smsauth.GetTLSConfig(smsConf.CAFile)
	if err != nil {
		log.Fatal(err)
	}

	httpServer := &http.Server{
		Handler:   httpRouter,
		Addr:      ":10443",
		TLSConfig: tlsConfig,
	}

	err = httpServer.ListenAndServeTLS(smsConf.ServerCert, smsConf.ServerKey)
	log.Fatal(err)
}

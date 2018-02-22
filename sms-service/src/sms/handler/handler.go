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

package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"

	smsbackend "sms/backend"
)

// handler stores two interface implementations that implement
// the backend functionality
type handler struct {
	secretBackend smsbackend.SecretBackend
	loginBackend  smsbackend.LoginBackend
}

// createSecretDomainHandler creates a secret domain with a name provided
func (h handler) createSecretDomainHandler(w http.ResponseWriter, r *http.Request) {
	var d smsbackend.SecretDomain

	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	dom, err := h.secretBackend.CreateSecretDomain(d.Name)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = json.NewEncoder(w).Encode(dom)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
}

// getSecretDomainHandler returns list of secret domains
func (h handler) getSecretDomainHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	domName := vars["domName"]

	fmt.Fprintf(w, "Got req %s", domName)
	//h.secretBackend.GetSecretDomain(domName)
	//encode data into json and return
}

// deleteSecretDomainHandler deletes a secret domain with the name provided
func (h handler) deleteSecretDomainHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	domName := vars["domName"]

	h.secretBackend.DeleteSecretDomain(domName)
}

// createSecretHandler handles creation of secrets on a given domain name
func (h handler) createSecretHandler(w http.ResponseWriter, r *http.Request) {
	// Get domain name from URL
	vars := mux.Vars(r)
	domName := vars["domName"]

	// Get secrets to be stored from body
	var b smsbackend.Secret
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	err = h.secretBackend.CreateSecret(domName, b)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// getSecretHandler handles reading a secret by given domain name and secret name
func (h handler) getSecretHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	domName := vars["domName"]
	secName := vars["secretName"]

	sec, err := h.secretBackend.GetSecret(domName, secName)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = json.NewEncoder(w).Encode(sec)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
}

// deleteSecretHandler handles deleting a secret by given domain name and secret name
func (h handler) deleteSecretHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	domName := vars["domName"]
	secName := vars["secretName"]

	h.secretBackend.DeleteSecret(domName, secName)
}

// struct that tracks various status items for SMS and backend
type backendStatus struct {
	Seal bool `json:"sealstatus"`
}

// statusHandler returns information related to SMS and SMS backend services
func (h handler) statusHandler(w http.ResponseWriter, r *http.Request) {
	s, err := h.secretBackend.GetStatus()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	status := backendStatus{Seal: s}
	err = json.NewEncoder(w).Encode(status)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

// loginHandler handles login via password and username
func (h handler) loginHandler(w http.ResponseWriter, r *http.Request) {

}

// initSMSHandler
func (h handler) initSMSHandler(w http.ResponseWriter, r *http.Request) {

}

// unsealHandler
func (h handler) unsealHandler(w http.ResponseWriter, r *http.Request) {

}

// CreateRouter returns an http.Handler for the registered URLs
// Takes an interface implementation as input
func CreateRouter(b smsbackend.SecretBackend) http.Handler {
	h := handler{secretBackend: b}

	// Create a new mux to handle URL endpoints
	router := mux.NewRouter()

	router.HandleFunc("/v1/sms/login", h.loginHandler).Methods("POST")

	// Initialization APIs which will be used by quorum client
	// to unseal and to provide root token to sms service
	router.HandleFunc("/v1/sms/status", h.statusHandler).Methods("GET")
	router.HandleFunc("/v1/sms/unseal", h.unsealHandler).Methods("POST")
	router.HandleFunc("/v1/sms/init", h.initSMSHandler).Methods("POST")

	router.HandleFunc("/v1/sms/domain", h.createSecretDomainHandler).Methods("POST")
	router.HandleFunc("/v1/sms/domain/{domName}", h.getSecretDomainHandler).Methods("GET")
	router.HandleFunc("/v1/sms/domain/{domName}", h.deleteSecretDomainHandler).Methods("DELETE")

	router.HandleFunc("/v1/sms/domain/{domName}/secret", h.createSecretHandler).Methods("POST")
	router.HandleFunc("/v1/sms/domain/{domName}/secret/{secretName}", h.getSecretHandler).Methods("GET")
	router.HandleFunc("/v1/sms/domain/{domName}/secret/{secretName}", h.deleteSecretHandler).Methods("DELETE")

	return router
}

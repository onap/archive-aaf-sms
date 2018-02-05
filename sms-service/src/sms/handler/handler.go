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
	"net/http"

	"sms/backend"

	"github.com/gorilla/mux"
)

type secretDomainJSON struct {
	name string
}

type secretKeyValue struct {
	name  string
	value string
}

type secretJSON struct {
	name   string
	values []secretKeyValue
}

type handler struct {
	secretBackend backend.SecretBackend
	loginBackend  backend.LoginBackend
}

// GetSecretDomainHandler returns list of secret domains
func (h handler) GetSecretDomainHandler(w http.ResponseWriter, r *http.Request) {

}

// CreateSecretDomainHandler creates a secret domain with a name provided
func (h handler) CreateSecretDomainHandler(w http.ResponseWriter, r *http.Request) {
	var d secretDomainJSON

	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
}

// DeleteSecretDomainHandler deletes a secret domain with the ID provided
func (h handler) DeleteSecretDomainHandler(w http.ResponseWriter, r *http.Request) {

}

// struct that tracks various status items for SMS and backend
type status struct {
	Seal bool `json:"sealstatus"`
}

// StatusHandler returns information related to SMS and SMS backend services
func (h handler) StatusHandler(w http.ResponseWriter, r *http.Request) {
	s := h.secretBackend.GetStatus()
	status := status{Seal: s}
	err := json.NewEncoder(w).Encode(status)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
}

// LoginHandler handles login via password and username
func (h handler) LoginHandler(w http.ResponseWriter, r *http.Request) {

}

// CreateRouter returns an http.Handler for the registered URLs
func CreateRouter(b backend.SecretBackend) http.Handler {
	h := handler{secretBackend: b}

	// Create a new mux to handle URL endpoints
	router := mux.NewRouter()

	router.HandleFunc("/v1/sms/login", h.LoginHandler).Methods("POST")

	router.HandleFunc("/v1/sms/status", h.StatusHandler).Methods("GET")

	router.HandleFunc("/v1/sms/domain", h.GetSecretDomainHandler).Methods("GET")
	router.HandleFunc("/v1/sms/domain", h.CreateSecretDomainHandler).Methods("POST")
	router.HandleFunc("/v1/sms/domain/{domName}", h.DeleteSecretDomainHandler).Methods("DELETE")

	return router
}

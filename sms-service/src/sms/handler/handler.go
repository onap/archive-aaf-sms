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
	"github.com/gorilla/mux"
	"net/http"

	smsbackend "sms/backend"
	smslogger "sms/log"
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
		smslogger.WriteError(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dom, err := h.secretBackend.CreateSecretDomain(d.Name)
	if err != nil {
		smslogger.WriteError(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jdata, err := json.Marshal(dom)
	if err != nil {
		smslogger.WriteError(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jdata)
}

// deleteSecretDomainHandler deletes a secret domain with the name provided
func (h handler) deleteSecretDomainHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	domName := vars["domName"]

	err := h.secretBackend.DeleteSecretDomain(domName)
	if err != nil {
		smslogger.WriteError(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
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
		smslogger.WriteError(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.secretBackend.CreateSecret(domName, b)
	if err != nil {
		smslogger.WriteError(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		smslogger.WriteError(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jdata, err := json.Marshal(sec)
	if err != nil {
		smslogger.WriteError(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jdata)
}

// listSecretHandler handles listing all secrets under a particular domain name
func (h handler) listSecretHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	domName := vars["domName"]

	secList, err := h.secretBackend.ListSecret(domName)
	if err != nil {
		smslogger.WriteError(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Creating an anonymous struct to store the returned list of data
	var retStruct = struct {
		SecretNames []string `json:"secretnames"`
	}{
		secList,
	}

	jdata, err := json.Marshal(retStruct)
	if err != nil {
		smslogger.WriteError(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jdata)
}

// deleteSecretHandler handles deleting a secret by given domain name and secret name
func (h handler) deleteSecretHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	domName := vars["domName"]
	secName := vars["secretName"]

	err := h.secretBackend.DeleteSecret(domName, secName)
	if err != nil {
		smslogger.WriteError(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// struct that tracks various status items for SMS and backend
type backendStatus struct {
	Seal bool `json:"sealstatus"`
}

// statusHandler returns information related to SMS and SMS backend services
func (h handler) statusHandler(w http.ResponseWriter, r *http.Request) {
	s, err := h.secretBackend.GetStatus()
	if err != nil {
		smslogger.WriteError(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	status := backendStatus{Seal: s}
	jdata, err := json.Marshal(status)
	if err != nil {
		smslogger.WriteError(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jdata)
}

// loginHandler handles login via password and username
func (h handler) loginHandler(w http.ResponseWriter, r *http.Request) {

}

// unsealHandler is a pass through that sends requests from quorum client
// to the backend.
func (h handler) unsealHandler(w http.ResponseWriter, r *http.Request) {
	// Get shards to be used for unseal
	type unsealStruct struct {
		UnsealShard string `json:"unsealshard"`
	}

	var inp unsealStruct
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&inp)
	if err != nil {
		smslogger.WriteError(err.Error())
		http.Error(w, "Bad input JSON", http.StatusBadRequest)
		return
	}

	err = h.secretBackend.Unseal(inp.UnsealShard)
	if err != nil {
		smslogger.WriteError(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// registerHandler allows the quorum clients to register with SMS
// with their PGP public keys that are then used by sms for backend
// initialization
func (h handler) registerHandler(w http.ResponseWriter, r *http.Request) {
	// Get shards to be used for unseal
	type registerStruct struct {
		PGPKey   string `json:"pgpkey"`
		QuorumID string `json:"quorumid"`
	}

	var inp registerStruct
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&inp)
	if err != nil {
		smslogger.WriteError(err.Error())
		http.Error(w, "Bad input JSON", http.StatusBadRequest)
		return
	}

	err = h.secretBackend.RegisterQuorum(inp.PGPKey, inp.QuorumID)
	if err != nil {
		smslogger.WriteError(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
	router.HandleFunc("/v1/sms/quorum/status", h.statusHandler).Methods("GET")
	router.HandleFunc("/v1/sms/quorum/unseal", h.unsealHandler).Methods("POST")
	router.HandleFunc("/v1/sms/quorum/register", h.registerHandler).Methods("POST")

	router.HandleFunc("/v1/sms/domain", h.createSecretDomainHandler).Methods("POST")
	router.HandleFunc("/v1/sms/domain/{domName}", h.deleteSecretDomainHandler).Methods("DELETE")

	router.HandleFunc("/v1/sms/domain/{domName}/secret", h.createSecretHandler).Methods("POST")
	router.HandleFunc("/v1/sms/domain/{domName}/secret", h.listSecretHandler).Methods("GET")
	router.HandleFunc("/v1/sms/domain/{domName}/secret/{secretName}", h.getSecretHandler).Methods("GET")
	router.HandleFunc("/v1/sms/domain/{domName}/secret/{secretName}", h.deleteSecretHandler).Methods("DELETE")

	return router
}

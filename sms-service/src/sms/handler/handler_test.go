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
	"net/http/httptest"
	"reflect"
	smsbackend "sms/backend"
	"strings"
	"testing"
)

var h handler

// Here we are using the anonymous variable feature of golang to
// override methods form an interface
type TestBackend struct {
	smsbackend.SecretBackend
}

func (b *TestBackend) Init() error {
	return nil
}

func (b *TestBackend) GetStatus() (bool, error) {
	return true, nil
}

func (b *TestBackend) GetSecretDomain(name string) (smsbackend.SecretDomain, error) {
	return smsbackend.SecretDomain{}, nil
}

func (b *TestBackend) GetSecret(dom string, sec string) (smsbackend.Secret, error) {
	return smsbackend.Secret{}, nil
}

func (b *TestBackend) ListSecret(dom string) ([]string, error) {
	return nil, nil
}

func (b *TestBackend) CreateSecretDomain(name string) (smsbackend.SecretDomain, error) {
	return smsbackend.SecretDomain{}, nil
}

func (b *TestBackend) CreateSecret(dom string, sec smsbackend.Secret) error {
	return nil
}

func (b *TestBackend) DeleteSecretDomain(name string) error {
	return nil
}

func (b *TestBackend) DeleteSecret(dom string, name string) error {
	return nil
}

func init() {
	testBackend := &TestBackend{}
	h = handler{secretBackend: testBackend}
}

func TestCreateRouter(t *testing.T) {
	router := CreateRouter(h.secretBackend)
	if router == nil {
		t.Fatal("CreateRouter: Got error when none expected")
	}
}

func TestStatusHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/v1/sms/status", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	hr := http.HandlerFunc(h.statusHandler)

	hr.ServeHTTP(rr, req)

	ret := rr.Code
	if ret != http.StatusOK {
		t.Errorf("statusHandler returned wrong status code: %v vs %v",
			ret, http.StatusOK)
	}

	expected := backendStatus{}
	got := backendStatus{}
	expectedStr := strings.NewReader(`{"sealstatus":true}`)
	json.NewDecoder(expectedStr).Decode(&expected)
	json.NewDecoder(rr.Body).Decode(&got)

	if reflect.DeepEqual(expected, got) == false {
		t.Errorf("statusHandler returned unexpected body: got %v vs %v",
			rr.Body.String(), expectedStr)
	}
}

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

package backend

import (
	smsconfig "sms/config"
	smslogger "sms/log"
)

// SecretDomain is where Secrets are stored.
// A single domain can have any number of secrets
type SecretDomain struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

// Secret is the struct that defines the structure of a secret
// It consists of a name and map containing key value pairs
type Secret struct {
	Name   string                 `json:"name"`
	Values map[string]interface{} `json:"values"`
}

// SecretBackend interface that will be implemented for various secret backends
type SecretBackend interface {
	Init() error
	GetStatus() (bool, error)
	Unseal(shard string) error
	RegisterQuorum(pgpkey string) (string, error)

	GetSecret(dom string, sec string) (Secret, error)
	ListSecret(dom string) ([]string, error)

	CreateSecretDomain(name string) (SecretDomain, error)
	CreateSecret(dom string, sec Secret) error

	DeleteSecretDomain(name string) error
	DeleteSecret(dom string, name string) error
}

// InitSecretBackend returns an interface implementation
func InitSecretBackend() (SecretBackend, error) {
	backendImpl := &Vault{
		vaultAddress: smsconfig.SMSConfig.BackendAddress,
	}

	err := backendImpl.Init()
	if err != nil {
		smslogger.WriteError(err.Error())
		return nil, err
	}

	return backendImpl, nil
}

// LoginBackend Interface that will be implemented for various login backends
type LoginBackend interface {
}

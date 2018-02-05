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
	vaultwrap "sms/backend/vault"
)

// SecretDomain struct that will be passed around between http handler
// and code that interfaces with vault
type SecretDomain struct {
	ID         int
	Name       string
	MountPoint string
}

// SecretBackend interface that will be implemented for various secret backends
type SecretBackend interface {
	Init()

	GetStatus() bool
}

// InitSecretBackend returns an interface implementation
func InitSecretBackend() SecretBackend {
	backendImpl := &vaultwrap.Vault{}
	backendImpl.Init()
	return backendImpl
}

// LoginBackend Interface that will be implemented for various login backends
type LoginBackend interface {
}

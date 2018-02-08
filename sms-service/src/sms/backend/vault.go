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
	"fmt"
	vaultapi "github.com/hashicorp/vault/api"

	smsconfig "sms/config"
)

// Vault is the main Struct used in Backend to initialize the struct
type Vault struct {
	vaultClient *vaultapi.Client
}

// Init will initialize the vault connection
// TODO: Check to see if we need to wait for vault to be running
func (v *Vault) Init() error {
	vaultCFG := vaultapi.DefaultConfig()
	vaultCFG.Address = smsconfig.SMSConfig.VaultAddress

	client, err := vaultapi.NewClient(vaultCFG)
	if err != nil {
		return err
	}

	v.vaultClient = client
	return nil
}

// GetStatus returns the current seal status of vault
func (v *Vault) GetStatus() (bool, error) {
	sys := v.vaultClient.Sys()
	sealStatus, err := sys.SealStatus()
	if err != nil {
		return false, err
	}

	return sealStatus.Sealed, nil
}

// GetSecretDomain returns any information related to the secretDomain
// More information can be added in the future with updates to the struct
func (v *Vault) GetSecretDomain(name string) (SecretDomain, error) {

	return SecretDomain{}, nil
}

// GetSecret returns a secret mounted on a particular domain name
// The secret itself is referenced via its name which translates to
// a mount path in vault
func (v *Vault) GetSecret(dom string, sec string) (Secret, error) {

	return Secret{}, nil
}

// CreateSecretDomain mounts the kv backend on a path with the given name
func (v *Vault) CreateSecretDomain(name string) (SecretDomain, error) {

	return SecretDomain{}, nil
}

// CreateSecret creates a secret mounted on a particular domain name
// The secret itself is mounted on a path specified by name
func (v *Vault) CreateSecret(dom string, sec Secret) (Secret, error) {

	return Secret{}, nil
}

// DeleteSecretDomain deletes a secret domain which translates to
// an unmount operation on the given path in Vault
func (v *Vault) DeleteSecretDomain(name string) error {

	return nil
}

// DeleteSecret deletes a secret mounted on the path provided
func (v *Vault) DeleteSecret(dom string, name string) error {

	return nil
}

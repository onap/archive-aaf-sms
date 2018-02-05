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

package vault

import (
	"fmt"
	"log"

	vaultapi "github.com/hashicorp/vault/api"
	smsConfig "sms/config"
)

// Vault is the main Struct used in Backend to initialize the struct
type Vault struct {
	vaultClient *vaultapi.Client
}

// Init will initialize the vault connection
// TODO: Check to see if we need to wait for vault to be running
func (v *Vault) Init() {
	vaultCFG := vaultapi.DefaultConfig()
	vaultCFG.Address = smsConfig.SMSConfig.VaultAddress

	client, err := vaultapi.NewClient(vaultCFG)
	if err != nil {
		log.Fatal(err)
	}

	v.vaultClient = client
}

// GetStatus returns the current seal status of vault
func (v *Vault) GetStatus() bool {
	sys := v.vaultClient.Sys()
	fmt.Println(v.vaultClient.Address())
	sealStatus, err := sys.SealStatus()
	if err != nil {
		log.Fatal(err)
	}

	return sealStatus.Sealed
}

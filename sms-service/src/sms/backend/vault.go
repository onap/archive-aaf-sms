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
	uuid "github.com/hashicorp/go-uuid"
	vaultapi "github.com/hashicorp/vault/api"

	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

// Vault is the main Struct used in Backend to initialize the struct
type Vault struct {
	vaultAddress   string
	vaultToken     string
	vaultMount     string
	vaultTempToken string

	vaultClient       *vaultapi.Client
	engineType        string
	policyName        string
	roleID            string
	secretID          string
	vaultTempTokenTTL time.Time

	tokenLock sync.Mutex
}

// Init will initialize the vault connection
// It will also create the initial policy if it does not exist
// TODO: Check to see if we need to wait for vault to be running
func (v *Vault) Init() error {
	vaultCFG := vaultapi.DefaultConfig()
	vaultCFG.Address = v.vaultAddress
	client, err := vaultapi.NewClient(vaultCFG)
	if err != nil {
		return err
	}

	v.engineType = "kv"
	v.policyName = "smsvaultpolicy"
	v.vaultMount = "sms"
	v.vaultClient = client

	// Check if vault is ready and unsealed
	seal, err := v.GetStatus()
	if err != nil {
		return err
	}
	if seal == true {
		return fmt.Errorf("Vault is still sealed. Unseal before use")
	}

	v.initRole()
	v.checkToken()
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
	// Check if token is still valid
	err := v.checkToken()
	if err != nil {
		return SecretDomain{}, err
	}

	name = strings.TrimSpace(name)
	mountPath := v.vaultMount + "/" + name
	mountInput := &vaultapi.MountInput{
		Type:        v.engineType,
		Description: "Mount point for domain: " + name,
		Local:       false,
		SealWrap:    false,
		Config:      vaultapi.MountConfigInput{},
	}

	err = v.vaultClient.Sys().Mount(mountPath, mountInput)
	if err != nil {
		return SecretDomain{}, err
	}

	uuid, _ := uuid.GenerateUUID()
	return SecretDomain{uuid, name}, nil
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

// initRole is called only once during the service bring up
func (v *Vault) initRole() error {
	// Use the root token once here
	v.vaultClient.SetToken(v.vaultToken)
	defer v.vaultClient.ClearToken()

	rules := `path "sms/*" { capabilities = ["create", "read", "update", "delete", "list"] }
			path "sys/mounts/sms*" { capabilities = ["update","delete","create"] }`
	v.vaultClient.Sys().PutPolicy(v.policyName, rules)

	rName := v.vaultMount + "-role"
	data := map[string]interface{}{
		"token_ttl": "60m",
		"policies":  [2]string{"default", v.policyName},
	}

	// Delete role if it already exists
	v.vaultClient.Logical().Delete("auth/approle/role/" + rName)

	// Mount approle in case its not already mounted
	v.vaultClient.Sys().EnableAuth("approle", "approle", "")

	// Create a role-id
	v.vaultClient.Logical().Write("auth/approle/role/"+rName, data)
	sec, err := v.vaultClient.Logical().Read("auth/approle/role/" + rName + "/role-id")
	if err != nil {
		log.Fatal(err)
	}
	v.roleID = sec.Data["role_id"].(string)

	// Create a secret-id to go with it
	sec, _ = v.vaultClient.Logical().Write("auth/approle/role/"+rName+"/secret-id",
		map[string]interface{}{})
	v.secretID = sec.Data["secret_id"].(string)

	return nil
}

// Function checkToken() gets called multiple times to create
// temporary tokens
func (v *Vault) checkToken() error {
	v.tokenLock.Lock()
	defer v.tokenLock.Unlock()

	// Return immediately if token still has life
	if v.vaultClient.Token() != "" &&
		time.Since(v.vaultTempTokenTTL) < time.Minute*50 {
		return nil
	}

	// Create a temporary token using our roleID and secretID
	out, err := v.vaultClient.Logical().Write("auth/approle/login",
		map[string]interface{}{"role_id": v.roleID, "secret_id": v.secretID})
	if err != nil {
		return err
	}

	tok, err := out.TokenID()

	v.vaultTempToken = tok
	v.vaultTempTokenTTL = time.Now()
	v.vaultClient.SetToken(v.vaultTempToken)
	return nil
}

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

	"errors"
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

	err = v.initRole()
	if err != nil {
		log.Fatalln("Unable to initRole in Vault. Exiting...")
	}

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
func (v *Vault) GetSecret(dom string, name string) (Secret, error) {
	err := v.checkToken()
	if err != nil {
		return Secret{}, errors.New("Token check returned error: " + err.Error())
	}

	dom = v.vaultMount + "/" + dom

	sec, err := v.vaultClient.Logical().Read(dom + "/" + name)
	if err != nil {
		return Secret{}, errors.New("Unable to read Secret at provided path")
	}

	// sec and err are nil in the case where a path does not exist
	if sec == nil {
		return Secret{}, errors.New("Secret not found at the provided path")
	}

	return Secret{Name: name, Values: sec.Data}, nil
}

// ListSecret returns a list of secret names on a particular domain
// The values of the secret are not returned
func (v *Vault) ListSecret(dom string) ([]string, error) {
	err := v.checkToken()
	if err != nil {
		return nil, errors.New("Token check returned error: " + err.Error())
	}

	dom = v.vaultMount + "/" + dom

	sec, err := v.vaultClient.Logical().List(dom)
	if err != nil {
		return nil, errors.New("Unable to read Secret at provided path")
	}

	// sec and err are nil in the case where a path does not exist
	if sec == nil {
		return nil, errors.New("Secret not found at the provided path")
	}

	val, ok := sec.Data["keys"].([]interface{})
	if !ok {
		return nil, errors.New("Secret not found at the provided path")
	}

	retval := make([]string, len(val))
	for i, v := range val {
		retval[i] = fmt.Sprint(v)
	}

	return retval, nil
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
func (v *Vault) CreateSecret(dom string, sec Secret) error {
	err := v.checkToken()
	if err != nil {
		return errors.New("Token checking returned an error" + err.Error())
	}

	dom = v.vaultMount + "/" + dom

	// Vault return is empty on successful write
	_, err = v.vaultClient.Logical().Write(dom+"/"+sec.Name, sec.Values)
	if err != nil {
		return errors.New("Unable to create Secret at provided path")
	}

	return nil
}

// DeleteSecretDomain deletes a secret domain which translates to
// an unmount operation on the given path in Vault
func (v *Vault) DeleteSecretDomain(name string) error {

	return nil
}

// DeleteSecret deletes a secret mounted on the path provided
func (v *Vault) DeleteSecret(dom string, name string) error {
	err := v.checkToken()
	if err != nil {
		return errors.New("Token checking returned an error" + err.Error())
	}

	dom = v.vaultMount + "/" + dom

	// Vault return is empty on successful delete
	_, err = v.vaultClient.Logical().Delete(dom + "/" + name)
	if err != nil {
		return errors.New("Unable to delete Secret at provided path")
	}

	return nil
}

// initRole is called only once during the service bring up
func (v *Vault) initRole() error {
	// Use the root token once here
	v.vaultClient.SetToken(v.vaultToken)
	defer v.vaultClient.ClearToken()

	rules := `path "sms/*" { capabilities = ["create", "read", "update", "delete", "list"] }
			path "sys/mounts/sms*" { capabilities = ["update","delete","create"] }`
	err := v.vaultClient.Sys().PutPolicy(v.policyName, rules)
	if err != nil {
		return errors.New("Unable to create policy for approle creation")
	}

	rName := v.vaultMount + "-role"
	data := map[string]interface{}{
		"token_ttl": "60m",
		"policies":  [2]string{"default", v.policyName},
	}

	// Delete role if it already exists
	_, err = v.vaultClient.Logical().Delete("auth/approle/role/" + rName)
	if err != nil {
		return errors.New("Unable to delete existing role")
	}

	//Check if approle is mounted
	authMounts, err := v.vaultClient.Sys().ListAuth()
	if err != nil {
		return errors.New("Unable to get mounted auth backends")
	}

	approleMounted := false
	for k, v := range authMounts {
		if v.Type == "approle" && k == "approle/" {
			approleMounted = true
			break
		}
	}

	// Mount approle in case its not already mounted
	if !approleMounted {
		v.vaultClient.Sys().EnableAuth("approle", "approle", "")
	}

	// Create a role-id
	v.vaultClient.Logical().Write("auth/approle/role/"+rName, data)
	sec, err := v.vaultClient.Logical().Read("auth/approle/role/" + rName + "/role-id")
	if err != nil {
		return errors.New("Unable to create role ID for approle")
	}
	v.roleID = sec.Data["role_id"].(string)

	// Create a secret-id to go with it
	sec, err = v.vaultClient.Logical().Write("auth/approle/role/"+rName+"/secret-id",
		map[string]interface{}{})
	if err != nil {
		return errors.New("Unable to create secret ID for role")
	}

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

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
	smsauth "sms/auth"
	smslogger "sms/log"

	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

// Vault is the main Struct used in Backend to initialize the struct
type Vault struct {
	sync.Mutex
	initRoleDone          bool
	policyName            string
	roleID                string
	secretID              string
	vaultAddress          string
	vaultClient           *vaultapi.Client
	vaultMountPrefix      string
	internalDomain        string
	internalDomainMounted bool
	vaultTempTokenTTL     time.Time
	vaultToken            string
	shards                []string
	prkey                 string
}

// initVaultClient will create the initial
// Vault strcuture and populate it with the
// right values and it will also create
// a vault client
func (v *Vault) initVaultClient() error {

	vaultCFG := vaultapi.DefaultConfig()
	vaultCFG.Address = v.vaultAddress
	client, err := vaultapi.NewClient(vaultCFG)
	if smslogger.CheckError(err, "Create new vault client") != nil {
		return err
	}

	v.initRoleDone = false
	v.policyName = "smsvaultpolicy"
	v.vaultClient = client
	v.vaultMountPrefix = "sms"
	v.internalDomain = "smsinternaldomain"
	v.internalDomainMounted = false
	v.prkey = ""
	return nil
}

// Init will initialize the vault connection
// It will also initialize vault if it is not
// already initialized.
// The initial policy will also be created
func (v *Vault) Init() error {

	v.initVaultClient()
	// Initialize vault if it is not already
	// Returns immediately if it is initialized
	v.initializeVault()

	err := v.initRole()
	if smslogger.CheckError(err, "InitRole First Attempt") != nil {
		smslogger.WriteInfo("InitRole will try again later")
	}

	return nil
}

// GetStatus returns the current seal status of vault
func (v *Vault) GetStatus() (bool, error) {

	sys := v.vaultClient.Sys()
	sealStatus, err := sys.SealStatus()
	if smslogger.CheckError(err, "Getting Status") != nil {
		return false, errors.New("Error getting status")
	}

	return sealStatus.Sealed, nil
}

// RegisterQuorum registers the PGP public key for a quorum client
// We will return a shard to the client that is registering
func (v *Vault) RegisterQuorum(pgpkey string) (string, error) {

	v.Lock()
	defer v.Unlock()

	if v.shards == nil {
		smslogger.WriteError("Invalid operation in RegisterQuorum")
		return "", errors.New("Invalid operation")
	}
	// Pop the slice
	var sh string
	sh, v.shards = v.shards[len(v.shards)-1], v.shards[:len(v.shards)-1]
	if len(v.shards) == 0 {
		v.shards = nil
	}

	// Decrypt with SMS pgp Key
	sh, _ = smsauth.DecryptPGPString(sh, v.prkey)
	// Encrypt with Quorum client pgp key
	sh, _ = smsauth.EncryptPGPString(sh, pgpkey)

	return sh, nil
}

// Unseal is a passthrough API that allows any
// unseal or initialization processes for the backend
func (v *Vault) Unseal(shard string) error {

	sys := v.vaultClient.Sys()
	_, err := sys.Unseal(shard)
	if smslogger.CheckError(err, "Unseal Operation") != nil {
		return errors.New("Unable to execute unseal operation with specified shard")
	}

	return nil
}

// GetSecret returns a secret mounted on a particular domain name
// The secret itself is referenced via its name which translates to
// a mount path in vault
func (v *Vault) GetSecret(dom string, name string) (Secret, error) {

	err := v.checkToken()
	if smslogger.CheckError(err, "Tocken Check") != nil {
		return Secret{}, errors.New("Token check failed")
	}

	dom = strings.TrimSpace(dom)
	dom = v.vaultMountPrefix + "/" + dom

	sec, err := v.vaultClient.Logical().Read(dom + "/" + name)
	if smslogger.CheckError(err, "Read Secret") != nil {
		return Secret{}, errors.New("Unable to read Secret at provided path")
	}

	// sec and err are nil in the case where a path does not exist
	if sec == nil {
		smslogger.WriteWarn("Vault read was empty. Invalid Path")
		return Secret{}, errors.New("Secret not found at the provided path")
	}

	return Secret{Name: name, Values: sec.Data}, nil
}

// ListSecret returns a list of secret names on a particular domain
// The values of the secret are not returned
func (v *Vault) ListSecret(dom string) ([]string, error) {

	err := v.checkToken()
	if smslogger.CheckError(err, "Token Check") != nil {
		return nil, errors.New("Token check failed")
	}

	dom = strings.TrimSpace(dom)
	dom = v.vaultMountPrefix + "/" + dom

	sec, err := v.vaultClient.Logical().List(dom)
	if smslogger.CheckError(err, "Read Secret") != nil {
		return nil, errors.New("Unable to read Secret at provided path")
	}

	// sec and err are nil in the case where a path does not exist
	if sec == nil {
		smslogger.WriteWarn("Vaultclient returned empty data")
		return nil, errors.New("Secret not found at the provided path")
	}

	val, ok := sec.Data["keys"].([]interface{})
	if !ok {
		smslogger.WriteError("Secret not found at the provided path")
		return nil, errors.New("Secret not found at the provided path")
	}

	retval := make([]string, len(val))
	for i, v := range val {
		retval[i] = fmt.Sprint(v)
	}

	return retval, nil
}

// Mounts the internal Domain if its not already mounted
func (v *Vault) mountInternalDomain(name string) error {

	if v.internalDomainMounted {
		return nil
	}

	name = strings.TrimSpace(name)
	mountPath := v.vaultMountPrefix + "/" + name
	mountInput := &vaultapi.MountInput{
		Type:        "kv",
		Description: "Mount point for domain: " + name,
		Local:       false,
		SealWrap:    false,
		Config:      vaultapi.MountConfigInput{},
	}

	err := v.vaultClient.Sys().Mount(mountPath, mountInput)
	if smslogger.CheckError(err, "Mount internal Domain") != nil {
		if strings.Contains(err.Error(), "existing mount") {
			// It is already mounted
			v.internalDomainMounted = true
			return nil
		}
		// Ran into some other error mounting it.
		return errors.New("Unable to mount internal Domain")
	}

	v.internalDomainMounted = true
	return nil
}

// Stores the UUID created for secretdomain in vault
// under v.vaultMountPrefix / smsinternal domain
func (v *Vault) storeUUID(uuid string, name string) error {

	// Check if token is still valid
	err := v.checkToken()
	if smslogger.CheckError(err, "Token Check") != nil {
		return errors.New("Token Check failed")
	}

	err = v.mountInternalDomain(v.internalDomain)
	if smslogger.CheckError(err, "Mount Internal Domain") != nil {
		return err
	}

	secret := Secret{
		Name: name,
		Values: map[string]interface{}{
			"uuid": uuid,
		},
	}

	err = v.CreateSecret(v.internalDomain, secret)
	if smslogger.CheckError(err, "Write UUID to domain") != nil {
		return err
	}

	return nil
}

// CreateSecretDomain mounts the kv backend on a path with the given name
func (v *Vault) CreateSecretDomain(name string) (SecretDomain, error) {

	// Check if token is still valid
	err := v.checkToken()
	if smslogger.CheckError(err, "Token Check") != nil {
		return SecretDomain{}, errors.New("Token Check failed")
	}

	name = strings.TrimSpace(name)
	mountPath := v.vaultMountPrefix + "/" + name
	mountInput := &vaultapi.MountInput{
		Type:        "kv",
		Description: "Mount point for domain: " + name,
		Local:       false,
		SealWrap:    false,
		Config:      vaultapi.MountConfigInput{},
	}

	err = v.vaultClient.Sys().Mount(mountPath, mountInput)
	if smslogger.CheckError(err, "Create Domain") != nil {
		return SecretDomain{}, errors.New("Unable to create Secret Domain")
	}

	uuid, _ := uuid.GenerateUUID()
	err = v.storeUUID(uuid, name)
	if smslogger.CheckError(err, "Store UUID") != nil {
		// Mount was successful at this point.
		// Rollback the mount operation since we could not
		// store the UUID for the mount.
		v.vaultClient.Sys().Unmount(mountPath)
		return SecretDomain{}, errors.New("Unable to store Secret Domain UUID. Retry")
	}

	return SecretDomain{uuid, name}, nil
}

// CreateSecret creates a secret mounted on a particular domain name
// The secret itself is mounted on a path specified by name
func (v *Vault) CreateSecret(dom string, sec Secret) error {

	err := v.checkToken()
	if smslogger.CheckError(err, "Token Check") != nil {
		return errors.New("Token check failed")
	}

	dom = strings.TrimSpace(dom)
	dom = v.vaultMountPrefix + "/" + dom

	// Vault return is empty on successful write
	// TODO: Check if values is not empty
	_, err = v.vaultClient.Logical().Write(dom+"/"+sec.Name, sec.Values)
	if smslogger.CheckError(err, "Create Secret") != nil {
		return errors.New("Unable to create Secret at provided path")
	}

	return nil
}

// DeleteSecretDomain deletes a secret domain which translates to
// an unmount operation on the given path in Vault
func (v *Vault) DeleteSecretDomain(dom string) error {

	err := v.checkToken()
	if smslogger.CheckError(err, "Token Check") != nil {
		return errors.New("Token Check Failed")
	}

	dom = strings.TrimSpace(dom)
	mountPath := v.vaultMountPrefix + "/" + dom

	err = v.vaultClient.Sys().Unmount(mountPath)
	if smslogger.CheckError(err, "Delete Domain") != nil {
		return errors.New("Unable to delete domain specified")
	}

	return nil
}

// DeleteSecret deletes a secret mounted on the path provided
func (v *Vault) DeleteSecret(dom string, name string) error {

	err := v.checkToken()
	if smslogger.CheckError(err, "Token Check") != nil {
		return errors.New("Token check failed")
	}

	dom = strings.TrimSpace(dom)
	dom = v.vaultMountPrefix + "/" + dom

	// Vault return is empty on successful delete
	_, err = v.vaultClient.Logical().Delete(dom + "/" + name)
	if smslogger.CheckError(err, "Delete Secret") != nil {
		return errors.New("Unable to delete Secret at provided path")
	}

	return nil
}

// initRole is called only once during SMS bring up
// It initially creates a role and secret id associated with
// that role. Later restarts will use the existing role-id
// and secret-id stored on disk
func (v *Vault) initRole() error {

	if v.initRoleDone {
		return nil
	}

	// Use the root token once here
	v.vaultClient.SetToken(v.vaultToken)
	defer v.vaultClient.ClearToken()

	// Check if roleID and secretID has already been created
	rID, error := smsauth.ReadFromFile("auth/role")
	if error != nil {
		smslogger.WriteWarn("Unable to find RoleID. Generating...")
	} else {
		sID, error := smsauth.ReadFromFile("auth/secret")
		if error != nil {
			smslogger.WriteWarn("Unable to find secretID. Generating...")
		} else {
			v.roleID = rID
			v.secretID = sID
			v.initRoleDone = true
			return nil
		}
	}

	rules := `path "sms/*" { capabilities = ["create", "read", "update", "delete", "list"] }
			path "sys/mounts/sms*" { capabilities = ["update","delete","create"] }`
	err := v.vaultClient.Sys().PutPolicy(v.policyName, rules)
	if smslogger.CheckError(err, "Creating Policy") != nil {
		return errors.New("Unable to create policy for approle creation")
	}

	//Check if applrole is mounted
	authMounts, err := v.vaultClient.Sys().ListAuth()
	if smslogger.CheckError(err, "Mount Auth Backend") != nil {
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

	rName := v.vaultMountPrefix + "-role"
	data := map[string]interface{}{
		"token_ttl": "60m",
		"policies":  [2]string{"default", v.policyName},
	}

	// Create a role-id
	v.vaultClient.Logical().Write("auth/approle/role/"+rName, data)
	sec, err := v.vaultClient.Logical().Read("auth/approle/role/" + rName + "/role-id")
	if smslogger.CheckError(err, "Create RoleID") != nil {
		return errors.New("Unable to create role ID for approle")
	}
	v.roleID = sec.Data["role_id"].(string)

	// Create a secret-id to go with it
	sec, err = v.vaultClient.Logical().Write("auth/approle/role/"+rName+"/secret-id",
		map[string]interface{}{})
	if smslogger.CheckError(err, "Create SecretID") != nil {
		return errors.New("Unable to create secret ID for role")
	}

	v.secretID = sec.Data["secret_id"].(string)
	v.initRoleDone = true
	/*
	* Revoke the Root token.
	* If a new Root Token is needed, it will need to be created
	* using the unseal shards.
	 */
	err = v.vaultClient.Auth().Token().RevokeSelf(v.vaultToken)
	if smslogger.CheckError(err, "Revoke Root Token") != nil {
		smslogger.WriteWarn("Unable to Revoke Token")
	} else {
		// Revoked successfully and clear it
		v.vaultToken = ""
	}

	// Store the role-id and secret-id
	// We will need this if SMS restarts
	smsauth.WriteToFile(v.roleID, "auth/role")
	smsauth.WriteToFile(v.secretID, "auth/secret")

	return nil
}

// Function checkToken() gets called multiple times to create
// temporary tokens
func (v *Vault) checkToken() error {

	v.Lock()
	defer v.Unlock()

	// Init Role if it is not yet done
	// Role needs to be created before token can be created
	err := v.initRole()
	if err != nil {
		smslogger.WriteError(err.Error())
		return errors.New("Unable to initRole in checkToken")
	}

	// Return immediately if token still has life
	if v.vaultClient.Token() != "" &&
		time.Since(v.vaultTempTokenTTL) < time.Minute*50 {
		return nil
	}

	// Create a temporary token using our roleID and secretID
	out, err := v.vaultClient.Logical().Write("auth/approle/login",
		map[string]interface{}{"role_id": v.roleID, "secret_id": v.secretID})
	if smslogger.CheckError(err, "Create Temp Token") != nil {
		return errors.New("Unable to create Temporary Token for Role")
	}

	tok, err := out.TokenID()

	v.vaultTempTokenTTL = time.Now()
	v.vaultClient.SetToken(tok)
	return nil
}

// vaultInit() is used to initialize the vault in cases where it is not
// initialized. This happens once during intial bring up.
func (v *Vault) initializeVault() error {

	// Check for vault init status and don't exit till it is initialized
	for {
		init, err := v.vaultClient.Sys().InitStatus()
		if smslogger.CheckError(err, "Get Vault Init Status") != nil {
			smslogger.WriteInfo("Trying again in 10s...")
			time.Sleep(time.Second * 10)
			continue
		}
		// Did not get any error
		if init == true {
			smslogger.WriteInfo("Vault is already Initialized")
			return nil
		}

		// init status is false
		// break out of loop and finish initialization
		smslogger.WriteInfo("Vault is not initialized. Initializing...")
		break
	}

	// Hardcoded this to 3. We should make this configurable
	// in the future
	initReq := &vaultapi.InitRequest{
		SecretShares:    3,
		SecretThreshold: 3,
	}

	pbkey, prkey, err := smsauth.GeneratePGPKeyPair()

	if smslogger.CheckError(err, "Generating PGP Keys") != nil {
		smslogger.WriteError("Error Generating PGP Keys. Vault Init will not use encryption!")
	} else {
		initReq.PGPKeys = []string{pbkey, pbkey, pbkey}
		initReq.RootTokenPGPKey = pbkey
	}

	resp, err := v.vaultClient.Sys().Init(initReq)
	if smslogger.CheckError(err, "Initialize Vault") != nil {
		return errors.New("FATAL: Unable to initialize Vault")
	}

	if resp != nil {
		v.prkey = prkey
		v.shards = resp.KeysB64
		v.vaultToken, _ = smsauth.DecryptPGPString(resp.RootToken, prkey)
		return nil
	}

	return errors.New("FATAL: Init response was empty")
}

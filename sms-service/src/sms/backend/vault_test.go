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
	vaultapi "github.com/hashicorp/vault/api"
	credAppRole "github.com/hashicorp/vault/builtin/credential/approle"
	vaulthttp "github.com/hashicorp/vault/http"
	vaultlogical "github.com/hashicorp/vault/logical"
	vaultinmem "github.com/hashicorp/vault/physical/inmem"
	vaulttesting "github.com/hashicorp/vault/vault"
	"reflect"
	smslog "sms/log"
	"testing"
)

var secret Secret

func init() {
	smslog.Init("")
	secret = Secret{
		Name: "testsecret",
		Values: map[string]interface{}{
			"name":    "john",
			"age":     "43",
			"isadmin": "true",
		},
	}
}

// Only needed when running tests against vault
func createLocalVaultServer(t *testing.T) (*vaulttesting.TestCluster, *Vault) {
	tc := vaulttesting.NewTestCluster(t,
		&vaulttesting.CoreConfig{
			DisableCache: true,
			DisableMlock: true,
			CredentialBackends: map[string]vaultlogical.Factory{
				"approle": credAppRole.Factory,
			},
		},
		&vaulttesting.TestClusterOptions{
			HandlerFunc: vaulthttp.Handler,
			NumCores:    1,
		})

	tc.Start()

	v := &Vault{}
	v.initVaultClient()
	v.vaultToken = tc.RootToken
	v.vaultClient = tc.Cores[0].Client

	return tc, v
}

func TestInitVaultClient(t *testing.T) {

	v := &Vault{}
	v.vaultAddress = "https://localhost:8200"
	err := v.initVaultClient()
	if err != nil || v.vaultClient == nil {
		t.Fatal("Init: Init() failed to create vaultClient")
	}
}

func TestInitRole(t *testing.T) {

	tc, v := createLocalVaultServer(t)
	defer tc.Cleanup()

	v.vaultToken = tc.RootToken
	v.vaultClient = tc.Cores[0].Client

	err := v.initRole()

	if err != nil {
		t.Fatal("InitRole: InitRole() failed to create roles")
	}
}

func TestGetStatus(t *testing.T) {

	tc, v := createLocalVaultServer(t)
	defer tc.Cleanup()

	st, err := v.GetStatus()

	if err != nil {
		t.Fatal("GetStatus: Returned error")
	}

	if st == true {
		t.Fatal("GetStatus: Returned true. Expected false")
	}
}

func TestCreateSecretDomain(t *testing.T) {

	tc, v := createLocalVaultServer(t)
	defer tc.Cleanup()

	sd, err := v.CreateSecretDomain("testdomain")

	if err != nil {
		t.Fatal("CreateSecretDomain: Returned error")
	}

	if sd.Name != "testdomain" {
		t.Fatal("CreateSecretDomain: Returned name does not match: " + sd.Name)
	}

	if sd.UUID == "" {
		t.Fatal("CreateSecretDomain: Returned UUID is empty")
	}
}

func TestDeleteSecretDomain(t *testing.T) {

	tc, v := createLocalVaultServer(t)
	defer tc.Cleanup()

	sd, err := v.CreateSecretDomain("testdomain")
	if err != nil {
		t.Fatal(err)
	}

	err = v.DeleteSecretDomain(sd.UUID)
	if err != nil {
		t.Fatal("DeleteSecretDomain: Unable to delete domain")
	}
}

func TestCreateSecret(t *testing.T) {

	tc, v := createLocalVaultServer(t)
	defer tc.Cleanup()

	sd, err := v.CreateSecretDomain("testdomain")
	if err != nil {
		t.Fatal(err)
	}

	err = v.CreateSecret(sd.UUID, secret)

	if err != nil {
		t.Fatal("CreateSecret: Error Creating secret")
	}
}

func TestGetSecret(t *testing.T) {

	tc, v := createLocalVaultServer(t)
	defer tc.Cleanup()

	sd, err := v.CreateSecretDomain("testdomain")
	if err != nil {
		t.Fatal(err)
	}

	err = v.CreateSecret(sd.UUID, secret)
	if err != nil {
		t.Fatal(err)
	}

	sec, err := v.GetSecret(sd.UUID, secret.Name)
	if err != nil {
		t.Fatal("GetSecret: Error Creating secret")
	}

	if sec.Name != secret.Name {
		t.Fatal("GetSecret: Returned incorrect name")
	}

	if reflect.DeepEqual(sec.Values, secret.Values) == false {
		t.Fatal("GetSecret: Returned incorrect Values")
	}
}

func TestListSecret(t *testing.T) {

	tc, v := createLocalVaultServer(t)
	defer tc.Cleanup()

	sd, err := v.CreateSecretDomain("testdomain")
	if err != nil {
		t.Fatal(err)
	}

	err = v.CreateSecret(sd.UUID, secret)
	if err != nil {
		t.Fatal(err)
	}

	_, err = v.ListSecret(sd.UUID)
	if err != nil {
		t.Fatal("ListSecret: Returned error")
	}
}

func TestDeleteSecret(t *testing.T) {

	tc, v := createLocalVaultServer(t)
	defer tc.Cleanup()

	sd, err := v.CreateSecretDomain("testdomain")
	if err != nil {
		t.Fatal(err)
	}

	err = v.CreateSecret(sd.UUID, secret)
	if err != nil {
		t.Fatal(err)
	}

	err = v.DeleteSecret(sd.UUID, secret.Name)
	if err != nil {
		t.Fatal("DeleteSecret: Error Creating secret")
	}
}

func TestInitializeVault(t *testing.T) {

	inm, err := vaultinmem.NewInmem(nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	core, err := vaulttesting.NewCore(&vaulttesting.CoreConfig{
		DisableMlock: true,
		DisableCache: true,
		Physical:     inm,
	})
	if err != nil {
		t.Fatal(err)
	}

	ln, addr := vaulthttp.TestServer(t, core)
	defer ln.Close()

	client, err := vaultapi.NewClient(&vaultapi.Config{
		Address: addr,
	})
	if err != nil {
		t.Fatal(err)
	}

	v := &Vault{}
	v.initVaultClient()
	v.vaultClient = client

	err = v.initializeVault()
	if err != nil {
		t.Fatal("InitializeVault: Error initializing Vault")
	}
}

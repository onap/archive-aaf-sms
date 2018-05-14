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

package config

import (
	"encoding/json"
	"os"
	smslogger "sms/log"
)

// SMSConfiguration loads up all the values that are used to configure
// backend implementations
// TODO: Review these and see if they can be created/discovered dynamically
type SMSConfiguration struct {
	CAFile     string `json:"cafile"`
	ServerCert string `json:"servercert"`
	ServerKey  string `json:"serverkey"`
	Password   string `json:"password"`

	BackendAddress            string `json:"smsdbaddress"`
	VaultToken                string `json:"vaulttoken"`
	DisableTLS                bool   `json:"disable_tls"`
	BackendAddressEnvVariable string `json:"smsdburlenv"`
}

// SMSConfig is the structure that stores the configuration
var SMSConfig *SMSConfiguration

// ReadConfigFile reads the specified smsConfig file to setup some env variables
func ReadConfigFile(file string) (*SMSConfiguration, error) {
	if SMSConfig == nil {
		f, err := os.Open(file)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		// Default behaviour is to enable TLS
		SMSConfig = &SMSConfiguration{DisableTLS: false}
		decoder := json.NewDecoder(f)
		err = decoder.Decode(SMSConfig)
		if err != nil {
			return nil, err
		}

		if SMSConfig.BackendAddress == "" && SMSConfig.BackendAddressEnvVariable != "" {
			// Get the value from ENV variable
			smslogger.WriteInfo("Using Environment Variable: " + SMSConfig.BackendAddressEnvVariable)
			SMSConfig.BackendAddress = os.Getenv(SMSConfig.BackendAddressEnvVariable)
		}
	}

	return SMSConfig, nil
}

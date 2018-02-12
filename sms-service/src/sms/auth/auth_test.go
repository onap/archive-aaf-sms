/*
* ============LICENSE_START=======================================================
* ONAP : AAF/SMS
* ================================================================================
* Copyright 2018 TechMahindra
*=================================================================================
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
* ============LICENSE_END=========================================================
 */

package auth

import (
   "testing"
)

/*
*Unit test to varify GetTLSconfig func and varify the tls config min version to be 771
*Assuming cert file name as server.cert
 */
func TestGetTLSConfig(t *testing.T) {
   a := 771
   tlsconfig := GetTLSConfig("server.cert")
   if tlsconfig != nil {
      if tlsconfig.MinVersion != uint16(a) {
         t.Errorf("Test Failed")
      }
   }
}

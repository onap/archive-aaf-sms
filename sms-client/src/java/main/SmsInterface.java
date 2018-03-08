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

package org.onap.aaf.sms;

import org.onap.aaf.sms.SmsResponse;
import java.util.Map;

public interface SmsInterface {
    /*
        Inputs dname - domain name
        Output - name and uuid
        Return SmsResponse object
            success or failure
            response code if connection succeeded, otherwise -1
            response string if expected.
    */
    public SmsResponse createDomain(String dname);

    /*
        Inputs dname - domain name
        Output - none
        Return SmsResponse object
            success or failure
            response code if connection succeeded, otherwise -1
            response string if expected.

    */
    public SmsResponse deleteDomain(String dname);

    /*
        Inputs dname - domain name
        Output - list of secret names
        Return SmsResponse object
            success or failure
            response code if connection succeeded, otherwise -1
            response string if expected.

    */
    public SmsResponse getSecretNames(String dname);

    /*
        Inputs dname - domain name
               sname - secret name
               values - list of key value pairs
        Output - none
        Return SmsResponse object
            success or failure
            response code if connection succeeded, otherwise -1
            response string if expected.

    */
    public SmsResponse storeSecret(String dname, String sname, Map<String, Object> values);

    /*
        Inputs dname - domain name
               sname - secret name
        Output values - list of value pairs
        Return SmsResponse object
            success or failure
            response code if connection succeeded, otherwise -1
            response string if expected.

    */
    public SmsResponse getSecret(String dname, String sname);

    /*
        Inputs dname - domain name
               sname - secret name
        Output - none
        Return SmsResponse object
            success or failure
            response code if connection succeeded, otherwise -1
            response string if expected.
    */
    public SmsResponse deleteSecret(String dname, String sname);
}

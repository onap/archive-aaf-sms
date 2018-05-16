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

import org.junit.runner.JUnitCore;
import org.junit.runner.Result;
import org.junit.runner.notification.Failure;

public class ClientTestRunner {
    public static void main(String[] args) {
        Result r = JUnitCore.runClasses(
            SmsCreateDomainTest.class,
            SmsDeleteDomainTest.class,
            SmsStoreSecretTest.class,
            SmsGetSecretNamesTest.class,
            SmsGetSecretTest.class,
            SmsDeleteSecretTest.class
        );

        for( Failure f : r.getFailures()) {
            System.out.println(f.toString());
        }
        System.out.println(r.wasSuccessful());
    }
}

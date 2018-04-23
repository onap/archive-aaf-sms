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

import junit.framework.*;
import org.onap.aaf.sms.SmsClient;
import org.onap.aaf.sms.SmsResponse;
import org.onap.aaf.sms.SmsSecureSocket;
import javax.net.ssl.SSLSocketFactory;
import java.util.HashMap;
import java.util.Map;

public class SmsDeleteSecretTest extends TestCase {

    public void testSmsDeleteSecret() {
        try {
            SmsTest sms = new SmsTest("otconap4.sc.intel.com", 10443, null);
            SmsResponse resp = sms.deleteSecret("onap.new.test.sms0", "testsec1");
            assertTrue(resp.getSuccess());
            if ( resp.getSuccess() ) {
                assertEquals(204, resp.getResponseCode());
            } else {
                fail("Unexpected response while deleting secret");
            }
        } catch ( Exception e ) {
            fail("Exception while deleting secret");
        }
    }
}

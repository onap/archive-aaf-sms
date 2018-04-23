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

public class SmsStoreSecretTest extends TestCase {

    public void testSmsStoreSecret() {
        try {
            SmsTest sms = new SmsTest("otconap4.sc.intel.com", 10443, null);
            Map<String, Object> v = new HashMap<String, Object>();
            v.put("username", "dbuser");
            v.put("password", "kmngjhk");
            SmsResponse resp = sms.storeSecret("onap.new.test.sms0", "dbs", v);
            assertTrue(resp.getSuccess());
            if ( resp.getSuccess() ) {
                assertEquals(201, resp.getResponseCode());
            } else {
                fail("Unexpected response while storing secret");
            }
        } catch ( Exception e ) {
            fail("Exception while storing secret");
        }
    }
}

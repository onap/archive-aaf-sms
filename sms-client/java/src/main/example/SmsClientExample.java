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

import java.io.FileInputStream;
import java.lang.Boolean;
import java.lang.Integer;
import java.net.URL;
import javax.net.ssl.HttpsURLConnection;
import javax.net.ssl.KeyManagerFactory;
import javax.net.ssl.SSLContext;
import javax.net.ssl.SSLSessionContext;
import javax.net.ssl.SSLSocketFactory;
import javax.net.ssl.TrustManagerFactory;
import java.security.KeyStore;
import java.security.Provider;
import java.security.SecureRandom;
import java.security.Security;
import java.util.HashMap;
import java.util.Map;
import org.onap.aaf.sms.SmsClient;
import org.onap.aaf.sms.SmsResponse;

/*
 * Sample application demonstrating various operations related
 * Secret Management Service's APIs
 */

public class SmsClientExample {

    public static SSLSocketFactory getSSLSocketFactory(String castore) {

        try {
            // Load the CA certificate
            // There are no private keys in the truststore
            FileInputStream tst = new FileInputStream("truststoreONAP.jks");
            KeyStore trustStore = KeyStore.getInstance("JKS");
            char[] password = "password".toCharArray();
            trustStore.load(tst, password);
            TrustManagerFactory trustManagerFactory =
                TrustManagerFactory.getInstance(TrustManagerFactory.getDefaultAlgorithm());
            trustManagerFactory.init(trustStore);

            //Create the context
            SSLContext context = SSLContext.getInstance("TLSv1.2");
            context.init(null, trustManagerFactory.getTrustManagers(), new SecureRandom());
            //Create a socket factory
            SSLSocketFactory ssf = context.getSocketFactory();
            return ssf;
        } catch (Exception e) {
            e.printStackTrace();
            return null;
        }

    }

    public static void main(String[] args) throws Exception {

        SSLSocketFactory ssf = SmsClientExample.getSSLSocketFactory("truststoreONAP.jks");

        // Create the SMSClient
        SmsClient sms = new SmsClient("aaf-sms.onap", 30243, ssf);

        // Create a test domain
        System.out.println("CREATE DOMAIN: ");
        SmsResponse resp = sms.createDomain("sms_test_domain");
        if ( resp.getSuccess() ) {
            System.out.println("-- Return Code: " + resp.getResponseCode());
            System.out.println("-- Return Data: " + resp.getResponse());
            System.out.println("");
        } else {
            System.out.println("-- Error String: " + resp.getErrorMessage());
            System.out.println("");
        }

        // Create secret data here
        Map<String, Object> data_1 = new HashMap<String, Object>();
        data_1.put("passwd", "gax6ChD0yft");

        // Store them in previously created domain
        System.out.println("STORE SECRET: " + "test_secret");
        resp = sms.storeSecret("sms_test_domain", "test_secret",  data_1);
        if ( resp.getSuccess() ) {
            System.out.println("-- Return Code: " + resp.getResponseCode());
            System.out.println("");
        }

        // A more complex data example on the same domain
        Map<String, Object> data_2 = new HashMap<String, Object>();
        data_2.put("username", "dbuser");
        data_2.put("isadmin", new Boolean(true));
        data_2.put("age", new Integer(40));
        data_2.put("secretkey", "asjdhkuhioeukadfjsadnfkjhsdukfhaskdjhfasdf");
        data_2.put("token", "2139084553458973452349230849234234908234342");

        // Store the secret
        System.out.println("STORE SECRET: " + "test_credentials");
        resp = sms.storeSecret("sms_test_domain", "test_credentials", data_2);
        if ( resp.getSuccess() ) {
            System.out.println("-- Return Code: " + resp.getResponseCode());
            System.out.println("");
        }

        // List all secret names stored in domain
        System.out.println("LIST SECRETS: ");
        resp = sms.getSecretNames("sms_test_domain");
        if ( resp.getSuccess() ) {
            System.out.println("-- Return Code: " + resp.getResponseCode());
            System.out.println("-- Return Data: " + resp.getResponse());
            System.out.println("");
        }

        // Retrieve a secret from stored domain
        System.out.println("GET SECRET: " + "test_secret");
        resp= sms.getSecret("sms_test_domain", "test_secret");
        if ( resp.getSuccess() ) {
            System.out.println("-- Return Code: " + resp.getResponseCode());
            System.out.println("-- Return Data: " + resp.getResponse());
            System.out.println("");
        }

        // Retrieve the second secret from stored domain
        // getResponse() on the return value retrieves the
        // map containing the key, values for the secret
        System.out.println("GET SECRET: " + "test_credentials");
        resp= sms.getSecret("sms_test_domain", "test_credentials");
        if ( resp.getSuccess() ) {
            System.out.println("-- Return Code: " + resp.getResponseCode());
            System.out.println("-- Return Data: " + resp.getResponse());

            //conditional processing of returned data
            Boolean b = (Boolean)resp.getResponse().get("isadmin");
            System.out.println("-- isadmin: " + b);
            if ( b )
                System.out.println("-- age: " + (Integer)resp.getResponse().get("age"));
            System.out.println("");
        }

        // Delete the secret
        System.out.println("DELETE SECRET: " + "test_credentials");
        resp=sms.deleteSecret("sms_test_domain", "test_credentials");
        if ( resp.getSuccess() ) {
            System.out.println("-- Return Code: " + resp.getResponseCode());
            System.out.println("");
        }

        // Delete the domain
        System.out.println("DELETE DOMAIN: " + "sms_test_domain");
        resp=sms.deleteDomain("sms_test_domain");
        if ( resp.getSuccess() ) {
            System.out.println("-- Return Code: " + resp.getResponseCode());
            System.out.println("");
        }
    }
}

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

import java.security.Provider;
import java.security.Security;
import java.security.KeyStore;
import java.security.SecureRandom;
import javax.net.ssl.KeyManagerFactory;
import javax.net.ssl.TrustManagerFactory;
import javax.net.ssl.SSLContext;
import javax.net.ssl.SSLSocketFactory;
import javax.net.ssl.SSLSessionContext;
import java.io.FileInputStream;
import java.net.URL;
import javax.net.ssl.HttpsURLConnection;
import org.onap.aaf.sms.SmsClient;
import org.onap.aaf.sms.SmsResponse;

public class SmsClientExample {
    public static void main(String[] args) throws Exception {
        // Set up the Sun PKCS 11 provider
        Provider p = Security.getProvider("SunPKCS11-pkcs11Test");
        if (p==null) {
            throw new RuntimeException("could not get security provider");
        }

        // Load the key store
        char[] pin = "45789654".toCharArray();
        KeyStore keyStore = KeyStore.getInstance("PKCS11", p);
        keyStore.load(null, pin);

        // Load the CA certificate
        FileInputStream tst = new FileInputStream("/ca.jks");
        KeyStore trustStore = KeyStore.getInstance("JKS");
        trustStore.load(tst, pin);

        KeyManagerFactory keyManagerFactory =
             KeyManagerFactory.getInstance(KeyManagerFactory.getDefaultAlgorithm());
        //Add to keystore to key manager
        keyManagerFactory.init(keyStore, pin);

        TrustManagerFactory trustManagerFactory =
             TrustManagerFactory.getInstance(TrustManagerFactory.getDefaultAlgorithm());
        trustManagerFactory.init(trustStore);

        //Create the context
        SSLContext context = SSLContext.getInstance("TLS");
        context.init(keyManagerFactory.getKeyManagers(),
             trustManagerFactory.getTrustManagers(), new SecureRandom());
        //Create a socket factory
        SSLSocketFactory ssf = context.getSocketFactory();
        SSLSessionContext sessCtx = context.getServerSessionContext();
        SmsClient sms = new SmsClient("onap.mydomain.com", 10443, ssf);
        SmsResponse resp1 = sms.createDomain("onap.new.test.sms0");
        if ( resp1.getSuccess() ) {
            System.out.println(resp1.getResponse());
            System.out.println(resp1.getResponseCode());
        }
        SmsResponse resp2 = sms.storeSecreat("onap.new.test.sms0",
            "testsec", "{\"passwd\": \"gax6ChD0yft\"}");
        if ( resp2.getSuccess() ) {
            System.out.println(resp2.getResponse());
            System.out.println(resp2.getResponseCode());
        }
        SmsResponse resp3 = sms.storeSecreat("onap.new.test.sms0",
            "user", "{\"username\": \"dbuser\"}");
        if ( resp3.getSuccess() ) {
            System.out.println(resp3.getResponse());
            System.out.println(resp3.getResponseCode());
        }
        SmsResponse resp4 = sms.getSecreatNames("onap.new.test.sms0");
        if ( resp4.getSuccess() ) {
            System.out.println(resp4.getResponse());
            System.out.println(resp4.getResponseCode());
        }
        SmsResponse resp5= sms.retrieveSecreat("onap.new.test.sms0", "testsec");
        if ( resp5.getSuccess() ) {
            System.out.println(resp5.getResponse());
            System.out.println(resp5.getResponseCode());
        }
        SmsResponse resp6= sms.retrieveSecreat("onap.new.test.sms0", "user");
        if ( resp6.getSuccess() ) {
            System.out.println(resp6.getResponse());
            System.out.println(resp6.getResponseCode());
        }
        SmsResponse resp7=sms.deleteDomain("onap.new.test.sms0");
        if ( resp7.getSuccess() ) {
            System.out.println(resp7.getResponse());
            System.out.println(resp7.getResponseCode());
        }
    }
}

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

package org.onap.aaf.sms.test;

import java.io.FileInputStream;
import javax.net.ssl.KeyManagerFactory;
import javax.net.ssl.SSLContext;
import javax.net.ssl.SSLSessionContext;
import javax.net.ssl.SSLSocketFactory;
import javax.net.ssl.TrustManagerFactory;
import java.security.KeyStore;
import java.security.Provider;
import java.security.SecureRandom;
import java.security.Security;

public class SmsSecureSocket {
    private SSLSocketFactory ssf = null;
    public SmsSecureSocket() throws Exception {
        // Set up the Sun PKCS 11 provider
        Provider p = Security.getProvider("SunPKCS11-pkcs11Test");
        if (p==null) {
            throw new RuntimeException("could not get security provider");
        }

        // Load the key store
        char[] pin = "123456789".toCharArray();
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
    }
    public SSLSocketFactory getSSF() {
        return(ssf);
    }
}

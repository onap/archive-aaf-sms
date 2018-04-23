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

import javax.net.ssl.SSLSocketFactory;
import java.net.URL;
import javax.net.ssl.HttpsURLConnection;
import org.onap.aaf.sms.SmsResponse;
import org.onap.aaf.sms.SmsClient;
import java.io.InputStream;
import java.io.OutputStream;
import java.io.InputStreamReader;
import java.io.BufferedReader;
import java.io.OutputStreamWriter;
import java.util.Map;
import java.util.HashMap;
import java.util.Iterator;
import java.util.List;
import java.util.ArrayList;
import org.json.JSONArray;
import org.json.JSONException;
import org.json.JSONObject;

public class SmsTest extends SmsClient {

    public SmsTest(String host, int port, SSLSocketFactory s) {
        super(host, port, s);
    }
    public SmsTest(String host, int port, String version, SSLSocketFactory s) {
        super(host, port, version, s);
    }
    public  SmsResponse execute(String reqtype, String t, String ins, boolean input, boolean output) {
        Map<String, Object> m;
        SmsResponse resp = new SmsResponse();
        System.out.println(t);

        switch ( reqtype ) {
            case "POST":
                if ( t.matches("(.*)/v1/sms/domain")) {
                    resp.setSuccess(true);
                    resp.setResponseCode(201);
                    try {
                        m = strtomap(ins);
                    } catch ( Exception e ) {
                        resp.setResponse(null);
                        return(resp);
                    }
                    resp.setResponse(m);
                } else {
                    if ( t.matches("(.*)/v1/sms/(.*)/secret") ) {
                        resp.setSuccess(true);
                        resp.setResponseCode(201);
                    }
                }
            break;
            case "GET":
                if ( t.matches("(.*)/v1/sms/(.*)/secret") ) {
                    resp.setSuccess(true);
                    resp.setResponseCode(200);
                    String jstr = "{\"secretnames\":[\"testsec1\",\"newtest\"]}";
                    try {
                        m = strtomap(jstr);
                    } catch ( Exception e ) {
                        resp.setResponse(null);
                        return(resp);
                    }
                    resp.setResponse(m);
                } else {
                    if ( t.matches("(.*)/v1/sms/(.*)/secret/testsec1")) {
                        resp.setSuccess(true);
                        resp.setResponseCode(200);
                        String js = "{\"name\":\"testsec1\",\"values\":{\"username\":\"dbuser\",\"passwd\":\"jdX784i-5k\"}}";
                        try {
                            m = strtomap(js);
                            Map<String, Object> sm = (Map<String, Object>)m.get("values");
                        } catch ( Exception e ) {
                            resp.setResponse(null);
                            return(resp);
                        }
                        resp.setResponse(m);
                    }
                }
            break;
            case "DELETE":
                if ( t.matches("(.*)/v1/sms/domain/(.*)") ) {
                    // for both delete domain & secret case
                    resp.setSuccess(true);
                    resp.setResponseCode(204);
                    resp.setResponse(null);
                } 
            break;

        }
        return resp;
    }
}

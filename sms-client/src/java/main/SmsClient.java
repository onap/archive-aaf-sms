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

public class SmsClient implements SmsInterface {

    private String baset;
    private SSLSocketFactory ssf;

    public SmsClient(String host, int port, SSLSocketFactory s) {
        baset = "https://"+ host + ":" + port + "/v1/sms";
        ssf = s;
    }
    public SmsClient(String host, int port, String version, SSLSocketFactory s) {
        baset = "https://"+ host + ":" + port + "/" + version + "/sms";
        ssf = s;
    }

    private  Map<String, Object> getSubmap(Map<String, Object> raw, String k) {
        Object v = raw.get(k);
        if ( v != null ) {
            Map<String, Object> r = (Map<String, Object>)v;
            return(r);
        }
        else {
            return(null);
        }
    }

    private List<Object> jsontolist(JSONArray a) throws JSONException {
        List<Object> l = new ArrayList<Object>();
        for(int i=0;i<a.length();i++) {
            Object v = a.get(i);
            if ( v instanceof JSONArray ) {
                v = jsontolist((JSONArray) v);
            } else if (v instanceof JSONObject) {
                v = jsontomap((JSONObject) v);
            }
            l.add(v);
        }
        return(l);
    }

    private Map<String, Object> jsontomap(JSONObject j) throws JSONException {
        Map<String, Object> m = new HashMap<String, Object>();

        Iterator<?> ks = j.keys();
        while( ks.hasNext() ) {
            String k = (String)ks.next();
            Object v = j.get(k);

            if ( v instanceof JSONArray ) {
                v = jsontolist((JSONArray) v);
            } else if ( v instanceof JSONObject ) {
                v = jsontomap((JSONObject) v);
            }
            m.put(k, v);
        }
        return(m);
    }

    private Map<String, Object> strtomap(String r) throws JSONException {
        JSONObject jobj = null;

        jobj = new JSONObject(r);
        return(jsontomap(jobj));

    }
    private SmsResponse execute(String reqtype, String t, String ins, boolean input, boolean output) {

        HttpsURLConnection conn;
        int errorcode = -1;
        SmsResponse resp = new SmsResponse();

        try {
            URL url = new URL(t);
            conn = (HttpsURLConnection)url.openConnection();
            conn.setSSLSocketFactory(ssf);
            conn.setRequestMethod(reqtype);
            conn.setDoOutput(true);
            conn.setDoInput(true);
            conn.setRequestProperty("Content-Type", "application/json");
            conn.setRequestProperty("Accept", "application/json");

            if ( input ) {
                OutputStream out = conn.getOutputStream();
                OutputStreamWriter wr = new OutputStreamWriter(out);
                wr.write(ins);
                wr.flush();
                wr.close();
            }
            errorcode = conn.getResponseCode();
            if ( output && errorcode > 0 ) {
                InputStream inputstream = conn.getInputStream();
                InputStreamReader inputstreamreader = new InputStreamReader(inputstream);
                BufferedReader bufferedreader = new BufferedReader(inputstreamreader);

                String response;
                String save = "";
                while ((response = bufferedreader.readLine()) != null) {
                    save = save + response;
                }
                if ( !save.isEmpty() ) {
                    if ( errorcode/100 == 2 ) {
                        resp.setResponse(strtomap(save));
                    } else {
                        resp.setErrorMessage(save);
                    }
                }
            }
        } catch ( Exception e ) {
            e.printStackTrace();
            resp.setResponseCode(errorcode);
            return(resp);
        }
        resp.setResponseCode(errorcode);
        return resp;
    }
    @Override
    public SmsResponse createDomain(String dname) {

        String t = baset + "/domain";
        String input = "{\"name\":\"" + dname + "\"}";

        SmsResponse resp = execute("POST", t, input, true, true);
        int errcode = resp.getResponseCode();

        if ( errcode > 0 && errcode/100 == 2 )
            resp.setSuccess(true);
        else
            resp.setSuccess(false);

        return(resp);
    }
    @Override
    public SmsResponse deleteDomain(String dname) {

        String t = baset + "/domain/" + dname;

        SmsResponse resp = execute("DELETE", t, null, false, true);
        int errcode = resp.getResponseCode();

        if ( errcode > 0 && errcode/100 == 2 )
            resp.setSuccess(true);
        else
            resp.setSuccess(false);

        return(resp);
    }
    @Override
    public SmsResponse storeSecret(String dname, String sname, Map<String, Object> values) {

        String t = baset + "/domain/" + dname + "/secret";
        Map<String, Object> cm = new HashMap<String, Object>();
        cm.put("name", sname);
        cm.put("values", values);
        JSONObject jobj = new JSONObject(cm);

        SmsResponse resp = execute("POST", t, jobj.toString(), true, false);
        int errcode = resp.getResponseCode();

        if ( errcode > 0 && errcode/100 == 2 )
            resp.setSuccess(true);
        else
            resp.setSuccess(false);

        return(resp);
    }
    @Override
    public SmsResponse getSecretNames(String dname) {

        String t = baset + "/domain/" + dname + "/secret";

        SmsResponse resp = execute("GET", t, null, false, true);
        int errcode = resp.getResponseCode();

        if ( errcode > 0 && errcode/100 == 2 )
            resp.setSuccess(true);
        else
            resp.setSuccess(false);

        return(resp);
    }
    @Override
    public SmsResponse getSecret(String dname, String sname) {

        String t = baset + "/domain/" + dname + "/secret/" + sname;

        SmsResponse resp = execute("GET", t, null, false, true);
        int errcode = resp.getResponseCode();

        if ( errcode > 0 && errcode/100 == 2 ) {
            resp.setSuccess(true);
            resp.setResponse(getSubmap(resp.getResponse(), "values"));
        }
        else {
            resp.setSuccess(false);
        }

        return(resp);

    }
    @Override
    public SmsResponse deleteSecret(String dname, String sname) {

        String t = baset + "/domain/" + dname + "/secret/" + sname;

        SmsResponse resp = execute("DELETE", t, null, false, true);
        int errcode = resp.getResponseCode();

        if ( errcode > 0 && errcode/100 == 2 )
            resp.setSuccess(true);
        else
            resp.setSuccess(false);

        return(resp);
    }
}

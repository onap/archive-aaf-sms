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

public class SmsResponse {
    private boolean success;
    private int responseCode;
    private String response;

    public SmsResponse() {
        success = false;
        responseCode = -1;
        response = "";
    }
    public void setResponseCode(int code) {
        responseCode = code;
    }
    public void setResponse(String res) {
        response = res;
    }
    public void setSuccess(boolean val) {
        success = val;
    }
    public int getResponseCode() {
        return responseCode;
    }
    public String getResponse() {
        return response;
    }
    public boolean getSuccess() {
        return success;
    }
}

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

import java.util.Map;

public class SmsResponse {
    private boolean success;
    private int responseCode;
    private String errorMessage;
    private Map<String, Object> response;

    public SmsResponse() {
        success = false;
        responseCode = -1;
        errorMessage = "";
        response = null;
    }
    public void setResponseCode(int code) {
        responseCode = code;
    }
    public void setResponse(Map<String, Object> res) {
        response = res;
    }
    public void setSuccess(boolean val) {
        success = val;
    }
    public int getResponseCode() {
        return responseCode;
    }
    public void setErrorMessage(String em) {
        errorMessage = em;
    }
    public String getErrorMessage() {
        return errorMessage;
    }
    public Map<String, Object> getResponse() {
        return response;
    }
    public boolean getSuccess() {
        return success;
    }
}

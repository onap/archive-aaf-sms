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

package log

import (
	"log"
	"os"
)

var errLogger *log.Logger
var warnLogger *log.Logger
var infoLogger *log.Logger

// Init will be called by sms.go before any other packages use it
func Init(filePath string) {
	if filePath == "" {
		errLogger = log.New(os.Stderr, "ERROR: ", log.Lshortfile|log.LstdFlags)
		warnLogger = log.New(os.Stdout, "WARNING: ", log.Lshortfile|log.LstdFlags)
		infoLogger = log.New(os.Stdout, "INFO: ", log.Lshortfile|log.LstdFlags)
		return
	}

	f, err := os.Create(filePath)
	if err != nil {
		log.Println("Unable to create a log file")
		log.Println(err)
		errLogger = log.New(os.Stderr, "ERROR: ", log.Lshortfile|log.LstdFlags)
		warnLogger = log.New(os.Stdout, "WARNING: ", log.Lshortfile|log.LstdFlags)
		infoLogger = log.New(os.Stdout, "INFO: ", log.Lshortfile|log.LstdFlags)
	} else {
		errLogger = log.New(f, "ERROR: ", log.Lshortfile|log.LstdFlags)
		warnLogger = log.New(f, "WARNING: ", log.Lshortfile|log.LstdFlags)
		infoLogger = log.New(f, "INFO: ", log.Lshortfile|log.LstdFlags)
	}
}

// WriteError writes output to the writer we have
// defined durint its creation with ERROR prefix
func WriteError(msg string) {
	if errLogger != nil {
		errLogger.Println(msg)
	}
}

// WriteWarn writes output to the writer we have
// defined durint its creation with WARNING prefix
func WriteWarn(msg string) {
	if warnLogger != nil {
		warnLogger.Println(msg)
	}
}

// WriteInfo writes output to the writer we have
// defined durint its creation with INFO prefix
func WriteInfo(msg string) {
	if infoLogger != nil {
		infoLogger.Println(msg)
	}
}

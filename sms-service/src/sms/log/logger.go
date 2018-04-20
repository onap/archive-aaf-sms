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
	"fmt"
	"log"
	"os"
)

var errL, warnL, infoL *log.Logger
var stdErr, stdWarn, stdInfo *log.Logger

// Init will be called by sms.go before any other packages use it
func Init(filePath string) {

	stdErr = log.New(os.Stderr, "ERROR: ", log.Lshortfile|log.LstdFlags)
	stdWarn = log.New(os.Stdout, "WARNING: ", log.Lshortfile|log.LstdFlags)
	stdInfo = log.New(os.Stdout, "INFO: ", log.Lshortfile|log.LstdFlags)

	if filePath == "" {
		// We will just to std streams
		return
	}

	f, err := os.Create(filePath)
	if err != nil {
		stdErr.Println("Unable to create log file: " + err.Error())
		return
	}

	errL = log.New(f, "ERROR: ", log.Lshortfile|log.LstdFlags)
	warnL = log.New(f, "WARNING: ", log.Lshortfile|log.LstdFlags)
	infoL = log.New(f, "INFO: ", log.Lshortfile|log.LstdFlags)
}

// WriteError writes output to the writer we have
// defined during its creation with ERROR prefix
func WriteError(msg string) {
	if errL != nil {
		errL.Output(2, fmt.Sprintln(msg))
	}
	if stdErr != nil {
		stdErr.Output(2, fmt.Sprintln(msg))
	}
}

// WriteWarn writes output to the writer we have
// defined during its creation with WARNING prefix
func WriteWarn(msg string) {
	if warnL != nil {
		warnL.Output(2, fmt.Sprintln(msg))
	}
	if stdWarn != nil {
		stdWarn.Output(2, fmt.Sprintln(msg))
	}
}

// WriteInfo writes output to the writer we have
// defined during its creation with INFO prefix
func WriteInfo(msg string) {
	if infoL != nil {
		infoL.Output(2, fmt.Sprintln(msg))
	}
	if stdInfo != nil {
		stdInfo.Output(2, fmt.Sprintln(msg))
	}
}

//CheckError is a helper function to reduce
//repitition of error checkign blocks of code
func CheckError(err error, topic string) error {
	if err != nil {
		msg := topic + ": " + err.Error()
		if errL != nil {
			errL.Output(2, fmt.Sprintln(msg))
		}
		if stdErr != nil {
			stdErr.Output(2, fmt.Sprintln(msg))
		}
		return err
	}
	return nil
}

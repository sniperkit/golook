//Copyright 2016-2017 Beate Ottenwälder
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
package api

import (
	"errors"
	"fmt"
	"net/http"

	. "github.com/ottenwbe/golook/broker/models"

	"encoding/json"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

/*
	Common (helper) functions and constants, required by all controllers.
*/

const (
	NACK = "{NACK}"
	ACK  = "{ACK}"

	systemPath = "system"
	FILE_PATH  = "file"
)

func isValidRequest(request *http.Request) bool {
	return (request != nil) && (request.Body != nil)
}

func returnAck(writer http.ResponseWriter) (int, error) {
	return fmt.Fprint(writer, ACK)
}

func returnNackAndLog(writer http.ResponseWriter, errorString string, status int) {
	log.Error(errorString)
	http.Error(writer, errors.New(NACK).Error(), status)
}

func returnNackAndLogError(writer http.ResponseWriter, errorString string, err error, status int) {
	log.WithError(err).Print(errorString)
	http.Error(writer, errors.New(NACK).Error(), status)
}

//func extractSystemFromPath(request *http.Request) string {
//	params := mux.Vars(request)
//	system := params[systemPath]
//	return system
//}

func extractFileFromPath(request *http.Request) string {
	params := mux.Vars(request)
	fileName := params[FILE_PATH]
	return fileName
}

func extractReport(request *http.Request) (*FileReport, error) {
	if !isValidRequest(request) {
		return nil, errors.New("No valid request")
	}

	var fileReport *FileReport
	err := json.NewDecoder(request.Body).Decode(fileReport)
	if err != nil {
		return nil, err
	}
	return fileReport, nil
}

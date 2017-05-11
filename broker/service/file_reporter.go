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

package service

import (
	"github.com/ottenwbe/golook/broker/models"
	"github.com/ottenwbe/golook/broker/repository"
	golook "github.com/ottenwbe/golook/broker/runtime/core"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

func broadcastLocalFiles(broadCastRouter *router) {
	files := repositories.GoLookRepository.GetFiles(golook.GolookSystem.UUID)
	broadcastFiles(files, broadCastRouter)
}

func broadcastFiles(files map[string]*models.File, broadCastRouter *router) {
	peerFileReport := &peerFileReport{Files: files, System: golook.GolookSystem.UUID}
	broadCastRouter.BroadCast(fileReport, peerFileReport)
}

func reportFileChanges(filePath string, broadCastRouter *router) {
	files := localFileReport(filePath, false)
	broadcastFiles(files, broadCastRouter)
}

/*
reportFileChangesLocal is a wrapper around localFileReport
*/
func reportFileChangesLocal(filePath string) {
	localFileReport(filePath, false)
}

func localFileReport(filePath string, _ bool) map[string]*models.File {

	var (
		files = map[string]*models.File{}
		err   error
	)

	file, err := models.NewFile(filePath)
	if err != nil {
		log.WithError(err).Error("Ignoring file report.")
		return files
	}

	if file.Directory {
		files, err = filesInFolder(file.Name)
		if err != nil {
			log.WithError(err).Error("Ignoring file report.")
			return files
		}
	}

	files[file.Name] = file

	repositories.GoLookRepository.UpdateFiles(golook.GolookSystem.UUID, files)

	return files
}

//filesInFolder generates a map with all files in the folder
func filesInFolder(folderPath string) (map[string]*models.File, error) {

	var (
		files  []os.FileInfo
		report = map[string]*models.File{}
		err    error
	)

	files, err = ioutil.ReadDir(folderPath)
	if err != nil {
		// return when errors like missing folder permissions disallow file report
		return report, err
	}

	for idx := range files {
		report = appendFile(files[idx], report)
	}
	return report, err
}

func appendFile(fileToAppend os.FileInfo, report map[string]*models.File) map[string]*models.File {
	if file, err := models.NewFile(fileToAppend.Name()); err == nil && !fileToAppend.IsDir() {
		report[file.Name] = file
	}
	return report
}
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
	"path/filepath"
)

func broadcastLocalFiles(broadCastRouter *router) {
	files := repositories.GoLookRepository.GetFiles(golook.GolookSystem.UUID)
	broadcastFiles(files, broadCastRouter)
}

func broadcastFiles(files map[string]map[string]*models.File, broadCastRouter *router) {
	peerFileReport := &peerFileReport{Files: files, SystemUUID: golook.GolookSystem.UUID}
	broadCastRouter.BroadCast(fileReport, peerFileReport)
}

func reportFileChanges(filePath string, broadCastRouter *router) {
	files := localFileReport(filePath)
	broadcastFiles(files, broadCastRouter)
}

/*
reportFileChangesLocal is a wrapper around localFileReport
*/
func reportFileChangesLocal(filePath string) {
	localFileReport(filePath)
}

func foreignFileReport(systemUUID string, files map[string]map[string]*models.File) {
	repositories.GoLookRepository.UpdateFiles(systemUUID, files, false)
}

func deleteLocalFiles(filePath string) map[string]map[string]*models.File {

	var (
		files = map[string]map[string]*models.File{}
		err   error
	)

	file, err := models.NewFile(filePath)
	if err != nil {
		log.WithError(err).Error("Ignoring file report.")
		return files
	}

	//Mark file as removed
	file.Meta.State = models.Removed

	files[filepath.Dir(filePath)] = map[string]*models.File{file.Name: file}

	repositories.GoLookRepository.UpdateFiles(golook.GolookSystem.UUID, files, true)

	return files
}

func localFileReport(filePath string) map[string]map[string]*models.File {

	var (
		files = map[string]map[string]*models.File{}
		err   error
	)

	file, err := models.NewFile(filePath)
	if err != nil {
		log.WithError(err).Error("Ignoring file report.")
		return files
	}

	if file.Directory {
		filesInFolder, err := filesInFolder(file.Name)
		if err != nil {
			log.WithError(err).Error("Ignoring file report.")
			return files
		}
		repositories.GoLookRepository.UpdateFiles(golook.GolookSystem.UUID, filesInFolder, false)

		// The folder is treated as special file for now. Therefore we do not break here.
		// This allows us to use the same mechanisms as for a file when handling it at the peers.
		// files[file.Name] = map[string]*models.File{file.Name: file}
	}

	files[filepath.Dir(filePath)] = map[string]*models.File{file.Name: file}

	repositories.GoLookRepository.UpdateFiles(golook.GolookSystem.UUID, files, true)

	return files
}

//filesInFolder generates a map with all files in the folder
func filesInFolder(folderPath string) (map[string]map[string]*models.File, error) {

	var (
		files  []os.FileInfo
		report = map[string]map[string]*models.File{folderPath: {}}
		err    error
	)

	files, err = ioutil.ReadDir(folderPath)
	if err != nil && !os.IsNotExist(err) {
		// return when errors like missing folder permissions disallow file report
		return report, err
	}

	for _, file := range files {
		if !file.IsDir() {
			report[folderPath] = appendFile(file, report[folderPath])
		}
	}
	return report, err
}

func appendFile(fileToAppend os.FileInfo, report map[string]*models.File) map[string]*models.File {
	if file, err := models.NewFile(fileToAppend.Name()); err == nil && !fileToAppend.IsDir() {
		report[file.Name] = file
	}
	return report
}

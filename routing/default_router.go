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
package routing

import (
	. "github.com/ottenwbe/golook/app"
	. "github.com/ottenwbe/golook/file_management"
	. "github.com/ottenwbe/golook/repository"
	"io/ioutil"
	"os"
)

type DefaultRouter struct {
}

func (DefaultRouter) handleQueryAllSystemsForFile(fileName string) (systems map[string]*System, err error) {
	if canAnswerQuery() {
		systems = GoLookRepository.FindSystemAndFiles(fileName)
	} else {
		systems, err = GolookClient.DoQuerySystemsAndFiles(fileName)
	}
	return
}

func (DefaultRouter) handleQueryFiles(systemName string) (files map[string]File, err error) {
	if canAnswerQuery() {
		files, _ = GoLookRepository.GetFilesOfSystem(systemName)
	} else {
		files, err = GolookClient.DoGetFiles(systemName)
	}
	return
}

func (DefaultRouter) handleReportFile(filePath string) error {
	if f, err := NewFile(filePath); err != nil {
		return err
	} else /* report file */ {
		GolookClient.DoPostFile(f)
	}
	return nil
}

// Report individual files
func (DefaultRouter) handleReportFileR(filePath string) error {
	if f, err := NewFile(filePath); err != nil {
		return err
	} else /* report file */ {
		GolookClient.DoPutFiles([]File{*f})
	}
	return nil
}

func (DefaultRouter) handleReportFolderR(folderPath string) error {
	report, err := generateReport(folderPath)
	GolookClient.DoPutFiles(report)
	return err
}

// Report files in a folder and replace all previously reported files
func (DefaultRouter) handleReportFolder(folderPath string) error {
	report, err := generateReport(folderPath)
	GolookClient.DoPostFiles(report)
	return err
}

// Report files in a folder
func generateReport(folderPath string) ([]File, error) {

	var (
		files     []os.FileInfo
		report    []File = make([]File, 0)
		returnErr error  = nil
	)

	files, returnErr = ioutil.ReadDir(folderPath)
	if returnErr != nil {
		return report, returnErr
	}

	for idx := range files {
		report, returnErr = appendFile(files[idx], report)
	}
	return report, returnErr
}

func appendFile(fileToAppend os.FileInfo, appendReport []File) (report []File, err error) {
	var file *File = nil
	if file, err = NewFile(fileToAppend.Name()); err == nil && !fileToAppend.IsDir() {
		report = append(appendReport, *file)
	}
	return
}

func canAnswerQuery() bool {
	return GolookClient == nil && GoLookRepository != nil
}

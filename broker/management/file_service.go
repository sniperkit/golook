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
package management

import (
	. "github.com/ottenwbe/golook/broker/models"

	log "github.com/sirupsen/logrus"
)

type ReportService interface {
	MakeFileReport(fileReport *FileReport)
	MakeFolderReport(folderReport *FileReport)
}

func NewReportService() ReportService {
	return &reportService{}
}

type reportService struct{}

func (*reportService) MakeFileReport(fileReport *FileReport) {

	if fileReport == nil {
		log.Error("Ignoring empty file api.")
		return
	}

	routeFile(fileReport.Path, fileReport.Replace)

	if fileReport.Monitor {
		AddFileMonitor(fileReport.Path)
	}
}

func (*reportService) MakeFolderReport(folderReport *FileReport) {

	if folderReport == nil {
		log.Error("Ignoring empty folder api.")
		return
	}

	routeFolder(folderReport.Path, folderReport.Replace)

	if folderReport.Monitor {
		AddFileMonitor(folderReport.Path)
	}
}

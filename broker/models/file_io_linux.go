// +build linux

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

package models

import (
	"github.com/fsnotify/fsnotify"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

func NewFile(filePath string) (f *File, err error) {

	//TODO refactor

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		f = &File{
			Name:      filePath,
			ShortName: filePath,
			Created:   time.Unix(0, 0),
			Modified:  time.Unix(0, 0),
			Accessed:  time.Unix(0, 0),
			Meta:      FileMeta{fsnotify.Remove},
		}
		return f, nil
	}

	var fi os.FileInfo
	var fileName string
	fi, err = os.Stat(filePath)
	if err != nil {
		return
	}
	var stat = fi.Sys().(*syscall.Stat_t)
	fileName, err = filepath.Abs(filePath)
	if err != nil {
		return
	}

	f = &File{
		Name:      fileName,
		ShortName: filepath.Base(filePath),
		Created:   time.Unix(stat.Ctim.Sec, stat.Ctim.Nsec),
		Modified:  time.Unix(stat.Mtim.Sec, stat.Mtim.Nsec),
		Accessed:  time.Unix(stat.Atim.Sec, stat.Atim.Nsec),
		Meta:      FileMeta{fsnotify.Create},
	}
	return
}
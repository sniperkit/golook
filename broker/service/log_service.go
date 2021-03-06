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
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"os"
)

/*
LogService interface for all log implementations. It basically is a facade that allows us to transparently change the log service.
*/
type LogService interface {
	RewriteLog(writer io.Writer) error
	//Init the configuration
	Init()
	//Apply the configuration
	Apply()
}

type viperLogService struct {
	logFileName string
}

const defaultLogFileName = "golook.log"

var logger = &viperLogService{}

/*
GetLogService returns the root log service.
*/
func GetLogService() LogService {
	return logger
}

func (logs *viperLogService) RewriteLog(writer io.Writer) error {

	f, err := os.Open(logs.logFileName)
	if os.IsNotExist(err) {
		// when no log file exists, write nothing to the io.Writer and return
		return nil
	} else if err != nil {
		return err
	}
	defer f.Close()

	buf := make([]byte, 32*1024) // default buffer size.

	for {
		n, err := f.Read(buf)

		if n > 0 {
			fmt.Fprint(writer, string(buf[:n]))
		}

		if err == io.EOF {
			break
		} else if err != nil {
			return errors.New("An unexpected error occurred while reading the log file")
		}
	}

	return nil
}

/*
Apply the log configuration
*/
func (logs *viperLogService) Apply() {
	lvl, err := log.ParseLevel(viper.GetString("log.level"))
	if err != nil {
		lvl = log.InfoLevel
		log.WithError(err).Infof("Failed to parse configured log level; falling back to default: '%s'", lvl.String())
	}
	log.SetLevel(lvl)

	logs.logFileName = viper.GetString("log.file")
	if logs.logFileName == "" {
		logs.logFileName = defaultLogFileName
		log.Infof("Failed to log's filename from configuration; falling back to default: '%s'", defaultLogFileName)
	}

	file, err := os.OpenFile(logs.logFileName, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Info("Failed to log to file, using default: 'stderr'")
	}
	log.SetOutput(file)
}

/*
Init initializes the log configuration
*/
func (logs *viperLogService) Init() {
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.file", defaultLogFileName)
}

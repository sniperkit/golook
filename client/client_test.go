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
package client

import (
	"encoding/json"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ottenwbe/golook/helper"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	sysName = "system"
)

var _ = Describe("The client", func() {

	var (
		server *httptest.Server
	)

	AfterEach(func() {
		server.Close()
	})

	It("should return a valid system", func() {
		server = httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			s := makeTestSystem()
			bytes, _ := json.Marshal(s)
			fmt.Fprintln(writer, string(bytes))
		}))
		defer server.Close()

		result := DoGetSystem(sysName)
		Expect(result).To(Not(BeNil()))
		//Expect(result.Name).To(Equal(sysName)) //TODO
	})

})

func TestDoGetHome(t *testing.T) {
	testString := "TestString"
	server := httptest.NewServer(
		http.HandlerFunc(
			func(writer http.ResponseWriter, _ *http.Request) {
				fmt.Fprintln(writer, testString)
			}))
	defer server.Close()

	if s := DoGetHome(); s != testString {
		t.Logf("TestString not successfully retrieved. Expected %s got %s", testString, s)
	}
}

func makeTestSystem() *helper.System {
	s := &helper.System{
		Name:  sysName,
		Files: nil,
		IP:    "1.1.1.1",
		OS:    "linux",
		UUID:  "1234"}
	return s
}

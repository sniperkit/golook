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
	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ottenwbe/golook/utils"
	"net/http"
	"net/http/httptest"
)

const (
	sysName = "system"
)

var _ = Describe("The client", func() {

	var (
		server *httptest.Server
		client *LookClient
	)

	BeforeEach(func() {
		client = NewLookClient()
	})

	AfterEach(func() {
		// ensure that the close method is executed and not forgotten
		server.Close()
		client = nil
	})

	Context(" System Methods ", func() {
		It("should return a valid system with Get", func() {
			server = httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
				s := newTestSystem()
				bytes, _ := json.Marshal(s)
				fmt.Fprintln(writer, string(bytes))
			}))
			client.serverUrl = server.URL

			result, err := client.DoGetSystem(sysName)
			Expect(err).To(BeNil())
			Expect(result).To(Not(BeNil()))
			Expect(result.Name).To(Equal(sysName))
		})

		It("should return a nil system with Get when the server does not exist", func() {
			client.serverUrl = "/"
			result, err := client.DoGetSystem(sysName)
			Expect(result).To(BeNil())
			Expect(err).To(Not(BeNil()))
		})

		It("should send a valid system to the server with Put", func() {

			testSystem := newTestSystem()

			server = httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
				receivedSystem, _ := utils.DecodeSystem(req.Body)
				Expect(receivedSystem.Name).To(Equal(testSystem.Name))
			}))
			client.serverUrl = server.URL

			result := client.DoPutSystem(testSystem)

			Expect(result).To(Not(BeNil()))
		})

		It("should transfer the delete request for a specific system to the server with DELETE", func() {

			testSystemName := "testSystem"

			server = httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
				params := mux.Vars(req)
				system := params["system"]
				Expect(system).To(Equal(testSystemName))
			}))
			server.URL = "/systems/{system}"
			client.serverUrl = server.URL

			result := client.DoDeleteSystem(testSystemName)

			Expect(result).To(Not(BeNil()))
		})
	})

	Context("Get Home", func() {

		const testString = "TestString"

		It("should pass the string which was sent by a server to the calle of DoGetHome()", func() {

			server := httptest.NewServer(
				http.HandlerFunc(
					func(writer http.ResponseWriter, _ *http.Request) {
						fmt.Fprintln(writer, testString)
					}))
			client.serverUrl = server.URL

			Expect(client.DoGetHome()).To(Equal(testString + "\n"))
		})
	})

})

func newTestSystem() *utils.System {
	s := &utils.System{
		Name:  sysName,
		Files: nil,
		IP:    "1.1.1.1",
		OS:    "linux",
		UUID:  "1234"}
	return s
}

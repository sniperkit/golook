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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ottenwbe/golook/broker/models"
	"github.com/ottenwbe/golook/broker/routing"
	"github.com/sirupsen/logrus"
	"reflect"
)

var _ = Describe("The query service", func() {
	It("creates a local query service by default", func() {
		q := newQueryService(LocalQueries, &router{})
		Expect(q).ToNot(BeNil())
		Expect(reflect.TypeOf(q)).To(Equal(reflect.TypeOf(&localQueryService{})))
	})

	It("creates broadcast query service when BCastQueries is specified", func() {
		q := newQueryService(BCastQueries, &router{})
		Expect(q).ToNot(BeNil())
		Expect(reflect.TypeOf(q)).To(Equal(reflect.TypeOf(&broadcastQueryService{})))
	})

	It("calls the routing service to initiate the query", func() {

		q := broadcastQueryService{router: &router{routing.NewMockedRouter()}}
		q.query("test.txt")

		Expect(q.router.Router.(*routing.MockRouter).Visited).To(Equal(1))
		Expect(q.router.Router.(*routing.MockRouter).VisitedMethod).To(Equal(fileQuery))

	})

	It("merges files responses", func() {

		f, err := models.NewFile("query_service_test")
		if err != nil {
			logrus.Fatal("Cannot create file.")
		}

		var expectedResult = FileQueryData{
			"a": {
				f,
			},
			"b": {
				f,
			},
		}

		var testResult = mergeFileQuery(
			testMerge{
				data: &FileQueryData{
					"a": {
						f,
					},
				},
			},
			testMerge{
				data: &FileQueryData{
					"b": {
						f,
					},
				},
			},
		)

		Expect(testResult.(FileQueryData)).To(Equal(expectedResult))

	})
})

type testMerge struct {
	data *FileQueryData
}

func (m testMerge) Unmarshal(v interface{}) error {
	tmp := reflect.ValueOf(v).Elem()
	tmp.Set(reflect.ValueOf(m.data).Elem())
	return nil
}

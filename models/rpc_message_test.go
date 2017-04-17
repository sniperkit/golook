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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("The encapsulated message", func() {
	It("should comprise a method name and the content after its creation", func() {
		m, err := NewRpcMessage("method", "msg")
		Expect(err).To(BeNil())
		Expect(m.Method).To(Equal("method"))
		Expect(len(m.Content)).ToNot(BeZero())
	})

	It("should return an error when the content cannot be encapsulated", func() {
		testChan := make(chan bool)
		m, err := NewRpcMessage("method", testChan)
		Expect(err).ToNot(BeNil())
		Expect(m.Method).To(Equal("invalid"))
	})

	It("should support to get the encapsulated method", func() {
		m, err := NewRpcMessage("method", "msg")

		var s string
		m.GetEncapsulated(&s)

		Expect(err).To(BeNil())
		Expect(s).To(Equal("msg"))
	})
})

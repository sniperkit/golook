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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ottenwbe/golook/broker/communication"
	"github.com/ottenwbe/golook/broker/models"
	"github.com/ottenwbe/golook/utils"
	"github.com/sirupsen/logrus"
	"reflect"
)

var _ = Describe("The broadcast router", func() {

	BeforeEach(func() {
		communication.ClientType = communication.MockRPC
	})

	It("implements the 'Router' interface", func() {
		r := newBroadcastRouter("test")

		Expect(r).ToNot(BeNil())
		Expect(reflect.TypeOf(r)).To(Equal(reflect.TypeOf(&BroadCastRouter{})))
	})

	It("broadcasts messages to one or more peerClients", func() {

		r := newBroadcastRouter("test")
		r.NewPeer(NewKey("peer1"), "1")
		r.NewPeer(NewKey("peer2"), "2")

		r.BroadCast("test", 123)

		for _, peer := range r.(*BroadCastRouter).routeTable.peers() {
			Expect(peer.(*communication.MockClient).VisitedCall).To(Equal(1))
			Expect(peer.(*communication.MockClient).Name).To(Equal("test"))
		}
	})

	It("should flood instead of sending directed messages via 'Route'", func() {

		r := newBroadcastRouter("test")
		r.NewPeer(NewKey("peer1"), "1")
		r.NewPeer(NewKey("peer2"), "2")

		r.Route(NewKey("peer2"), "test", 123)

		for _, peer := range r.(*BroadCastRouter).routeTable.peers() {
			Expect(peer.(*communication.MockClient).VisitedCall).To(Equal(1))
			Expect(peer.(*communication.MockClient).Name).To(Equal("test"))
		}
	})

	It("routes messages to corresponding handlers.", func() {

		const testHandlerName = "testHandler"
		r := newBroadcastRouter("test")
		testMsgHandler := &testHandler{}

		r.AddHandler(testHandlerName, &Handler{testMsgHandler.testMsgHandle, testMsgHandler.testMerge})

		testRequestMessage, _ := NewRequestMessage(NilKey(), 0, testHandlerName, "test")
		r.Handle(testHandlerName, Params(utils.MarshalSD(testRequestMessage)))

		Expect(testMsgHandler.message).To(Equal("test"))

	})
})

type testHandler struct {
	message string
}

func (t *testHandler) testMsgHandle(params models.EncapsulatedValues) interface{} {
	err := params.Unmarshal(&t.message)
	if err != nil {
		logrus.Fatal("Error handling message")
	}
	return nil
}

func (t *testHandler) testMerge(raw1 models.EncapsulatedValues, raw2 models.EncapsulatedValues) interface{} {
	return nil
}

type mockPeer struct {
	visitedCall int
	request     *RequestMessage
}

func (p *mockPeer) Call(index string, message interface{}) (models.EncapsulatedValues, error) {
	p.visitedCall++
	p.request = message.(*RequestMessage)
	return nil, nil
}

func (mockPeer) URL() string {
	return "test"
}

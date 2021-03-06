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

package communication

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/osamingo/jsonrpc"
	golook "github.com/ottenwbe/golook/broker/runtime/core"
	log "github.com/sirupsen/logrus"
)

type (
	/*JSONRPCServerStub implements the server stub for a particular handler function. The handler function is called when a valid message is received from a client.*/
	JSONRPCServerStub struct {
		handler string
		active  bool
	}

	/*JSONRPCParams implements the interface EncapsulatedValues.*/
	JSONRPCParams struct {
		params json.RawMessage
	}
)

var (
	/*HTTPRPCServer is the http server for accepting json rpc messages*/
	HTTPRPCServer golook.Server
)

var _ (jsonrpc.Handler) = (*JSONRPCServerStub)(nil)

/*
ServeJSONRPC handles json rpc messages which arrive for registered handlers
*/
func (rpc *JSONRPCServerStub) ServeJSONRPC(_ context.Context, params *json.RawMessage) (interface{}, *jsonrpc.Error) {

	// if the interface is not active, return an error
	if !rpc.active {
		return nil, jsonrpc.ErrMethodNotFound()
	}

	jsonRPCLogger().Debug("Received RPC message: %s", string(*params))

	p := &JSONRPCParams{params: *params}

	response, err := MessageDispatcher.handleMessage(rpc.handler, p)
	if err != nil {
		jsonRPCLogger().WithError(err).Error("Error when dispatching Json RPC call.")
		return response, jsonrpc.ErrMethodNotFound()
	}
	return response, nil

}

/*
Associate a handler with a json rpc server
*/
func (rpc *JSONRPCServerStub) Associate(handlerName string, request interface{}, response interface{}) {
	rpc.handler = handlerName
	rpc.active = true
	jsonrpc.RegisterMethod(handlerName, rpc, request, response)
}

/*
Finalize ensures that a server stub does not forward messages to a handler
*/
func (rpc *JSONRPCServerStub) Finalize() {
	//Unfortunately jsonrpc has no method for removing a registered function, therefore we only mark it as deleted
	rpc.active = false
}

/*
Unmarshal has to be called by the receiver of the parameter in order to unmarshal the params.
*/
func (p *JSONRPCParams) Unmarshal(v interface{}) error {

	var interfaceParams []json.RawMessage
	if err := jsonrpc.Unmarshal(&p.params, &interfaceParams); err != nil {
		return err
	}

	if len(interfaceParams) == 1 {
		if err := jsonrpc.Unmarshal(&interfaceParams[0], v); err != nil {
			return err
		}
	} else {
		return errors.New("Slices are not supported")

	}

	return nil
}

func jsonRPCLogger() *log.Entry {
	return log.WithField("com", "jsonRPCServerStub")
}

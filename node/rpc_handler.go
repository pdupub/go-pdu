// Copyright 2024 The PDU Authors
// This file is part of the PDU library.
//
// The PDU library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The PDU library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the PDU library. If not, see <http://www.gnu.org/licenses/>.

package node

import (
	"encoding/json"
	"io"
	"net/http"
)

// JSON-RPC request structure
type JSONRPCRequest struct {
	Jsonrpc string          `json:"jsonrpc"`
	Method  string          `json:"method,omitempty"`
	Params  json.RawMessage `json:"params"`
	ID      json.RawMessage `json:"id,omitempty"`
}

// JSON-RPC response structure
type JSONRPCResponse struct {
	Jsonrpc string          `json:"jsonrpc"`
	Result  interface{}     `json:"result,omitempty"`
	Error   *RPCError       `json:"error,omitempty"`
	ID      json.RawMessage `json:"id"`
}

// RPCError structure
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func withCORS(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		handler.ServeHTTP(w, r)
	})
}

func (n *Node) handleRPCRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is supported", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Log the incoming request
	// log.Printf("Received request: %s\n", string(body))

	var req JSONRPCRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Handle standard JSON-RPC requests
	var result interface{}
	var rpcErr *RPCError

	switch req.Method {
	case "eth_chainId":
		result = "0x2304" // Example chain ID, replace with actual chain ID
	case "net_version":
		result = "1" // Example network version, replace with actual network version
	case "eth_blockNumber":
		result = "0xBC614E" // Example block number, replace with actual block number
	case "eth_getBlockByNumber":
		result = "0xBC614E" // Example block data, replace with actual block data
	case "eth_gasPrice":
		result = "0x09184e72a000" // Example gas price, replace with actual gas price
	case "eth_getBalance":
		result = "0x8AC7230489E80000" // Example balance, replace with actual logic to fetch balance
	case "eth_getTransactionCount":
		result = "0x1" // Example transaction count, replace with actual logic to fetch transaction count
	case "pdu_sendQuantums":
		n.handleSendQuamtumsRequest(req.Params)
		result = map[string]string{"status": "success"}
	case "pdu_getQuantums":
		quantums := n.handleQueryQuamtumsRequest(req.Params)
		result = map[string]interface{}{"quantums": quantums, "status": "success"}

	// case "eth_call":
	// 	result, rpcErr = handleEthCall(req.Params)
	default:
		rpcErr = &RPCError{
			Code:    -32601,
			Message: "Method not found",
		}
	}

	resp := JSONRPCResponse{
		Jsonrpc: "2.0",
		Result:  result,
		Error:   rpcErr,
		ID:      req.ID,
	}

	responseBody, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Log the outgoing response
	// log.Printf("Sending response: %s\n", string(responseBody))

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBody)
}

package node

import (
	"encoding/json"
	"log"

	"github.com/pdupub/go-pdu/core"
)

func handleRecvQuamtumsRequest(params json.RawMessage) {

	var paramsData []interface{}
	err := json.Unmarshal(params, &paramsData)
	if err != nil {
		return
	}

	// Log the JSON data received
	log.Printf("Received custom JSON data: %+s\n", params)
	if len(paramsData) == 0 || paramsData[0] == nil {
		return
	}

	body, err := json.Marshal(paramsData[0])
	if err != nil {
		return
	}

	quantum, err := core.JsonToQuantum(body)
	if err != nil {
		return
	}

	// // Log the quantum
	// log.Printf("Quantum: %+v\n", quantum)
	// log.Printf("Quantum signature: %s\n", core.Sig2Hex(quantum.Signature))
	// log.Printf("Quantum content[0] data: %+s\n", quantum.Contents[0].Data)
	// log.Printf("Quantum content[0] format: %+s\n", quantum.Contents[0].Format)
	// log.Printf("Quantum ref[0]: %s\n", core.Sig2Hex(quantum.References[0]))

	addr, err := quantum.Ecrecover()
	if err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		// return
		log.Printf("Failed to recover address: %s\n", err.Error())
	}
	log.Printf("Quantum address: %s\n", addr.Hex())

}

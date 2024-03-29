// Copyright 2021 The PDU Authors
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
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/pdupub/go-pdu/core"
	"github.com/pdupub/go-pdu/identity"
	"github.com/pdupub/go-pdu/udb/fb"

	"github.com/labstack/echo/v4"
)

var (
	errQuantumAlreadyExist = RespErr{ErrCode: 1, ErrMsg: core.ErrQuantumAlreadyExist.Error()}
	errSignatureIncorrect  = RespErr{ErrCode: 2, ErrMsg: core.ErrSignatureIncorrect.Error()}
	errSelfRefMissing      = RespErr{ErrCode: 3, ErrMsg: core.ErrSelfRefMissing.Error()}
	errOffspringDuplicate  = RespErr{ErrCode: 4, ErrMsg: core.ErrOffspringDuplicate.Error()}
	errAncestorMissing     = RespErr{ErrCode: 5, ErrMsg: core.ErrAncestorMissing.Error()}

	errCodeRequestHeader = 101
	errCodeRequestBody   = 102
	errCodeRequestJSON   = 110
	errCodeUnknown       = 200
)

// Node
type Node struct {
	interval int64
	univ     core.Universe
	qChan    chan *core.Quantum
	e        *echo.Echo
}

type Resp struct {
	Data  interface{} `json:"data,omitempty"`
	Error RespErr     `json:"error,omitempty"`
}

type RespErr struct {
	ErrCode int           `json:"code"`
	ErrMsg  string        `json:"message"`
	Params  []interface{} `json:"params,omitempty"` // key-value
}

func New(interval int64, firebaseKeyPath, firebaseProjectID string) (*Node, error) {
	ctx := context.Background()
	fbu, err := fb.NewFBUniverse(ctx, firebaseKeyPath, firebaseProjectID)
	if err != nil {
		return nil, err
	}
	return &Node{interval: interval, univ: fbu, qChan: make(chan *core.Quantum), e: echo.New()}, nil
}

func (n *Node) RunEcho(port int64, c <-chan os.Signal) {
	n.e.HideBanner = true
	n.e.Static("/resource", "resource")
	n.e.POST("/rec", n.receiverHandler)
	n.e.GET("/individual/:address", n.getIndividualHandler)
	n.e.POST("/report", n.reportHandler)

	go n.Run(c)
	n.e.Logger.Fatal(n.e.Start(fmt.Sprintf(":%d", port)))
}

func (n *Node) getIndividualHandler(c echo.Context) error {
	resp := Resp{}

	addrHex := c.Param("address")
	individual, err := n.univ.GetIndividual(identity.HexToAddress(addrHex))
	if err != nil {
		resp.Error = RespErr{ErrCode: errCodeRequestJSON, ErrMsg: err.Error()}
		return c.JSON(http.StatusOK, resp)
	}
	resp.Data = map[string]interface{}{"status": "ok", "individual": individual}
	return c.JSON(http.StatusOK, resp)
}

func (n *Node) reportHandler(c echo.Context) error {
	resp := Resp{}
	if c.Request().Header.Get("Content-Type") != echo.MIMEApplicationJSON {
		resp.Error = RespErr{ErrCode: errCodeRequestHeader, ErrMsg: ""}
		return c.JSON(http.StatusOK, resp)
	}

	data, err := io.ReadAll(c.Request().Body)
	if err != nil {
		resp.Error = RespErr{ErrCode: errCodeRequestBody, ErrMsg: err.Error()}
		return c.JSON(http.StatusOK, resp)
	}
	// report is not part of pdu, only function for iOS or Android platform
	// user can report quantum, user or species. Node will record if the signer exist,
	// but should not do anything auto. report should not use quantum struct to avoid
	// be used as bad behavir proof of signer.

	var report fb.Report
	err = json.Unmarshal(data, &report)
	if err != nil {
		resp.Error = RespErr{ErrCode: errCodeRequestJSON, ErrMsg: err.Error()}
		return c.JSON(http.StatusOK, resp)
	}

	err = n.univ.(*fb.FBUniverse).ReceiveReport(&report)
	if err != nil {
		resp.Error = RespErr{ErrCode: errCodeUnknown, ErrMsg: err.Error()}
		return c.JSON(http.StatusOK, resp)
	}
	resp.Data = map[string]string{"status": "ok"}
	return c.JSON(http.StatusOK, resp)
}

func (n *Node) receiverHandler(c echo.Context) error {
	resp := Resp{}
	if c.Request().Header.Get("Content-Type") != echo.MIMEApplicationJSON {
		resp.Error = RespErr{ErrCode: errCodeRequestHeader, ErrMsg: ""}
		return c.JSON(http.StatusOK, resp)
	}

	data, err := io.ReadAll(c.Request().Body)
	if err != nil {
		resp.Error = RespErr{ErrCode: errCodeRequestBody, ErrMsg: err.Error()}
		return c.JSON(http.StatusOK, resp)
	}

	var quantum core.Quantum
	err = json.Unmarshal(data, &quantum)
	if err != nil {
		resp.Error = RespErr{ErrCode: errCodeRequestJSON, ErrMsg: err.Error()}
		return c.JSON(http.StatusOK, resp)
	}

	// check if exist already
	// GetQuantum return err if not exist
	_, err = n.univ.GetQuantum(quantum.Signature)
	if err == nil {
		resp.Error = errQuantumAlreadyExist
		return c.JSON(http.StatusOK, resp)
	}

	// check if refs[0] is individual.last
	addr, err := quantum.Ecrecover()
	if err != nil {
		resp.Error = errSignatureIncorrect
		return c.JSON(http.StatusOK, resp)
	}

	// get individual will return err if address not exist, so ignore if err return
	indv, err := n.univ.GetIndividual(addr)
	if err == nil {
		// check if refs[0] is exist
		if len(quantum.References) == 0 {
			resp.Error = errSelfRefMissing
			return c.JSON(http.StatusOK, resp)
		}

		if core.Sig2Hex(indv.LastSig) != core.Sig2Hex(quantum.References[0]) {
			params := []interface{}{map[string]string{"last": core.Sig2Hex(indv.LastSig)}}
			if _, err = n.univ.GetQuantum(quantum.References[0]); err == nil || core.Sig2Hex(quantum.References[0]) == core.Sig2Hex(core.FirstQuantumReference) {
				// if exist return duplicate ref use, return err with sigHex of last
				resp.Error = errOffspringDuplicate
				resp.Error.Params = params
			} else {
				// if not exist return missing quantums since last err with sigHex of last
				// !!! return err, but still accept the quantum
				n.qChan <- &quantum
				resp.Error = errAncestorMissing
				resp.Error.Params = params
			}
			return c.JSON(http.StatusOK, resp)

		}
	}

	n.qChan <- &quantum
	resp.Data = map[string]string{"status": "ok"}
	return c.JSON(http.StatusOK, resp)
}

// Run
func (n *Node) Run(c <-chan os.Signal) {
	var runCnt uint64
	busy := false
	for {
		select {
		case <-c:
			log.Println("main closed")
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			if err := n.e.Shutdown(ctx); err != nil {
				n.e.Logger.Fatal(err)
			}
			return
		case <-time.After(time.Second * time.Duration(n.interval)):
			log.Println("execution", runCnt)
			runCnt++
			if !busy {
				busy = true
				n.univ.ProcessQuantums(10, 0)
				busy = false
			}
		case q := <-n.qChan:
			log.Println("execution", core.Sig2Hex(q.Signature))
			n.univ.ReceiveQuantums([]*core.Quantum{q})
		}
	}
}

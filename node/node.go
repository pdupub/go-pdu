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
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/pdupub/go-pdu/core"
	"github.com/pdupub/go-pdu/udb/fb"

	"github.com/labstack/echo/v4"
)

// Node
type Node struct {
	interval int64
	univ     core.Universe
	qChan    chan *core.Quantum
	e        *echo.Echo
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
	n.e.POST("/rec", n.receiverHandler)

	go n.Run(c)
	n.e.Logger.Fatal(n.e.Start(fmt.Sprintf(":%d", port)))
}

func (n *Node) receiverHandler(c echo.Context) error {
	if c.Request().Header.Get("Content-Type") != echo.MIMEApplicationJSON {
		return c.String(http.StatusBadRequest, "")
	}

	data, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	var quantum core.Quantum
	err = json.Unmarshal(data, &quantum)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	// c.Logger().Info(string(data))

	// check if exist already
	// n.univ.GetQuantum(quantum.Signature)

	// check if refs[0] is individual.last
	// if not
	// check if refs[0] is exist
	// if exist return duplicate ref use, return err with sigHex of last
	// if not exist return missing quantums since last err with sigHex of last

	n.qChan <- &quantum
	return c.JSON(http.StatusOK, nil)
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

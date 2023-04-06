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
	"net"
	"net/http"
	"os"
	"time"

	"github.com/pdupub/go-pdu/core"
	"github.com/pdupub/go-pdu/udb/fb"

	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "8123"
	CONN_TYPE = "tcp"
)

// Node
type Node struct {
	interval int64
	univ     core.Universe
}

func New(interval int64, firebaseKeyPath, firebaseProjectID string) (*Node, error) {
	ctx := context.Background()
	fbu, err := fb.NewFBUniverse(ctx, firebaseKeyPath, firebaseProjectID)
	if err != nil {
		return nil, err
	}
	return &Node{interval: interval, univ: fbu}, nil
}

func (n *Node) RunListenSample() {
	// Listen for incoming connections.
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go n.handleRequest(conn)
	}
}

// Handles incoming requests.
func (n *Node) handleRequest(conn net.Conn) {
	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)
	// Read the incoming connection into the buffer.
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	// Send a response back to person contacting us.
	conn.Write([]byte("Message received."))
	// Close the connection when you're done with it.
	conn.Close()
}

func (n *Node) RunGin(port int64) {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run()
}

func (n *Node) RunEcho(port int64) {
	e := echo.New()
	e.HideBanner = true
	e.POST("/rec", n.receiverHandler)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
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

	n.univ.ReceiveQuantums([]*core.Quantum{&quantum})
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
			return
		case <-time.After(time.Second * time.Duration(n.interval)):
			log.Println("execution", runCnt)
			runCnt++
			if !busy {
				busy = true
				n.univ.ProcessQuantums(10, 0)
				busy = false
			}
		}
	}
}

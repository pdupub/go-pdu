// Copyright 2018 The PDU Authors
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

package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/TATAUFO/PDU/mydb"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type NatureTime struct {
	Timestamp int64  `json:"timestamp"` // Timestamp
	Proof     string `json:"proof"`     // The Proof of Not Before Timestamp
}

var (
	err    error
	memory *mydb.MyDB
)

func initDB() error {
	memory, err = mydb.Open()
	if err != nil {
		return err
	}
	return nil
}

func timeResponse(c echo.Context) error {

	t := time.Now().UnixNano()
	p := MD5(strconv.Itoa(int(t)))
	natureTime := NatureTime{t, p}
	res, err := json.Marshal(natureTime)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	memory.SaveTime(t, p)
	return c.String(http.StatusOK, string(res))
}

func MD5(text string) string {
	ctx := md5.New()
	ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))
}

func main() {
	// Echo instance
	e := echo.New()

	// init DB
	err = initDB()
	if err != nil {
		e.Logger.Fatal(err)
	}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", timeResponse)

	// Port
	port := "1323"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}
	// Start server
	e.Logger.Fatal(e.Start(":" + port))
}

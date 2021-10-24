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

package p2p

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"path"
	"time"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/pdupub/go-pdu/core"
	"go.etcd.io/bbolt"
)

func New(templatePath, dbPath string, g *core.Genesis, port int, ignoreUnknownSource bool, peers []string) {
	// Setup
	e := echo.New()
	// Hide banner
	e.HideBanner = true
	// e.IPExtractor = echo.ExtractIPDirect()
	// e.IPExtractor = echo.ExtractIPFromXFFHeader()
	// e.IPExtractor = echo.ExtractIPFromRealIPHeader()
	e.Logger.SetLevel(log.INFO)

	tpl, err := template.ParseGlob(path.Join(templatePath, "views/*.html"))
	if err != nil {
		e.Logger.Error(err)
		e.Logger.Error("web display disabled")
	} else {
		renderer := &Template{
			templates: tpl,
		}
		e.Renderer = renderer
		e.Static("/js", path.Join(templatePath, "js"))
	}

	db, err := bbolt.Open(dbPath, 0600, &bbolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		e.Logger.Fatal(err)
	}
	// Close database connection
	defer db.Close()

	// Stats
	bn := NewBootNode(db, e.Logger)
	bn.SetIgnoreUnknownSource(ignoreUnknownSource)
	if err := bn.SetUniverse(g.GetUniverse()); err != nil {
		e.Logger.Fatal(err)
	}
	if err := bn.UpsertDBMsgs(g.GetMsgs()); err != nil {
		e.Logger.Fatal(err)
	}
	if err := bn.LoadUniverse(); err != nil {
		e.Logger.Fatal(err)
	}

	if err := bn.AddPeers(peers); err != nil {
		e.Logger.Error(err)
	}
	e.Use(bn.Process)
	// Routes

	// node information
	e.GET("/", bn.welcome)

	// receive message from client & peers
	e.POST("/", bn.newMsg)

	// upload image for test
	e.Static("/files", "tmp/files")
	e.POST("/upload", bn.upload)

	// dag topology
	e.GET("/society", bn.getSociety)
	e.GET("/entropy", bn.getEntropy)
	e.GET("/message/:sig", bn.getMessage)
	e.GET("/profile/:uid", bn.getProfile)

	// k-v database
	e.GET("/info/detail/:uid", bn.getDetail)
	e.GET("/info/latest/:uid", bn.getLatest)
	e.GET("/info/full/:sig", bn.getFullMsg)
	e.GET("/info/photons", bn.getPhotons)

	// graphql
	schema := graphql.MustParseSchema(schema, &query{})
	gQLQuery := echo.WrapHandler(&relay.Handler{Schema: schema})
	e.POST("/query", gQLQuery)

	// Start server
	go func() {
		if err := e.Start(fmt.Sprintf(":%d", port)); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

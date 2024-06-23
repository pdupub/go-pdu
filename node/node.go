package node

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/libp2p/go-libp2p/p2p/transport/tcp"
)

func createNode() (host.Host, context.Context) {
	ctx := context.Background()
	h, err := libp2p.New(
		libp2p.Transport(tcp.NewTCPTransport),
	)
	if err != nil {
		log.Fatal(err)
	}
	return h, ctx
}

func handleInterrupt(ctx context.Context, h host.Host) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		fmt.Println("\nShutting down...")
		if err := h.Close(); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
}

func Run() {
	h, ctx := createNode()
	handleInterrupt(ctx, h)

	h.SetStreamHandler(protocol.ID("/p2p/1.0.0"), func(s network.Stream) {
		fmt.Println("Got a new stream!")
		// 处理流的代码
	})

	fmt.Printf("Node ID: %s\n", h.ID().String())
	for _, addr := range h.Addrs() {
		fmt.Printf("Node Address: %s\n", addr.String())
	}

	<-ctx.Done() // 保持程序运行
}

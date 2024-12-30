package main

import (
	// "context"
	"context"
	"fmt"
	"log"

	"github.com/libp2p/go-libp2p"
	kaddht "github.com/libp2p/go-libp2p-kad-dht"
)

func main() {
	ctx := context.Background()

	host, err := libp2p.New()

	if err != nil {
		log.Fatalf("创建主机失败: %v", err)
	}
	defer host.Close()

	_, err = kaddht.New(ctx, host)
	if err != nil {
		log.Fatalf("启动 DHT 失败: %v", err)
	}

	fmt.Printf("DHT 节点已启动: %s\n", host.ID())
	// fmt.Println("Hello world!")
	select {}
}

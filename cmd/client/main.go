package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	shell "github.com/ipfs/go-ipfs-api"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

func main() {
	ctx := context.Background()

	// 创建一个新的 libp2p Host
	host, err := libp2p.New(
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"))
	if err != nil {
		log.Fatalf("Failed to create host: %s", err)
	}

	// 打印节点的地址
	fmt.Println("This node's addresses:")
	for _, addr := range host.Addrs() {
		fmt.Printf("%s/p2p/%s\n", addr, host.ID().String())
	}

	// 创建一个新的 DHT 实例
	kademliaDHT, err := dht.New(ctx, host)
	if err != nil {
		log.Fatalf("Failed to create DHT: %s", err)
	}

	// 启动 DHT
	if err := kademliaDHT.Bootstrap(ctx); err != nil {
		log.Fatalf("Failed to bootstrap DHT: %s", err)
	}

	// 连接到 IPFS 引导节点
	bootstrapNodes := []string{

		"/dnsaddr/bootstrap.libp2p.io/p2p/QmNnooDu7bfjPFoTZYxMNLWUQJyrVwtbZg5gBMjTezGAJN",
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmQCU2EcMqAqQPR2i9bChDtGNJchTbq5TbXJJ16u19uLTa",
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmbLHAnMoJPWSCR5Zhtx6BHJX9KiKNN6tpvbUcqanj75Nb",
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmcZf59bWwK5XFi76CZX8cbJ4BhTzzA3gU1ZjYZcYW3dwt",
		"/ip4/104.131.131.82/tcp/4001/p2p/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ",
	}

	for _, addr := range bootstrapNodes {
		maddr, err := multiaddr.NewMultiaddr(addr)
		if err != nil {
			log.Printf("Failed to parse multiaddr %s: %s", addr, err)
			continue
		}

		info, err := peer.AddrInfoFromP2pAddr(maddr)
		if err != nil {
			log.Printf("Failed to get peer info from multiaddr %s: %s", addr, err)
			continue
		}

		if err := host.Connect(ctx, *info); err != nil {
			log.Printf("Failed to connect to bootstrap node %s: %s", info.ID, err)
		} else {
			fmt.Printf("Connected to bootstrap node %s\n", info.ID.String())
		}
	}

	// 使用 DHT 查找对等方
	peerID := "QmTzQ1Npf7BZw2BzGiCBzBvPxT7Q3xtFYsHwo77G5n98RN" // 示例对等方ID
	if pi, err := kademliaDHT.FindPeer(ctx, peer.ID(peerID)); err != nil {
		log.Printf("Failed to find peer %s: %s", peerID, err)
	} else {
		fmt.Printf("Found peer %s at %s\n", peerID, pi.Addrs)
	}

	// 添加内容到 IPFS
	ipfsAddr := "/ip4/127.0.0.1/tcp/5001" // 确保 IPFS 守护程序在这个地址运行
	sh := shell.NewShell(ipfsAddr)

	content := "Hello, IPFS!"
	cid, err := addContentToIPFS(sh, content)
	if err != nil {
		log.Fatalf("Failed to add content to IPFS: %s", err)
	} else {
		fmt.Printf("Content added to IPFS with CID: %s\n", cid)
	}

	// 阻塞主线程，避免程序退出
	select {}
}

// addContentToIPFS 向 IPFS 添加内容并返回内容 ID (CID)
func addContentToIPFS(sh *shell.Shell, content string) (string, error) {
	reader := strings.NewReader(content)
	cid, err := sh.Add(reader)
	if err != nil {
		return "", err
	}
	return cid, nil
}

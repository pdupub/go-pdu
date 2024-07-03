# [ParaDigi Universe | PDU](https://pdu.pub) &nbsp; [![Tweet](https://img.shields.io/twitter/url/http/shields.io.svg?style=social)](https://x.com/pdupub) &nbsp; [![Telegram](https://img.shields.io/badge/-telegram-red?color=white&logo=telegram)](https://t.me/pdugroup)

[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/pdupub/go-pdu)
[![GoReport](https://goreportcard.com/badge/github.com/pdupub/go-pdu)](https://goreportcard.com/report/github.com/pdupub/go-pdu)
[![License](https://img.shields.io/badge/license-GPL%20v3-blue.svg)](LICENSE)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go#distributed-systems)

Golang implementation of PDU.

- [What is PDU?](#what-is-pdu)
- [Usage](#usage)
- [Development](#development)
- [Contributing](#contributing)

## Join Test

iOS : [https://testflight.apple.com/join/FqQGxhbn](https://testflight.apple.com/join/FqQGxhbn)

## What is PDU?

PDU is short for ParaDigi Universe, is a fully peer-to-peer (P2P) social network system designed to enable participants to freely publish and efficiently access information without relying on any third-party services. Traditional systems that do not use centralized verification methods, such as phone numbers, are vulnerable to Sybil attacks, where the cost-free creation of new accounts can overwhelm the network with spam, undermining reward and punishment mechanisms. PDU addresses this issue by establishing trusted publisher identities through a sequence of messages signed by the same private key. Interactions such as reposts, comments, and likes create associations between publishers, allowing participants to form a custom set of visible publisher identities. This relatively stable scope enables an identity-based reward and punishment mechanism to effectively filter information.

Please read the WhitePaper on [https://pdu.pub](https://pdu.pub/white_paper.html) for more details.


## Usage

```
ParaDigi Universe
	A Peer-to-Peer Social Network Service
	Website: https://pdu.pub

Usage:
  pdu [flags]

Flags:
  -d, --dbName string    Database name (default "pdu.db")
  -h, --help             help for pdu
  -n, --nodeKey string   Node key (default "node.key")
  -p, --peerPort int     Port to listen for P2P connections (default 4001)
  -r, --rpcPort int      Port for the RPC server (default 8545)
  -t, --test             Run in test mode
  -w, --webPort int      Port for the web server (default 8546)
```


## Development

To copy the repository:

```
get clone https://github.com/pdupub/go-pdu.git

```

To build:
```
make build 
```


## Contributing

1. Fork the repository on GitHub to start making your changes to the master branch
2. Write a test which shows that the bug was fixed or that the feature works as expected
3. Send a pull request and bug the maintainer until it gets merged and published



<a href="https://pdu.pub"><img height="32" align="right" src="https://pdu.pub/assets/images/logo.png"></a>

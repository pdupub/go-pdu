# [ParaDigi Universe](https://pdu.pub) &nbsp; [![Tweet](https://img.shields.io/twitter/url/http/shields.io.svg?style=social)](https://twitter.com/intent/tweet?text=A%20decentralized%20identity-based%20social%20network&url=https://pdu.pub&via=PDUPUB&hashtags=P2P,SocialNetwork,decentralized,identity,Blockchain) &nbsp; [![Telegram](https://img.shields.io/badge/-telegram-red?color=white&logo=telegram)](https://t.me/pdugroup)

[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/pdupub/go-pdu)
[![GoReport](https://goreportcard.com/badge/github.com/pdupub/go-pdu)](https://goreportcard.com/report/github.com/pdupub/go-pdu)
[![Travis](https://travis-ci.org/pdupub/go-pdu.svg?branch=master)](https://travis-ci.org/pdupub/go-pdu)
[![License](https://img.shields.io/badge/license-GPL%20v3-blue.svg)](LICENSE)
[![Chat](https://img.shields.io/badge/gitter-Docs%20chat-4AB495.svg)](https://gitter.im/pdupub/go-pdu)
[![Coverage Status](https://coveralls.io/repos/github/pdupub/go-pdu/badge.svg?branch=master)](https://coveralls.io/github/pdupub/go-pdu?branch=master)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go#distributed-systems)

Golang implementation of PDU.

PDU的Go语言版本，以Firebase作为数据存储。


- [What is PDU?](#what-is-pdu)
- [Usage](#usage)
- [Development](#development)
- [Contributing](#contributing)

## Join Test

iOS : [https://testflight.apple.com/join/FqQGxhbn](https://testflight.apple.com/join/FqQGxhbn)

## What is PDU?

PDU is a social network service (SNS) system based on peer-to-peer (P2P) methods. All information in the system is identified by signature, called a message. Through references between messages, messages signed by the same key can form a total order relationship, and all messages can form a partial order relationship. Accounts are then defined on the basis of homologous total sequence messages. Any account can freely create a species, define its rules, and confirm other accounts to join its own species based on self-identity. In any species, accounts and information rows can be filtered according to the identification association relationship to achieve effective information acquisition. Please read the WhitePaper on [https://pdu.pub](https://pdu.pub/docs/en/WhitePaperV2.html) for more details.

PDU是基于点对点的方式构建社交网络服务，系统中所有信息均由签名确定其来源，称为消息。通过消息间的引用，同源消息可以构成全序关系，异源消息构成偏序关系。又以同源全序消息为基础，定义账户。任何账户都可以自由的创建社区规则，并基于自我认同，邀请其他账户加入自身所属社区。社区内的账户均存在偏序关联，可以基于这种关联关系，对账户和信息行筛选，以实现信息的有效获取。更多内容，详见[PDU白皮书(v2)](https://pdu.pub/docs/zh/WhitePaperV2.html)。

## Usage

```
Run node daemon

Usage:
  pdu run [loop interval] [flags]

Flags:
  -h, --help           help for run
      --interval int   time interval between consecutive processing on node (default 5)

Global Flags:
      --fbKeyPath string     path of firebase json key (default "udb/fb/test-firebase-adminsdk.json")
      --fbProjectID string   project ID (default "pdupub-a2bdd")
      --projectPath string   project root path (default "./")
```

```
Operations on node

Usage:
  pdu node [command]

Available Commands:
  backup      Backup processed quantums to local
  exe         Do process quantum once on node
  hide        Hide processed Quantum in node
  judge       Judge Individual & Species on your own node
  truncate    Clear up all data on firebase collections

Flags:
  -h, --help   help for node

Global Flags:
      --fbKeyPath string     path of firebase json key (default "udb/fb/test-firebase-adminsdk.json")
      --fbProjectID string   project ID (default "pdupub-a2bdd")
      --projectPath string   project root path (default "./")

Use "pdu node [command] --help" for more information about a command.
```

```
Create and Broadcast Message (For test your own node)

Usage:
  pdu msg [flags]

Flags:
  -h, --help   help for msg

Global Flags:
      --fbKeyPath string     path of firebase json key (default "udb/fb/test-firebase-adminsdk.json")
      --fbProjectID string   project ID (default "pdupub-a2bdd")
      --projectPath string   project root path (default "./")
```


## Development

To copy the repository:

```
get clone https://github.com/pdupub/go-pdu.git

```

To build:
```
make install 
```


## Contributing

1. Fork the repository on GitHub to start making your changes to the master branch
2. Write a test which shows that the bug was fixed or that the feature works as expected
3. Send a pull request and bug the maintainer until it gets merged and published


<a href="https://pdu.pub"><img height="32" align="right" src="https://pdu.pub/assets/img/logo.png"></a>

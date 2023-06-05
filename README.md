# [ParaDigi Universe](https://pdu.pub) &nbsp; [![Tweet](https://img.shields.io/twitter/url/http/shields.io.svg?style=social)](https://twitter.com/intent/tweet?text=A%20decentralized%20identity-based%20social%20network&url=https://pdu.pub&via=PDUPUB&hashtags=P2P,SocialNetwork,decentralized,identity,Blockchain) &nbsp; [![Telegram](https://img.shields.io/badge/-telegram-red?color=white&logo=telegram)](https://t.me/pdugroup)

[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/pdupub/go-pdu)
[![GoReport](https://goreportcard.com/badge/github.com/pdupub/go-pdu)](https://goreportcard.com/report/github.com/pdupub/go-pdu)
[![License](https://img.shields.io/badge/license-GPL%20v3-blue.svg)](LICENSE)
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

PDU is a social network service (SNS) system based on peer-to-peer and help users to effectively screen information publishers without relying on third-party authentication. All messages in the system determine the ordered relationship through mutual reference, and then determine their source by signature. Homologous total order message sequence is regarded as an information publisher identity, and all messages in the system can constitute one or more message sets with partial order relationship. Any information publisher is free to create a new species or identify other identities as belonging to a certain species. The user constructs the species range based on the obtained species identification information, and further filters suspicious information publishers according to the identification relationship. This process changes the unified verification and consistent user range in centralized services into a scalable species-based identity range determined by users themselves.

Please read the WhitePaper on [https://pdu.pub](https://pdu.pub/docs/en/WhitePaperV3.html) for more details.

PDU是基于点对点的方式构建社交网络服务。帮助使用者能够在不依赖第三方认证的情况下，实现对信息发布者的有效筛选。系统中所有消息通过相互的引用确定有序关系，再由签名确定其来源。同源的全序消息序列被视为一个信息发布者身份，而所有的消息在系统中可构成一个或者多个有偏序关系的消息集合。任何信息发布者都可以自由的创建新族群或对其他身份做出属于某族群的认定。使用者根据已获取的族群认定消息来构建族群范围，并依据认定关系进一步过滤可疑的信息发布者。此过程将中心化服务中统一验证和一致的用户范围，变为基于可扩展的族群，由使用者自行决定的身份范围。

更多内容，详见[PDU白皮书(v3)](https://pdu.pub/docs/zh/WhitePaperV3.html)。

## Usage

```
ParaDigi Universe
	A decentralized social networking service
	Website: https://pdu.pub

Usage:
  pdu [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  key         Create keystores (For test)
  msg         Create and Broadcast Message (For test your own node)
  node        Operations on node
  run         Run node daemon

Flags:
      --fbKeyPath string     path of firebase json key (default "udb/fb/test-firebase-adminsdk.json")
      --fbProjectID string   project ID (default "pdu-dev-1")
  -h, --help                 help for pdu
      --projectPath string   project root path (default "./")
  -v, --version              version for pdu

Use "pdu [command] --help" for more information about a command.
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

# [ParaDigi Universe](https://pdu.pub) &nbsp; [![Tweet](https://img.shields.io/twitter/url/http/shields.io.svg?style=social)](https://twitter.com/intent/tweet?text=A%20decentralized%20identity-based%20social%20network&url=https://pdu.pub&via=PDUPUB&hashtags=P2P,SocialNetwork,decentralized,identity,Blockchain)

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

## What is PDU?

PDU is a decentralized social networking service, all information in the system is determined by the signature of its source, referred to as a **message**. Through references between messages, same source messages can form a total order relationship, and different source messages form a partial order relationship. Based on the source total order messages, the definition of an account is established. Any account can freely create community rules and invite other accounts to join the community they belong to, based on self-identification. Accounts within the community have partial order relationships, so we can filter accounts and information to achieve effective information retrieval. Please read the WhitePaper on [https://pdu.pub](https://pdu.pub/docs/en/WhitePaperV2.html) for more details.

PDU是基于点对点的方式构建社交网络服务，系统中所有信息均由签名确定其来源，称为消息。通过消息间的引用，同源消息可以构成全序关系，异源消息构成偏序关系。又以同源全序消息为基础，定义账户。任何账户都可以自由的创建社区规则，并基于自我认同，邀请其他账户加入自身所属社区。社区内的账户均存在偏序关联，可以基于这种关联关系，对账户和信息行筛选，以实现信息的有效获取。更多内容，详见[PDU白皮书(v2)](https://pdu.pub/docs/zh/WhitePaperV2.html)。

## Usage

```
Parallel Digital Universe
	A decentralized social networking service
	Website: https://pdu.pub

Usage:
  pdu [command]

Available Commands:
  ck          Create keystores
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  send        Send sampel msg to node
  start       Start run node
  test        Test some methods

Flags:
  -h, --help                 help for pdu
      --projectPath string   project root path (default "./")

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

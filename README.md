# [Parallel Digital Universe](https://pdu.pub) &nbsp; [![Tweet](https://img.shields.io/twitter/url/http/shields.io.svg?style=social)](https://twitter.com/intent/tweet?text=A%20decentralized%20identity-based%20social%20network&url=https://pdu.pub&via=PDUPUB&hashtags=P2P,SocialNetwork,decentralized,identity,Blockchain)

[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/pdupub/go-pdu)
[![GoReport](https://goreportcard.com/badge/github.com/pdupub/go-pdu)](https://goreportcard.com/report/github.com/pdupub/go-pdu)
[![Travis](https://travis-ci.org/pdupub/go-pdu.svg?branch=master)](https://travis-ci.org/pdupub/go-pdu)
[![License](https://img.shields.io/badge/license-GPL%20v3-blue.svg)](LICENSE)
[![Chat](https://img.shields.io/badge/gitter-Docs%20chat-4AB495.svg)](https://gitter.im/pdupub/go-pdu)
[![Coverage Status](https://coveralls.io/repos/github/pdupub/go-pdu/badge.svg?branch=master)](https://coveralls.io/github/pdupub/go-pdu?branch=master)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go#distributed-systems)

Golang implementation of PDU.


- [What is PDU?](#pdu)
- [Usage](#usage)
- [Development](#development)
- [Contributing](#contributing)

## PDU
PDU is a decentralized identity-based social network, please read the WhitePaper on [github.com/pdupub/Documentation](https://github.com/pdupub/Documentation) for more details.


## Usage

```
Parallel Digital Universe
	A decentralized identity-based social network
	Website: https://pdu.pub

Usage:
  pdu [command]

Available Commands:
  auto        Auto Initialize PDU
  ck          Create keystores
  help        Help about any command
  init        Initialize PDU
  send        Send hello msg to node
  start       Start run node
  test        Test some tmp func
  upload      Upload file to node

Flags:
      --db string            path of database (default "pdu.db")
  -h, --help                 help for pdu
      --nodes string         node list
      --port int             port to start server or send msg (default 1323)
      --projectPath string   project root path (default "./")
      --url string           target url (default "http://127.0.0.1")

Use "pdu [command] --help" for more information about a command.
```



## Development

To copy the repository:

```
go get github.com/pdupub/go-pdu

```
OR 
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


<a href="https://pdu.pub"><img height="32" align="right" src="https://pdu.pub/images/icon.svg"></a>


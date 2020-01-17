go-pdu
====
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/pdupub/go-pdu)
[![GoReport](https://goreportcard.com/badge/github.com/pdupub/go-pdu)](https://goreportcard.com/report/github.com/pdupub/go-pdu)
[![Travis](https://travis-ci.org/pdupub/go-pdu.svg?branch=master)](https://travis-ci.org/pdupub/go-pdu)
[![License](https://img.shields.io/badge/license-GPL%20v3-blue.svg)](LICENSE)
[![Chat](https://img.shields.io/badge/gitter-Docs%20chat-4AB495.svg)](https://gitter.im/pdupub/go-pdu)
[![Coverage Status](https://coveralls.io/repos/github/pdupub/go-pdu/badge.svg?branch=master)](https://coveralls.io/github/pdupub/go-pdu?branch=master)

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
  account     Account generate or inspect
  create      Create a new PDU Universe
  help        Help about any command
  start       Start to run PDU Universe

Flags:
  -c, --config string   config file
  -h, --help            help for pdu

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

To build and run:
```
make install && pdu start
```


## Contributing

1. Fork the repository on GitHub to start making your changes to the master branch
2. Write a test which shows that the bug was fixed or that the feature works as expected
3. Send a pull request and bug the maintainer until it gets merged and published


<a href="https://pdu.pub"><img height="32" align="right" src="https://pdu.pub/images/icon.svg"></a>


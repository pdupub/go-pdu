# PDU:一种点对点的社交网络
Parallel Digital Universe - A Peer-to-Peer Social Network

email: liupeng@tataufo.com

[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/TATAUFO/PDU)
[![GoReport](https://goreportcard.com/badge/github.com/TATAUFO/PDU)](https://goreportcard.com/report/github.com/TATAUFO/PDU)
[![Travis](https://travis-ci.org/TATAUFO/PDU.svg?branch=master)](https://travis-ci.org/TATAUFO/PDU)
[![License](https://img.shields.io/badge/license-GPL%20v3-blue.svg)](LICENSE)


## Abstract

SNS，即社交网络服务，如Facebook，用户可以在其上创建身份，维护社交关系并进行信息传播，交互。但现有的SNS均依赖于某个第三方的网络服务，而这个第三方（如Facebook）则越来越不被信任。BitTorrent协议，能够实现P2P的信息传播，但其根本目的是提高对于已知内容的传播效率，缺乏在P2P网络下的账号系统，所以无法对未知内容有所判断。即便有数字签名，能够证明每个信息的来源，但是因为缺少第三方验证（如手机号注册），创建账号的成本为零，所以当无用（虚假）的信息会充斥整个网络且无法信息来源进行惩罚。

我们提出一种在纯粹P2P的环境下增加创建成本的方式，并基于这种账户系统，构建完整的P2P社交网络形态。首先我们引入的时间证明，用以证明某个特定行为发生于某时刻之后。新账号的创建必须由多个（通常2个）合法账号签名，同一账号的此类签名操作需满足时间间隔。每个网络的参与者（用户），都在本地以DAG的结构维护所有账号之间的关系拓扑，并随时可以根据自己获知的新消息，对新的账号进行验证增补，同时也可因作恶行为对某些账号及关联账号进行惩罚。

与比特币的共识不同，在PDU中你只相信你所相信的。



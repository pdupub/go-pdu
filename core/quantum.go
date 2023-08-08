// Copyright 2021 The PDU Authors
// This file is part of the PDU library.
//
// The PDU library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The PDU library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the PDU library. If not, see <http://www.gnu.org/licenses/>.

package core

// Quantum is the only structure for the information released by the publisher in the PDU,
// and it is also the only structure for information exchange between nodes in this P2P system.
// Each Quantum consists of three parts: signature, reference list and data.
// Signatures are used to determine the source of information and integrity.
// References are used to determine the order of the information, the data content is composed of QContent.
// Quantum是PDU中发布者所发布信息的唯一结构，也是这个P2P的系统中，各节点之间信息交流的唯一结构。
// 每个Quantum都由三部分构成：签名，引用和数据内容。签名用以确定信息源即完整性，引用用来确定信息
// 间的顺序，数据内容由QContent数组构成。

import (
	"encoding/json"
	"errors"

	"github.com/pdupub/go-pdu/identity"
)

const (
	maxReferencesCnt = 16
	maxContentsCnt   = 16
	maxContentSize   = 1024 * 512
)

var (
	errQuantumRefsCntOutOfLimit     = errors.New("quantum references count out of limit")
	errQuantumContentsCntOutOfLimit = errors.New("quantum contents count out of limit")
	errQuantumContentSizeOutOfLimit = errors.New("quantum content size out of limit")
	errQuantumSignedAlready         = errors.New("quantum have been signed already")
	errQuantumTypeNotFit            = errors.New("quantum type is not fit")
)

type Sig []byte

const (
	// QuantumTypeInformation specifies quantum which user just want to share information, such as
	// post articles, reply to others or chat message with others (encrypted by receiver's public key)
	// QuantumTypeInformation 表示这个Quantum包含的内容只应被作为信息处理，如发布文章，图片等媒体内容或对于他人的回复，
	// 不同于后面所介绍的四种类型，此类型没有任何会被系统默认处理的特殊功能。
	QuantumTypeInformation = 0

	// QuantumTypeIntegration specifies the quantum to update user's (signer's) profile.
	// contents = [key1(QCFmtStringTEXT), value1, key2, value2 ...]
	// if user want to update key1, send new QuantumTypeIntegration quantum with
	// contents = [key1, newValue ...]
	// if user want to delete key1, send new QuantumTypeIntegration quantum with
	// contents = [key1, newValue which content with empty data]
	// QuantumTypeIntegration 表示这个Quantum为整合类型，此类型的quantum主要用于维护信息发布者的公开字典，可用户维护
	// 信息发布者的profile或者作为账本维护每个地址的资金情况。QContent数组是由key和value循环组成，如果这个数组的长度为
	// 奇数则忽略最后一个key
	QuantumTypeIntegration = 1

	// QuantumTypeSpeciation specifies the quantum of rule to build new species
	// contents[0] is the display information of current species
	// {fmt:QCFmtStringJSON/QCFmtStringTEXT..., data: ...}
	// contents[1] is the number of identification (co-signature) from users in current species
	// {fmt:QCFmtStringInt, data:1} at least 1. (creater is absolutely same with others)
	// contents[2] is the max number of identification by one user
	// {fmt:QCFmtStringInt, data:1} -1 means no limit, 0 means not allowed
	// contents[3] ~ contents[15] is the initial users in this species
	// {fmt:QCFmtBytesAddress/QCFmtStringAddressHex, data:0x1232...}
	// signer of this species is also the initial user in this species
	// QuantumTypeSpeciation特指创建新的族群，QContent数组中的第一项是这个族群的说明，二三两项为认定的限制条件
	// 从第四项开始是初始地址，族群默认包含创建者属于这个族群且创建者无需将自己的地址写入到族群创建的quantum当中。
	// 需要注意：第二三两项的认定限制条件有可能在后续的开发中消失，转而将族群的判定条件完全交由信息使用者确定
	QuantumTypeSpeciation = 2

	// QuantumTypeIdentification specifies the quantum of identification
	// contents[0] is the signature of target species rule quantum
	// {fmt:QCFmtBytesSignature/QCFmtStringSignatureHex, data:signature of target species rule quantum}
	// contents[1] ~ contents[n] is the address of be identified
	// {fmt:QCFmtBytesAddress/QCFmtStringAddressHex, data:0x123...}
	// no matter all identification send in same quantum of different quantum
	// only first n address (rule.contents[2]) will be accepted
	// user can not quit species, but any user can block any other user (or self) from any species
	// accepted by any species is decided by user in that species feel about u, not opposite.
	// User belong to species, quantum belong to user. (On trandation forum, posts usually belong to
	// one topic and have lots of tag, is just the function easy to implemnt not base struct here)
	// QuantumTypeIdentification 特指表达对于其他信息发布者属于某个特定族群的认定类Quantum。
	// 需要注意：后续的认定消息发布不仅限于属于某特定族群的地址，而可以由任何地址发布，因为判定的某地址是否属于某族群的权利
	// 将完全交于信息使用者，此改变和QuantumTypeSpeciation将会同步进行。
	QuantumTypeIdentification = 3

	// QuantumTypeTermination specifies the quantum of ending life cycle
	// When receive this quantum, no more quantum from same address should be received or broadcast.
	// The decision of end from any identity should be respected, no matter for security reason or just want to leave.
	// QuantumTypeTermination 特指废弃这个地址的quantum，表达发布者希望系统不在接受和处理由这个私钥签发的后续任何信息。
	QuantumTypeTermination = 4
)

var (
	FirstQuantumReference = Hex2Sig("0x00")
	SameQuantumReference  = Hex2Sig("0x01")
)

// UnsignedQuantum defines the single message from user without signature,
// all variables should be in alphabetical order.
type UnsignedQuantum struct {
	// Contents contain all data in this quantum
	Contents []*QContent `json:"cs,omitempty"`

	// References must from the exist signature or 0x00
	// References[0] is the last signature by user-self, 0x00 if this quantum is the first quantum by user
	// References[1] is the last signature of public or important quantum by user-self, which quantum user
	// want other user to view most. (refs[1] usually contains article and forward, but not comments, like or chat)
	// if References[1] same with References[0], References[1] can be SameQuantumReference(0x01)
	// References[2] ~ References[n] is optional, recommend to use new & valid quantum
	// If two quantums by same user with same References[0], these two quantums will cause conflict, and
	// this user maybe block by others. The reason to do that punishment is user should act like individual,
	// all proactive event from one user should be sequence should be total order (全序关系). References[1~n]
	// do not need follow this restriction, because all other references show the partial order (偏序关系).
	// all the quantums which be reference as References[1] should be total order too.
	References []Sig `json:"refs"`

	// Type specifies the type of this quantum
	Type int `json:"type"`
}

// Quantum defines the single message signed by user.
type Quantum struct {
	UnsignedQuantum
	Signature Sig `json:"sig,omitempty"`
}

// NewQuantum try to build Quantum without signature
func NewQuantum(t int, cs []*QContent, refs ...Sig) (*Quantum, error) {
	if len(cs) > maxContentsCnt {
		return nil, errQuantumContentsCntOutOfLimit
	}
	if len(refs) > maxReferencesCnt || len(refs) < 1 {
		return nil, errQuantumRefsCntOutOfLimit
	}
	for _, v := range cs {
		if len(v.Data) > maxContentSize {
			return nil, errQuantumContentSizeOutOfLimit
		}
		// TODO: check fmt of each content
	}

	uq := UnsignedQuantum{
		Contents:   cs,
		References: refs,
		Type:       t}

	return &Quantum{UnsignedQuantum: uq}, nil
}

// Sign try to add signature to Quantum
func (q *Quantum) Sign(did *identity.DID) error {
	if q.Signature != nil {
		return errQuantumSignedAlready
	}

	b, err := json.Marshal(q.UnsignedQuantum)
	if err != nil {
		return err
	}
	sig, err := did.Sign(b)
	if err != nil {
		return err
	}
	q.Signature = sig
	return nil
}

// Ecrecover recover
func (q *Quantum) Ecrecover() (identity.Address, error) {
	b, err := json.Marshal(q.UnsignedQuantum)
	if err != nil {
		return identity.Address{}, err
	}
	return identity.Ecrecover(b, q.Signature)
}

func CreateInfoQuantum(qcs []*QContent, refs ...Sig) (*Quantum, error) {
	return NewQuantum(QuantumTypeInformation, qcs, refs...)
}

func CreateIntegrationQuantum(data map[string]interface{}, refs ...Sig) (*Quantum, error) {
	var qcs []*QContent
	for k, v := range data {
		key := CreateTextContent(k)
		var value *QContent
		switch v := v.(type) {
		case string:
			value = CreateTextContent(v)
		case float64:
			value = CreateFloatContent(v)
		case int:
			value = CreateIntContent(int64(v))
		case int64:
			value = CreateIntContent(v)
		default:
			// img for avator will be deal later
			return nil, errContentFmtNotFit
		}
		qcs = append(qcs, key, value)
	}

	return NewQuantum(QuantumTypeIntegration, qcs, refs...)
}

func CreateSpeciesQuantum(note string, minCosignCnt int, maxIdentifyCnt int, initAddrsHex []string, refs ...Sig) (*Quantum, error) {
	var qcs []*QContent

	qcs = append(qcs, CreateTextContent(note))
	qcs = append(qcs, CreateIntContent(int64(minCosignCnt)))
	qcs = append(qcs, CreateIntContent(int64(maxIdentifyCnt)))
	for _, addrHex := range initAddrsHex {
		qcs = append(qcs, &QContent{Format: QCFmtStringAddressHex, Data: []byte(addrHex)})
	}

	return NewQuantum(QuantumTypeSpeciation, qcs, refs...)
}

func CreateIdentificationQuantum(target Sig, addrsHex []string, refs ...Sig) (*Quantum, error) {
	var qcs []*QContent
	qcs = append(qcs, &QContent{Format: QCFmtBytesSignature, Data: target})
	for _, addrHex := range addrsHex {
		qcs = append(qcs, &QContent{Format: QCFmtStringAddressHex, Data: []byte(addrHex)})
	}
	return NewQuantum(QuantumTypeIdentification, qcs, refs...)
}

func CreateTerminationQuantum(refs ...Sig) (*Quantum, error) {
	return NewQuantum(QuantumTypeTermination, []*QContent{}, refs...)
}

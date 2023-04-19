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

package fb

import (
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"context"

	firebase "firebase.google.com/go"
	// "firebase.google.com/go/auth"
	"cloud.google.com/go/firestore"

	"github.com/pdupub/go-pdu/core"
	"github.com/pdupub/go-pdu/identity"
	"github.com/pdupub/go-pdu/params"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

var testKeyJSON = "../../" + params.TestFirebaseAdminSDKPath
var testProjectID = params.TestFirebaseProjectID

const (
	collectionQuantum    = "quantum"
	collectionUniverse   = "universe"
	collectionCommunity  = "community"
	collectionIndividual = "individual"
)

func testCreateEndQuantum(t *testing.T, ctx context.Context, client *firestore.Client, did *identity.DID,
	refs ...core.Sig) (*core.Quantum, *firestore.DocumentRef) {

	q, err := core.CreateEndQuantum(refs...)
	if err != nil {
		t.Error(err)
	}
	q.Sign(did)
	return testUploadQuantum(t, ctx, client, q)
}

func testCreateInviteQuantum(t *testing.T, ctx context.Context, client *firestore.Client, did *identity.DID,
	target core.Sig, addrsHex []string,
	refs ...core.Sig) (*core.Quantum, *firestore.DocumentRef) {
	q, err := core.CreateInvitationQuantum(target, addrsHex, refs...)
	if err != nil {
		t.Error(err)
	}
	q.Sign(did)

	return testUploadQuantum(t, ctx, client, q)
}

func testCreateCommunityQuantum(t *testing.T, ctx context.Context, client *firestore.Client, did *identity.DID,
	note string, minCosignCnt int, maxInviteCnt int, initAddrsHex []string,
	refs ...core.Sig) (*core.Quantum, *firestore.DocumentRef) {

	q, err := core.CreateCommunityQuantum(note, minCosignCnt, maxInviteCnt, initAddrsHex, refs...)
	if err != nil {
		t.Error(err)
	}
	q.Sign(did)
	return testUploadQuantum(t, ctx, client, q)
}

func testCreateProfileQuantum(t *testing.T, ctx context.Context, client *firestore.Client, did *identity.DID, profiles map[string]interface{}, refs ...core.Sig) (*core.Quantum, *firestore.DocumentRef) {
	q, err := core.CreateProfileQuantum(profiles, refs...)
	if err != nil {
		t.Error(err)
	}
	q.Sign(did)
	return testUploadQuantum(t, ctx, client, q)
}

func testCreateInfoQuantum(t *testing.T, ctx context.Context, client *firestore.Client, did *identity.DID, qcs []*core.QContent, refs ...core.Sig) (*core.Quantum, *firestore.DocumentRef) {
	q, err := core.CreateInfoQuantum(qcs, refs...)
	if err != nil {
		t.Error(err)
	}
	q.Sign(did)
	return testUploadQuantum(t, ctx, client, q)
}

func testClearQuantum(t *testing.T) {

	ctx := context.Background()
	opt := option.WithCredentialsFile(testKeyJSON)
	config := &firebase.Config{ProjectID: testProjectID}
	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		t.Error(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		t.Error(err)
	}
	defer client.Close()

	testCollection := client.Collection(collectionQuantum)
	docRefs, err := testCollection.DocumentRefs(ctx).GetAll()
	if err != nil {
		t.Error(err)
		return
	}
	for _, docRef := range docRefs {
		docRef.Delete(ctx)
	}

	individualCollection := client.Collection(collectionIndividual)
	docRefs, err = individualCollection.DocumentRefs(ctx).GetAll()
	if err != nil {
		t.Error(err)
		return
	}
	for _, docRef := range docRefs {
		docRef.Delete(ctx)
	}

	communityCollection := client.Collection(collectionCommunity)
	docRefs, err = communityCollection.DocumentRefs(ctx).GetAll()
	if err != nil {
		t.Error(err)
		return
	}
	for _, docRef := range docRefs {
		docRef.Delete(ctx)
	}

	configCollection := client.Collection("universe")
	configDocRef := configCollection.Doc("status")
	configMap := make(map[string]interface{})
	configMap["lastSequence"] = 0
	configMap["lastSigHex"] = ""
	configMap["updateTime"] = time.Now().UnixMilli()
	configDocRef.Set(ctx, configMap, firestore.Merge([]string{"lastSequence"}, []string{"lastSigHex"}, []string{"updateTime"}))
}

func testUploadQuantum(t *testing.T, ctx context.Context, client *firestore.Client, q *core.Quantum) (*core.Quantum, *firestore.DocumentRef) {
	testCollection := client.Collection(collectionQuantum)

	docID := core.Sig2Hex(q.Signature)

	docRef := testCollection.Doc(docID)

	dMap := make(map[string]interface{})

	qBytes, err := json.Marshal(q)
	if err != nil {
		t.Error(err)
	}

	dMap["recv"] = qBytes

	_, err = docRef.Set(ctx, dMap)
	if err != nil {
		t.Error(err)
	}

	return q, docRef
}

func testUploadQuantumBack(t *testing.T, ctx context.Context, client *firestore.Client, q *core.Quantum) (*core.Quantum, *firestore.DocumentRef) {

	testCollection := client.Collection(collectionQuantum)

	fbq, _ := NewFBQuantum(q)
	dMap, _ := FBStruct2Data(fbq)

	docID := fbq.SigHex

	// add document
	docRef := testCollection.Doc(docID)
	_, err := docRef.Set(ctx, dMap)

	if err != nil {
		t.Error(err)
	}
	// get doc snapshot
	docSnapshot, err := docRef.Get(ctx)
	if err != nil {
		t.Error(err)
	}

	// get data of snapshot
	mapRes := docSnapshot.Data()
	fbqRes, err := Data2FBQuantum(mapRes)

	if err != nil {
		t.Error(err)
	}

	qRes, err := fbqRes.GetOriginQuantum()
	if err != nil {
		t.Error(err)
	}

	addr, err := qRes.Ecrecover()
	if err != nil {
		t.Error(err)
	}
	t.Log("addr", addr.Hex())

	return q, docRef
}

func testCreateQuantums(t *testing.T) {
	ctx := context.Background()
	opt := option.WithCredentialsFile(testKeyJSON)
	config := &firebase.Config{ProjectID: testProjectID}
	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		t.Error(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		t.Error(err)
	}
	defer client.Close()

	did1, _ := identity.New()
	did1.UnlockWallet("../../"+params.TestKeystore(0), params.TestPassword)
	did2, _ := identity.New()
	did2.UnlockWallet("../../"+params.TestKeystore(1), params.TestPassword)
	did3, _ := identity.New()
	did3.UnlockWallet("../../"+params.TestKeystore(2), params.TestPassword)
	did4, _ := identity.New()
	did4.UnlockWallet("../../"+params.TestKeystore(3), params.TestPassword)

	ref1 := []core.Sig{}
	ref2 := []core.Sig{}
	ref3 := []core.Sig{}
	// ref4 := []core.Sig{}

	c0 := core.CreateTextContent("Hello! ")
	c1 := core.CreateIntContent(100)
	c2 := core.CreateTextContent(">")
	c3 := core.CreateFloatContent(99.9)

	q1, _ := testCreateInfoQuantum(t, ctx, client, did1, []*core.QContent{c0}, core.FirstQuantumReference)
	ref1 = append(ref1, q1.Signature)

	q2, _ := testCreateInfoQuantum(t, ctx, client, did2, []*core.QContent{c1, c2, c3}, core.FirstQuantumReference, q1.Signature)
	ref2 = append(ref2, q2.Signature)

	q3, _ := testCreateInfoQuantum(t, ctx, client, did3, []*core.QContent{c1, c3, c0}, core.FirstQuantumReference, q1.Signature, q2.Signature)
	ref3 = append(ref3, q3.Signature)

	for i := 2; i < 6; i++ {
		q11, _ := testCreateInfoQuantum(t, ctx, client, did1, []*core.QContent{core.CreateTextContent("Hello! A " + strconv.Itoa(i))}, ref1[len(ref1)-1], ref2[len(ref2)-1])
		ref1 = append(ref1, q11.Signature)
		q22, _ := testCreateInfoQuantum(t, ctx, client, did2, []*core.QContent{core.CreateTextContent("Hello! B " + strconv.Itoa(i))}, ref2[len(ref2)-1], ref3[len(ref3)-1], ref1[len(ref1)-1])
		ref2 = append(ref2, q22.Signature)
		q33, _ := testCreateInfoQuantum(t, ctx, client, did3, []*core.QContent{core.CreateTextContent("Hello! C " + strconv.Itoa(i))}, ref3[len(ref3)-1], ref2[len(ref2)-1])
		ref3 = append(ref3, q33.Signature)
	}

	profile1 := make(map[string]interface{})
	profile1["name"] = "hello"
	profile1["age"] = 100
	profile1["temp"] = 12.3

	profile2 := make(map[string]interface{})
	profile2["name"] = "hahaha AAA"
	profile2["city"] = "BeiJing"

	profile22 := make(map[string]interface{})
	profile22["city"] = "NONON"
	profile22["temp"] = 12.3

	q4, _ := testCreateProfileQuantum(t, ctx, client, did1, profile1, ref1[len(ref1)-1], ref2[len(ref2)-1], q2.Signature)
	ref1 = append(ref1, q4.Signature)

	q5, _ := testCreateProfileQuantum(t, ctx, client, did2, profile2, ref2[len(ref2)-1], ref3[len(ref3)-1], q2.Signature)
	ref2 = append(ref2, q5.Signature)

	q6, _ := testCreateCommunityQuantum(t, ctx, client, did3, "Tody is Great", 2, 3, []string{did1.GetAddress().Hex(), did2.GetAddress().Hex()}, ref3[len(ref3)-1], q5.Signature)
	ref3 = append(ref3, q6.Signature)

	q7, _ := testCreateInviteQuantum(t, ctx, client, did3, q6.Signature, []string{did4.GetAddress().Hex()}, ref3[len(ref3)-1], q6.Signature, q5.Signature)
	ref3 = append(ref3, q7.Signature)

	q8, _ := testCreateEndQuantum(t, ctx, client, did2, ref2[len(ref2)-1], q6.Signature, q7.Signature)
	ref2 = append(ref2, q8.Signature)

	q9, _ := testCreateProfileQuantum(t, ctx, client, did2, profile22, ref2[len(ref2)-1], ref3[len(ref3)-1], q8.Signature)
	ref2 = append(ref2, q9.Signature)

	t.Log("did1 last sig: ", core.Sig2Hex(ref1[len(ref1)-1]))
	t.Log("did2 last sig: ", core.Sig2Hex(ref2[len(ref2)-1]))
	t.Log("did3 last sig: ", core.Sig2Hex(ref3[len(ref3)-1]))

	t.Log("community sig: ", core.Sig2Hex(q6.Signature))

}

func testManualInviteQuantums(t *testing.T) {
	ctx := context.Background()
	opt := option.WithCredentialsFile(testKeyJSON)
	config := &firebase.Config{ProjectID: testProjectID}
	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		t.Error(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		t.Error(err)
	}
	defer client.Close()

	did1, _ := identity.New()
	did1.UnlockWallet("../../"+params.TestKeystore(0), params.TestPassword)
	did2, _ := identity.New()
	did2.UnlockWallet("../../"+params.TestKeystore(1), params.TestPassword)
	did3, _ := identity.New()
	did3.UnlockWallet("../../"+params.TestKeystore(2), params.TestPassword)
	did4, _ := identity.New()
	did4.UnlockWallet("../../"+params.TestKeystore(3), params.TestPassword)

	// did1, did2, did3 exist in communtiy
	extraInviteAddrHexList := []string{"0x00008Bd373Ac9f168f087E976d0068732fbD6835"}
	communityDefineSig := core.Hex2Sig("0x4ac617c2ead08dd4ae046c048200e74c64d7b93f22f12c78d3f85b539afb94982157837dae9e4720193f64b480509a7815cda6ca356025555c1df631d995a62e01")
	did1SelfRef := core.Hex2Sig("0x5c9110f071ee589a918f26b988dba3b990300e372392ceb83e36515d8beeef4a7bf5838123dc0b1f0717ed8b68e7be270538a954f5b17529bae44c412714d0f000")
	did3SelfRef := core.Hex2Sig("0xadd9248fdeeb795bf0b6758418f28570075547019324eac7fb0d4e686709106a158a0929956bac853f1e70bf1de28e03c929e46514c6d0a85927fb6ec8d8948000")

	q1, _ := testCreateInviteQuantum(t, ctx, client, did1, communityDefineSig, extraInviteAddrHexList, did1SelfRef, did3SelfRef)
	q2, _ := testCreateInviteQuantum(t, ctx, client, did3, communityDefineSig, extraInviteAddrHexList, did3SelfRef, did1SelfRef)

	t.Log("did1 last sig: ", core.Sig2Hex(q1.Signature))
	t.Log("did3 last sig: ", core.Sig2Hex(q2.Signature))

}

func testCustomQuantum(t *testing.T) {
	ctx := context.Background()
	fbu, err := NewFBUniverse(ctx, testKeyJSON, testProjectID)
	if err != nil {
		t.Error(err)
		return
	}
	sigHex := "0x0d6c6fbeb40c817e74c4890d2652f7edf610274a58c399bbb6eadf671b80ea430fafdd1a38bf5880e1e2c44dad7e8a4a62d6ae77e5daac061e956e16c29e2ccc01"
	field := "origin"
	if _, _, _, err := fbu.customProcess(sigHex, field); err != nil {
		t.Error(err)
	}
}

func testDealQuantums(t *testing.T) {
	ctx := context.Background()
	fbu, err := NewFBUniverse(ctx, testKeyJSON, testProjectID)
	if err != nil {
		t.Error(err)
		return
	}
	if _, _, _, err := fbu.ProcessQuantums(100, 0); err != nil {
		t.Error(err)
	}
	if err := fbu.Close(); err != nil {
		t.Error(err)
	}
}

func testGetQuantums(t *testing.T) {
	didMap := make(map[string]*identity.DID)

	for i := 0; i < 4; i++ {
		did, _ := identity.New()
		did.UnlockWallet("../../"+params.TestKeystore(i), params.TestPassword)
		didMap[did.GetAddress().Hex()] = did
	}

	ctx := context.Background()
	fbu, err := NewFBUniverse(ctx, testKeyJSON, testProjectID)
	if err != nil {
		t.Error(err)
	}
	if qs, err := fbu.QueryQuantums(identity.Address{}, 0, 0, 100, false); err != nil {
		t.Error(err)
	} else {
		var sigs []string
		for _, q := range qs {
			resMap := make(map[string]interface{})
			resMap["type"] = q.Type
			resMap["signature"] = core.Sig2Hex(q.Signature)
			var refs []string
			for _, ref := range q.References {
				refs = append(refs, core.Sig2Hex(ref))
			}
			resMap["refs"] = refs
			resMap["cs"], err = CS2Readable(q.Contents)
			if err != nil {
				t.Error(err)
			}
			id, err := q.Ecrecover()
			if err != nil {
				t.Error(err)
			}
			resMap["address"] = id.Hex()
			if _, ok := didMap[id.Hex()]; !ok {
				t.Error("signer not exist")
			}

			b, err := json.Marshal(resMap)
			if err != nil {
				t.Error(err)
			}
			t.Log(string(b))
			sigs = append(sigs, core.Sig2Hex(q.Signature))
		}
		for _, sig := range sigs {
			t.Log(sig)
		}
	}

	if err := fbu.Close(); err != nil {
		t.Error(err)
	}
}

func testTemp(t *testing.T) {
	t.Log("tmp test")
	ctx := context.Background()
	var err error
	fbu, err := NewFBUniverse(ctx, testKeyJSON, testProjectID)
	if err != nil {
		t.Error(err)
	}
	snapshots, err := fbu.individualC.Where("rp.email.data", "==", "Hi").Documents(fbu.ctx).GetAll()
	if err != nil {
		t.Error(err)
	}
	for _, snap := range snapshots {
		fbi, err := Data2FBIndividual(snap.Data())
		if err != nil {
			t.Error(err)
		}
		t.Log(fbi.AddrHex)
	}
}

func testCommunityInfo(t *testing.T) {
	t.Log("community test")
	ctx := context.Background()
	var fbu core.Universe
	var err error
	fbu, err = NewFBUniverse(ctx, testKeyJSON, testProjectID)
	if err != nil {
		t.Error(err)
	}
	// quantums, err := fbu.QueryQuantums(identity.Address{}, 0, 0, 3, false)
	// if err != nil {
	// 	t.Error(err)
	// }
	// for _, q := range quantums {
	// 	t.Log(q)
	// }
	comSig := core.Hex2Sig("0x4ac617c2ead08dd4ae046c048200e74c64d7b93f22f12c78d3f85b539afb94982157837dae9e4720193f64b480509a7815cda6ca356025555c1df631d995a62e01")
	// comSig := core.Hex2Sig("0xe2f00788cd5a4ba91c6c0a1e6e0944da1631a236ed7df74b83521dba9397dcd44f4c46792f0235099cc7352e33d82be29da1b08f23e25276a7095100862e2ac601")
	// community, err := fbu.GetCommunity(comSig)
	// if err != nil {
	// 	t.Error(err)
	// }
	// t.Log(community)

	individuals, err := fbu.QueryIndividuals(comSig, 1, 4, false)
	if err != nil {
		t.Error(err)
	}
	for _, ind := range individuals {
		t.Log(ind.Address.Hex())
	}
}

func testCheckQuantum(t *testing.T) {

	ctx := context.Background()
	opt := option.WithCredentialsFile(testKeyJSON)
	config := &firebase.Config{ProjectID: testProjectID}
	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		t.Error(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		t.Error(err)
	}
	defer client.Close()

	testCollection := client.Collection(collectionQuantum)

	// checkKey := "recv"
	checkKey := "origin"
	iter := testCollection.Where(checkKey, "!=", []byte{}).Documents(ctx)
	for docSnapshot, err := iter.Next(); err != iterator.Done; docSnapshot, err = iter.Next() {
		if qBytes, ok := docSnapshot.Data()[checkKey]; ok {
			t.Log(docSnapshot.Ref.ID)
			// t.Log(qBytes)
			q := core.Quantum{}
			json.Unmarshal(qBytes.([]byte), &q)
			addr, err := q.Ecrecover()
			if err != nil {
				t.Error(err)
			}
			t.Log("Address", addr.Hex())
			for _, ref := range q.References {
				t.Log("Ref", core.Sig2Hex(ref))
			}
		}
	}
}

func testProcessOriginQuantum(t *testing.T) {
	ctx := context.Background()
	opt := option.WithCredentialsFile(testKeyJSON)
	config := &firebase.Config{ProjectID: testProjectID}
	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		t.Error(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		t.Error(err)
	}
	defer client.Close()

	testCollection := client.Collection(collectionQuantum)
	docRefs, err := testCollection.DocumentRefs(ctx).GetAll()
	if err != nil {
		t.Error(err)
		return
	}
	for _, docRef := range docRefs {
		t.Log(docRef.ID)
		snap, err := docRef.Get(ctx)
		if err != nil {
			t.Error(err)
			continue
		}
		if v, ok := snap.Data()["origin"]; ok {
			value := v.([]byte)
			var q core.Quantum
			if err := json.Unmarshal(value, &q); err != nil {
				t.Error(err)
			} else {
				t.Log("origin quantums", core.Sig2Hex(q.Signature))
				addr, err := q.Ecrecover()
				if err != nil {
					t.Error(err)
				} else {
					t.Log("address", addr.Hex())
				}
			}

		}
	}
	t.Log("$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$")
	iter := testCollection.Where("origin", "!=", "").Documents(ctx)
	for snap, err := iter.Next(); err != iterator.Done; snap, err = iter.Next() {
		if snap != nil {
			if _, ok := snap.Data()["origin"]; ok {
				t.Log(snap.Ref.ID)
			}
		}
	}
}

func testShowPrivateKey(t *testing.T) {
	t.Log("address / privateKey ")
	t.Log("======================")

	for i := 0; i < 4; i++ {
		did, _ := identity.New()
		did.UnlockWallet("../../"+params.TestKeystore(i), params.TestPassword)
		addr, _, priv, _ := did.Inspect(true)
		t.Log("address\t", addr)
		// t.Log("pubKey\t", pub)
		t.Log("privKey\t", priv)
		t.Log("======================")
	}
}

func TestMain(t *testing.T) {
	// testClearQuantum(t)
	// testCreateQuantums(t)
	// testManualInviteQuantums(t)
	testDealQuantums(t)
	// testCustomQuantum(t)
	// testCheckQuantum(t)
	// testGetQuantums(t)
	// testCommunityInfo(t)
	// testTemp(t)
	// testShowPrivateKey(t)
	// testProcessOriginQuantum(t)
}

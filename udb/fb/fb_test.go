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
	"strconv"
	"testing"

	"context"

	firebase "firebase.google.com/go"
	// "firebase.google.com/go/auth"
	"cloud.google.com/go/firestore"

	"github.com/pdupub/go-pdu/core"
	"github.com/pdupub/go-pdu/identity"
	"github.com/pdupub/go-pdu/params"
	"google.golang.org/api/option"
)

var testKeyJSON = "./firebase-adminsdk.json"
var testProjectID = "tweetsample-201fd"

var clearBeforeTest = true

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

func testClearQuantum(t *testing.T, ctx context.Context, client *firestore.Client) {
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

	configCollection := client.Collection("config")
	configDocRef := configCollection.Doc("system")
	configMap := make(map[string]int64)
	configMap["sequence"] = 0
	configDocRef.Set(ctx, configMap, firestore.Merge([]string{"sequence"}))
}

func testUploadQuantum(t *testing.T, ctx context.Context, client *firestore.Client, q *core.Quantum) (*core.Quantum, *firestore.DocumentRef) {

	testCollection := client.Collection(collectionQuantum)

	docID, fbq := Quantum2FBQuantum(q)
	dMap, _ := FBStruct2Data(fbq)

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

	qRes, err := FBQuantum2Quantum(docRef.ID, fbqRes)
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

	if clearBeforeTest {
		testClearQuantum(t, ctx, client)
	}

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

	c0 := core.NewTextC("Hello! ")
	c1 := core.NewIntC(100)
	c2 := core.NewTextC(">")
	c3 := core.NewFloatC(99.9)

	q1, _ := testCreateInfoQuantum(t, ctx, client, did1, []*core.QContent{c0}, core.FirstQuantumReference)
	ref1 = append(ref1, q1.Signature)

	q2, _ := testCreateInfoQuantum(t, ctx, client, did2, []*core.QContent{c1, c2, c3}, core.FirstQuantumReference, q1.Signature)
	ref2 = append(ref2, q2.Signature)

	q3, _ := testCreateInfoQuantum(t, ctx, client, did3, []*core.QContent{c1, c3, c0}, core.FirstQuantumReference, q1.Signature, q2.Signature)
	ref3 = append(ref3, q3.Signature)

	for i := 2; i < 6; i++ {
		q11, _ := testCreateInfoQuantum(t, ctx, client, did1, []*core.QContent{core.NewTextC("Hello! A " + strconv.Itoa(i))}, ref1[len(ref1)-1], ref2[len(ref2)-1])
		ref1 = append(ref1, q11.Signature)
		q22, _ := testCreateInfoQuantum(t, ctx, client, did2, []*core.QContent{core.NewTextC("Hello! B " + strconv.Itoa(i))}, ref2[len(ref2)-1], ref3[len(ref3)-1], ref1[len(ref1)-1])
		ref2 = append(ref2, q22.Signature)
		q33, _ := testCreateInfoQuantum(t, ctx, client, did3, []*core.QContent{core.NewTextC("Hello! C " + strconv.Itoa(i))}, ref3[len(ref3)-1], ref2[len(ref2)-1])
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

	t.Log(len(ref1), len(ref2), len(ref3))
}

func testDealQuantums(t *testing.T) {
	ctx := context.Background()
	fbs, err := NewFBS(ctx, testKeyJSON, testProjectID)
	if err != nil {
		t.Error(err)
	}
	if err := fbs.DealNewQuantums(); err != nil {
		t.Error(err)
	}
	if err := fbs.Close(); err != nil {
		t.Error(err)
	}
}

func TestMain(t *testing.T) {
	testCreateQuantums(t)
	testDealQuantums(t)
}

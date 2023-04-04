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

package params

import (
	"fmt"
)

// test for udb/firebase

const (
	TestFirebaseAdminSDKPath = "udb/fb/test-firebase-adminsdk.json"
	TestFirebaseProjectID    = "pdu-dev-1" // "pdupub-a2bdd"
)

// test for identity
const (
	TestPassword = "1"
)

func TestKeystore(i int) string {
	return fmt.Sprintf("identity/testdata/keyfile_%d.json", i)
}

func TestAddrsHex() []string {
	return []string{"0xaf040ed5498f9808550402ebb6c193e2a73b860a",
		"0x4c62002051d76ba82e43e1a5785590dfca261c10",
		"0x50cc750695f0701e193b95f833ab04256d3d9daa",
		"0xcab7ecfd63f93d405656633db6d0d8a236b74aff",
		"0x7402c45f230361e2e397dfdf15e264d55ac56da1",
		"0x9a27f8de1c88fa4c76d7322956de736520477b2c",
		"0x608302944ed7d3e0f136bede44b13992b26a9a86",
		"0xb85caf6c2648985657366e9ad37257cdc5f5a238",
		"0x646ec7961710cdbb554dc442fc10d108eafa0959",
		"0xf2f16ad0af506a2525fe924d4b09a108c2999c6a",
		"0xfb113012790c4b151319b41c756b42190bb96ec1",
		"0x1a3d54a80db0a80a6226e7d3088e639edf56525d",
		"0x587bf58f61462c4e9b8dd8355336ac75052dec1e",
		"0x1b2181ae4012faa9518ac3bc2c98e6b278c7ab86",
		"0xcae3edfaaa9707577d63808bb535e04f34066085",
		"0x2342ea371796e94e120de6835b38053cfd9279df",
		"0x4c98e8f2872bda14db1883428bf5ca1d2451c0b3",
		"0x4b725bb13bbc0dcb8acfb61222267a9bb9ad7ab3",
		"0x557aead4e0f3b0f25d1bc2bb9de36decba7124da",
		"0xe4adb46f5f67c4851110a199f75380b27609a57f",
		"0x4e13c4633fc8673cbc70ac9380d9fd658037ed73",
		"0x0e0fd41f921458d7b712a27db0112336595c3a86",
		"0x18311e0d647f5d154166d84cd5f79f2ea70c4c97",
		"0xf9882a169b982bdc65a007f5a1a8b64e146f26e9",
		"0x60c058b8563317a6eafc72dc65a6c88e739691cd",
		"0xf322dfb5a08a3791b6b179d654e8536fc31cc202",
		"0xaf91d663c43573ce9d8c5b7088033c71fc9b2a48",
		"0xd95ba8d39e244083da112e0b062ee53663a2a4c4",
		"0x8686d84d971696fb9d4e0def40d384aada7c670f",
		"0xc42682dad4dcc4b2ca859cac52b693c489815026",
		"0xfc58f7ffa8920fdafd2d697e4e467174aa5cddf0",
		"0x04b8b140df06685d8c1547156e0856044a9ba0a7",
		"0x62d64728be066685fbeb4bcc115f1081f6da4123",
		"0xf38c262fb9d6bc7e96b5ae0c689ce8b78c16c965",
		"0xbb7d8af4c34a83c32cd0f07d34635c7b538d58fb",
		"0x3421794c3b99b10b5e87e18201a3ad48c1c2d42c",
		"0xa251b3c4ffd7481b80e1fa7641384e2bb2fdcd81",
		"0x70bb8f86d6de728a58eb62c1242f9c6fa1c4401d",
		"0x3996fb6aab6d9ebb3ab9b3abe6e8482171c70aef",
		"0x6f20cf92bdea49962a779dcf63a086f9a63b66b4",
		"0xa49cdb9cfbeec0c144f5847763fe15c54fc0bd6f",
		"0x0a2be25c0df4bfdd647edb3da1eb613838239593",
		"0xc1df41f94ecb37ed6f713f06557c8616a54e8c2f",
		"0x61b0302d755c3f970a4806caf6a4f71a81e1fd37",
		"0xd30fcc4ceed4709d2fb3c3fcde7684bba7b6f3c0",
		"0x702af1dbf465bb2de13af94f23371ff84260e016",
		"0x4783a275546f3da5fa48b0162eeb8cbe941b12c0",
		"0x0d05f12a82ad4c8b00e276540ed0887bfecb9033",
		"0x5d4db511d4ce264d85a9d0ca85bf80e6b0078447",
		"0xc630263900404300d92f455d75610dd22acb1643",
		"0x810c1016b3339583722250a4bb728f090e821868",
		"0xb42efa83bd7198a13d17dff577de225df1b3a1bb",
		"0x0b0c66b0ba060928087e82394c19cd410d191fe7",
		"0x955b4edd4e51dc3cb54efc886f5cae9081e0f2d2",
		"0x1526ad23cc8e1e20a48b0b45d6124b9db404a8d3",
		"0x4366b15a25df5ad9e95621ee19b8f947739c722a",
		"0x79b238541af4100f2cdd056582520b4a4b29b573",
		"0x9b484f78e7ab9135e86ee489b6fb655620e7e538",
		"0xce584faada6c7f0e5b6784bbb5411af1c638e0b4",
		"0x793c0b6dbec447d97075d70602249e1fa018c32e",
		"0x9a032e23e3c4efd02421a48d683604df9c53dae7",
		"0x1ce8c2c440c225e70891c3a55b0e681358351501",
		"0x93d351aa7604a1465cf950682f964db0dd49c312",
		"0x1a7374fec40e7406e06885021898acd33776d543",
		"0x96f2faf4b51b90439b4e47f79e3d64a301780cb3",
		"0x19f471ad3ae6d6626238ed9673f21e003d5c4ff9",
		"0x00a33f55e96cbddb28707bbcc463244103ee624c",
		"0x9e0a4d5aa4b70ac415ce18b4f6f5cb00e65bbc49",
		"0x3f83a5a82889a5c22ceadd22e351de10b8558017",
		"0x664579dad3e1da41c1a50497d9d1e53bf7c8c493",
		"0xd2d9197eb5ab2e7040bb5d9215482068e634bc63",
		"0x8562980b3c91084e393f8c8ad53b6804d882d4a5",
		"0x44d501c4c6d3acbff637e4d76313d8c6e27428ab",
		"0x19b3969cdea171c10f9f6ed366de6dc7ce294390",
		"0xce09f79e42f2bc09cce1637717eb33509691486c",
		"0x0fbc7a821b7130064096e8e79ee4994ab1b9f918",
		"0x768186c50f82c4edd2bd507ae63ea8c473e36ead",
		"0x57aa20eb6cf85a48785d9b9135c17b8b8b5cbd8a",
		"0x53a1b2c543830dabd61b72a96409d868dfcadee2",
		"0x9a7ce733d30e9c6c7487ffce21fea9638bee500c",
		"0xc44269cf4ed91e6b77544ee01e9f78f09b257579",
		"0xb135d63c5b179c5789d91c36d3cecd7f00c76d6c",
		"0x37eaf34df3f2b7fbeb3ec180fe96dc76f6da443e",
		"0x03ceae3c2bb8a262e4219e7069de1f05f54e5b92",
		"0x09269cd9fccf7774cd743ef3b7c3d5e057752cba",
		"0x3d9338f4dfc25d5d3a61a25ee7e0107af6a6d035",
		"0x681f4a1509e1c47d0eee8cbb5bba8eafdc33b5d7",
		"0xbbc298b534b67406fce069b3c7ae29ed433031c0",
		"0x77568ae3453e02d6cd5c917d362a712984d59ce0",
		"0xb6b8845fc56e4605f2f082a827b08f0e2e15ef1e",
		"0x56f6afe355e64240febbd56c541e426cbe5dbe1c",
		"0x0ee1fc873de326c637285f2ccd7971e4c91509cc",
		"0x050bc73399c21a3f5fd0fd55ac6036939842e952",
		"0xeb19200f5ce221b9fb343e6508922514d8b15746",
		"0x73e640f271119aafb63c9589071817bb8cdb8f97",
		"0x6751ff2ec8015e78a1cebd7c8d6c2948110aaae9",
		"0xd32059287e944e09895dfe73992df05ff0a55732",
		"0xa179605f8f03522e031e0ab416a41e153faa565a",
		"0xe991899b8c0b36de12613974653753e546fae778",
		"0xd6eaa55f1d61e304dc1eb7aee0e80b35b617d9b5",
	}
}

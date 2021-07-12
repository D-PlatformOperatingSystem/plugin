// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tests

import (
	"encoding/hex"
	"testing"

	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/evm/executor/vm/common"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/evm/executor/vm/model"
)

//
func TestCreateContract1(t *testing.T) {
	deployCode, _ := hex.DecodeString("608060405260358060116000396000f3006080604052600080fd00a165627a7a723058203f5c7a16b3fd4fb82c8b466dd5a3f43773e41cc9c0acb98f83640880a39a68080029")
	execCode, _ := hex.DecodeString("6080604052600080fd00a165627a7a723058203f5c7a16b3fd4fb82c8b466dd5a3f43773e41cc9c0acb98f83640880a39a68080029")

	privKey := getPrivKey()

	gas := uint64(210000)
	gasLimit := gas
	tx := createTx(privKey, deployCode, gas, 10000000)
	mdb := buildStateDB(getAddr(privKey).String(), 500000000)
	ret, addr, leftGas, statedb, err := createContract(mdb, tx, 0)

	test := NewTester(t)
	test.assertNil(err)

	test.assertEqualsB(ret, execCode)
	test.assertBigger(int(gasLimit), int(leftGas))
	test.assertNotEqualsI(addr, common.EmptyAddress())

	//
	test.assertEqualsV(statedb.GetLastSnapshot().GetID(), 0)
}

//     gas
func TestCreateContract2(t *testing.T) {
	deployCode, _ := hex.DecodeString("60606040523415600e57600080fd5b603580601b6000396000f3006060604052600080fd00a165627a7a723058204bf1accefb2526a5077bcdfeaeb8020162814272245a9741cc2fddd89191af1c0029")
	//execCode, _ := hex.DecodeString("6060604052600080fd00a165627a7a723058204bf1accefb2526a5077bcdfeaeb8020162814272245a9741cc2fddd89191af1c0029")

	privKey := getPrivKey()
	//               61 Gas，        10600 Gas
	gas := uint64(30)
	tx := createTx(privKey, deployCode, gas, 0)
	mdb := buildStateDB(getAddr(privKey).String(), 100000000)
	ret, _, leftGas, _, err := createContract(mdb, tx, 0)

	test := NewTester(t)

	//    gas  ，
	test.assertNilB(ret)

	test.assertEqualsE(err, model.ErrOutOfGas)

	// gas   ，
	test.assertEqualsV(int(leftGas), 0)

}

//     gas
func TestCreateContract3(t *testing.T) {
	deployCode, _ := hex.DecodeString("60606040523415600e57600080fd5b603580601b6000396000f3006060604052600080fd00a165627a7a723058204bf1accefb2526a5077bcdfeaeb8020162814272245a9741cc2fddd89191af1c0029")
	execCode, _ := hex.DecodeString("6060604052600080fd00a165627a7a723058204bf1accefb2526a5077bcdfeaeb8020162814272245a9741cc2fddd89191af1c0029")

	privKey := getPrivKey()
	//               61 Gas，        10600 Gas
	usedGas := uint64(61)
	gas := uint64(100)
	gasLimit := gas
	tx := createTx(privKey, deployCode, gas, 0)
	mdb := buildStateDB(getAddr(privKey).String(), 100000000)
	ret, _, leftGas, _, err := createContract(mdb, tx, 0)

	test := NewTester(t)

	//            ，       StateDb
	test.assertEqualsB(ret, execCode)

	//  gas
	test.assertEqualsE(err, model.ErrCodeStoreOutOfGas)

	//
	// Gas      61 ，
	test.assertEqualsV(int(leftGas), int(gasLimit-usedGas))
}

// Gas  ，         （            ）
func TestCreateContract4(t *testing.T) {
	deployCode, _ := hex.DecodeString("60606040523415600e57600080fd5b603580601b6000396000f3006060604052600080fd00a165627a7a723058204bf1accefb2526a5077bcdfeaeb8020162814272245a9741cc2fddd89191af1c0029")
	execCode, _ := hex.DecodeString("6060604052600080fd00a165627a7a723058204bf1accefb2526a5077bcdfeaeb8020162814272245a9741cc2fddd89191af1c0029")

	privKey := getPrivKey()
	//               61 Gas，        10600 Gas
	gas := uint64(210000)
	gasLimit := gas
	tx := createTx(privKey, deployCode, gas, 0)
	mdb := buildStateDB(getAddr(privKey).String(), 100000000)
	ret, _, leftGas, _, err := createContract(mdb, tx, 50)

	test := NewTester(t)

	//
	test.assertNotNil(ret)
	test.assertEqualsB(ret, execCode)

	test.assertBigger(int(gasLimit), int(leftGas))

	//
	test.assertNotNil(err)
	test.assertEqualsS(err.Error(), "evm: max code size exceeded")
}

//
//      ：608060405234801561001057600080fd5b506298967f60008190555060df806100296000396000f3006080604052600436106049576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806360fe47b114604e5780636d4ce63c146078575b600080fd5b348015605957600080fd5b5060766004803603810190808035906020019092919050505060a0565b005b348015608357600080fd5b50608a60aa565b6040518082815260200191505060405180910390f35b8060008190555050565b600080549050905600a165627a7a72305820b3ccec4d8cbe393844da31834b7464f23d3b81b24f36ce7e18bb09601f2eb8660029
//contract MyStore {
//    uint value;
//    constructor() public{
//        value=9999999;
//    }
//    function set(uint x) public {
//        value = x;
//    }
//
//    function get() public constant returns (uint){
//        return value;
//    }
//}

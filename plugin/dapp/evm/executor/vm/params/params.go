// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package params

import (
	"math/big"

	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/evm/executor/vm/common"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/evm/executor/vm/state"
)

// GasParam       Gas           Gas
//        ，       、   、Gas
//
type GasParam struct {
	// Gas         Gas（            ）
	Gas uint64

	// Address
	//   ，     CallCode   ，                  ，
	Address common.Address
}

// EVMParam       Gas           EVM
//        ，       、Gas
//          StateDB CallGasTemp  ，     Gas       ：
// 1.   EVM  EVMParam（   CallGasTemp ，      ）；
// 2.  EVMParam   ，  Gas  ；
// 3.      ，  EVMParam      EVM ；
type EVMParam struct {

	// EVMStateDB
	StateDB state.EVMStateDB

	// CallGasTemp               Gas
	//       ，      gasCost  ，           Gas，
	//      opCall ，         Gas
	CallGasTemp uint64

	// BlockNumber NUMBER   ，
	BlockNumber *big.Int
}

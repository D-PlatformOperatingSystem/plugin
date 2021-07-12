// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gas

import (
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/evm/executor/vm/model"
	"github.com/holiman/uint256"
)

const (
	// GasQuickStep
	GasQuickStep uint64 = 2
	// GasFastestStep
	GasFastestStep uint64 = 3
	// GasFastStep
	GasFastStep uint64 = 5
	// GasMidStep
	GasMidStep uint64 = 8
	// GasSlowStep
	GasSlowStep uint64 = 10
	// GasExtStep
	GasExtStep uint64 = 20

	// MaxNewMemSize              ，
	MaxNewMemSize uint64 = 0xffffffffe0
)

//        Gas
//  availableGas - base * 63 / 64.
func callGas(gasTable Table, availableGas, base uint64, callCost *uint256.Int) (uint64, error) {
	if availableGas == callCost.Uint64() {
		availableGas = availableGas - base
		gas := availableGas - availableGas/64

		//      callCost  ，           ，        gas
		if callCost.BitLen() > 64 || gas < callCost.Uint64() {
			return gas, nil
		}
	}

	//   Gas  ，
	if callCost.BitLen() > 64 {
		return 0, model.ErrGasUintOverflow
	}

	return callCost.Uint64(), nil
}

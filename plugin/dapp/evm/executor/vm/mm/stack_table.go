// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mm

import (
	"fmt"

	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/evm/executor/vm/params"
)

type (
	// StackValidationFunc
	StackValidationFunc func(*Stack) error
)

// MakeStackFunc           （                 ）
func MakeStackFunc(pop, push int) StackValidationFunc {
	return func(stack *Stack) error {
		if err := stack.Require(pop); err != nil {
			return err
		}

		if stack.Len()+push-pop > int(params.StackLimit) {
			return fmt.Errorf("stack limit reached %d (%d)", stack.Len(), params.StackLimit)
		}
		return nil
	}
}

// MakeDupStackFunc
func MakeDupStackFunc(n int) StackValidationFunc {
	return MakeStackFunc(n, n+1)
}

// MakeSwapStackFunc
func MakeSwapStackFunc(n int) StackValidationFunc {
	return MakeStackFunc(n, n)
}

//func MinSwapStack(n int) int {
//	return MinStack(n, n)
//}
//func MaxSwapStack(n int) int {
//	return MaxStack(n, n)
//}
//
//func MinDupStack(n int) int {
//	return MinStack(n, n+1)
//}
//func MaxDupStack(n int) int {
//	return MaxStack(n, n+1)
//}

//func MaxStack(pop, push int) int {
//	return int(params.StackLimit) + pop - push
//}
//func MinStack(pops, push int) int {
//	return pops
//}

package executor

import (
	"github.com/D-PlatformOperatingSystem/dpos/common"
	"github.com/D-PlatformOperatingSystem/dpos/common/address"
	"github.com/robertkrimen/otto"
)

//  js
func execaddressFunc(vm *otto.Otto) {
	vm.Set("execaddress", func(call otto.FunctionCall) otto.Value {
		key, err := call.Argument(0).ToString()
		if err != nil {
			return errReturn(vm, err)
		}
		addr := address.ExecAddress(key)
		return okReturn(vm, addr)
	})
}

func sha256Func(vm *otto.Otto) {
	vm.Set("sha256", func(call otto.FunctionCall) otto.Value {
		key, err := call.Argument(0).ToString()
		if err != nil {
			return errReturn(vm, err)
		}
		var hash = common.Sha256([]byte(key))
		return okReturn(vm, common.ToHex(hash))
	})
}

/*
//
//randnum

//       hash
//prev_blockhash()
*/

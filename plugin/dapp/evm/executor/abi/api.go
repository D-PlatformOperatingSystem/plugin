package abi

import (
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"strings"

	"errors"

	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/evm/executor/vm/common"
	"github.com/golang-collections/collections/stack"
)

// Pack   ABI     ，        EVM
// abiData    ABI
// param
// readOnly     ，           ，
//     ： foo(param1,param2)
func Pack(param, abiData string, readOnly bool) (methodName string, packData []byte, err error) {
	//          ，
	methodName, params, err := procFuncCall(param)
	if err != nil {
		return methodName, packData, err
	}

	//   ABI    ，
	abi, err := JSON(strings.NewReader(abiData))
	if err != nil {
		return methodName, packData, err
	}

	var method Method
	var ok bool
	if method, ok = abi.Methods[methodName]; !ok {
		err = fmt.Errorf("function %v not exists", methodName)
		return methodName, packData, err
	}

	if readOnly && !method.IsConstant() {
		return methodName, packData, errors.New("method is not readonly")
	}
	if len(params) != method.Inputs.LengthNonIndexed() {
		err = fmt.Errorf("function params error:%v", params)
		return methodName, packData, err
	}
	//         ，       ，     Go
	paramVals := []interface{}{}
	if len(params) != 0 {
		//          ABI
		if method.Inputs.LengthNonIndexed() != len(params) {
			err = fmt.Errorf("function Params count error: %v", param)
			return methodName, packData, err
		}

		for i, v := range method.Inputs.NonIndexed() {
			paramVal, err := str2GoValue(v.Type, params[i])
			if err != nil {
				return methodName, packData, err
			}
			paramVals = append(paramVals, paramVal)
		}
	}

	//   Abi
	packData, err = abi.Pack(methodName, paramVals...)
	return methodName, packData, err
}

// Unpack          ABI       json
// data
// abiData    ABI
func Unpack(data []byte, methodName, abiData string) (output string, err error) {
	if len(data) == 0 {
		return output, err
	}
	//   ABI    ，
	abi, err := JSON(strings.NewReader(abiData))
	if err != nil {
		return output, err
	}

	var method Method
	var ok bool
	if method, ok = abi.Methods[methodName]; !ok {
		return output, fmt.Errorf("function %v not exists", methodName)
	}

	if method.Outputs.LengthNonIndexed() == 0 {
		return output, err
	}

	values, err := method.Outputs.UnpackValues(data)
	if err != nil {
		return output, err
	}

	outputs := []*Param{}

	for i, v := range values {
		arg := method.Outputs[i]
		pval := &Param{Name: arg.Name, Type: arg.Type.String(), Value: v}
		outputs = append(outputs, pval)
	}

	jsondata, err := json.Marshal(outputs)
	if err != nil {
		return output, err
	}
	return string(jsondata), err
}

// Param
type Param struct {
	// Name
	Name string `json:"name"`
	// Type
	Type string `json:"type"`
	// Value
	Value interface{} `json:"value"`
}

func convertUint(val uint64, kind reflect.Kind) interface{} {
	switch kind {
	case reflect.Uint:
		return uint(val)
	case reflect.Uint8:
		return uint8(val)
	case reflect.Uint16:
		return uint16(val)
	case reflect.Uint32:
		return uint32(val)
	case reflect.Uint64:
		return val
	}
	return val
}

func convertInt(val int64, kind reflect.Kind) interface{} {
	switch kind {
	case reflect.Int:
		return int(val)
	case reflect.Int8:
		return int8(val)
	case reflect.Int16:
		return int16(val)
	case reflect.Int32:
		return int32(val)
	case reflect.Int64:
		return val
	}
	return val
}

//              （  ），  Go
func str2GoValue(typ Type, val string) (res interface{}, err error) {
	switch typ.T {
	case IntTy:
		if typ.Size < 256 {
			x, err := strconv.ParseInt(val, 10, typ.Size)
			if err != nil {
				return res, err
			}
			return convertInt(x, typ.GetType().Kind()), nil
		}
		b := new(big.Int)
		b.SetString(val, 10)
		return b, err
	case UintTy:
		if typ.Size < 256 {
			x, err := strconv.ParseUint(val, 10, typ.Size)
			if err != nil {
				return res, err
			}
			return convertUint(x, typ.GetType().Kind()), nil
		}
		b := new(big.Int)
		b.SetString(val, 10)
		return b, err
	case BoolTy:
		x, err := strconv.ParseBool(val)
		if err != nil {
			return res, err
		}
		return x, nil
	case StringTy:
		return val, nil
	case SliceTy:
		subs, err := procArrayItem(val)
		if err != nil {
			return res, err
		}
		rval := reflect.MakeSlice(typ.GetType(), len(subs), len(subs))
		for idx, sub := range subs {
			subVal, er := str2GoValue(*typ.Elem, sub)
			if er != nil {
				return res, er
			}
			rval.Index(idx).Set(reflect.ValueOf(subVal))
		}
		return rval.Interface(), nil
	case ArrayTy:
		rval := reflect.New(typ.GetType()).Elem()
		subs, err := procArrayItem(val)
		if err != nil {
			return res, err
		}
		for idx, sub := range subs {
			subVal, er := str2GoValue(*typ.Elem, sub)
			if er != nil {
				return res, er
			}
			rval.Index(idx).Set(reflect.ValueOf(subVal))
		}
		return rval.Interface(), nil
	case AddressTy:
		addr := common.StringToAddress(val)
		if addr == nil {
			return res, fmt.Errorf("invalid  address: %v", val)
		}
		return addr.ToHash160(), nil
	case FixedBytesTy:
		//        ，            ，  0xabcd00ff
		x, err := common.HexToBytes(val)
		if err != nil {
			return res, err
		}
		rval := reflect.New(typ.GetType()).Elem()
		for i, b := range x {
			rval.Index(i).Set(reflect.ValueOf(b))
		}
		return rval.Interface(), nil
	case BytesTy:
		//     ，            ，  0xab
		x, err := common.HexToBytes(val)
		if err != nil {
			return res, err
		}
		return x, nil
	case HashTy:
		//     ，          ， ：0xabcdef
		x, err := common.HexToBytes(val)
		if err != nil {
			return res, err
		}
		return common.BytesToHash(x), nil
	default:
		return res, fmt.Errorf("not support type: %v", typ.stringKind)
	}
}

//                 ，     ，          ；
//        ，
//   ："[a,b,c]" -> "a","b","c"
//   ："[[a,b],[c,d]]" -> "[a,b]", "[c,d]"
//         ，          ，
func procArrayItem(val string) (res []string, err error) {
	ss := stack.New()
	data := []rune{}
	for _, b := range val {
		switch b {
		case ' ':
			//
			if ss.Len() > 0 && peekRune(ss) == '"' {
				data = append(data, b)
			}
		case ',':
			//                 ，             ，
			//   ，              '['，
			if ss.Len() == 1 && peekRune(ss) == '[' {
				//
				res = append(res, string(data))
				data = []rune{}

			} else {
				data = append(data, b)
			}
		case '"':
			//             ，
			if ss.Peek() == b {
				ss.Pop()
			} else {
				ss.Push(b)
			}
			//data = append(data, b)
		case '[':
			//        ，'['         ，
			if ss.Len() == 0 {
				data = []rune{}
			} else {
				data = append(data, b)
			}
			ss.Push(b)
		case ']':
			//          ']' ，         ，
			if ss.Len() == 1 && peekRune(ss) == '[' {
				//
				res = append(res, string(data))
			} else {
				data = append(data, b)
			}
			ss.Pop()
		default:
			//
			data = append(data, b)
		}
	}

	if ss.Len() != 0 {
		return nil, fmt.Errorf("invalid array format:%v", val)
	}
	return res, err
}

func peekRune(ss *stack.Stack) rune {
	return ss.Peek().(rune)
}

//          ，
//   ：foo(param1,param2) -> [foo,param1,param2]
func procFuncCall(param string) (funcName string, res []string, err error) {
	lidx := strings.Index(param, "(")
	ridx := strings.LastIndex(param, ")")

	if lidx == -1 || ridx == -1 {
		return funcName, res, fmt.Errorf("invalid function signature:%v", param)
	}

	funcName = strings.TrimSpace(param[:lidx])
	params := strings.TrimSpace(param[lidx+1 : ridx])

	//             ，          ，
	if len(params) > 0 {
		res, err = procArrayItem(fmt.Sprintf("[%v]", params))
	}

	return funcName, res, err
}

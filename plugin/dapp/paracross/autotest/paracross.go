// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package autotest

import (
	"reflect"

	"github.com/D-PlatformOperatingSystem/dpos/cmd/autotest/types"
	coinautotest "github.com/D-PlatformOperatingSystem/dpos/system/dapp/coins/autotest"
	tokenautotest "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/token/autotest"
)

type paracrossAutoTest struct {
	SimpleCaseArr            []types.SimpleCase                    `toml:"SimpleCase,omitempty"`
	TokenPreCreateCaseArr    []tokenautotest.TokenPreCreateCase    `toml:"TokenPreCreateCase,omitempty"`
	TokenFinishCreateCaseArr []tokenautotest.TokenFinishCreateCase `toml:"TokenFinishCreateCase,omitempty"`
	TransferCaseArr          []coinautotest.TransferCase           `toml:"TransferCase,omitempty"`
}

func init() {

	types.RegisterAutoTest(paracrossAutoTest{})

}

func (config paracrossAutoTest) GetName() string {

	return "paracross"
}

func (config paracrossAutoTest) GetTestConfigType() reflect.Type {

	return reflect.TypeOf(config)
}

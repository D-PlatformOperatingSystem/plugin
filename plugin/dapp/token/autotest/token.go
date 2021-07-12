// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package autotest

import (
	"reflect"

	"github.com/D-PlatformOperatingSystem/dpos/cmd/autotest/types"
	"github.com/D-PlatformOperatingSystem/dpos/system/dapp/coins/autotest"
)

type tokenAutoTest struct {
	SimpleCaseArr            []types.SimpleCase      `toml:"SimpleCase,omitempty"`
	TokenPreCreateCaseArr    []TokenPreCreateCase    `toml:"TokenPreCreateCase,omitempty"`
	TokenFinishCreateCaseArr []TokenFinishCreateCase `toml:"TokenFinishCreateCase,omitempty"`
	TransferCaseArr          []autotest.TransferCase `toml:"TransferCase,omitempty"`
	WithdrawCaseArr          []autotest.WithdrawCase `toml:"WithdrawCase,omitempty"`
	TokenRevokeCaseArr       []TokenRevokeCase       `toml:"TokenRevokeCase,omitempty"`
	TokenMintCaseArr         []TokenMintCase         `toml:"TokenMintCase,omitempty"`
	TokenBurnCaseArr         []TokenBurnCase         `toml:"TokenBurnCase,omitempty"`
}

func init() {

	types.RegisterAutoTest(tokenAutoTest{})

}

func (config tokenAutoTest) GetName() string {

	return "token"
}

func (config tokenAutoTest) GetTestConfigType() reflect.Type {

	return reflect.TypeOf(config)
}

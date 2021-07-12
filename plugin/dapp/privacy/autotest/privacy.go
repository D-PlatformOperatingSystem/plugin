// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package autotest

import (
	"reflect"

	"github.com/D-PlatformOperatingSystem/dpos/cmd/autotest/types"
	ca "github.com/D-PlatformOperatingSystem/dpos/system/dapp/coins/autotest"
	ta "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/token/autotest"
)

type privacyAutoTest struct {
	SimpleCaseArr            []types.SimpleCase         `toml:"SimpleCase,omitempty"`
	TokenPreCreateCaseArr    []ta.TokenPreCreateCase    `toml:"TokenPreCreateCase,omitempty"`
	TokenFinishCreateCaseArr []ta.TokenFinishCreateCase `toml:"TokenFinishCreateCase,omitempty"`
	TransferCaseArr          []ca.TransferCase          `toml:"TransferCase,omitempty"`
	PubToPrivCaseArr         []PubToPrivCase            `toml:"PubToPrivCase,omitempty"`
	PrivToPrivCaseArr        []PrivToPrivCase           `toml:"PrivToPrivCase,omitempty"`
	PrivToPubCaseArr         []PrivToPubCase            `toml:"PrivToPubCase,omitempty"`
	CreateUtxosCaseArr       []CreateUtxosCase          `toml:"CreateUtxosCase,omitempty"`
}

func init() {

	types.RegisterAutoTest(privacyAutoTest{})

}

func (config privacyAutoTest) GetName() string {

	return "privacy"
}

func (config privacyAutoTest) GetTestConfigType() reflect.Type {

	return reflect.TypeOf(config)
}

// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddressBig(t *testing.T) {
	saddr := "16ERTbYtKKQ64wMthAY9J4La4nAiidG45A"
	addr := StringToAddress(saddr)
	baddr := addr.Big()
	naddr := BigToAddress(baddr)
	if saddr != naddr.String() {
		t.Fail()
	}
}
func TestExecAddress2(t *testing.T) {
	addr := ExecAddress("autonomy")
	t.Log(addr)
}
func TestAddressBytes(t *testing.T) {
	addr := BytesToAddress([]byte{1})
	assert.Equal(t, addr.String(), "11111111111111111111BZbvjr")
}
func TestExecAddress(t *testing.T) {
	addr := ExecAddress("1CQXE6TxaYCG5mADtWij4AxhZCUTpoABb3")
	t.Log(addr)
}

// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testnode

import (
	"testing"

	"github.com/D-PlatformOperatingSystem/dpos/util/testnode"
	ty "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/ticket/types"
	ticketwallet "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/ticket/wallet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "github.com/D-PlatformOperatingSystem/dpos/system"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	_ "github.com/D-PlatformOperatingSystem/plugin/plugin"
)

func TestWalletTicket(t *testing.T) {
	minerAddr := "12oupcayRT7LvaC4qW4avxsTE7U41cKSio"
	t.Log("Begin wallet ticket test")

	cfg := types.NewDplatformOSConfig(types.GetDefaultCfgstring())
	cfg.GetModuleConfig().Consensus.Name = "ticket"
	mockDOM := testnode.NewWithConfig(cfg, nil)
	defer mockDOM.Close()
	err := mockDOM.WaitHeight(0)
	assert.Nil(t, err)
	msg, err := mockDOM.GetAPI().Query(ty.TicketX, "TicketList", &ty.TicketList{Addr: minerAddr, Status: 1})
	assert.Nil(t, err)
	ticketList := msg.(*ty.ReplyTicketList)
	assert.NotNil(t, ticketList)
	//return
	ticketwallet.FlushTicket(mockDOM.GetAPI())
	err = mockDOM.WaitHeight(2)
	assert.Nil(t, err)
	header, err := mockDOM.GetAPI().GetLastHeader()
	require.Equal(t, err, nil)
	require.Equal(t, header.Height >= 2, true)

	in := &ty.TicketClose{MinerAddress: minerAddr}
	msg, err = mockDOM.GetAPI().ExecWalletFunc(ty.TicketX, "CloseTickets", in)
	assert.Nil(t, err)
	hashes := msg.(*types.ReplyHashes)
	assert.NotNil(t, hashes)

	in = &ty.TicketClose{}
	msg, err = mockDOM.GetAPI().ExecWalletFunc(ty.TicketX, "CloseTickets", in)
	assert.Nil(t, err)
	hashes = msg.(*types.ReplyHashes)
	assert.NotNil(t, hashes)
	t.Log("End wallet ticket test")
}

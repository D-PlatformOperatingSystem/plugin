// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"github.com/D-PlatformOperatingSystem/dpos/types"
	pkt "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/pokerbull/types"
)

// Exec_Start
func (c *PokerBull) Exec_Start(payload *pkt.PBGameStart, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(c, tx, index)
	return action.GameStart(payload)
}

// Exec_Continue
func (c *PokerBull) Exec_Continue(payload *pkt.PBGameContinue, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(c, tx, index)
	return action.GameContinue(payload)
}

// Exec_Quit
func (c *PokerBull) Exec_Quit(payload *pkt.PBGameQuit, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(c, tx, index)
	return action.GameQuit(payload)
}

// Exec_Play
func (c *PokerBull) Exec_Play(payload *pkt.PBGamePlay, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(c, tx, index)
	return action.GamePlay(payload)
}

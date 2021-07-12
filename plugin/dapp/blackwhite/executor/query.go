// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"github.com/D-PlatformOperatingSystem/dpos/types"
	gt "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/blackwhite/types"
)

// Query_GetBlackwhiteRoundInfo
func (c *Blackwhite) Query_GetBlackwhiteRoundInfo(in *gt.ReqBlackwhiteRoundInfo) (types.Message, error) {
	if in == nil {
		return nil, types.ErrInvalidParam
	}
	return c.GetBlackwhiteRoundInfo(in)
}

// Query_GetBlackwhiteByStatusAndAddr
func (c *Blackwhite) Query_GetBlackwhiteByStatusAndAddr(in *gt.ReqBlackwhiteRoundList) (types.Message, error) {
	if in == nil {
		return nil, types.ErrInvalidParam
	}
	return c.GetBwRoundListInfo(in)
}

// Query_GetBlackwhiteloopResult
func (c *Blackwhite) Query_GetBlackwhiteloopResult(in *gt.ReqLoopResult) (types.Message, error) {
	if in == nil {
		return nil, types.ErrInvalidParam
	}
	return c.GetBwRoundLoopResult(in)
}

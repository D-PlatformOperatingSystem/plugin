// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"github.com/D-PlatformOperatingSystem/dpos/types"
	pty "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/norm/types"
)

// Query_NormGet get value
func (n *Norm) Query_NormGet(in *pty.NormGetKey) (types.Message, error) {
	value, err := n.GetStateDB().Get(Key(in.Key))
	if err != nil {
		return nil, types.ErrNotFound
	}
	return &types.ReplyString{Data: string(value)}, nil
}

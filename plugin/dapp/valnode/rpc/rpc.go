// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rpc

import (
	"context"

	"github.com/D-PlatformOperatingSystem/dpos/types"
	vt "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/valnode/types"
)

// IsSync query is sync
func (c *channelClient) IsSync(ctx context.Context, req *types.ReqNil) (*vt.IsHealthy, error) {
	data, err := c.QueryConsensusFunc("tendermint", "IsHealthy", &types.ReqNil{})
	if err != nil {
		return nil, err
	}
	if resp, ok := data.(*vt.IsHealthy); ok {
		return resp, nil
	}
	return nil, types.ErrDecode
}

// IsSync query is sync
func (c *Jrpc) IsSync(req *types.ReqNil, result *interface{}) error {
	data, err := c.cli.IsSync(context.Background(), req)
	if err != nil {
		return err
	}
	*result = data.IsHealthy
	return nil
}

// GetNodeInfo query node info
func (c *channelClient) GetNodeInfo(ctx context.Context, req *types.ReqNil) (*vt.ValidatorSet, error) {
	data, err := c.QueryConsensusFunc("tendermint", "NodeInfo", &types.ReqNil{})
	if err != nil {
		return nil, err
	}
	if resp, ok := data.(*vt.ValidatorSet); ok {
		return resp, nil
	}
	return nil, types.ErrDecode
}

// GetNodeInfo query node info
func (c *Jrpc) GetNodeInfo(req *types.ReqNil, result *interface{}) error {
	data, err := c.cli.GetNodeInfo(context.Background(), req)
	if err != nil {
		return err
	}
	*result = data.Validators
	return nil
}

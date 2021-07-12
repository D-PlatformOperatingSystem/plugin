// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rpc

import (
	"context"
	"encoding/hex"

	"github.com/D-PlatformOperatingSystem/dpos/types"
	pty "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/unfreeze/types"
)

// GetUnfreeze
func (c *channelClient) GetUnfreeze(ctx context.Context, in *types.ReqString) (*pty.Unfreeze, error) {
	v, err := c.Query(pty.UnfreezeX, "GetUnfreeze", in)
	if err != nil {
		return nil, err
	}
	if resp, ok := v.(*pty.Unfreeze); ok {
		return resp, nil
	}
	return nil, types.ErrDecode
}

// GetUnfreezeWithdraw
func (c *channelClient) GetUnfreezeWithdraw(ctx context.Context, in *types.ReqString) (*pty.ReplyQueryUnfreezeWithdraw, error) {
	v, err := c.Query(pty.UnfreezeX, "GetUnfreezeWithdraw", in)
	if err != nil {
		return nil, err
	}
	if resp, ok := v.(*pty.ReplyQueryUnfreezeWithdraw); ok {
		return resp, nil
	}
	return nil, types.ErrDecode
}

// GetUnfreeze
func (c *Jrpc) GetUnfreeze(in *types.ReqString, result *interface{}) error {
	v, err := c.cli.GetUnfreeze(context.Background(), in)
	if err != nil {
		return err
	}

	*result = v
	return nil
}

// GetUnfreezeWithdraw
func (c *Jrpc) GetUnfreezeWithdraw(in *types.ReqString, result *interface{}) error {
	v, err := c.cli.GetUnfreezeWithdraw(context.Background(), in)
	if err != nil {
		return err
	}

	*result = v
	return nil
}

// CreateRawUnfreezeCreate
func (c *Jrpc) CreateRawUnfreezeCreate(param *pty.UnfreezeCreate, result *interface{}) error {
	if param == nil {
		return types.ErrInvalidParam
	}
	cfg := c.cli.GetConfig()
	data, err := types.CallCreateTx(cfg, cfg.ExecName(pty.UnfreezeX), "Create", param)
	if err != nil {
		return err
	}
	*result = hex.EncodeToString(data)
	return nil
}

// CreateRawUnfreezeWithdraw
func (c *Jrpc) CreateRawUnfreezeWithdraw(param *pty.UnfreezeWithdraw, result *interface{}) error {
	if param == nil {
		return types.ErrInvalidParam
	}
	cfg := c.cli.GetConfig()
	data, err := types.CallCreateTx(cfg, cfg.ExecName(pty.UnfreezeX), "Withdraw", param)
	if err != nil {
		return err
	}
	*result = hex.EncodeToString(data)
	return nil
}

// CreateRawUnfreezeTerminate
func (c *Jrpc) CreateRawUnfreezeTerminate(param *pty.UnfreezeTerminate, result *interface{}) error {
	if param == nil {
		return types.ErrInvalidParam
	}
	cfg := c.cli.GetConfig()
	data, err := types.CallCreateTx(cfg, cfg.ExecName(pty.UnfreezeX), "Terminate", param)
	if err != nil {
		return err
	}
	*result = hex.EncodeToString(data)
	return nil
}

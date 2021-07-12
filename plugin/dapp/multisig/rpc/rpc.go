// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rpc

import (
	"encoding/hex"

	"github.com/D-PlatformOperatingSystem/dpos/types"
	mty "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/multisig/types"
)

// MultiSigAccCreateTx :
func (c *Jrpc) MultiSigAccCreateTx(param *mty.MultiSigAccCreate, result *interface{}) error {
	if param == nil {
		return types.ErrInvalidParam
	}
	cfg := c.cli.GetConfig()
	data, err := types.CallCreateTx(cfg, cfg.ExecName(mty.MultiSigX), "MultiSigAccCreate", param)
	if err != nil {
		return err
	}
	*result = hex.EncodeToString(data)
	return nil
}

// MultiSigOwnerOperateTx :          owner
func (c *Jrpc) MultiSigOwnerOperateTx(param *mty.MultiSigOwnerOperate, result *interface{}) error {
	if param == nil {
		return types.ErrInvalidParam
	}
	cfg := c.cli.GetConfig()
	data, err := types.CallCreateTx(cfg, cfg.ExecName(mty.MultiSigX), "MultiSigOwnerOperate", param)
	if err != nil {
		return err
	}
	*result = hex.EncodeToString(data)
	return nil
}

// MultiSigAccOperateTx :
func (c *Jrpc) MultiSigAccOperateTx(param *mty.MultiSigAccOperate, result *interface{}) error {
	if param == nil {
		return types.ErrInvalidParam
	}
	cfg := c.cli.GetConfig()
	data, err := types.CallCreateTx(cfg, cfg.ExecName(mty.MultiSigX), "MultiSigAccOperate", param)
	if err != nil {
		return err
	}
	*result = hex.EncodeToString(data)
	return nil
}

// MultiSigConfirmTx :
func (c *Jrpc) MultiSigConfirmTx(param *mty.MultiSigConfirmTx, result *interface{}) error {
	if param == nil {
		return types.ErrInvalidParam
	}
	cfg := c.cli.GetConfig()
	data, err := types.CallCreateTx(cfg, cfg.ExecName(mty.MultiSigX), "MultiSigConfirmTx", param)
	if err != nil {
		return err
	}
	*result = hex.EncodeToString(data)
	return nil
}

// MultiSigAccTransferInTx :
func (c *Jrpc) MultiSigAccTransferInTx(param *mty.MultiSigExecTransferTo, result *interface{}) error {
	if param == nil {
		return types.ErrInvalidParam
	}
	v := *param
	cfg := c.cli.GetConfig()
	data, err := types.CallCreateTx(cfg, cfg.ExecName(mty.MultiSigX), "MultiSigExecTransferTo", &v)
	if err != nil {
		return err
	}
	*result = hex.EncodeToString(data)
	return nil
}

// MultiSigAccTransferOutTx :
func (c *Jrpc) MultiSigAccTransferOutTx(param *mty.MultiSigExecTransferFrom, result *interface{}) error {
	if param == nil {
		return types.ErrInvalidParam
	}
	v := *param
	cfg := c.cli.GetConfig()
	data, err := types.CallCreateTx(cfg, cfg.ExecName(mty.MultiSigX), "MultiSigExecTransferFrom", &v)
	if err != nil {
		return err
	}
	*result = hex.EncodeToString(data)
	return nil
}

// MultiSigAddresList   owner            {multiSigAddr，owneraddr，weight}
func (c *Jrpc) MultiSigAddresList(in *types.ReqString, result *interface{}) error {
	v := *in
	data, err := c.cli.ExecWalletFunc(mty.MultiSigX, "MultiSigAddresList", &v)
	if err != nil {
		return err
	}
	ownerAttrs := data.(*mty.OwnerAttrs)
	*result = ownerAttrs
	return nil
}

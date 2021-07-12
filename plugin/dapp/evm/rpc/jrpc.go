// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rpc

import (
	"context"
	"encoding/hex"

	"github.com/D-PlatformOperatingSystem/dpos/types"
	evm "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/evm/types"
)

// EvmCreateTx   Evm
func (c *Jrpc) EvmCreateTx(parm *evm.EvmContractCreateReq, result *interface{}) error {
	if parm == nil {
		return types.ErrInvalidParam
	}

	reply, err := c.cli.Create(context.Background(), *parm)
	if err != nil {
		return err
	}
	*result = hex.EncodeToString(reply.Data)
	return nil
}

// EvmCallTx   Evm
func (c *Jrpc) EvmCallTx(parm *evm.EvmContractCallReq, result *interface{}) error {
	if parm == nil {
		return types.ErrInvalidParam
	}

	reply, err := c.cli.Call(context.Background(), *parm)
	if err != nil {
		return err
	}
	*result = hex.EncodeToString(reply.Data)
	return nil
}

// EvmTransferTx Evm
func (c *Jrpc) EvmTransferTx(parm *evm.EvmContractTransferReq, result *interface{}) error {
	if parm == nil {
		return types.ErrInvalidParam
	}

	reply, err := c.cli.Transfer(context.Background(), *parm, false)
	if err != nil {
		return err
	}

	*result = hex.EncodeToString(reply.Data)
	return nil
}

// EvmWithdrawTx Evm
func (c *Jrpc) EvmWithdrawTx(parm *evm.EvmContractTransferReq, result *interface{}) error {
	if parm == nil {
		return types.ErrInvalidParam
	}

	reply, err := c.cli.Transfer(context.Background(), *parm, true)
	if err != nil {
		return err
	}

	*result = hex.EncodeToString(reply.Data)
	return nil
}

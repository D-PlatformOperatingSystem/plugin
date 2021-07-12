// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rpc

import (
	"encoding/hex"

	"github.com/D-PlatformOperatingSystem/dpos/account"
	"github.com/D-PlatformOperatingSystem/dpos/common/address"
	rpctypes "github.com/D-PlatformOperatingSystem/dpos/rpc/types"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	tokenty "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/token/types"
	context "golang.org/x/net/context"
)

//TODO: GetBalance      ，  LoadAccounts LoadExecAccountQueue         , added by hzj
func (c *channelClient) getTokenBalance(in *tokenty.ReqTokenBalance) ([]*types.Account, error) {
	cfg := c.GetConfig()
	accountTokendb, err := account.NewAccountDB(cfg, tokenty.TokenX, in.GetTokenSymbol(), nil)
	if err != nil {
		return nil, err
	}

	switch in.GetExecer() {
	case cfg.ExecName(tokenty.TokenX):
		addrs := in.GetAddresses()
		var queryAddrs []string
		queryAddrs = append(queryAddrs, addrs...)

		accounts, err := accountTokendb.LoadAccounts(c.QueueProtocolAPI, queryAddrs)
		if err != nil {
			log.Error("GetTokenBalance", "err", err.Error(), "token symbol", in.GetTokenSymbol(), "address", queryAddrs)
			return nil, err
		}
		return accounts, nil

	default: //trade
		execaddress := address.ExecAddress(in.GetExecer())
		addrs := in.GetAddresses()
		var accounts []*types.Account
		for _, addr := range addrs {
			acc, err := accountTokendb.LoadExecAccountQueue(c.QueueProtocolAPI, addr, execaddress)
			if err != nil {
				log.Error("GetTokenBalance for exector", "err", err.Error(), "token symbol", in.GetTokenSymbol(),
					"address", addr)
				continue
			}
			accounts = append(accounts, acc)
		}

		return accounts, nil
	}
}

// GetTokenBalance   token  （channelClient）
func (c *channelClient) GetTokenBalance(ctx context.Context, in *tokenty.ReqTokenBalance) (*types.Accounts, error) {
	reply, err := c.getTokenBalance(in)
	if err != nil {
		return nil, err
	}
	return &types.Accounts{Acc: reply}, nil
}

// GetTokenBalance   token   (Jrpc)
func (c *Jrpc) GetTokenBalance(in tokenty.ReqTokenBalance, result *interface{}) error {
	balances, err := c.cli.GetTokenBalance(context.Background(), &in)
	if err != nil {
		return err
	}
	var accounts []*rpctypes.Account
	for _, balance := range balances.Acc {
		accounts = append(accounts, &rpctypes.Account{Addr: balance.GetAddr(),
			Balance:  balance.GetBalance(),
			Currency: balance.GetCurrency(),
			Frozen:   balance.GetFrozen()})
	}
	*result = accounts
	return nil
}

// CreateRawTokenPreCreateTx         Token
func (c *Jrpc) CreateRawTokenPreCreateTx(param *tokenty.TokenPreCreate, result *interface{}) error {
	if param == nil || param.Symbol == "" {
		return types.ErrInvalidParam
	}
	cfg := c.cli.GetConfig()
	data, err := types.CallCreateTx(cfg, cfg.ExecName(tokenty.TokenX), "TokenPreCreate", param)
	if err != nil {
		return err
	}
	*result = hex.EncodeToString(data)
	return nil
}

// CreateRawTokenFinishTx         Token
func (c *Jrpc) CreateRawTokenFinishTx(param *tokenty.TokenFinishCreate, result *interface{}) error {
	if param == nil || param.Symbol == "" {
		return types.ErrInvalidParam
	}
	cfg := c.cli.GetConfig()
	data, err := types.CallCreateTx(cfg, cfg.ExecName(tokenty.TokenX), "TokenFinishCreate", param)
	if err != nil {
		return err
	}
	*result = hex.EncodeToString(data)
	return nil
}

// CreateRawTokenRevokeTx         Token
func (c *Jrpc) CreateRawTokenRevokeTx(param *tokenty.TokenRevokeCreate, result *interface{}) error {
	if param == nil || param.Symbol == "" {
		return types.ErrInvalidParam
	}
	cfg := c.cli.GetConfig()
	data, err := types.CallCreateTx(cfg, cfg.ExecName(tokenty.TokenX), "TokenRevokeCreate", param)
	if err != nil {
		return err
	}
	*result = hex.EncodeToString(data)
	return nil
}

// CreateRawTokenMintTx       mint Token
func (c *Jrpc) CreateRawTokenMintTx(param *tokenty.TokenMint, result *interface{}) error {
	if param == nil || param.Symbol == "" || param.Amount <= 0 {
		return types.ErrInvalidParam
	}
	cfg := c.cli.GetConfig()
	data, err := types.CallCreateTx(cfg, cfg.ExecName(tokenty.TokenX), "TokenMint", param)
	if err != nil {
		return err
	}
	*result = hex.EncodeToString(data)
	return nil
}

// CreateRawTokenBurnTx        burn Token
func (c *Jrpc) CreateRawTokenBurnTx(param *tokenty.TokenBurn, result *interface{}) error {
	if param == nil || param.Symbol == "" || param.Amount <= 0 {
		return types.ErrInvalidParam
	}
	cfg := c.cli.GetConfig()
	data, err := types.CallCreateTx(cfg, cfg.ExecName(tokenty.TokenX), "TokenBurn", param)
	if err != nil {
		return err
	}
	*result = hex.EncodeToString(data)
	return nil
}

// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"github.com/D-PlatformOperatingSystem/dpos/types"
	tokenty "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/token/types"
)

// Query_GetTokens   token
func (t *token) Query_GetTokens(in *tokenty.ReqTokens) (types.Message, error) {
	if in == nil {
		return nil, types.ErrInvalidParam
	}
	return t.getTokens(in)
}

// Query_GetTokenInfo   token
func (t *token) Query_GetTokenInfo(in *types.ReqString) (types.Message, error) {
	if in == nil {
		return nil, types.ErrInvalidParam
	}
	return t.getTokenInfo(in.GetData())
}

// Query_GetTotalAmount   token
func (t *token) Query_GetTotalAmount(in *types.ReqString) (types.Message, error) {
	if in == nil {
		return nil, types.ErrInvalidParam
	}
	ret, err := t.getTokenInfo(in.GetData())
	if err != nil {
		return nil, err
	}
	tokenInfo, ok := ret.(*tokenty.LocalToken)
	if !ok {
		return nil, types.ErrTypeAsset
	}
	return &types.TotalAmount{
		Total: tokenInfo.Total,
	}, nil
}

// Query_GetAddrReceiverforTokens   token
func (t *token) Query_GetAddrReceiverforTokens(in *tokenty.ReqAddrTokens) (types.Message, error) {
	if in == nil {
		return nil, types.ErrInvalidParam
	}
	return t.getAddrReceiverforTokens(in)
}

// Query_GetAccountTokenAssets      token
func (t *token) Query_GetAccountTokenAssets(in *tokenty.ReqAccountTokenAssets) (types.Message, error) {
	if in == nil {
		return nil, types.ErrInvalidParam
	}
	return t.getAccountTokenAssets(in)
}

// Query_GetTxByToken   token
func (t *token) Query_GetTxByToken(in *tokenty.ReqTokenTx) (types.Message, error) {
	if in == nil {
		return nil, types.ErrInvalidParam
	}
	if !subCfg.SaveTokenTxList {
		return nil, types.ErrActionNotSupport
	}
	return t.getTxByToken(in)
}

// Query_GetTokenHistory   token
func (t *token) Query_GetTokenHistory(in *types.ReqString) (types.Message, error) {
	if in == nil {
		return nil, types.ErrInvalidParam
	}
	rows, err := list(t.GetLocalDB(), "symbol", &tokenty.LocalLogs{Symbol: in.Data}, -1, 0)
	if err != nil {
		tokenlog.Error("Query_GetTokenHistory", "err", err)
		return nil, err
	}
	var replys tokenty.ReplyTokenLogs
	for _, row := range rows {
		o, ok := row.Data.(*tokenty.LocalLogs)
		if !ok {
			tokenlog.Error("Query_GetTokenHistory", "err", "bad row type")
			return nil, types.ErrTypeAsset
		}
		replys.Logs = append(replys.Logs, o)
	}
	return &replys, nil
}

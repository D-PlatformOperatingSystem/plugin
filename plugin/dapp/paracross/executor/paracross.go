// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"bytes"
	"encoding/hex"

	"github.com/D-PlatformOperatingSystem/dpos/common"
	"github.com/D-PlatformOperatingSystem/dpos/common/crypto"
	log "github.com/D-PlatformOperatingSystem/dpos/common/log/log15"
	drivers "github.com/D-PlatformOperatingSystem/dpos/system/dapp"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	"github.com/D-PlatformOperatingSystem/dpos/util"
	pt "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/paracross/types"
)

var (
	clog                    = log.New("module", "execs.paracross")
	enableParacrossTransfer = true
	driverName              = pt.ParaX
)

// Paracross exec
type Paracross struct {
	cryptoCli crypto.Crypto
	drivers.DriverBase
}

//Init paracross exec register
func Init(name string, cfg *types.DplatformOSConfig, sub []byte) {
	drivers.Register(cfg, GetName(), newParacross, cfg.GetDappFork(driverName, "Enable"))
	InitExecType()
	setPrefix()
}

//InitExecType ...
func InitExecType() {
	ety := types.LoadExecutorType(driverName)
	ety.InitFuncList(types.ListMethod(&Paracross{}))
}

//GetName return paracross name
func GetName() string {
	return newParacross().GetName()
}

func newParacross() drivers.Driver {
	c := &Paracross{}
	c.SetChild(c)
	c.SetExecutorType(types.LoadExecutorType(driverName))
	cli, err := crypto.New("bls")
	if err != nil {
		panic("paracross need bls sign register")
	}
	c.cryptoCli = cli
	return c
}

// GetDriverName return paracross driver name
func (c *Paracross) GetDriverName() string {
	return pt.ParaX
}

func (c *Paracross) checkTxGroup(tx *types.Transaction, index int) ([]*types.Transaction, error) {
	if tx.GroupCount >= 2 {
		txs, err := c.GetTxGroup(index)
		if err != nil {
			clog.Error("ParacrossActionAssetTransfer", "get tx group failed", err, "hash", hex.EncodeToString(tx.Hash()))
			return nil, err
		}
		return txs, nil
	}
	return nil, nil
}

func (c *Paracross) saveLocalParaTxs(tx *types.Transaction, isDel bool) (*types.LocalDBSet, error) {

	var payload pt.ParacrossAction
	err := types.Decode(tx.Payload, &payload)
	if err != nil {
		return nil, err
	}
	if payload.Ty != pt.ParacrossActionCommit || payload.GetCommit() == nil {
		return nil, nil
	}

	commit := payload.GetCommit()
	crossTxHashs, crossTxResult, err := getCrossTxHashs(c.GetAPI(), commit.Status)
	if err != nil {
		return nil, err
	}

	return c.udpateLocalParaTxs(commit.Status.Title, commit.Status.Height, crossTxHashs, crossTxResult, isDel)

}

//     commit tx  ， commitDone
func (c *Paracross) saveLocalParaTxsFork(commitDone *pt.ReceiptParacrossDone, isDel bool) (*types.LocalDBSet, error) {
	status := &pt.ParacrossNodeStatus{
		MainBlockHash:   commitDone.MainBlockHash,
		MainBlockHeight: commitDone.MainBlockHeight,
		Title:           commitDone.Title,
		Height:          commitDone.Height,
		BlockHash:       commitDone.BlockHash,
		TxResult:        commitDone.TxResult,
	}

	crossTxHashs, crossTxResult, err := getCrossTxHashs(c.GetAPI(), status)
	if err != nil {
		return nil, err
	}

	return c.udpateLocalParaTxs(commitDone.Title, commitDone.Height, crossTxHashs, crossTxResult, isDel)

}

func (c *Paracross) udpateLocalParaTxs(paraTitle string, paraHeight int64, crossTxHashs [][]byte, crossTxResult []byte, isDel bool) (*types.LocalDBSet, error) {
	var set types.LocalDBSet

	if len(crossTxHashs) == 0 {
		return &set, nil
	}

	for i := 0; i < len(crossTxHashs); i++ {
		success := util.BitMapBit(crossTxResult, uint32(i))

		paraTx, err := GetTx(c.GetAPI(), crossTxHashs[i])
		if err != nil {
			clog.Crit("paracross.Commit Load Tx failed", "para title", paraTitle,
				"para height", paraHeight, "para tx index", i, "error", err, "txHash",
				hex.EncodeToString(crossTxHashs[i]))
			return nil, err
		}

		var payload pt.ParacrossAction
		err = types.Decode(paraTx.Tx.Payload, &payload)
		if err != nil {
			clog.Crit("paracross.Commit Decode Tx failed", "para title", paraTitle,
				"para height", paraHeight, "para tx index", i, "error", err, "txHash",
				hex.EncodeToString(crossTxHashs[i]))
			return nil, err
		}
		if payload.Ty == pt.ParacrossActionCrossAssetTransfer {
			act, err := getCrossAction(payload.GetCrossAssetTransfer(), string(paraTx.Tx.Execer))
			if err != nil {
				clog.Crit("udpateLocalParaTxs getCrossAction failed", "error", err)
				return nil, err
			}
			//     ，            transfer
			if act == pt.ParacrossMainAssetTransfer || act == pt.ParacrossParaAssetWithdraw {
				kv, err := c.updateLocalAssetTransfer(paraTx.Tx, paraHeight, success, isDel)
				if err != nil {
					return nil, err
				}
				set.KV = append(set.KV, kv)
			}
			//     ，             withdraw
			if act == pt.ParacrossMainAssetWithdraw || act == pt.ParacrossParaAssetTransfer {
				asset, err := c.getCrossAssetTransferInfo(payload.GetCrossAssetTransfer(), paraTx.Tx, act)
				if err != nil {
					return nil, err
				}
				kv, err := c.initLocalAssetTransferDone(paraTx.Tx, asset, paraHeight, success, isDel)
				if err != nil {
					return nil, err
				}
				set.KV = append(set.KV, kv)
			}
		}

		if payload.Ty == pt.ParacrossActionAssetTransfer {
			kv, err := c.updateLocalAssetTransfer(paraTx.Tx, paraHeight, success, isDel)
			if err != nil {
				return nil, err
			}
			set.KV = append(set.KV, kv)
		} else if payload.Ty == pt.ParacrossActionAssetWithdraw {
			asset, err := c.getAssetTransferInfo(paraTx.Tx, payload.GetAssetWithdraw().Cointoken, true)
			if err != nil {
				return nil, err
			}
			kv, err := c.initLocalAssetTransferDone(paraTx.Tx, asset, paraHeight, success, isDel)
			if err != nil {
				return nil, err
			}
			set.KV = append(set.KV, kv)
		}
	}

	return &set, nil
}

func (c *Paracross) getAssetTransferInfo(tx *types.Transaction, coinToken string, isWithdraw bool) (*pt.ParacrossAsset, error) {
	exec := "coins"
	symbol := types.DOM
	if coinToken != "" {
		exec = "token"
		symbol = coinToken
	}
	amount, err := tx.Amount()
	if err != nil {
		return nil, err
	}

	asset := &pt.ParacrossAsset{
		From:       tx.From(),
		To:         tx.To,
		Amount:     amount,
		IsWithdraw: isWithdraw,
		TxHash:     common.ToHex(tx.Hash()),
		Height:     c.GetHeight(),
		Exec:       exec,
		Symbol:     symbol,
	}
	return asset, nil
}

func (c *Paracross) getCrossAssetTransferInfo(payload *pt.CrossAssetTransfer, tx *types.Transaction, act int64) (*pt.ParacrossAsset, error) {
	exec := payload.AssetExec
	symbol := payload.AssetSymbol
	if payload.AssetSymbol == "" {
		symbol = types.DOM
		exec = "coins"
	}

	amount, err := tx.Amount()
	if err != nil {
		return nil, err
	}

	isWithDraw := true
	if act == pt.ParacrossMainAssetTransfer || act == pt.ParacrossParaAssetTransfer {
		isWithDraw = false
	}

	asset := &pt.ParacrossAsset{
		From:       tx.From(),
		To:         tx.To,
		IsWithdraw: isWithDraw,
		Amount:     amount,
		TxHash:     common.ToHex(tx.Hash()),
		Height:     c.GetHeight(),
		Exec:       exec,
		Symbol:     symbol,
	}
	return asset, nil
}

func (c *Paracross) initLocalAssetTransfer(tx *types.Transaction, isDel bool, asset *pt.ParacrossAsset) (*types.KeyValue, error) {
	clog.Debug("para execLocal", "tx hash", hex.EncodeToString(tx.Hash()), "action name", log.Lazy{Fn: tx.ActionName})
	key := calcLocalAssetKey(tx.Hash())
	if isDel {
		c.GetLocalDB().Set(key, nil)
		return &types.KeyValue{Key: key, Value: nil}, nil
	}

	err := c.GetLocalDB().Set(key, types.Encode(asset))
	if err != nil {
		clog.Error("para execLocal", "set", hex.EncodeToString(tx.Hash()), "failed", err)
	}
	return &types.KeyValue{Key: key, Value: types.Encode(asset)}, nil
}

func (c *Paracross) initLocalAssetTransferDone(tx *types.Transaction, asset *pt.ParacrossAsset, paraHeight int64, success, isDel bool) (*types.KeyValue, error) {
	key := calcLocalAssetKey(tx.Hash())
	if isDel {
		c.GetLocalDB().Set(key, nil)
		return &types.KeyValue{Key: key, Value: nil}, nil
	}

	asset.ParaHeight = paraHeight
	asset.CommitDoneHeight = c.GetHeight()
	asset.Success = success

	err := c.GetLocalDB().Set(key, types.Encode(asset))
	if err != nil {
		clog.Error("para execLocal", "set", "", "failed", err)
	}
	return &types.KeyValue{Key: key, Value: types.Encode(asset)}, nil
}

func (c *Paracross) updateLocalAssetTransfer(tx *types.Transaction, paraHeight int64, success, isDel bool) (*types.KeyValue, error) {
	clog.Debug("para execLocal", "tx hash", hex.EncodeToString(tx.Hash()))
	key := calcLocalAssetKey(tx.Hash())

	var asset pt.ParacrossAsset
	v, err := c.GetLocalDB().Get(key)
	if err != nil {
		return nil, err
	}
	err = types.Decode(v, &asset)
	if err != nil {
		panic(err)
	}
	if !isDel {
		asset.ParaHeight = paraHeight

		asset.CommitDoneHeight = c.GetHeight()
		asset.Success = success
	} else {
		asset.ParaHeight = 0
		asset.CommitDoneHeight = 0
		asset.Success = false
	}
	c.GetLocalDB().Set(key, types.Encode(&asset))
	return &types.KeyValue{Key: key, Value: types.Encode(&asset)}, nil
}

//IsFriend call exec is same seariase exec
func (c *Paracross) IsFriend(myexec, writekey []byte, tx *types.Transaction) bool {
	//
	cfg := c.GetAPI().GetConfig()
	if cfg.IsPara() {
		return false
	}
	//friend
	if string(myexec) != c.GetDriverName() {
		return false
	}
	//          （tx      paracross）
	if string(types.GetRealExecName(tx.Execer)) != c.GetDriverName() {
		return false
	}
	//
	return c.allow(tx, 0) == nil
}

func (c *Paracross) allow(tx *types.Transaction, index int) error {
	//       :         title  asset-transfer/asset-withdraw
	// 1. user.p.${tilte}.${paraX}
	// 1. payload   actionType = t/w
	cfg := c.GetAPI().GetConfig()
	if !cfg.IsPara() && c.allowIsParaTx(tx.Execer) {
		var payload pt.ParacrossAction
		err := types.Decode(tx.Payload, &payload)
		if err != nil {
			return err
		}
		if payload.Ty == pt.ParacrossActionAssetTransfer || payload.Ty == pt.ParacrossActionAssetWithdraw {
			return nil
		}
		//       feature，           ，   none   ，            none
		//          ，                 ，         fork  ForkRootHash，             ，
		if cfg.IsDappFork(c.GetHeight(), pt.ParaX, pt.ForkCommitTx) {
			if payload.Ty == pt.ParacrossActionCommit || payload.Ty == pt.ParacrossActionNodeConfig ||
				payload.Ty == pt.ParacrossActionNodeGroupApply {
				return nil
			}
		}
		if cfg.IsDappFork(c.GetHeight(), pt.ParaX, pt.ForkParaAssetTransferRbk) {
			if payload.Ty == pt.ParacrossActionCrossAssetTransfer {
				return nil
			}
		}
	}
	return types.ErrNotAllow
}

// Allow add paracross allow rule
func (c *Paracross) Allow(tx *types.Transaction, index int) error {
	//
	err := c.DriverBase.Allow(tx, index)
	if err == nil {
		return nil
	}
	//paracross
	return c.allow(tx, index)
}

func (c *Paracross) allowIsParaTx(execer []byte) bool {
	if !bytes.HasPrefix(execer, types.ParaKey) {
		return false
	}
	count := 0
	index := 0
	s := len(types.ParaKey)
	for i := s; i < len(execer); i++ {
		if execer[i] == '.' {
			count++
			index = i
		}
	}
	if count == 1 && c.AllowIsSame(execer[index+1:]) {
		return true
	}
	return false
}

// CheckReceiptExecOk return true to check if receipt ty is ok
func (c *Paracross) CheckReceiptExecOk() bool {
	return true
}

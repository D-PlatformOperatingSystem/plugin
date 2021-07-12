// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import (
	"encoding/json"
	"reflect"

	log "github.com/D-PlatformOperatingSystem/dpos/common/log/log15"
	"github.com/D-PlatformOperatingSystem/dpos/types"
)

var (
	// ParaX paracross exec name
	ParaX = "paracross"
	glog  = log.New("module", ParaX)
	// ForkCommitTx main chain support paracross commit tx
	ForkCommitTx = "ForkParacrossCommitTx"
	// MainForkParacrossCommitTx            ForkCommitTx
	MainForkParacrossCommitTx = "mainForkParacrossCommitTx"
	// ForkLoopCheckCommitTxDone         done fork
	ForkLoopCheckCommitTxDone = "ForkLoopCheckCommitTxDone"
	// MainLoopCheckCommitTxDoneForkHeight        ，     ForkLoopCheckCommitTxDone
	MainLoopCheckCommitTxDoneForkHeight = "mainLoopCheckCommitTxDoneForkHeight"
	// ForkParaSelfConsStages
	ForkParaSelfConsStages = "ForkParaSelfConsStages"
	// ForkParaAssetTransferRbk
	ForkParaAssetTransferRbk = "ForkParaAssetTransferRbk"

	// ParaConsSubConf sub
	ParaConsSubConf = "consensus.sub.para"
	//ParaPrefixConsSubConf prefix
	ParaPrefixConsSubConf = "config." + ParaConsSubConf
	//ParaSelfConsInitConf self stage init config
	ParaSelfConsInitConf = "paraSelfConsInitDisable"
	//ParaSelfConsConfPreContract self consens enable string as ["0-100"] config pre stage contract
	ParaSelfConsConfPreContract = "selfConsensEnablePreContract"
	//ParaFilterIgnoreTxGroup adapt 6.1.0 to check para tx in group
	ParaFilterIgnoreTxGroup = "filterIgnoreParaTxGroup"
)

func init() {
	// init executor type
	types.AllowUserExec = append(types.AllowUserExec, []byte(ParaX))
	types.RegFork(ParaX, InitFork)
	types.RegExec(ParaX, InitExecutor)

}

//InitFork ...
func InitFork(cfg *types.DplatformOSConfig) {
	cfg.RegisterDappFork(ParaX, "Enable", 0)
	cfg.RegisterDappFork(ParaX, "ForkParacrossWithdrawFromParachain", 1298600)
	cfg.RegisterDappFork(ParaX, ForkCommitTx, 1850000)
	cfg.RegisterDappFork(ParaX, ForkLoopCheckCommitTxDone, 3230000)
	cfg.RegisterDappFork(ParaX, ForkParaAssetTransferRbk, 4500000)

	//
	cfg.RegisterDappFork(ParaX, ForkParaSelfConsStages, types.MaxHeight)
}

//InitExecutor ...
func InitExecutor(cfg *types.DplatformOSConfig) {
	types.RegistorExecutor(ParaX, NewType(cfg))
}

// GetExecName get para exec name
func GetExecName(cfg *types.DplatformOSConfig) string {
	return cfg.ExecName(ParaX)
}

// ParacrossType base paracross type
type ParacrossType struct {
	types.ExecTypeBase
}

// NewType get paracross type
func NewType(cfg *types.DplatformOSConfig) *ParacrossType {
	c := &ParacrossType{}
	c.SetChild(c)
	c.SetConfig(cfg)
	return c
}

// GetName
func (p *ParacrossType) GetName() string {
	return ParaX
}

// GetLogMap get receipt log map
func (p *ParacrossType) GetLogMap() map[int64]*types.LogInfo {
	return map[int64]*types.LogInfo{
		TyLogParacrossCommit:           {Ty: reflect.TypeOf(ReceiptParacrossCommit{}), Name: "LogParacrossCommit"},
		TyLogParacrossCommitDone:       {Ty: reflect.TypeOf(ReceiptParacrossDone{}), Name: "LogParacrossCommitDone"},
		TyLogParacrossCommitRecord:     {Ty: reflect.TypeOf(ReceiptParacrossRecord{}), Name: "LogParacrossCommitRecord"},
		TyLogParaAssetWithdraw:         {Ty: reflect.TypeOf(types.ReceiptAccountTransfer{}), Name: "LogParaAssetWithdraw"},
		TyLogParaAssetTransfer:         {Ty: reflect.TypeOf(types.ReceiptAccountTransfer{}), Name: "LogParaAssetTransfer"},
		TyLogParaAssetDeposit:          {Ty: reflect.TypeOf(types.ReceiptAccountTransfer{}), Name: "LogParaAssetDeposit"},
		TyLogParaCrossAssetTransfer:    {Ty: reflect.TypeOf(types.ReceiptAccountTransfer{}), Name: "LogParaCrossAssetTransfer"},
		TyLogParacrossMiner:            {Ty: reflect.TypeOf(ReceiptParacrossMiner{}), Name: "LogParacrossMiner"},
		TyLogParaNodeConfig:            {Ty: reflect.TypeOf(ReceiptParaNodeConfig{}), Name: "LogParaNodeConfig"},
		TyLogParaNodeStatusUpdate:      {Ty: reflect.TypeOf(ReceiptParaNodeAddrStatUpdate{}), Name: "LogParaNodeAddrStatUpdate"},
		TyLogParaNodeGroupAddrsUpdate:  {Ty: reflect.TypeOf(types.ReceiptConfig{}), Name: "LogParaNodeGroupAddrsUpdate"},
		TyLogParaNodeVoteDone:          {Ty: reflect.TypeOf(ReceiptParaNodeVoteDone{}), Name: "LogParaNodeVoteDone"},
		TyLogParaNodeGroupConfig:       {Ty: reflect.TypeOf(ReceiptParaNodeGroupConfig{}), Name: "LogParaNodeGroupConfig"},
		TyLogParaNodeGroupStatusUpdate: {Ty: reflect.TypeOf(ReceiptParaNodeGroupConfig{}), Name: "LogParaNodeGroupStatusUpdate"},
		TyLogParaSelfConsStageConfig:   {Ty: reflect.TypeOf(ReceiptSelfConsStageConfig{}), Name: "LogParaSelfConsStageConfig"},
		TyLogParaStageVoteDone:         {Ty: reflect.TypeOf(ReceiptSelfConsStageVoteDone{}), Name: "LogParaSelfConfStageVoteDoen"},
		TyLogParaStageGroupUpdate:      {Ty: reflect.TypeOf(ReceiptSelfConsStagesUpdate{}), Name: "LogParaSelfConfStagesUpdate"},
		TyLogParaBindMinerAddr:         {Ty: reflect.TypeOf(ReceiptParaBindMinerInfo{}), Name: "TyLogParaBindMinerAddrUpdate"},
		TyLogParaBindMinerNode:         {Ty: reflect.TypeOf(ReceiptParaNodeBindListUpdate{}), Name: "TyLogParaBindNodeListUpdate"},
	}
}

// GetTypeMap get action type
func (p *ParacrossType) GetTypeMap() map[string]int32 {
	return map[string]int32{
		"Commit":             ParacrossActionCommit,
		"Miner":              ParacrossActionMiner,
		"AssetTransfer":      ParacrossActionAssetTransfer,
		"AssetWithdraw":      ParacrossActionAssetWithdraw,
		"Transfer":           ParacrossActionTransfer,
		"Withdraw":           ParacrossActionWithdraw,
		"TransferToExec":     ParacrossActionTransferToExec,
		"CrossAssetTransfer": ParacrossActionCrossAssetTransfer,
		"NodeConfig":         ParacrossActionNodeConfig,
		"NodeGroupConfig":    ParacrossActionNodeGroupApply,
		"SelfStageConfig":    ParacrossActionSelfStageConfig,
		"ParaBindMiner":      ParacrossActionParaBindMiner,
	}
}

// GetPayload paracross get action payload
func (p *ParacrossType) GetPayload() types.Message {
	return &ParacrossAction{}
}

// CreateTx paracross create tx by different action
func (p ParacrossType) CreateTx(action string, message json.RawMessage) (*types.Transaction, error) {
	cfg := p.GetConfig()
	//    ParacrossAssetTransfer  ，   AssetTransfer　
	if action == "ParacrossAssetTransfer" || action == "ParacrossAssetWithdraw" {
		var param types.CreateTx
		err := json.Unmarshal(message, &param)
		if err != nil {
			glog.Error("CreateTx", "Error", err)
			return nil, types.ErrInvalidParam
		}
		return CreateRawAssetTransferTx(cfg, &param)
	} else if action == "Transfer" || action == "Withdraw" || action == "TransferToExec" {
		//transfer/withdraw/toExec           tx.to
		return p.CreateRawTransferTx(action, message)
	}
	return p.ExecTypeBase.CreateTx(action, message)
}

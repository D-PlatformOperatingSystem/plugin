// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"fmt"

	log "github.com/D-PlatformOperatingSystem/dpos/common/log/log15"
	drivers "github.com/D-PlatformOperatingSystem/dpos/system/dapp"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	pkt "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/pokerbull/types"
)

var logger = log.New("module", "execs.pokerbull")

// Init
func Init(name string, cfg *types.DplatformOSConfig, sub []byte) {
	drivers.Register(cfg, newPBGame().GetName(), newPBGame, cfg.GetDappFork(driverName, "Enable"))
	InitExecType()
}

var driverName = pkt.PokerBullX

//InitExecType ...
func InitExecType() {
	ety := types.LoadExecutorType(driverName)
	ety.InitFuncList(types.ListMethod(&PokerBull{}))
}

// PokerBull
type PokerBull struct {
	drivers.DriverBase
}

func newPBGame() drivers.Driver {
	t := &PokerBull{}
	t.SetChild(t)
	t.SetExecutorType(types.LoadExecutorType(driverName))
	return t
}

// GetName
func GetName() string {
	return newPBGame().GetName()
}

// GetDriverName
func (g *PokerBull) GetDriverName() string {
	return pkt.PokerBullX
}

func calcPBGameAddrPrefix(addr string) []byte {
	key := fmt.Sprintf("LODB-pokerbull-addr:%s:", addr)
	return []byte(key)
}

func calcPBGameAddrKey(addr string, index int64) []byte {
	key := fmt.Sprintf("LODB-pokerbull-addr:%s:%018d", addr, index)
	return []byte(key)
}

func calcPBGameStatusPrefix(status int32) []byte {
	key := fmt.Sprintf("LODB-pokerbull-status-index:%d:", status)
	return []byte(key)
}

func calcPBGameStatusKey(status int32, index int64) []byte {
	key := fmt.Sprintf("LODB-pokerbull-status-index:%d:%018d", status, index)
	return []byte(key)
}

func calcPBGameStatusAndPlayerKey(status, player int32, value, index int64) []byte {
	key := fmt.Sprintf("LODB-pokerbull-status:%d:%d:%015d:%018d", status, player, value, index)
	return []byte(key)
}

func calcPBGameStatusAndPlayerPrefix(status, player int32, value int64) []byte {
	var key string
	if value == 0 {
		key = fmt.Sprintf("LODB-pokerbull-status:%d:%d:", status, player)
	} else {
		key = fmt.Sprintf("LODB-pokerbull-status:%d:%d:%015d", status, player, value)
	}

	return []byte(key)
}

func addPBGameStatusIndexKey(status int32, gameID string, index int64) *types.KeyValue {
	kv := &types.KeyValue{}
	kv.Key = calcPBGameStatusKey(status, index)
	record := &pkt.PBGameIndexRecord{
		GameId: gameID,
		Index:  index,
	}
	kv.Value = types.Encode(record)
	return kv
}

func delPBGameStatusIndexKey(status int32, index int64) *types.KeyValue {
	kv := &types.KeyValue{}
	kv.Key = calcPBGameStatusKey(status, index)
	kv.Value = nil
	return kv
}

func addPBGameAddrIndexKey(status int32, addr, gameID string, index int64) *types.KeyValue {
	kv := &types.KeyValue{}
	kv.Key = calcPBGameAddrKey(addr, index)
	record := &pkt.PBGameRecord{
		GameId: gameID,
		Status: status,
		Index:  index,
	}
	kv.Value = types.Encode(record)
	return kv
}

func delPBGameAddrIndexKey(addr string, index int64) *types.KeyValue {
	kv := &types.KeyValue{}
	kv.Key = calcPBGameAddrKey(addr, index)
	kv.Value = nil
	return kv
}

func addPBGameStatusAndPlayer(status int32, player int32, value, index int64, gameID string) *types.KeyValue {
	kv := &types.KeyValue{}
	kv.Key = calcPBGameStatusAndPlayerKey(status, player, value, index)
	record := &pkt.PBGameIndexRecord{
		GameId: gameID,
		Index:  index,
	}
	kv.Value = types.Encode(record)
	return kv
}

func delPBGameStatusAndPlayer(status int32, player int32, value, index int64) *types.KeyValue {
	kv := &types.KeyValue{}
	kv.Key = calcPBGameStatusAndPlayerKey(status, player, value, index)
	kv.Value = nil
	return kv
}

// CheckReceiptExecOk return true to check if receipt ty is ok
func (g *PokerBull) CheckReceiptExecOk() bool {
	return true
}

// ExecutorOrder   localdb EnableRead
func (g *PokerBull) ExecutorOrder() int64 {
	cfg := g.GetAPI().GetConfig()
	if cfg.IsFork(g.GetHeight(), "ForkLocalDBAccess") {
		return drivers.ExecLocalSameTime
	}
	return g.DriverBase.ExecutorOrder()
}

// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kvmvccmavl

import (
	"sync"

	"github.com/D-PlatformOperatingSystem/dpos/common"
	dbm "github.com/D-PlatformOperatingSystem/dpos/common/db"
	"github.com/D-PlatformOperatingSystem/dpos/queue"
	mavl "github.com/D-PlatformOperatingSystem/dpos/system/store/mavl/db"
	"github.com/D-PlatformOperatingSystem/dpos/types"
)

const (
	//  store/mavl       ，
	hashNodePrefix        = "_mh_"
	leafNodePrefix        = "_mb_"
	leafKeyCountPrefix    = "..mk.."
	oldLeafKeyCountPrefix = "..mok.."
)

// MavlStore mavl store struct
type MavlStore struct {
	db      dbm.DB
	trees   *sync.Map
	treeCfg *mavl.TreeConfig
}

// NewMavl new mavl store module
func NewMavl(sub *subMavlConfig, db dbm.DB) *MavlStore {
	var subcfg subMavlConfig
	if sub != nil {
		subcfg.EnableMavlPrefix = sub.EnableMavlPrefix
		subcfg.EnableMVCC = sub.EnableMVCC
		subcfg.EnableMavlPrune = sub.EnableMavlPrune
		subcfg.PruneHeight = sub.PruneHeight
		subcfg.EnableMemTree = sub.EnableMemTree
		subcfg.EnableMemVal = sub.EnableMemVal
		subcfg.TkCloseCacheLen = sub.TkCloseCacheLen
	}

	//
	if subcfg.EnableMavlPrune {
		subcfg.EnableMavlPrefix = subcfg.EnableMavlPrune
	}
	treeCfg := &mavl.TreeConfig{
		EnableMavlPrefix: subcfg.EnableMavlPrefix,
		EnableMVCC:       subcfg.EnableMVCC,
		EnableMavlPrune:  subcfg.EnableMavlPrune,
		PruneHeight:      subcfg.PruneHeight,
		EnableMemTree:    subcfg.EnableMemTree,
		EnableMemVal:     subcfg.EnableMemVal,
		TkCloseCacheLen:  subcfg.TkCloseCacheLen,
	}
	mavls := &MavlStore{db, &sync.Map{}, treeCfg}
	return mavls
}

// Close close mavl store
func (mavls *MavlStore) Close() {
	mavl.ClosePrune()
	kmlog.Info("store mavl closed")
}

// Set set k v to mavl store db; sync is true represent write sync
func (mavls *MavlStore) Set(datas *types.StoreSet, sync bool) ([]byte, error) {
	return mavl.SetKVPair(mavls.db, datas, sync, mavls.treeCfg)
}

// Get get values by keys
func (mavls *MavlStore) Get(datas *types.StoreGet) [][]byte {
	var tree *mavl.Tree
	var err error
	values := make([][]byte, len(datas.Keys))
	search := string(datas.StateHash)
	if data, ok := mavls.trees.Load(search); ok {
		tree = data.(*mavl.Tree)
	} else {
		tree = mavl.NewTree(mavls.db, true, mavls.treeCfg)
		//get
		//tree.SetBlockHeight(datas.Height)
		err = tree.Load(datas.StateHash)
		kmlog.Debug("store mavl get tree", "err", err, "StateHash", common.ToHex(datas.StateHash))
	}
	if err == nil {
		for i := 0; i < len(datas.Keys); i++ {
			_, value, exit := tree.Get(datas.Keys[i])
			if exit {
				values[i] = value
			}
		}
	}
	return values
}

// MemSet set keys values to memcory mavl, return root hash and error
func (mavls *MavlStore) MemSet(datas *types.StoreSet, sync bool) ([]byte, error) {
	beg := types.Now()
	defer func() {
		kmlog.Debug("mavl MemSet", "cost", types.Since(beg))
	}()
	if len(datas.KV) == 0 {
		kmlog.Info("store mavl memset,use preStateHash as stateHash for kvset is null")
		mavls.trees.Store(string(datas.StateHash), nil)
		return datas.StateHash, nil
	}
	tree := mavl.NewTree(mavls.db, sync, mavls.treeCfg)
	tree.SetBlockHeight(datas.Height)
	err := tree.Load(datas.StateHash)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(datas.KV); i++ {
		tree.Set(datas.KV[i].Key, datas.KV[i].Value)
	}
	hash := tree.Hash()
	mavls.trees.Store(string(hash), tree)
	return hash, nil
}

// MemSetUpgrade   hash
func (mavls *MavlStore) MemSetUpgrade(datas *types.StoreSet, sync bool) ([]byte, error) {
	beg := types.Now()
	defer func() {
		kmlog.Debug("mavl MemSet", "cost", types.Since(beg))
	}()
	if len(datas.KV) == 0 {
		kmlog.Info("store mavl memset,use preStateHash as stateHash for kvset is null")
		return datas.StateHash, nil
	}
	tree := mavl.NewTree(mavls.db, sync, mavls.treeCfg)
	tree.SetBlockHeight(datas.Height)
	err := tree.Load(datas.StateHash)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(datas.KV); i++ {
		tree.Set(datas.KV[i].Key, datas.KV[i].Value)
	}
	hash := tree.Hash()
	return hash, nil
}

// Commit convert memcory mavl to storage db
func (mavls *MavlStore) Commit(req *types.ReqHash) ([]byte, error) {
	beg := types.Now()
	defer func() {
		kmlog.Debug("mavl Commit", "cost", types.Since(beg))
	}()
	tree, ok := mavls.trees.Load(string(req.Hash))
	if !ok {
		kmlog.Error("store mavl commit", "err", types.ErrHashNotFound)
		return nil, types.ErrHashNotFound
	}
	if tree == nil {
		kmlog.Info("store mavl commit,do nothing for kvset is null")
		mavls.trees.Delete(string(req.Hash))
		return req.Hash, nil
	}
	hash := tree.(*mavl.Tree).Save()
	if hash == nil {
		kmlog.Error("store mavl commit", "err", types.ErrHashNotFound)
		return nil, types.ErrDataBaseDamage
	}
	mavls.trees.Delete(string(req.Hash))
	return req.Hash, nil
}

// Rollback       mavl
func (mavls *MavlStore) Rollback(req *types.ReqHash) ([]byte, error) {
	beg := types.Now()
	defer func() {
		kmlog.Debug("Rollback", "cost", types.Since(beg))
	}()
	_, ok := mavls.trees.Load(string(req.Hash))
	if !ok {
		kmlog.Error("store mavl rollback", "err", types.ErrHashNotFound)
		return nil, types.ErrHashNotFound
	}
	mavls.trees.Delete(string(req.Hash))
	return req.Hash, nil
}

// IterateRangeByStateHash       ； statehash：    hash, start：     key, end:    key, ascending：  ，  , fn
func (mavls *MavlStore) IterateRangeByStateHash(statehash []byte, start []byte, end []byte, ascending bool, fn func(key, value []byte) bool) {
	mavl.IterateRangeByStateHash(mavls.db, statehash, start, end, ascending, mavls.treeCfg, fn)
}

// ProcEvent not support message
func (mavls *MavlStore) ProcEvent(msg queue.Message) {
	msg.ReplyErr("Store", types.ErrActionNotSupport)
}

// Del ...
func (mavls *MavlStore) Del(req *types.StoreDel) ([]byte, error) {
	//not support
	return nil, nil
}

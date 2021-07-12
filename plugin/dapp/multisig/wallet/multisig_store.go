// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wallet

import (
	"github.com/D-PlatformOperatingSystem/dpos/common/db"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	wcom "github.com/D-PlatformOperatingSystem/dpos/wallet/common"
	mtypes "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/multisig/types"
)

func newStore(db db.DB) *multisigStore {
	return &multisigStore{Store: wcom.NewStore(db)}
}

// multisigStore
type multisigStore struct {
	*wcom.Store
}

//    owner
func (store *multisigStore) listOwnerAttrsByAddr(addr string) (*mtypes.OwnerAttrs, error) {
	if len(addr) == 0 {
		bizlog.Error("listMultisigAddrByOwnerAddr addr is nil")
		return nil, types.ErrInvalidParam
	}

	ownerAttrByte, err := store.Get(calcMultisigAddr(addr))
	if err != nil {
		bizlog.Error("listMultisigAddrByOwnerAddr", "addr", addr, "db Get error ", err)
		if err == db.ErrNotFoundInDb {
			return nil, types.ErrNotFound
		}
		return nil, err
	}
	if nil == ownerAttrByte || len(ownerAttrByte) == 0 {
		return nil, types.ErrNotFound
	}
	var ownerAttrs mtypes.OwnerAttrs
	err = types.Decode(ownerAttrByte, &ownerAttrs)
	if err != nil {
		bizlog.Error("listMultisigAddrByOwnerAddr", "proto.Unmarshal err:", err)
		return nil, types.ErrUnmarshal
	}
	return &ownerAttrs, nil
}

//
func (store *multisigStore) listOwnerAttrs() (*mtypes.OwnerAttrs, error) {

	list := store.NewListHelper()
	ownerbytes := list.PrefixScan(calcPrefixMultisigAddr())
	if len(ownerbytes) == 0 {
		bizlog.Error("listOwnerAttrs is null")
		return nil, types.ErrNotFound
	}
	var replayOwnerAttrs mtypes.OwnerAttrs
	for _, ownerattrbytes := range ownerbytes {
		var ownerAttrs mtypes.OwnerAttrs
		err := types.Decode(ownerattrbytes, &ownerAttrs)
		if err != nil {
			bizlog.Error("listOwnerAttrs", "Decode err", err)
			continue
		}
		replayOwnerAttrs.Items = append(replayOwnerAttrs.Items, ownerAttrs.Items...)
	}
	return &replayOwnerAttrs, nil
}

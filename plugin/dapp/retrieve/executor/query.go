// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"github.com/D-PlatformOperatingSystem/dpos/types"
	rt "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/retrieve/types"
)

// Query_GetRetrieveInfo get retrieve state
func (r *Retrieve) Query_GetRetrieveInfo(in *rt.ReqRetrieveInfo) (types.Message, error) {
	rlog.Debug("Retrieve Query", "backupaddr", in.BackupAddress, "defaddr", in.DefaultAddress)
	info, err := getRetrieveInfo(r.GetLocalDB(), in.BackupAddress, in.DefaultAddress)
	if info == nil {
		return nil, err
	}
	if info.Status == retrievePrepare {
		info.RemainTime = info.DelayPeriod - (r.GetBlockTime() - info.PrepareTime)
		if info.RemainTime < 0 {
			info.RemainTime = 0
		}
	}

	//    asset     ，     asset
	if info.Status == retrievePerform && in.GetAssetExec() != "" {
		// retrievePerform   ，
		// 1     , 2 fork      coins
		// 2 fork       coins      ,
		// localdb not support PrefixCount
		//              ，

		asset, _ := getRetrieveAsset(r.GetLocalDB(), in.BackupAddress, in.DefaultAddress, in.AssetExec, in.AssetSymbol)
		if asset != nil {
			return asset, nil
		}

		// 1
		info.Status = retrievePrepare
		info.RemainTime = zeroRemainTime
		return info, nil

	}
	return info, nil
}

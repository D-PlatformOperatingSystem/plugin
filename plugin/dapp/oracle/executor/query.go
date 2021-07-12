/*
 * Copyright D-Platform Corp. 2018 All Rights Reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package executor

import (
	"github.com/D-PlatformOperatingSystem/dpos/types"
	oty "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/oracle/types"
)

// statedb
func (o *oracle) Query_QueryOraclesByIDs(in *oty.QueryOracleInfos) (types.Message, error) {
	return getOracleLisByIDs(o.GetStateDB(), in)
}

//      ids
func (o *oracle) Query_QueryEventIDsByStatus(in *oty.QueryEventID) (types.Message, error) {
	eventIds, err := getEventIDListByStatus(o.GetLocalDB(), in.Status, in.EventID)
	if err != nil {
		return nil, err
	}

	return eventIds, nil
}

//
func (o *oracle) Query_QueryEventIDsByAddrAndStatus(in *oty.QueryEventID) (types.Message, error) {
	eventIds, err := getEventIDListByAddrAndStatus(o.GetLocalDB(), in.Addr, in.Status, in.EventID)
	if err != nil {
		return nil, err
	}

	return eventIds, nil
}

//
func (o *oracle) Query_QueryEventIDsByTypeAndStatus(in *oty.QueryEventID) (types.Message, error) {
	eventIds, err := getEventIDListByTypeAndStatus(o.GetLocalDB(), in.Type, in.Status, in.EventID)
	if err != nil {
		return nil, err
	}
	return eventIds, nil
}

// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

//Issuance op
const (
	IssuanceActionCreate = 1 + iota //
	IssuanceActionDebt              //
	IssuanceActionRepay             //
	IssuanceActionFeed              //
	IssuanceActionClose             //
	IssuanceActionManage            //

	//log for Issuance
	TyLogIssuanceCreate = 741
	TyLogIssuanceDebt   = 742
	TyLogIssuanceRepay  = 743
	TyLogIssuanceFeed   = 745
	TyLogIssuanceClose  = 746
)

// Issuance name
const (
	IssuanceX                   = "issuance"
	CCNYTokenName               = "CCNY"
	IssuancePreLiquidationRatio = 11000 //TODO      ，         ccny  110%
)

//Issuance status
const (
	IssuanceStatusCreated = 1 + iota
	IssuanceStatusClose
)

//status ...
const (
	IssuanceUserStatusCreate = 1 + iota
	IssuanceUserStatusWarning
	IssuanceUserStatusSystemLiquidate
	IssuanceUserStatusExpire
	IssuanceUserStatusExpireLiquidate
	IssuanceUserStatusClose
)

//type ...
const (
	PriceFeedKey = "issuance-price-feed"
	GuarantorKey = "issuance-guarantor"
	ManageKey    = "issuance-manage"
	FundKey      = "issuance-fund"
)

//fork ...
var (
	ForkIssuanceTableUpdate = "ForkIssuanceTableUpdate"
)

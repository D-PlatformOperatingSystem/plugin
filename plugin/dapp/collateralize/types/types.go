// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

//Collateralize op
const (
	CollateralizeActionCreate = 1 + iota
	CollateralizeActionBorrow
	CollateralizeActionRepay
	CollateralizeActionAppend
	CollateralizeActionFeed
	CollateralizeActionRetrieve
	CollateralizeActionManage

	//log for Collateralize
	TyLogCollateralizeCreate   = 731
	TyLogCollateralizeBorrow   = 732
	TyLogCollateralizeRepay    = 733
	TyLogCollateralizeAppend   = 734
	TyLogCollateralizeFeed     = 735
	TyLogCollateralizeRetrieve = 736
)

// Collateralize name
const (
	CollateralizeX                   = "collateralize"
	CCNYTokenName                    = "CCNY"
	CollateralizePreLiquidationRatio = 1.1 * 1e4 //TODO      ，         ccny  110%
)

//Collateralize status
const (
	CollateralizeStatusCreated = 1 + iota
	CollateralizeStatusClose
)

//     dpos
//const (
//	CollateralizeAssetTypeDpos = 1 + iota
//	CollateralizeAssetTypeBtc
//	CollateralizeAssetTypeEth
//)

//collater ...
const (
	CollateralizeUserStatusCreate = 1 + iota
	CollateralizeUserStatusWarning
	CollateralizeUserStatusSystemLiquidate
	CollateralizeUserStatusExpire
	CollateralizeUserStatusExpireLiquidate
	CollateralizeUserStatusClose
)

//fork ...
var (
	ForkCollateralizeTableUpdate = "ForkCollateralizeTableUpdate"
)

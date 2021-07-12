// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

//game action ty
const (
	GuessGameActionStart   = 5
	GuessGameActionBet     = 6
	GuessGameActionStopBet = 7
	GuessGameActionAbort   = 8
	GuessGameActionPublish = 9
	GuessGameActionQuery   = 10

	GuessGameStatusStart   = 11
	GuessGameStatusBet     = 12
	GuessGameStatusStopBet = 13
	GuessGameStatusAbort   = 14
	GuessGameStatusPublish = 15
	GuessGameStatusTimeOut = 16
)

//game log ty
const (
	TyLogGuessGameStart   = 901
	TyLogGuessGameBet     = 902
	TyLogGuessGameStopBet = 903
	TyLogGuessGameAbort   = 904
	TyLogGuessGamePublish = 905
	TyLogGuessGameTimeout = 906
)

//
//   github     ，        ,
//      ，
var (
	GuessX      = "guess"
	ExecerGuess = []byte(GuessX)
)

const (
	//FuncNameQueryGamesByIDs func name
	FuncNameQueryGamesByIDs = "QueryGamesByIDs"

	//FuncNameQueryGameByID func name
	FuncNameQueryGameByID = "QueryGameByID"

	//FuncNameQueryGameByAddr func name
	FuncNameQueryGameByAddr = "QueryGamesByAddr"

	//FuncNameQueryGameByStatus func name
	FuncNameQueryGameByStatus = "QueryGamesByStatus"

	//FuncNameQueryGameByAdminAddr func name
	FuncNameQueryGameByAdminAddr = "QueryGamesByAdminAddr"

	//FuncNameQueryGameByAddrStatus func name
	FuncNameQueryGameByAddrStatus = "QueryGamesByAddrStatus"

	//FuncNameQueryGameByAdminStatus func name
	FuncNameQueryGameByAdminStatus = "QueryGamesByAdminStatus"

	//FuncNameQueryGameByCategoryStatus func name
	FuncNameQueryGameByCategoryStatus = "QueryGamesByCategoryStatus"

	//CreateStartTx
	CreateStartTx = "Start"

	//CreateBetTx
	CreateBetTx = "Bet"

	//CreateStopBetTx
	CreateStopBetTx = "StopBet"

	//CreatePublishTx
	CreatePublishTx = "Publish"

	//CreateAbortTx
	CreateAbortTx = "Abort"
)

const (
	//DevShareAddr default value
	DevShareAddr = "1D6RFZNp2rh6QdbcZ1d7RWuBUz61We6SD7"

	//PlatformShareAddr default value
	PlatformShareAddr = "1PHtChNt3UcfssR7v7trKSk3WJtAWjKjjX"
)

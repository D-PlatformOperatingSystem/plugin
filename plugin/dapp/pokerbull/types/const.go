// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import "github.com/D-PlatformOperatingSystem/dpos/types"

//game action ty
const (
	PBGameActionStart = iota + 1
	PBGameActionContinue
	PBGameActionQuit
	PBGameActionQuery
	PBGameActionPlay
)

const (
	// PlayStyleDefault
	PlayStyleDefault = iota + 1
	// PlayStyleDealer
	PlayStyleDealer
)

const (
	// TyLogPBGameStart log for start PBgame
	TyLogPBGameStart = 721
	// TyLogPBGameContinue log for continue PBgame
	TyLogPBGameContinue = 722
	// TyLogPBGameQuit log for quit PBgame
	TyLogPBGameQuit = 723
	// TyLogPBGameQuery log for query PBgame
	TyLogPBGameQuery = 724
	// TyLogPBGamePlay log for play PBgame
	TyLogPBGamePlay = 725
)

//
//   github     ，        ,
//      ，
var (
	JRPCName        = "pokerbull"
	PokerBullX      = "pokerbull"
	ExecerPokerBull = []byte(PokerBullX)
)

const (
	// FuncNameQueryGameListByIDs   id    game
	FuncNameQueryGameListByIDs = "QueryGameListByIDs"
	// FuncNameQueryGameByID   id  game
	FuncNameQueryGameByID = "QueryGameByID"
	// FuncNameQueryGameByAddr       game
	FuncNameQueryGameByAddr = "QueryGameByAddr"
	// FuncNameQueryGameByStatus   status  game
	FuncNameQueryGameByStatus = "QueryGameByStatus"
	// FuncNameQueryGameByRound
	FuncNameQueryGameByRound = "QueryGameByRound"
	// CreateStartTx
	CreateStartTx = "Start"
	// CreateContinueTx
	CreateContinueTx = "Continue"
	// CreateQuitTx
	CreateQuitTx = "Quit"
	// CreatePlayTx
	CreatePlayTx = "Play"
)

const (
	// ListDESC
	ListDESC = int32(0)
	// DefaultCount
	DefaultCount = int32(20)
	// MaxPlayerNum
	MaxPlayerNum = 5
	// MinPlayerNum
	MinPlayerNum = 2
	// MinPlayValue
	MinPlayValue = 10 * types.Coin
	// DefaultStyle
	DefaultStyle = PlayStyleDefault
	// PlatformAddress
	PlatformAddress = "1PHtChNt3UcfssR7v7trKSk3WJtAWjKjjX"
	// PlatformFee
	PlatformFee = int64(0.005 * float64(types.Coin))
	// DeveloperAddress
	DeveloperAddress = "1D6RFZNp2rh6QdbcZ1d7RWuBUz61We6SD7"
	// DeveloperFee
	DeveloperFee = int64(0.005 * float64(types.Coin))
	// WinnerReturn
	WinnerReturn = types.Coin - DeveloperFee - PlatformFee
	// PlatformSignAddress
	PlatformSignAddress = "1Geb4ppNiAwMKKyrJgcis3JA57FkqsXvdR"
)

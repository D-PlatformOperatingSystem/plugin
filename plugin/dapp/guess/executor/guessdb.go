// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/D-PlatformOperatingSystem/dpos/account"
	"github.com/D-PlatformOperatingSystem/dpos/common"
	"github.com/D-PlatformOperatingSystem/dpos/common/db"
	dbm "github.com/D-PlatformOperatingSystem/dpos/common/db"
	"github.com/D-PlatformOperatingSystem/dpos/common/db/table"
	"github.com/D-PlatformOperatingSystem/dpos/system/dapp"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	gty "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/guess/types"
)

const (
	//ListDESC
	ListDESC = int32(0)

	//ListASC
	ListASC = int32(1)

	//DefaultCount
	DefaultCount = int32(10)

	//DefaultCategory
	DefaultCategory = "default"

	//MaxBetsOneTime
	MaxBetsOneTime = 10000e8

	//MaxBetsNumber
	MaxBetsNumber = 10000000e8

	//MaxBetHeight
	MaxBetHeight = 1000000

	//MaxExpireHeight
	MaxExpireHeight = 1000000
)

//Action
type Action struct {
	coinsAccount *account.DB
	db           dbm.KV
	txhash       []byte
	fromaddr     string
	blocktime    int64
	height       int64
	execaddr     string
	localDB      dbm.KVDB
	index        int
	mainHeight   int64
}

//NewAction   Action
func NewAction(guess *Guess, tx *types.Transaction, index int) *Action {
	hash := tx.Hash()
	fromAddr := tx.From()

	return &Action{
		coinsAccount: guess.GetCoinsAccount(),
		db:           guess.GetStateDB(),
		txhash:       hash,
		fromaddr:     fromAddr,
		blocktime:    guess.GetBlockTime(),
		height:       guess.GetHeight(),
		execaddr:     dapp.ExecAddress(string(tx.Execer)),
		localDB:      guess.GetLocalDB(),
		index:        index,
		mainHeight:   guess.GetMainHeight(),
	}
}

//CheckExecAccountBalance      Guess
func (action *Action) CheckExecAccountBalance(fromAddr string, ToFrozen, ToActive int64) bool {
	acc := action.coinsAccount.LoadExecAccount(fromAddr, action.execaddr)
	if acc.GetBalance() >= ToFrozen && acc.GetFrozen() >= ToActive {
		return true
	}
	return false
}

//Key State         Key
func Key(id string) (key []byte) {
	//key = append(key, []byte("mavl-"+types.ExecName(pkt.GuessX)+"-")...)
	key = append(key, []byte("mavl-"+gty.GuessX+"-")...)
	key = append(key, []byte(id)...)
	return key
}

//queryGameInfos     id
func queryGameInfos(kvdb db.KVDB, infos *gty.QueryGuessGameInfos) (types.Message, error) {
	var games []*gty.GuessGame
	gameTable := gty.NewGuessGameTable(kvdb)
	query := gameTable.GetQuery(kvdb)

	for i := 0; i < len(infos.GameIDs); i++ {
		rows, err := query.ListIndex("gameid", []byte(infos.GameIDs[i]), nil, 1, 0)
		if err != nil {
			return nil, err
		}

		game := rows[0].Data.(*gty.GuessGame)
		games = append(games, game)
	}
	return &gty.ReplyGuessGameInfos{Games: games}, nil
}

//queryGameInfo   gameid  game
func queryGameInfo(kvdb db.KVDB, gameID []byte) (*gty.GuessGame, error) {
	gameTable := gty.NewGuessGameTable(kvdb)
	query := gameTable.GetQuery(kvdb)
	rows, err := query.ListIndex("gameid", gameID, nil, 1, 0)
	if err != nil {
		return nil, err
	}

	game := rows[0].Data.(*gty.GuessGame)

	return game, nil
}

//queryUserTableData   user
func queryUserTableData(query *table.Query, indexName string, prefix, primaryKey []byte) (types.Message, error) {
	rows, err := query.ListIndex(indexName, prefix, primaryKey, DefaultCount, 0)
	if err != nil {
		return nil, err
	}

	var records []*gty.GuessGameRecord

	for i := 0; i < len(rows); i++ {
		userBet := rows[i].Data.(*gty.UserBet)
		var record gty.GuessGameRecord
		record.GameID = userBet.GameID
		record.StartIndex = userBet.StartIndex
		records = append(records, &record)
	}

	var primary string
	if len(rows) == int(DefaultCount) {
		primary = string(rows[len(rows)-1].Primary)
	}

	return &gty.GuessGameRecords{Records: records, PrimaryKey: primary}, nil
}

//queryGameTableData   game
func queryGameTableData(query *table.Query, indexName string, prefix, primaryKey []byte) (types.Message, error) {
	rows, err := query.ListIndex(indexName, prefix, primaryKey, DefaultCount, 0)
	if err != nil {
		return nil, err
	}

	var records []*gty.GuessGameRecord

	for i := 0; i < len(rows); i++ {
		game := rows[i].Data.(*gty.GuessGame)
		var record gty.GuessGameRecord
		record.GameID = game.GameID
		record.StartIndex = game.StartIndex
		records = append(records, &record)
	}

	var primary string
	if len(rows) == int(DefaultCount) {
		primary = string(rows[len(rows)-1].Primary)
	}

	return &gty.GuessGameRecords{Records: records, PrimaryKey: primary}, nil
}

//queryJoinTableData   join
func queryJoinTableData(talbeJoin *table.JoinTable, indexName string, prefix, primaryKey []byte) (types.Message, error) {
	rows, err := talbeJoin.ListIndex(indexName, prefix, primaryKey, DefaultCount, 0)
	if err != nil {
		return nil, err
	}

	var records []*gty.GuessGameRecord

	for i := 0; i < len(rows); i++ {
		game := rows[i].Data.(*table.JoinData).Right.(*gty.GuessGame)
		var record gty.GuessGameRecord
		record.GameID = game.GameID
		record.StartIndex = game.StartIndex
		records = append(records, &record)
	}

	var primary string
	if len(rows) == int(DefaultCount) {
		primary = fmt.Sprintf("%018d", rows[len(rows)-1].Data.(*table.JoinData).Left.(*gty.UserBet).Index)
	}

	return &gty.GuessGameRecords{Records: records, PrimaryKey: primary}, nil
}

func (action *Action) saveGame(game *gty.GuessGame) (kvset []*types.KeyValue) {
	value := types.Encode(game)
	err := action.db.Set(Key(game.GetGameID()), value)
	if err != nil {
		logger.Error("saveGame have err:", err.Error())
	}
	kvset = append(kvset, &types.KeyValue{Key: Key(game.GameID), Value: value})
	return kvset
}

func (action *Action) getIndex() int64 {
	return action.height*types.MaxTxsPerBlock + int64(action.index)
}

//getReceiptLog
func (action *Action) getReceiptLog(game *gty.GuessGame, statusChange bool, bet *gty.GuessGameBet) *types.ReceiptLog {
	log := &types.ReceiptLog{}
	r := &gty.ReceiptGuessGame{}
	r.Addr = action.fromaddr
	if game.Status == gty.GuessGameStatusStart {
		log.Ty = gty.TyLogGuessGameStart
	} else if game.Status == gty.GuessGameStatusBet {
		log.Ty = gty.TyLogGuessGameBet
	} else if game.Status == gty.GuessGameStatusStopBet {
		log.Ty = gty.TyLogGuessGameStopBet
	} else if game.Status == gty.GuessGameStatusAbort {
		log.Ty = gty.TyLogGuessGameAbort
	} else if game.Status == gty.GuessGameStatusPublish {
		log.Ty = gty.TyLogGuessGamePublish
	} else if game.Status == gty.GuessGameStatusTimeOut {
		log.Ty = gty.TyLogGuessGameTimeout
	}

	r.StartIndex = game.StartIndex
	r.Index = action.getIndex()
	r.GameID = game.GameID
	r.Status = game.Status
	r.AdminAddr = game.AdminAddr
	r.PreStatus = game.PreStatus
	r.StatusChange = statusChange
	r.PreIndex = game.PreIndex
	r.Category = game.Category
	if nil != bet {
		r.Bet = true
		r.Option = bet.Option
		r.BetsNumber = bet.BetsNum
	}
	r.Game = game
	log.Log = types.Encode(r)
	return log
}

func (action *Action) readGame(id string) (*gty.GuessGame, error) {
	data, err := action.db.Get(Key(id))
	if err != nil {
		logger.Error("readGame have err", "err", err.Error())
		return nil, err
	}
	var game gty.GuessGame
	//decode
	err = types.Decode(data, &game)
	if err != nil {
		logger.Error("decode game have err:", err.Error())
		return nil, err
	}
	return &game, nil
}

//
func (action *Action) newGame(gameID string, start *gty.GuessGameStart) *gty.GuessGame {
	game := &gty.GuessGame{
		GameID: gameID,
		Status: gty.GuessGameStatusStart,
		//StartTime:   action.blocktime,
		StartTxHash:    gameID,
		Topic:          start.Topic,
		Category:       start.Category,
		Options:        start.Options,
		MaxBetHeight:   start.MaxBetHeight,
		MaxBetsOneTime: start.MaxBetsOneTime,
		MaxBetsNumber:  start.MaxBetsNumber,
		DevFeeFactor:   start.DevFeeFactor,
		DevFeeAddr:     start.DevFeeAddr,
		PlatFeeFactor:  start.PlatFeeFactor,
		PlatFeeAddr:    start.PlatFeeAddr,
		ExpireHeight:   start.ExpireHeight,
		//AdminAddr: action.fromaddr,
		BetsNumber: 0,
		//Index:       action.getIndex(game),
		DrivenByAdmin: start.DrivenByAdmin,
	}

	return game
}

//GameStart
func (action *Action) GameStart(start *gty.GuessGameStart) (*types.Receipt, error) {
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue

	if start.MaxBetHeight >= MaxBetHeight {
		logger.Error("GameStart", "addr", action.fromaddr, "execaddr", action.execaddr,
			"err", fmt.Sprintf("The maximum height diff number is %d which is less than start.MaxBetHeight %d", MaxBetHeight, start.MaxBetHeight))
		return nil, types.ErrInvalidParam
	}

	if start.ExpireHeight >= MaxExpireHeight {
		logger.Error("GameStart", "addr", action.fromaddr, "execaddr", action.execaddr,
			"err", fmt.Sprintf("The maximum height diff number is %d which is less than start.MaxBetHeight %d", MaxBetHeight, start.MaxBetHeight))
		return nil, types.ErrInvalidParam
	}

	if start.MaxBetsNumber >= MaxBetsNumber {
		logger.Error("GameStart", "addr", action.fromaddr, "execaddr", action.execaddr,
			"err", fmt.Sprintf("The maximum bets number is %d which is less than start.MaxBetsNumber %d", int64(MaxBetsNumber), start.MaxBetsNumber))
		return nil, gty.ErrOverBetsLimit
	}

	if len(start.Topic) == 0 || len(start.Options) == 0 {
		logger.Error("GameStart", "addr", action.fromaddr, "execaddr", action.execaddr,
			"err", fmt.Sprintf("Illegal parameters,Topic:%s | options: %s | category: %s", start.Topic, start.Options, start.Category))
		return nil, types.ErrInvalidParam
	}

	options, ok := getOptions(start.Options)
	if !ok {
		logger.Error("GameStart", "addr", action.fromaddr, "execaddr", action.execaddr,
			"err", fmt.Sprintf("The options is illegal:%s", start.Options))
		return nil, types.ErrInvalidParam
	}

	if !action.checkTime(start) {
		logger.Error("GameStart", "addr", action.fromaddr, "execaddr", action.execaddr,
			"err", fmt.Sprintf("The height and time parameters are illegal:MaxHeight %d ,ExpireHeight %d", start.MaxBetHeight, start.ExpireHeight))
		return nil, types.ErrInvalidParam
	}

	if len(start.Category) == 0 {
		start.Category = DefaultCategory
	}

	if start.MaxBetsOneTime >= MaxBetsOneTime {
		start.MaxBetsOneTime = MaxBetsOneTime
	}

	gameID := common.ToHex(action.txhash)
	game := action.newGame(gameID, start)
	game.StartTime = action.blocktime
	game.StartHeight = action.mainHeight
	game.AdminAddr = action.fromaddr
	game.PreIndex = 0
	game.Index = action.getIndex()
	game.StartIndex = game.Index
	game.Status = gty.GuessGameStatusStart
	game.BetStat = &gty.GuessBetStat{TotalBetTimes: 0, TotalBetsNumber: 0}
	for i := 0; i < len(options); i++ {
		item := &gty.GuessBetStatItem{Option: options[i], BetsNumber: 0, BetsTimes: 0}
		game.BetStat.Items = append(game.BetStat.Items, item)
	}

	receiptLog := action.getReceiptLog(game, false, nil)
	logs = append(logs, receiptLog)
	kv = append(kv, action.saveGame(game)...)

	return &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}, nil
}

//GameBet
func (action *Action) GameBet(pbBet *gty.GuessGameBet) (*types.Receipt, error) {
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue

	game, err := action.readGame(pbBet.GetGameID())
	if err != nil || game == nil {
		logger.Error("GameBet", "addr", action.fromaddr, "execaddr", action.execaddr, "get game failed",
			pbBet.GetGameID(), "err", err)
		return nil, err
	}

	prevStatus := game.Status
	if game.Status != gty.GuessGameStatusStart && game.Status != gty.GuessGameStatusBet {
		logger.Error("GameBet", "addr", action.fromaddr, "execaddr", action.execaddr, "Status error",
			game.GetStatus())
		return nil, gty.ErrGuessStatus
	}

	canBet := action.refreshStatusByTime(game)

	if !canBet {
		var receiptLog *types.ReceiptLog
		if prevStatus != game.Status {
			//       ，            ，         addr  ， addr:status
			action.changeAllAddrIndex(game)
			receiptLog = action.getReceiptLog(game, true, nil)
		} else {
			receiptLog = action.getReceiptLog(game, false, nil)
		}

		logs = append(logs, receiptLog)
		kv = append(kv, action.saveGame(game)...)

		return &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}, nil
	}

	//
	options, legal := getOptions(game.GetOptions())
	if !legal || len(options) == 0 {
		logger.Error("GameBet", "addr", action.fromaddr, "execaddr", action.execaddr, "Game Options illegal",
			game.GetOptions())
		return nil, types.ErrInvalidParam
	}

	if !isLegalOption(options, pbBet.GetOption()) {
		logger.Error("GameBet", "addr", action.fromaddr, "execaddr", action.execaddr, "Option illegal",
			pbBet.GetOption())
		return nil, types.ErrInvalidParam
	}

	//          ，    ，
	if pbBet.GetBetsNum() > game.GetMaxBetsOneTime() {
		pbBet.BetsNum = game.GetMaxBetsOneTime()
	}

	if game.BetsNumber+pbBet.GetBetsNum() > game.MaxBetsNumber {
		logger.Error("GameBet", "addr", action.fromaddr, "execaddr", action.execaddr, "MaxBetsNumber over limit",
			game.MaxBetsNumber, "current Bets Number", game.BetsNumber)
		return nil, types.ErrInvalidParam
	}

	//
	checkValue := pbBet.BetsNum
	if !action.CheckExecAccountBalance(action.fromaddr, checkValue, 0) {
		logger.Error("GameBet", "addr", action.fromaddr, "execaddr", action.execaddr, "id",
			pbBet.GetGameID(), "err", types.ErrNoBalance)
		return nil, types.ErrNoBalance
	}

	receipt, err := action.coinsAccount.ExecFrozen(action.fromaddr, action.execaddr, checkValue)
	if err != nil {
		logger.Error("GameCreate.ExecFrozen", "addr", action.fromaddr, "execaddr", action.execaddr, "amount", checkValue, "err", err.Error())
		return nil, err
	}
	logs = append(logs, receipt.Logs...)
	kv = append(kv, receipt.KV...)

	var receiptLog *types.ReceiptLog
	if prevStatus != gty.GuessGameStatusBet {
		action.changeStatus(game, gty.GuessGameStatusBet)
		action.addGuessBet(game, pbBet)
		receiptLog = action.getReceiptLog(game, true, pbBet)
	} else {
		action.addGuessBet(game, pbBet)
		receiptLog = action.getReceiptLog(game, false, pbBet)
	}

	logs = append(logs, receiptLog)
	kv = append(kv, action.saveGame(game)...)

	return &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}, nil
}

//GameStopBet
func (action *Action) GameStopBet(pbBet *gty.GuessGameStopBet) (*types.Receipt, error) {
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue

	game, err := action.readGame(pbBet.GetGameID())
	if err != nil || game == nil {
		logger.Error("GameStopBet", "addr", action.fromaddr, "execaddr", action.execaddr, "get game failed",
			pbBet.GetGameID(), "err", err)
		return nil, err
	}

	if game.Status != gty.GuessGameStatusStart && game.Status != gty.GuessGameStatusBet {
		logger.Error("GameBet", "addr", action.fromaddr, "execaddr", action.execaddr, "Status error",
			game.GetStatus())
		return nil, gty.ErrGuessStatus
	}

	//  adminAddr    stopBet
	if game.AdminAddr != action.fromaddr {
		logger.Error("GameStopBet", "addr", action.fromaddr, "execaddr", action.execaddr, "fromAddr is not adminAddr",
			action.fromaddr, "adminAddr", game.AdminAddr)
		return nil, gty.ErrNoPrivilege
	}

	action.changeStatus(game, gty.GuessGameStatusStopBet)

	var receiptLog *types.ReceiptLog
	//      ，    addr     index
	action.changeAllAddrIndex(game)
	receiptLog = action.getReceiptLog(game, true, nil)

	logs = append(logs, receiptLog)
	kv = append(kv, action.saveGame(game)...)

	return &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}, nil
}

//addGuessBet
func (action *Action) addGuessBet(game *gty.GuessGame, pbBet *gty.GuessGameBet) {
	bet := &gty.GuessBet{Option: pbBet.GetOption(), BetsNumber: pbBet.BetsNum, Index: action.getIndex()}
	player := &gty.GuessPlayer{Addr: action.fromaddr, Bet: bet}
	game.Plays = append(game.Plays, player)

	for i := 0; i < len(game.BetStat.Items); i++ {
		if game.BetStat.Items[i].Option == trimStr(pbBet.GetOption()) {
			//
			game.BetStat.Items[i].BetsNumber += pbBet.GetBetsNum()
			game.BetStat.Items[i].BetsTimes++

			//
			game.BetStat.TotalBetsNumber += pbBet.GetBetsNum()
			game.BetStat.TotalBetTimes++
			break
		}
	}

	game.BetsNumber += pbBet.GetBetsNum()
}

//GamePublish
func (action *Action) GamePublish(publish *gty.GuessGamePublish) (*types.Receipt, error) {
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue

	game, err := action.readGame(publish.GetGameID())
	if err != nil || game == nil {
		logger.Error("GamePublish", "addr", action.fromaddr, "execaddr", action.execaddr, "get game failed",
			publish.GetGameID(), "err", err)
		return nil, err
	}

	//  adminAddr    publish
	if game.AdminAddr != action.fromaddr {
		logger.Error("GamePublish", "addr", action.fromaddr, "execaddr", action.execaddr, "fromAddr is not adminAddr",
			action.fromaddr, "adminAddr", game.AdminAddr)
		return nil, gty.ErrNoPrivilege
	}

	if game.Status != gty.GuessGameStatusStart && game.Status != gty.GuessGameStatusBet && game.Status != gty.GuessGameStatusStopBet {
		logger.Error("GamePublish", "addr", action.fromaddr, "execaddr", action.execaddr, "Status error",
			game.GetStatus())
		return nil, gty.ErrGuessStatus
	}

	//
	options, legal := getOptions(game.GetOptions())
	if !legal || len(options) == 0 {
		logger.Error("GamePublish", "addr", action.fromaddr, "execaddr", action.execaddr, "Game Options illegal",
			game.GetOptions())
		return nil, types.ErrInvalidParam
	}

	if !isLegalOption(options, publish.GetResult()) {
		logger.Error("GamePublish", "addr", action.fromaddr, "execaddr", action.execaddr, "Option illegal",
			publish.GetResult())
		return nil, types.ErrInvalidParam
	}

	game.Result = trimStr(publish.Result)

	//         ，     Admin      ；
	for i := 0; i < len(game.Plays); i++ {
		player := game.Plays[i]
		value := player.Bet.BetsNumber
		receipt, err := action.coinsAccount.ExecActive(player.Addr, action.execaddr, value)
		if err != nil {
			logger.Error("GamePublish.ExecActive", "addr", player.Addr, "execaddr", action.execaddr, "amount", value,
				"err", err)
			return nil, err
		}
		logs = append(logs, receipt.Logs...)
		kv = append(kv, receipt.KV...)

		receipt, err = action.coinsAccount.ExecTransfer(player.Addr, game.AdminAddr, action.execaddr, value)
		if err != nil {
			//action.coinsAccount.ExecFrozen(game.AdminAddr, action.execaddr, value) // rollback
			logger.Error("GamePublish", "addr", player.Addr, "execaddr", action.execaddr,
				"amount", value, "err", err)
			return nil, err
		}
		logs = append(logs, receipt.Logs...)
		kv = append(kv, receipt.KV...)
	}

	action.changeStatus(game, gty.GuessGameStatusPublish)
	//
	totalBetsNumber := game.BetStat.TotalBetsNumber
	winBetsNumber := int64(0)
	for j := 0; j < len(game.BetStat.Items); j++ {
		if game.BetStat.Items[j].Option == game.Result {
			winBetsNumber = game.BetStat.Items[j].BetsNumber
		}
	}

	//           ，
	devAddr := gty.DevShareAddr
	platAddr := gty.PlatformShareAddr
	devFee := int64(0)
	platFee := int64(0)
	if len(game.DevFeeAddr) > 0 {
		devAddr = game.DevFeeAddr
	}

	if len(game.PlatFeeAddr) > 0 {
		platAddr = game.PlatFeeAddr
	}

	if game.DevFeeFactor > 0 {
		fee := big.NewInt(totalBetsNumber)
		factor := big.NewInt(game.DevFeeFactor)
		thousand := big.NewInt(1000)
		devFee = fee.Mul(fee, factor).Div(fee, thousand).Int64()
		receipt, err := action.coinsAccount.ExecTransfer(game.AdminAddr, devAddr, action.execaddr, devFee)
		if err != nil {
			//action.coinsAccount.ExecFrozen(game.AdminAddr, action.execaddr, devFee) // rollback
			logger.Error("GamePublish", "adminAddr", game.AdminAddr, "execaddr", action.execaddr,
				"amount", devFee, "err", err)
			return nil, err
		}
		logs = append(logs, receipt.Logs...)
		kv = append(kv, receipt.KV...)
	}

	if game.PlatFeeFactor > 0 {
		fee := big.NewInt(totalBetsNumber)
		factor := big.NewInt(game.PlatFeeFactor)
		thousand := big.NewInt(1000)
		platFee = fee.Mul(fee, factor).Div(fee, thousand).Int64()
		receipt, err := action.coinsAccount.ExecTransfer(game.AdminAddr, platAddr, action.execaddr, platFee)
		if err != nil {
			//action.coinsAccount.ExecFrozen(game.AdminAddr, action.execaddr, platFee) // rollback
			logger.Error("GamePublish", "adminAddr", game.AdminAddr, "execaddr", action.execaddr,
				"amount", platFee, "err", err)
			return nil, err
		}
		logs = append(logs, receipt.Logs...)
		kv = append(kv, receipt.KV...)
	}

	//     ，
	winValue := totalBetsNumber - devFee - platFee
	for j := 0; j < len(game.Plays); j++ {
		player := game.Plays[j]
		if trimStr(player.Bet.Option) == trimStr(game.Result) {
			betsNumber := big.NewInt(player.Bet.BetsNumber)
			totalWinBetsNumber := big.NewInt(winBetsNumber)
			leftWinBetsNumber := big.NewInt(winValue)

			value := betsNumber.Mul(betsNumber, leftWinBetsNumber).Div(betsNumber, totalWinBetsNumber).Int64()
			receipt, err := action.coinsAccount.ExecTransfer(game.AdminAddr, player.Addr, action.execaddr, value)
			if err != nil {
				//action.coinsAccount.ExecFrozen(player.Addr, action.execaddr, value) // rollback
				logger.Error("GamePublish", "addr", player.Addr, "execaddr", action.execaddr,
					"amount", value, "err", err)
				return nil, err
			}
			logs = append(logs, receipt.Logs...)
			kv = append(kv, receipt.KV...)
			player.Bet.IsWinner = true
			player.Bet.Profit = value
		}
	}

	var receiptLog *types.ReceiptLog
	action.changeAllAddrIndex(game)
	receiptLog = action.getReceiptLog(game, true, nil)

	logs = append(logs, receiptLog)
	kv = append(kv, action.saveGame(game)...)

	return &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}, nil
}

//GameAbort
func (action *Action) GameAbort(pbend *gty.GuessGameAbort) (*types.Receipt, error) {
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue

	game, err := action.readGame(pbend.GetGameID())
	if err != nil || game == nil {
		logger.Error("GameAbort", "addr", action.fromaddr, "execaddr", action.execaddr, "get game failed",
			pbend.GetGameID(), "err", err)
		return nil, err
	}

	if game.Status == gty.GuessGameStatusPublish || game.Status == gty.GuessGameStatusAbort {

		logger.Error("GameAbort", "addr", action.fromaddr, "execaddr", action.execaddr, "game status not allow abort",
			game.Status)
		return nil, gty.ErrGuessStatus
	}

	preStatus := game.Status
	//                。
	action.refreshStatusByTime(game)

	//      ，        Abort，             Abort
	if game.Status != gty.GuessGameStatusTimeOut {
		if game.AdminAddr != action.fromaddr {
			logger.Error("GameAbort", "addr", action.fromaddr, "execaddr", action.execaddr, "Only admin can abort",
				action.fromaddr, "status", game.Status)
			return nil, err
		}
	}

	//
	for i := 0; i < len(game.Plays); i++ {
		player := game.Plays[i]
		value := player.Bet.BetsNumber
		receipt, err := action.coinsAccount.ExecActive(player.Addr, action.execaddr, value)
		if err != nil {
			logger.Error("GameAbort", "addr", player.Addr, "execaddr", action.execaddr, "amount", value, "err", err)
			continue
		}

		player.Bet.IsWinner = true
		player.Bet.Profit = value

		logs = append(logs, receipt.Logs...)
		kv = append(kv, receipt.KV...)
	}

	if game.Status != preStatus {
		//  action.RefreshStatusByTime(game)           index ，           。
		game.Status = gty.GuessGameStatusAbort
	} else {
		action.changeStatus(game, gty.GuessGameStatusAbort)
	}

	//      ，      addr   index
	action.changeAllAddrIndex(game)

	receiptLog := action.getReceiptLog(game, true, nil)
	logs = append(logs, receiptLog)
	kv = append(kv, action.saveGame(game)...)
	return &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}, nil
}

//getOptions       ，           ，  "A:xxxx;B:xxxx;C:xxx"，“：”      ，    ，":"      。
func getOptions(strOptions string) (options []string, legal bool) {
	if len(strOptions) == 0 {
		return nil, false
	}

	legal = true
	items := strings.Split(strOptions, ";")
	for i := 0; i < len(items); i++ {
		item := strings.Split(items[i], ":")
		for j := 0; j < len(options); j++ {
			if item[0] == options[j] {
				legal = false
				return
			}
		}

		options = append(options, trimStr(item[0]))
	}

	return options, legal
}

//trimStr          、   、
func trimStr(str string) string {
	str = strings.Replace(str, " ", "", -1)
	str = strings.Replace(str, "\t", "", -1)
	str = strings.Replace(str, "\n", "", -1)

	return str
}

//isLegalOption
func isLegalOption(options []string, option string) bool {
	option = trimStr(option)
	for i := 0; i < len(options); i++ {
		if options[i] == option {
			return true
		}
	}

	return false
}

//changeStatus       ，
func (action *Action) changeStatus(game *gty.GuessGame, destStatus int32) {
	if game.Status != destStatus {
		game.PreStatus = game.Status
		game.PreIndex = game.Index
		game.Status = destStatus
		game.Index = action.getIndex()
	}
}

//changeAllAddrIndex      ，
func (action *Action) changeAllAddrIndex(game *gty.GuessGame) {
	for i := 0; i < len(game.Plays); i++ {
		player := game.Plays[i]
		player.Bet.PreIndex = player.Bet.Index
		player.Bet.Index = action.getIndex()
	}
}

//refreshStatusByTime         ，
func (action *Action) refreshStatusByTime(game *gty.GuessGame) (canBet bool) {
	mainHeight := action.mainHeight
	//              ，           ，        。
	if game.DrivenByAdmin {

		if (mainHeight - game.StartHeight) >= game.ExpireHeight {
			action.changeStatus(game, gty.GuessGameStatusTimeOut)
			canBet = false
			return canBet
		}

		return true
	}

	//                    ，
	heightDiff := mainHeight - game.StartHeight
	if heightDiff >= game.MaxBetHeight {
		logger.Error("GameBet", "addr", action.fromaddr, "execaddr", action.execaddr, "Height over limit",
			mainHeight, "startHeight", game.StartHeight, "MaxHeightDiff", game.GetMaxBetHeight())
		if game.ExpireHeight > heightDiff {
			action.changeStatus(game, gty.GuessGameStatusStopBet)
		} else {
			action.changeStatus(game, gty.GuessGameStatusTimeOut)
		}

		canBet = false
		return canBet
	}

	canBet = true
	return canBet
}

//checkTime          。
func (action *Action) checkTime(start *gty.GuessGameStart) bool {
	if start.MaxBetHeight == 0 && start.ExpireHeight == 0 {
		//          ，      admin     。
		start.DrivenByAdmin = true

		//           ，
		start.ExpireHeight = MaxExpireHeight
		return true
	}

	if start.MaxBetHeight == 0 {
		start.MaxBetHeight = MaxBetHeight
	}

	if start.ExpireHeight == 0 {
		start.ExpireHeight = MaxExpireHeight
	}

	if start.MaxBetHeight <= start.ExpireHeight {
		return true
	}

	return false
}

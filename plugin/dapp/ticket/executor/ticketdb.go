// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

//database opeartion for execs ticket
import (

	//"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/D-PlatformOperatingSystem/dpos/account"
	"github.com/D-PlatformOperatingSystem/dpos/client"
	"github.com/D-PlatformOperatingSystem/dpos/common"
	"github.com/D-PlatformOperatingSystem/dpos/common/address"
	dbm "github.com/D-PlatformOperatingSystem/dpos/common/db"
	log "github.com/D-PlatformOperatingSystem/dpos/common/log/log15"
	"github.com/D-PlatformOperatingSystem/dpos/system/dapp"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	ty "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/ticket/types"
)

var tlog = log.New("module", "ticket.db")

//var genesisKey = []byte("mavl-acc-genesis")
//var addrSeed = []byte("address seed bytes for public key")

// DB db
type DB struct {
	ty.Ticket
	prevstatus int32
}

//GetRealPrice
func (t *DB) GetRealPrice(cfg *types.DplatformOSConfig) int64 {
	if t.GetPrice() == 0 {
		cfg := ty.GetTicketMinerParam(cfg, cfg.GetFork("ForkChainParamV1"))
		return cfg.TicketPrice
	}
	return t.GetPrice()
}

// NewDB new instance
func NewDB(cfg *types.DplatformOSConfig, id, minerAddress, returnWallet string, blocktime, height, price int64, isGenesis bool) *DB {
	t := &DB{}
	t.TicketId = id
	t.MinerAddress = minerAddress
	t.ReturnAddress = returnWallet
	t.CreateTime = blocktime
	t.Status = ty.TicketOpened
	t.IsGenesis = isGenesis
	t.prevstatus = 0
	//height == 0     ，     genesis block
	if cfg.IsFork(height, "ForkChainParamV2") && height > 0 {
		t.Price = price
	}
	return t
}

//ticket      ：
//1. status == 1 (NewTicket   )
//2. status == 2 (       )
//3. status == 3 (Close   )

//add prevStatus:        ，
//list      :
//minerAddress:status:ticketId=ticketId

// GetReceiptLog get receipt
func (t *DB) GetReceiptLog() *types.ReceiptLog {
	log := &types.ReceiptLog{}
	if t.Status == ty.TicketOpened {
		log.Ty = ty.TyLogNewTicket
	} else if t.Status == ty.TicketMined {
		log.Ty = ty.TyLogMinerTicket
	} else if t.Status == ty.TicketClosed {
		log.Ty = ty.TyLogCloseTicket
	}
	r := &ty.ReceiptTicket{}
	r.TicketId = t.TicketId
	r.Status = t.Status
	r.PrevStatus = t.prevstatus
	r.Addr = t.MinerAddress
	log.Log = types.Encode(r)
	return log
}

// GetKVSet get kv set
func (t *DB) GetKVSet() (kvset []*types.KeyValue) {
	value := types.Encode(&t.Ticket)
	kvset = append(kvset, &types.KeyValue{Key: Key(t.TicketId), Value: value})
	return kvset
}

// Save save
func (t *DB) Save(db dbm.KV) {
	set := t.GetKVSet()
	for i := 0; i < len(set); i++ {
		db.Set(set[i].GetKey(), set[i].Value)
	}
}

//Key address to save key
func Key(id string) (key []byte) {
	key = append(key, []byte("mavl-ticket-")...)
	key = append(key, []byte(id)...)
	return key
}

// BindKey bind key
func BindKey(id string) (key []byte) {
	key = append(key, []byte("mavl-ticket-tbind-")...)
	key = append(key, []byte(id)...)
	return key
}

// Action action type
type Action struct {
	coinsAccount *account.DB
	db           dbm.KV
	txhash       []byte
	fromaddr     string
	blocktime    int64
	height       int64
	execaddr     string
	api          client.QueueProtocolAPI
}

// NewAction new action type
func NewAction(t *Ticket, tx *types.Transaction) *Action {
	hash := tx.Hash()
	fromaddr := tx.From()
	return &Action{t.GetCoinsAccount(), t.GetStateDB(), hash, fromaddr,
		t.GetBlockTime(), t.GetHeight(), dapp.ExecAddress(string(tx.Execer)), t.GetAPI()}
}

// GenesisInit init genesis
func (action *Action) GenesisInit(genesis *ty.TicketGenesis) (*types.Receipt, error) {
	dplatformosCfg := action.api.GetConfig()
	prefix := common.ToHex(action.txhash)
	prefix = genesis.MinerAddress + ":" + prefix + ":"
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	cfg := ty.GetTicketMinerParam(dplatformosCfg, action.height)
	for i := 0; i < int(genesis.Count); i++ {
		id := prefix + fmt.Sprintf("%010d", i)
		t := NewDB(dplatformosCfg, id, genesis.MinerAddress, genesis.ReturnAddress, action.blocktime, action.height, cfg.TicketPrice, true)
		//
		receipt, err := action.coinsAccount.ExecFrozen(genesis.ReturnAddress, action.execaddr, cfg.TicketPrice)
		if err != nil {
			tlog.Error("GenesisInit.Frozen", "addr", genesis.ReturnAddress, "execaddr", action.execaddr)
			panic(err)
		}
		t.Save(action.db)
		logs = append(logs, t.GetReceiptLog())
		kv = append(kv, t.GetKVSet()...)
		logs = append(logs, receipt.Logs...)
		kv = append(kv, receipt.KV...)
	}
	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil
}

func saveBind(db dbm.KV, tbind *ty.TicketBind) {
	set := getBindKV(tbind)
	for i := 0; i < len(set); i++ {
		db.Set(set[i].GetKey(), set[i].Value)
	}
}

func getBindKV(tbind *ty.TicketBind) (kvset []*types.KeyValue) {
	value := types.Encode(tbind)
	kvset = append(kvset, &types.KeyValue{Key: BindKey(tbind.ReturnAddress), Value: value})
	return kvset
}

func getBindLog(tbind *ty.TicketBind, old string) *types.ReceiptLog {
	log := &types.ReceiptLog{}
	log.Ty = ty.TyLogTicketBind
	r := &ty.ReceiptTicketBind{}
	r.ReturnAddress = tbind.ReturnAddress
	r.OldMinerAddress = old
	r.NewMinerAddress = tbind.MinerAddress
	log.Log = types.Encode(r)
	return log
}

func (action *Action) getBind(addr string) string {
	value, err := action.db.Get(BindKey(addr))
	if err != nil || value == nil {
		return ""
	}
	var bind ty.TicketBind
	err = types.Decode(value, &bind)
	if err != nil {
		panic(err)
	}
	return bind.MinerAddress
}

//TicketBind
func (action *Action) TicketBind(tbind *ty.TicketBind) (*types.Receipt, error) {
	//todo: query address is a minered address
	if action.fromaddr != tbind.ReturnAddress {
		return nil, types.ErrFromAddr
	}
	//""
	if len(tbind.MinerAddress) > 0 {
		if err := address.CheckAddress(tbind.MinerAddress); err != nil {
			return nil, err
		}
	}
	var logs []*types.ReceiptLog
	var kvs []*types.KeyValue
	oldbind := action.getBind(tbind.ReturnAddress)
	log := getBindLog(tbind, oldbind)
	logs = append(logs, log)
	saveBind(action.db, tbind)
	kv := getBindKV(tbind)
	kvs = append(kvs, kv...)
	receipt := &types.Receipt{Ty: types.ExecOk, KV: kvs, Logs: logs}
	return receipt, nil
}

// TicketOpen ticket open
func (action *Action) TicketOpen(topen *ty.TicketOpen) (*types.Receipt, error) {
	dplatformosCfg := action.api.GetConfig()
	prefix := common.ToHex(action.txhash)
	prefix = topen.MinerAddress + ":" + prefix + ":"
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	//addr from
	if action.fromaddr != topen.ReturnAddress {
		mineraddr := action.getBind(topen.ReturnAddress)
		if mineraddr != action.fromaddr {
			return nil, ty.ErrMinerNotPermit
		}
		if topen.MinerAddress != mineraddr {
			return nil, ty.ErrMinerAddr
		}
	}
	//action.fromaddr == topen.ReturnAddress or mineraddr == action.fromaddr
	cfg := ty.GetTicketMinerParam(dplatformosCfg, action.height)
	for i := 0; i < int(topen.Count); i++ {
		id := prefix + fmt.Sprintf("%010d", i)
		//add pubHash
		if dplatformosCfg.IsDappFork(action.height, ty.TicketX, "ForkTicketId") {
			if len(topen.PubHashes) == 0 {
				return nil, ty.ErrOpenTicketPubHash
			}
			id = id + ":" + fmt.Sprintf("%x:%d", topen.PubHashes[i], topen.RandSeed)
		}
		t := NewDB(dplatformosCfg, id, topen.MinerAddress, topen.ReturnAddress, action.blocktime, action.height, cfg.TicketPrice, false)

		//
		receipt, err := action.coinsAccount.ExecFrozen(topen.ReturnAddress, action.execaddr, cfg.TicketPrice)
		if err != nil {
			tlog.Error("TicketOpen.Frozen", "addr", topen.ReturnAddress, "execaddr", action.execaddr, "n", topen.Count)
			return nil, err
		}
		t.Save(action.db)
		logs = append(logs, t.GetReceiptLog())
		kv = append(kv, t.GetKVSet()...)
		logs = append(logs, receipt.Logs...)
		kv = append(kv, receipt.KV...)
	}
	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil
}

func readTicket(db dbm.KV, id string) (*ty.Ticket, error) {
	data, err := db.Get(Key(id))
	if err != nil {
		return nil, err
	}
	var ticket ty.Ticket
	//decode
	err = types.Decode(data, &ticket)
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

func genPubHash(tid string) string {
	var pubHash string
	parts := strings.Split(tid, ":")
	if len(parts) > ty.TicketOldParts {
		pubHash = parts[ty.TicketOldParts]
	}
	return pubHash
}

// TicketMiner ticket miner
func (action *Action) TicketMiner(miner *ty.TicketMiner, index int) (*types.Receipt, error) {
	if index != 0 {
		return nil, types.ErrCoinBaseIndex
	}
	dplatformosCfg := action.api.GetConfig()
	ticket, err := readTicket(action.db, miner.TicketId)
	if err != nil {
		return nil, err
	}
	if ticket.Status != ty.TicketOpened {
		return nil, types.ErrCoinBaseTicketStatus
	}
	cfg := ty.GetTicketMinerParam(dplatformosCfg, action.height)
	if !ticket.IsGenesis {
		if action.blocktime-ticket.GetCreateTime() < cfg.TicketFrozenTime {
			return nil, ty.ErrTime
		}
	}
	//check from address
	if action.fromaddr != ticket.MinerAddress && action.fromaddr != ticket.ReturnAddress {
		return nil, types.ErrFromAddr
	}
	//check pubHash and privHash
	if !dplatformosCfg.IsDappFork(action.height, ty.TicketX, "ForkTicketId") {
		miner.PrivHash = nil
	}
	if len(miner.PrivHash) != 0 {
		pubHash := genPubHash(ticket.TicketId)
		if len(pubHash) == 0 || hex.EncodeToString(common.Sha256(miner.PrivHash)) != pubHash {
			tlog.Error("TicketMiner", "pubHash", pubHash, "privHashHash", common.Sha256(miner.PrivHash), "ticketId", ticket.TicketId)
			return nil, errors.New("ErrCheckPubHash")
		}
	}
	prevstatus := ticket.Status
	ticket.Status = ty.TicketMined
	ticket.MinerValue = miner.Reward
	if dplatformosCfg.IsFork(action.height, "ForkMinerTime") {
		ticket.MinerTime = action.blocktime
	}
	t := &DB{*ticket, prevstatus}
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	//user
	receipt1, err := action.coinsAccount.ExecDepositFrozen(t.ReturnAddress, action.execaddr, ticket.MinerValue)
	if err != nil {
		tlog.Error("TicketMiner.ExecDepositFrozen user", "addr", t.ReturnAddress, "execaddr", action.execaddr)
		return nil, err
	}
	//fund
	var receipt2 *types.Receipt
	receipt2, err = action.coinsAccount.ExecDepositFrozen(dplatformosCfg.GetFundAddr(), action.execaddr, cfg.CoinDevFund)
	if err != nil {
		tlog.Error("TicketMiner.ExecDepositFrozen fund", "addr", dplatformosCfg.GetFundAddr(), "execaddr", action.execaddr, "error", err)
		return nil, err
	}
	/*if dplatformosCfg.IsFork(action.height, "ForkTicketFundAddrV1") {
		// issue coins to exec addr
		addr := dplatformosCfg.MGStr("mver.consensus.fundKeyAddr", action.height)
		receipt2, err = action.coinsAccount.ExecIssueCoins(addr, cfg.CoinDevFund)
		if err != nil {
			tlog.Error("TicketMiner.ExecDepositFrozen fund to autonomy fund", "addr", addr, "error", err)
			return nil, err
		}
	} else {

	}*/

	t.Save(action.db)
	logs = append(logs, t.GetReceiptLog())
	kv = append(kv, t.GetKVSet()...)
	logs = append(logs, receipt1.Logs...)
	kv = append(kv, receipt1.KV...)
	logs = append(logs, receipt2.Logs...)
	kv = append(kv, receipt2.KV...)
	return &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}, nil
}

// TicketClose close tick
func (action *Action) TicketClose(tclose *ty.TicketClose) (*types.Receipt, error) {
	dplatformosCfg := action.api.GetConfig()
	tickets := make([]*DB, len(tclose.TicketId))
	cfg := ty.GetTicketMinerParam(dplatformosCfg, action.height)
	for i := 0; i < len(tclose.TicketId); i++ {
		ticket, err := readTicket(action.db, tclose.TicketId[i])
		if err != nil {
			return nil, err
		}
		//ticket         2 ,
		if ticket.Status != ty.TicketMined && ticket.Status != ty.TicketOpened {
			tlog.Error("ticket", "id", ticket.GetTicketId(), "status", ticket.GetStatus())
			return nil, ty.ErrTicketClosed
		}
		if !ticket.IsGenesis {
			//
			if ticket.Status == ty.TicketOpened && action.blocktime-ticket.GetCreateTime() < cfg.TicketWithdrawTime {
				return nil, ty.ErrTime
			}
			//
			if ticket.Status == ty.TicketMined && action.blocktime-ticket.GetCreateTime() < cfg.TicketWithdrawTime {
				return nil, ty.ErrTime
			}
			if ticket.Status == ty.TicketMined && action.blocktime-ticket.GetMinerTime() < cfg.TicketMinerWaitTime {
				return nil, ty.ErrTime
			}
		}
		//check from address
		if action.fromaddr != ticket.MinerAddress && action.fromaddr != ticket.ReturnAddress {
			return nil, types.ErrFromAddr
		}
		prevstatus := ticket.Status
		ticket.Status = ty.TicketClosed
		tickets[i] = &DB{*ticket, prevstatus}
	}
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	for i := 0; i < len(tickets); i++ {
		t := tickets[i]
		if t.prevstatus == 1 {
			t.MinerValue = 0
		}
		retValue := t.GetRealPrice(dplatformosCfg) + t.MinerValue
		receipt1, err := action.coinsAccount.ExecActive(t.ReturnAddress, action.execaddr, retValue)
		if err != nil {
			tlog.Error("TicketClose.ExecActive user", "addr", t.ReturnAddress, "execaddr", action.execaddr, "value", retValue)
			return nil, err
		}
		logs = append(logs, t.GetReceiptLog())
		kv = append(kv, t.GetKVSet()...)
		logs = append(logs, receipt1.Logs...)
		kv = append(kv, receipt1.KV...)
		//  ticket        ，
		if t.prevstatus == 2 {
			if !dplatformosCfg.IsFork(action.height, "ForkTicketFundAddrV1") {
				receipt2, err := action.coinsAccount.ExecActive(dplatformosCfg.GetFundAddr(), action.execaddr, cfg.CoinDevFund)
				if err != nil {
					tlog.Error("TicketClose.ExecActive fund", "addr", dplatformosCfg.GetFundAddr(), "execaddr", action.execaddr, "value", retValue)
					return nil, err
				}
				logs = append(logs, receipt2.Logs...)
				kv = append(kv, receipt2.KV...)
			}
		}
		t.Save(action.db)
	}
	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil
}

// List list db
func List(db dbm.Lister, db2 dbm.KV, tlist *ty.TicketList) (types.Message, error) {
	values, err := db.List(calcTicketPrefix(tlist.Addr, tlist.Status), nil, 0, 0)
	if err != nil {
		return nil, err
	}
	if len(values) == 0 {
		return &ty.ReplyTicketList{}, nil
	}
	var ids ty.TicketInfos
	for i := 0; i < len(values); i++ {
		ids.TicketIds = append(ids.TicketIds, string(values[i]))
	}
	return Infos(db2, &ids)
}

// Infos info
func Infos(db dbm.KV, tinfos *ty.TicketInfos) (types.Message, error) {
	var tickets []*ty.Ticket
	for i := 0; i < len(tinfos.TicketIds); i++ {
		id := tinfos.TicketIds[i]
		ticket, err := readTicket(db, id)
		//         ，
		if err != nil {
			continue
		}
		tickets = append(tickets, ticket)
	}
	return &ty.ReplyTicketList{Tickets: tickets}, nil
}

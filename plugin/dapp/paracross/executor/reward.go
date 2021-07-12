package executor

import (
	"bytes"

	"github.com/pkg/errors"

	dbm "github.com/D-PlatformOperatingSystem/dpos/common/db"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	pt "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/paracross/types"
	"github.com/golang/protobuf/proto"
)

const (
	opBind   = 1
	opUnBind = 2
)

//                       ，
func (a *action) getBindAddrs(nodes []string, statusHeight int64) (*pt.ParaNodeBindList, error) {
	nodesMap := make(map[string]bool)
	for _, n := range nodes {
		nodesMap[n] = true
	}

	var newLists pt.ParaNodeBindList
	list, err := getBindNodeInfo(a.db)
	if err != nil {
		clog.Error("paracross getBindAddrs err", "height", statusHeight)
		return nil, err
	}
	//       list     ，    nodes   (      )
	for _, m := range list.Miners {
		if nodesMap[m.SuperNode] {
			newLists.Miners = append(newLists.Miners, m)
		}
	}

	return &newLists, nil

}

func (a *action) rewardSuperNode(coinReward int64, miners []string, statusHeight int64) (*types.Receipt, int64, error) {
	//
	minerUnit := coinReward / int64(len(miners))
	var change int64
	receipt := &types.Receipt{Ty: types.ExecOk}
	if minerUnit > 0 {
		//
		change = coinReward % minerUnit
		for _, addr := range miners {
			rep, err := a.coinsAccount.ExecDeposit(addr, a.execaddr, minerUnit)

			if err != nil {
				clog.Error("paracross super node reward deposit err", "height", statusHeight,
					"execAddr", a.execaddr, "minerAddr", addr, "amount", minerUnit, "err", err)
				return nil, 0, err
			}
			receipt = mergeReceipt(receipt, rep)
		}
	}
	return receipt, change, nil
}

//
func (a *action) rewardBindAddr(coinReward int64, bindList *pt.ParaNodeBindList, statusHeight int64) (*types.Receipt, int64, error) {
	if coinReward <= 0 {
		return nil, 0, nil
	}

	//     bindAddr    node  ，
	var bindAddrList []*pt.ParaBindMinerInfo
	for _, node := range bindList.Miners {
		info, err := getBindAddrInfo(a.db, node.SuperNode, node.Miner)
		if err != nil {
			return nil, 0, err
		}
		bindAddrList = append(bindAddrList, info)
	}

	var totalCoins int64
	for _, addr := range bindAddrList {
		totalCoins += addr.BindCoins
	}

	//
	minerUnit := coinReward / totalCoins
	var change int64
	receipt := &types.Receipt{Ty: types.ExecOk}
	if minerUnit > 0 {
		//
		change = coinReward % minerUnit
		for _, miner := range bindAddrList {
			rep, err := a.coinsAccount.ExecDeposit(miner.Addr, a.execaddr, minerUnit*miner.BindCoins)
			if err != nil {
				clog.Error("paracross bind miner reward deposit err", "height", statusHeight,
					"execAddr", a.execaddr, "minerAddr", miner.Addr, "amount", minerUnit*miner.BindCoins, "err", err)
				return nil, 0, err
			}
			receipt = mergeReceipt(receipt, rep)
		}
	}
	return receipt, change, nil
}

// reward     ，           ，       ，
func (a *action) reward(nodeStatus *pt.ParacrossNodeStatus, stat *pt.ParacrossHeightStatus) (*types.Receipt, error) {
	//        ，           ，
	cfg := a.api.GetConfig()
	coinReward := cfg.MGInt("mver.consensus.paracross.coinReward", nodeStatus.Height)
	fundReward := cfg.MGInt("mver.consensus.paracross.coinDevFund", nodeStatus.Height)
	coinBaseReward := cfg.MGInt("mver.consensus.paracross.coinBaseReward", nodeStatus.Height)

	decimalMode := cfg.MIsEnable("mver.consensus.paracross.decimalMode", nodeStatus.Height)
	if !decimalMode {
		coinReward *= types.Coin
		fundReward *= types.Coin
		coinBaseReward *= types.Coin
	}

	fundAddr := cfg.MGStr("mver.consensus.fundKeyAddr", nodeStatus.Height)

	//  coinBaseReward       ， coinBaseReward     coinReward
	if coinBaseReward >= coinReward {
		coinBaseReward = coinReward / 10
	}

	//
	nodeAddrs := getSuperNodes(stat.Details, nodeStatus.BlockHash)
	//
	bindAddrs, err := a.getBindAddrs(nodeAddrs, nodeStatus.Height)
	if err != nil {
		return nil, err
	}

	//
	minderRewards := coinReward
	//         ，      baseReward  ，
	if len(bindAddrs.Miners) > 0 {
		minderRewards = coinBaseReward
	}
	receipt := &types.Receipt{Ty: types.ExecOk}
	r, change, err := a.rewardSuperNode(minderRewards, nodeAddrs, nodeStatus.Height)
	if err != nil {
		return nil, err
	}
	fundReward += change
	mergeReceipt(receipt, r)

	//
	r, change, err = a.rewardBindAddr(coinReward-minderRewards, bindAddrs, nodeStatus.Height)
	if err != nil {
		return nil, err
	}
	fundReward += change
	mergeReceipt(receipt, r)

	//
	if fundReward > 0 {
		rep, err := a.coinsAccount.ExecDeposit(fundAddr, a.execaddr, fundReward)
		if err != nil {
			clog.Error("paracross fund reward deposit err", "height", nodeStatus.Height,
				"execAddr", a.execaddr, "fundAddr", fundAddr, "amount", fundReward, "err", err)
			return nil, err
		}
		receipt = mergeReceipt(receipt, rep)
	}

	return receipt, nil
}

// getSuperNodes
func getSuperNodes(detail *pt.ParacrossStatusDetails, blockHash []byte) []string {
	addrs := make([]string, 0)
	for i, hash := range detail.BlockHash {
		if bytes.Equal(hash, blockHash) {
			addrs = append(addrs, detail.Addrs[i])
		}
	}
	return addrs
}

//
func mergeReceipt(receipt1, receipt2 *types.Receipt) *types.Receipt {
	if receipt2 != nil {
		receipt1.KV = append(receipt1.KV, receipt2.KV...)
		receipt1.Logs = append(receipt1.Logs, receipt2.Logs...)
	}

	return receipt1
}

func makeAddrBindReceipt(node, addr string, prev, current *pt.ParaBindMinerInfo) *types.Receipt {
	key := calcParaBindMinerAddr(node, addr)
	log := &pt.ReceiptParaBindMinerInfo{
		Addr:    addr,
		Prev:    prev,
		Current: current,
	}

	return &types.Receipt{
		Ty: types.ExecOk,
		KV: []*types.KeyValue{
			{Key: key, Value: types.Encode(current)},
		},
		Logs: []*types.ReceiptLog{
			{
				Ty:  pt.TyLogParaBindMinerAddr,
				Log: types.Encode(log),
			},
		},
	}
}

func makeNodeBindReceipt(prev, current *pt.ParaNodeBindList) *types.Receipt {
	key := calcParaBindMinerNode()
	log := &pt.ReceiptParaNodeBindListUpdate{
		Prev:    prev,
		Current: current,
	}

	return &types.Receipt{
		Ty: types.ExecOk,
		KV: []*types.KeyValue{
			{Key: key, Value: types.Encode(current)},
		},
		Logs: []*types.ReceiptLog{
			{
				Ty:  pt.TyLogParaBindMinerNode,
				Log: types.Encode(log),
			},
		},
	}
}

//
func (a *action) bind2Node(node string) (*types.Receipt, error) {
	list, err := getBindNodeInfo(a.db)
	if err != nil {
		return nil, errors.Wrap(err, "bind2Node")
	}

	//  kvmvcc    ，       nil，     ，           ，unbind ，           ，    ，title
	if len(list.Title) <= 0 {
		list.Title = a.api.GetConfig().GetTitle()
	}

	old := proto.Clone(list).(*pt.ParaNodeBindList)
	list.Miners = append(list.Miners, &pt.ParaNodeBindOne{SuperNode: node, Miner: a.fromaddr})

	return makeNodeBindReceipt(old, list), nil

}

//
func (a *action) unbind2Node(node string) (*types.Receipt, error) {
	list, err := getBindNodeInfo(a.db)
	if err != nil {
		return nil, errors.Wrap(err, "unbind2Node")
	}
	newList := &pt.ParaNodeBindList{Title: a.api.GetConfig().GetTitle()}
	old := proto.Clone(list).(*pt.ParaNodeBindList)

	for _, m := range list.Miners {
		if m.SuperNode == node && m.Miner == a.fromaddr {
			continue
		}
		newList.Miners = append(newList.Miners, m)
	}
	return makeNodeBindReceipt(old, newList), nil

}

func getBindNodeInfo(db dbm.KV) (*pt.ParaNodeBindList, error) {
	var list pt.ParaNodeBindList
	key := calcParaBindMinerNode()
	data, err := db.Get(key)
	if isNotFound(err) {
		return &list, nil
	}
	if err != nil {
		return nil, errors.Wrapf(err, "get key failed")
	}

	err = types.Decode(data, &list)
	if err != nil {
		return nil, errors.Wrapf(err, "decode failed")
	}
	return &list, nil
}

func getBindAddrInfo(db dbm.KV, node, addr string) (*pt.ParaBindMinerInfo, error) {
	key := calcParaBindMinerAddr(node, addr)
	data, err := db.Get(key)
	if err != nil {
		return nil, errors.Wrapf(err, "get key failed node=%s,addr=%s", node, addr)
	}

	var info pt.ParaBindMinerInfo
	err = types.Decode(data, &info)
	if err != nil {
		return nil, errors.Wrapf(err, "decode failed node=%s,addr=%s", node, addr)
	}
	return &info, nil
}

func (a *action) bindOp(cmd *pt.ParaBindMinerCmd) (*types.Receipt, error) {
	if cmd.BindCoins <= 0 {
		return nil, errors.Wrapf(types.ErrInvalidParam, "bindMiner BindCoins nil from addr %s", a.fromaddr)
	}

	err := a.isValidSuperNode(cmd.TargetNode)
	if err != nil {
		return nil, err
	}

	current, err := getBindAddrInfo(a.db, cmd.TargetNode, a.fromaddr)
	if err != nil && !isNotFound(errors.Cause(err)) {
		return nil, errors.Wrap(err, "getBindAddrInfo")
	}

	//found,
	if current != nil && current.BindStatus == opBind {
		var receipt *types.Receipt

		if cmd.BindCoins == current.BindCoins {
			return nil, errors.Wrapf(types.ErrInvalidParam, "bind coins same current=%d, cmd=%d", current.BindCoins, cmd.BindCoins)
		}

		//     coins
		if cmd.BindCoins < current.BindCoins {
			receipt, err = a.coinsAccount.ExecActive(a.fromaddr, a.execaddr, (current.BindCoins-cmd.BindCoins)*types.Coin)
			if err != nil {
				return nil, errors.Wrapf(err, "bindOp Active addr=%s,execaddr=%s,coins=%d", a.fromaddr, a.execaddr, current.BindCoins-cmd.BindCoins)
			}
		} else {
			//
			receipt, err = a.coinsAccount.ExecFrozen(a.fromaddr, a.execaddr, (cmd.BindCoins-current.BindCoins)*types.Coin)
			if err != nil {
				return nil, errors.Wrapf(err, "bindOp frozen more addr=%s,execaddr=%s,coins=%d", a.fromaddr, a.execaddr, cmd.BindCoins-current.BindCoins)
			}
		}

		acctCopy := *current
		current.BindCoins = cmd.BindCoins
		r := makeAddrBindReceipt(cmd.TargetNode, a.fromaddr, &acctCopy, current)
		return mergeReceipt(receipt, r), nil
	}

	//not bind,
	receipt, err := a.coinsAccount.ExecFrozen(a.fromaddr, a.execaddr, cmd.BindCoins*types.Coin)
	if err != nil {
		return nil, errors.Wrapf(err, "bindOp frozen addr=%s,execaddr=%s,count=%d", a.fromaddr, a.execaddr, cmd.BindCoins)
	}

	//bind addr
	newer := &pt.ParaBindMinerInfo{
		Addr:        a.fromaddr,
		BindStatus:  opBind,
		BindCoins:   cmd.BindCoins,
		BlockTime:   a.blocktime,
		BlockHeight: a.height,
		TargetNode:  cmd.TargetNode,
	}
	rBind := makeAddrBindReceipt(cmd.TargetNode, a.fromaddr, current, newer)
	mergeReceipt(receipt, rBind)

	//
	rList, err := a.bind2Node(cmd.TargetNode)
	if err != nil {
		return nil, err
	}
	mergeReceipt(receipt, rList)
	return receipt, nil

}

func (a *action) unBindOp(cmd *pt.ParaBindMinerCmd) (*types.Receipt, error) {
	acct, err := getBindAddrInfo(a.db, cmd.TargetNode, a.fromaddr)
	if err != nil {
		return nil, err
	}

	cfg := a.api.GetConfig()
	unBindHours := cfg.MGInt("mver.consensus.paracross.unBindTime", a.height)
	if a.blocktime-acct.BlockTime < unBindHours*60*60 {
		return nil, errors.Wrapf(types.ErrNotAllow, "unBindOp unbind time=%d less %d hours than bind time =%d", a.blocktime, unBindHours, acct.BlockTime)
	}

	if acct.BindStatus != opBind {
		return nil, errors.Wrapf(types.ErrNotAllow, "unBindOp,current addr is unbind status")
	}

	//unfrozen
	receipt, err := a.coinsAccount.ExecActive(a.fromaddr, a.execaddr, acct.BindCoins*types.Coin)
	if err != nil {
		return nil, errors.Wrapf(err, "unBindOp addr=%s,execaddr=%s,count=%d", a.fromaddr, a.execaddr, acct.BindCoins)
	}

	//   bind addr
	//  kvmvcc   ，       key =nil     ，kvmvcc          ，         ，&struct{}   ，len=0
	acctCopy := *acct
	acct.BindStatus = opUnBind
	acct.BlockHeight = a.height
	acct.BlockTime = a.blocktime
	rUnBind := makeAddrBindReceipt(cmd.TargetNode, a.fromaddr, &acctCopy, acct)
	mergeReceipt(receipt, rUnBind)

	//
	rUnList, err := a.unbind2Node(cmd.TargetNode)
	if err != nil {
		return nil, err
	}
	mergeReceipt(receipt, rUnList)

	return receipt, nil
}

func (a *action) bindMiner(info *pt.ParaBindMinerCmd) (*types.Receipt, error) {
	if len(info.TargetNode) <= 0 {
		return nil, errors.Wrapf(types.ErrInvalidParam, "bindMiner TargetNode should not be nil to addr %s", a.fromaddr)
	}

	//
	if !types.IsParaExecName(string(a.tx.Execer)) {
		return nil, errors.Wrapf(types.ErrInvalidParam, "exec=%s,should prefix with user.p.", string(a.tx.Execer))
	}

	if info.BindAction != opBind && info.BindAction != opUnBind {
		return nil, errors.Wrapf(types.ErrInvalidParam, "bindMiner action=%d not correct", info.BindAction)
	}

	if info.BindAction == opBind {
		return a.bindOp(info)
	}
	return a.unBindOp(info)
}

package executor

import (
	dbm "github.com/D-PlatformOperatingSystem/dpos/common/db"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	pty "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/trade/types"
	"github.com/pkg/errors"
)

const (
	tradeLocaldbVersioin = "LODB-trade-version"
)

//
//    key -> id           ，
//
const (
	sellOrderSHTAS = "LODB-trade-sellorder-shtas:"
	sellOrderASTS  = "LODB-trade-sellorder-asts:"
	sellOrderATSS  = "LODB-trade-sellorder-atss:"
	sellOrderTSPAS = "LODB-trade-sellorder-tspas:"
	buyOrderSHTAS  = "LODB-trade-buyorder-shtas:"
	buyOrderASTS   = "LODB-trade-buyorder-asts:"
	buyOrderATSS   = "LODB-trade-buyorder-atss:"
	buyOrderTSPAS  = "LODB-trade-buyorder-tspas:"
	// Addr-Status-Type-Height-Key
	orderASTHK = "LODB-trade-order-asthk:"
)

// Upgrade
func (t *trade) Upgrade() (*types.LocalDBSet, error) {
	localDB := t.GetLocalDB()
	//      coins symbol，
	coinSymbol := t.GetAPI().GetConfig().GetCoinSymbol()
	kvs, err := UpgradeLocalDBV2(localDB, coinSymbol)
	if err != nil {
		tradelog.Error("Upgrade failed", "err", err)
		return nil, errors.Cause(err)
	}
	return kvs, nil
}

// UpgradeLocalDBV2 trade
// from 1 to 2
func UpgradeLocalDBV2(localDB dbm.KVDB, coinSymbol string) (*types.LocalDBSet, error) {
	toVersion := 2
	tradelog.Info("UpgradeLocalDBV2 upgrade start", "to_version", toVersion)
	version, err := getVersion(localDB)
	if err != nil {
		return nil, errors.Wrap(err, "UpgradeLocalDBV2 get version")
	}
	if version >= toVersion {
		tradelog.Debug("UpgradeLocalDBV2 not need to upgrade", "current_version", version, "to_version", toVersion)
		return nil, nil
	}

	var kvset types.LocalDBSet
	kvs, err := UpgradeLocalDBPart2(localDB, coinSymbol)
	if err != nil {
		return nil, errors.Wrap(err, "UpgradeLocalDBV2 UpgradeLocalDBPart2")
	}
	if len(kvs) > 0 {
		kvset.KV = append(kvset.KV, kvs...)
	}

	kvs, err = UpgradeLocalDBPart1(localDB)
	if err != nil {
		return nil, errors.Wrap(err, "UpgradeLocalDBV2 UpgradeLocalDBPart1")
	}
	if len(kvs) > 0 {
		kvset.KV = append(kvset.KV, kvs...)
	}

	kvs, err = setVersion(localDB, toVersion)
	if err != nil {
		return nil, errors.Wrap(err, "UpgradeLocalDBV2 setVersion")
	}
	if len(kvs) > 0 {
		kvset.KV = append(kvset.KV, kvs...)
	}

	tradelog.Info("UpgradeLocalDBV2 upgrade done")
	return &kvset, nil
}

// UpgradeLocalDBPart1     KV，
func UpgradeLocalDBPart1(localDB dbm.KVDB) ([]*types.KeyValue, error) {
	prefixes := []string{
		sellOrderSHTAS,
		sellOrderASTS,
		sellOrderATSS,
		sellOrderTSPAS,
		buyOrderSHTAS,
		buyOrderASTS,
		buyOrderATSS,
		buyOrderTSPAS,
		orderASTHK,
	}

	var allKvs []*types.KeyValue
	for _, prefix := range prefixes {
		kvs, err := delOnePrefix(localDB, prefix)
		if err != nil {
			return nil, errors.Wrapf(err, "UpdateLocalDBPart1 delOnePrefix: %s", prefix)
		}
		if len(kvs) > 0 {
			allKvs = append(allKvs, kvs...)
		}

	}
	return allKvs, nil
}

// delOnePrefix
func delOnePrefix(localDB dbm.KVDB, prefix string) ([]*types.KeyValue, error) {
	start := []byte(prefix)
	keys, err := localDB.List(start, nil, 0, dbm.ListASC|dbm.ListKeyOnly)
	if err != nil {
		if err == types.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}

	var kvs []*types.KeyValue
	tradelog.Debug("delOnePrefix", "len", len(keys), "prefix", prefix)
	for _, key := range keys {
		err = localDB.Set(key, nil)
		if err != nil {
			return nil, err
		}
		kvs = append(kvs, &types.KeyValue{Key: key, Value: nil})
	}

	return kvs, nil
}

// UpgradeLocalDBPart2   order
// order   v1     v2
//   tableV1   ，   tableV2   ,
func UpgradeLocalDBPart2(kvdb dbm.KVDB, coinSymbol string) ([]*types.KeyValue, error) {
	return upgradeOrder(kvdb, coinSymbol)
}

func upgradeOrder(kvdb dbm.KVDB, coinSymbol string) ([]*types.KeyValue, error) {
	tab2 := NewOrderTableV2(kvdb)
	tab := NewOrderTable(kvdb)
	q1 := tab.GetQuery(kvdb)

	var order1 pty.LocalOrder
	rows, err := q1.List("key", &order1, []byte(""), 0, 0)
	if err != nil {
		if err == types.ErrNotFound {
			return nil, nil
		}
		return nil, errors.Wrap(err, "upgradeOrder list from order v1 table")
	}

	tradelog.Debug("upgradeOrder", "len", len(rows))
	for _, row := range rows {
		o1, ok := row.Data.(*pty.LocalOrder)
		if !ok {
			return nil, errors.Wrap(types.ErrTypeAsset, "decode order v1")
		}

		o2 := types.Clone(o1).(*pty.LocalOrder)
		upgradeLocalOrder(o2, coinSymbol)
		err = tab2.Add(o2)
		if err != nil {
			return nil, errors.Wrap(err, "upgradeOrder add to order v2 table")
		}

		err = tab.Del([]byte(o1.TxIndex))
		if err != nil {
			return nil, errors.Wrapf(err, "upgradeOrder del from order v1 table, key: %s", o1.TxIndex)
		}
	}

	kvs, err := tab2.Save()
	if err != nil {
		return nil, errors.Wrap(err, "upgradeOrder save-add to order v2 table")
	}
	kvs2, err := tab.Save()
	if err != nil {
		return nil, errors.Wrap(err, "upgradeOrder save-del to order v1 table")
	}
	kvs = append(kvs, kvs2...)

	for _, kv := range kvs {
		tradelog.Debug("upgradeOrder", "KEY", string(kv.GetKey()))
		err = kvdb.Set(kv.GetKey(), kv.GetValue())
		if err != nil {
			return nil, errors.Wrapf(err, "upgradeOrder set localdb key: %s", string(kv.GetKey()))
		}
	}

	return kvs, nil
}

// upgradeLocalOrder     fork
// 1.
// 2.
func upgradeLocalOrder(order *pty.LocalOrder, coinSymbol string) {
	if order.AssetExec == "" {
		order.AssetExec = defaultAssetExec
	}
	if order.PriceExec == "" {
		order.PriceExec = defaultPriceExec
		order.PriceSymbol = coinSymbol
	}
}

// localdb Version
func getVersion(kvdb dbm.KV) (int, error) {
	value, err := kvdb.Get([]byte(tradeLocaldbVersioin))
	if err != nil && err != types.ErrNotFound {
		return 1, err
	}
	if err == types.ErrNotFound {
		return 1, nil
	}
	var v types.Int32
	err = types.Decode(value, &v)
	if err != nil {
		return 1, err
	}
	return int(v.Data), nil
}

func setVersion(kvdb dbm.KV, version int) ([]*types.KeyValue, error) {
	v := types.Int32{Data: int32(version)}
	x := types.Encode(&v)
	err := kvdb.Set([]byte(tradeLocaldbVersioin), x)
	return []*types.KeyValue{{Key: []byte(tradeLocaldbVersioin), Value: x}}, err
}

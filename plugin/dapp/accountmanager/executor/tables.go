package executor

import (
	"fmt"

	"github.com/D-PlatformOperatingSystem/dpos/common/db"
	"github.com/D-PlatformOperatingSystem/dpos/common/db/table"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	aty "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/accountmanager/types"
)

/*
 *       kv   ，key
 *  key = keyPrefix + userKey
 *          ，  ’-‘
 */

const (
	//KeyPrefixStateDB state db key
	KeyPrefixStateDB = "mavl-accountmanager-"
	//KeyPrefixLocalDB local db key
	KeyPrefixLocalDB = "LODB-accountmanager"
)

var opt_account = &table.Option{
	Prefix:  KeyPrefixLocalDB,
	Name:    "account",
	Primary: "index",
	Index:   []string{"status", "accountID", "addr"},
}

//
func calcAccountKey(accountID string) []byte {
	key := fmt.Sprintf("%s"+"accountID:%s", KeyPrefixStateDB, accountID)
	return []byte(key)
}

//NewAccountTable ...
func NewAccountTable(kvdb db.KV) *table.Table {
	rowmeta := NewAccountRow()
	table, err := table.NewTable(rowmeta, kvdb, opt_account)
	if err != nil {
		panic(err)
	}
	return table
}

//AccountRow account table meta
type AccountRow struct {
	*aty.Account
}

//NewAccountRow     meta
func NewAccountRow() *AccountRow {
	return &AccountRow{Account: &aty.Account{}}
}

//CreateRow      (  index             ,     eventid)
func (m *AccountRow) CreateRow() *table.Row {
	return &table.Row{Data: &aty.Account{}}
}

//SetPayload
func (m *AccountRow) SetPayload(data types.Message) error {
	if txdata, ok := data.(*aty.Account); ok {
		m.Account = txdata
		return nil
	}
	return types.ErrTypeAsset
}

//Get   indexName    indexValue
func (m *AccountRow) Get(key string) ([]byte, error) {
	if key == "accountID" {
		return []byte(m.AccountID), nil
	} else if key == "status" {
		return []byte(fmt.Sprintf("%d", m.Status)), nil
	} else if key == "index" {
		return []byte(fmt.Sprintf("%015d", m.GetIndex())), nil
	} else if key == "addr" {
		return []byte(m.GetAddr()), nil
	}
	return nil, types.ErrNotFound
}

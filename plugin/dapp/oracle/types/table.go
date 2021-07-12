package types

import (
	"fmt"

	"github.com/D-PlatformOperatingSystem/dpos/common/db"

	"github.com/D-PlatformOperatingSystem/dpos/common/db/table"
	"github.com/D-PlatformOperatingSystem/dpos/types"
)

/*
table  struct
data:  oracle
index: addr,status,index,type
*/

var opt = &table.Option{
	Prefix:  "LODB",
	Name:    "oracle",
	Primary: "eventid",
	Index:   []string{"status", "addr_status", "type_status"},
}

//NewTable
func NewTable(kvdb db.KV) *table.Table {
	rowmeta := NewOracleRow()
	table, err := table.NewTable(rowmeta, kvdb, opt)
	if err != nil {
		panic(err)
	}
	return table
}

//OracleRow table meta
type OracleRow struct {
	*ReceiptOracle
}

//NewOracleRow     meta
func NewOracleRow() *OracleRow {
	return &OracleRow{ReceiptOracle: &ReceiptOracle{}}
}

//CreateRow      (  index             ,     eventid)
func (tx *OracleRow) CreateRow() *table.Row {
	return &table.Row{Data: &ReceiptOracle{}}
}

//SetPayload
func (tx *OracleRow) SetPayload(data types.Message) error {
	if txdata, ok := data.(*ReceiptOracle); ok {
		tx.ReceiptOracle = txdata
		return nil
	}
	return types.ErrTypeAsset
}

//Get   indexName    indexValue
func (tx *OracleRow) Get(key string) ([]byte, error) {
	if key == "eventid" {
		return []byte(tx.EventID), nil
	} else if key == "status" {
		return []byte(fmt.Sprintf("%2d", tx.Status)), nil
	} else if key == "addr_status" {
		return []byte(fmt.Sprintf("%s:%2d", tx.Addr, tx.Status)), nil
	} else if key == "type_status" {
		return []byte(fmt.Sprintf("%s:%2d", tx.Type, tx.Status)), nil
	}
	return nil, types.ErrNotFound
}

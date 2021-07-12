package types

import (
	"fmt"

	"github.com/D-PlatformOperatingSystem/dpos/common/db"

	"github.com/D-PlatformOperatingSystem/dpos/common/db/table"
	"github.com/D-PlatformOperatingSystem/dpos/types"
)

var opt = &table.Option{
	Prefix:  "LODB-collateralize",
	Name:    "coller",
	Primary: "collateralizeid",
	Index:   []string{"status", "addr", "addr_status"},
}

//NewCollateralizeTable
func NewCollateralizeTable(kvdb db.KV) *table.Table {
	rowmeta := NewCollatetalizeRow()
	table, err := table.NewTable(rowmeta, kvdb, opt)
	if err != nil {
		panic(err)
	}
	return table
}

//CollatetalizeRow table meta
type CollatetalizeRow struct {
	*ReceiptCollateralize
}

//NewCollatetalizeRow     meta
func NewCollatetalizeRow() *CollatetalizeRow {
	return &CollatetalizeRow{ReceiptCollateralize: &ReceiptCollateralize{}}
}

//CreateRow
func (tx *CollatetalizeRow) CreateRow() *table.Row {
	return &table.Row{Data: &ReceiptCollateralize{}}
}

//SetPayload
func (tx *CollatetalizeRow) SetPayload(data types.Message) error {
	if txdata, ok := data.(*ReceiptCollateralize); ok {
		tx.ReceiptCollateralize = txdata
		return nil
	}
	return types.ErrTypeAsset
}

//Get   indexName    indexValue
func (tx *CollatetalizeRow) Get(key string) ([]byte, error) {
	if key == "collateralizeid" {
		return []byte(tx.CollateralizeId), nil
	} else if key == "status" {
		return []byte(fmt.Sprintf("%2d", tx.Status)), nil
	} else if key == "addr" {
		return []byte(tx.AccountAddr), nil
	} else if key == "addr_status" {
		return []byte(fmt.Sprintf("%s:%2d", tx.AccountAddr, tx.Status)), nil
	}
	return nil, types.ErrNotFound
}

var optRecord = &table.Option{
	Prefix:  "LODB-collateralize",
	Name:    "borrow",
	Primary: "borrowid",
	Index:   []string{"status", "addr", "addr_status", "id_status", "id_addr"},
}

// NewRecordTable
func NewRecordTable(kvdb db.KV) *table.Table {
	rowmeta := NewRecordRow()
	table, err := table.NewTable(rowmeta, kvdb, optRecord)
	if err != nil {
		panic(err)
	}
	return table
}

//CollateralizeRecordRow table meta
type CollateralizeRecordRow struct {
	*ReceiptCollateralize
}

//NewRecordRow     meta
func NewRecordRow() *CollateralizeRecordRow {
	return &CollateralizeRecordRow{ReceiptCollateralize: &ReceiptCollateralize{}}
}

//CreateRow
func (tx *CollateralizeRecordRow) CreateRow() *table.Row {
	return &table.Row{Data: &ReceiptCollateralize{}}
}

//SetPayload
func (tx *CollateralizeRecordRow) SetPayload(data types.Message) error {
	if txdata, ok := data.(*ReceiptCollateralize); ok {
		tx.ReceiptCollateralize = txdata
		return nil
	}
	return types.ErrTypeAsset
}

//Get   indexName    indexValue
func (tx *CollateralizeRecordRow) Get(key string) ([]byte, error) {
	if key == "borrowid" {
		return []byte(tx.RecordId), nil
	} else if key == "status" {
		return []byte(fmt.Sprintf("%2d", tx.Status)), nil
	} else if key == "addr" {
		return []byte(tx.AccountAddr), nil
	} else if key == "addr_status" {
		return []byte(fmt.Sprintf("%s:%2d", tx.AccountAddr, tx.Status)), nil
	} else if key == "id_status" {
		return []byte(fmt.Sprintf("%s:%2d", tx.CollateralizeId, tx.Status)), nil
	} else if key == "id_addr" {
		return []byte(fmt.Sprintf("%s:%s", tx.CollateralizeId, tx.AccountAddr)), nil
	}
	return nil, types.ErrNotFound
}

package executor

import (
	"fmt"

	"github.com/D-PlatformOperatingSystem/dpos/common/db"

	"github.com/D-PlatformOperatingSystem/dpos/common/db/table"
	"github.com/D-PlatformOperatingSystem/dpos/system/dapp"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	auty "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/autonomy/types"
)

/*
table  struct
data:  autonomy board
index: status, addr
*/

var boardOpt = &table.Option{
	Prefix:  "LODB-autonomy",
	Name:    "board",
	Primary: "heightindex",
	Index:   []string{"addr", "status", "addr_status"},
}

//NewBoardTable
func NewBoardTable(kvdb db.KV) *table.Table {
	rowmeta := NewBoardRow()
	table, err := table.NewTable(rowmeta, kvdb, boardOpt)
	if err != nil {
		panic(err)
	}
	return table
}

//BoardRow table meta
type BoardRow struct {
	*auty.AutonomyProposalBoard
}

//NewBoardRow     meta
func NewBoardRow() *BoardRow {
	return &BoardRow{AutonomyProposalBoard: &auty.AutonomyProposalBoard{}}
}

//CreateRow      (  index             ,     heightindex)
func (r *BoardRow) CreateRow() *table.Row {
	return &table.Row{Data: &auty.AutonomyProposalBoard{}}
}

//SetPayload
func (r *BoardRow) SetPayload(data types.Message) error {
	if d, ok := data.(*auty.AutonomyProposalBoard); ok {
		r.AutonomyProposalBoard = d
		return nil
	}
	return types.ErrTypeAsset
}

//Get   indexName    indexValue
func (r *BoardRow) Get(key string) ([]byte, error) {
	if key == "heightindex" {
		return []byte(dapp.HeightIndexStr(r.Height, int64(r.Index))), nil
	} else if key == "status" {
		return []byte(fmt.Sprintf("%2d", r.Status)), nil
	} else if key == "addr" {
		return []byte(r.Address), nil
	} else if key == "addr_status" {
		return []byte(fmt.Sprintf("%s:%2d", r.Address, r.Status)), nil
	}
	return nil, types.ErrNotFound
}

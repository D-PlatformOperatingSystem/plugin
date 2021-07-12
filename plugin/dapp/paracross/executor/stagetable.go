package executor

import (
	"fmt"

	"github.com/D-PlatformOperatingSystem/dpos/common/db"
	"github.com/D-PlatformOperatingSystem/dpos/common/db/table"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	pt "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/paracross/types"
)

/*
table  struct
data:  self consens stage
index: status
*/

var boardOpt = &table.Option{
	Prefix:  "LODB-paracross",
	Name:    "stage",
	Primary: "heightindex",
	Index:   []string{"id", "status"},
}

//NewStageTable
func NewStageTable(kvdb db.KV) *table.Table {
	rowmeta := NewStageRow()
	table, err := table.NewTable(rowmeta, kvdb, boardOpt)
	if err != nil {
		panic(err)
	}
	return table
}

//StageRow table meta
type StageRow struct {
	*pt.LocalSelfConsStageInfo
}

//NewStageRow     meta
func NewStageRow() *StageRow {
	return &StageRow{LocalSelfConsStageInfo: &pt.LocalSelfConsStageInfo{}}
}

//CreateRow      (  index             ,     heightindex)
func (r *StageRow) CreateRow() *table.Row {
	return &table.Row{Data: &pt.LocalSelfConsStageInfo{}}
}

//SetPayload
func (r *StageRow) SetPayload(data types.Message) error {
	if d, ok := data.(*pt.LocalSelfConsStageInfo); ok {
		r.LocalSelfConsStageInfo = d
		return nil
	}
	return types.ErrTypeAsset
}

//Get   indexName    indexValue
func (r *StageRow) Get(key string) ([]byte, error) {
	if key == "heightindex" {
		return []byte(r.TxIndex), nil
	} else if key == "id" {
		return []byte(r.Stage.Id), nil
	} else if key == "status" {
		return []byte(fmt.Sprintf("%2d", r.Stage.Status)), nil
	}

	return nil, types.ErrNotFound
}

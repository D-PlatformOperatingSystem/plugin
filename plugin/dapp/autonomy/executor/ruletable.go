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
data:  autonomy rule
index: status, addr
*/

var ruleOpt = &table.Option{
	Prefix:  "LODB-autonomy",
	Name:    "rule",
	Primary: "heightindex",
	Index:   []string{"addr", "status", "addr_status"},
}

//NewRuleTable
func NewRuleTable(kvdb db.KV) *table.Table {
	rowmeta := NewRuleRow()
	table, err := table.NewTable(rowmeta, kvdb, ruleOpt)
	if err != nil {
		panic(err)
	}
	return table
}

//RuleRow table meta
type RuleRow struct {
	*auty.AutonomyProposalRule
}

//NewRuleRow     meta
func NewRuleRow() *RuleRow {
	return &RuleRow{AutonomyProposalRule: &auty.AutonomyProposalRule{}}
}

//CreateRow      (  index             ,     heightindex)
func (r *RuleRow) CreateRow() *table.Row {
	return &table.Row{Data: &auty.AutonomyProposalRule{}}
}

//SetPayload
func (r *RuleRow) SetPayload(data types.Message) error {
	if d, ok := data.(*auty.AutonomyProposalRule); ok {
		r.AutonomyProposalRule = d
		return nil
	}
	return types.ErrTypeAsset
}

//Get   indexName    indexValue
func (r *RuleRow) Get(key string) ([]byte, error) {
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

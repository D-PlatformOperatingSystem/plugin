package executor

import (
	"github.com/D-PlatformOperatingSystem/dpos/types"
)

/*
 *
 */

// ExecDelLocal       ，
func (e *exchange) ExecDelLocal(tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	kvs, err := e.DelRollbackKV(tx, tx.Execer)
	if err != nil {
		return nil, err
	}
	dbSet := &types.LocalDBSet{}
	dbSet.KV = append(dbSet.KV, kvs...)
	return dbSet, nil
}

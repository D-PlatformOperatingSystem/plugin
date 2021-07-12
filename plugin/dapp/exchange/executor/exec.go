package executor

import (
	"github.com/D-PlatformOperatingSystem/dpos/types"
	exchangetypes "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/exchange/types"
)

/*
 *
 *       （statedb）       （log）
 */

func (e *exchange) Exec_LimitOrder(payload *exchangetypes.LimitOrder, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(e, tx, index)
	return action.LimitOrder(payload)
}

func (e *exchange) Exec_MarketOrder(payload *exchangetypes.MarketOrder, tx *types.Transaction, index int) (*types.Receipt, error) {
	//TODO marketOrder
	return nil, types.ErrActionNotSupport
}

func (e *exchange) Exec_RevokeOrder(payload *exchangetypes.RevokeOrder, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(e, tx, index)
	return action.RevokeOrder(payload)
}

package executor

import (
	"github.com/D-PlatformOperatingSystem/dpos/types"
	storagetypes "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/storage/types"
)

/*
 *
 *       （statedb）       （log）
 */

func (s *storage) Exec_ContentStorage(payload *storagetypes.ContentOnlyNotaryStorage, tx *types.Transaction, index int) (*types.Receipt, error) {
	s.GetAPI()
	action := newStorageAction(s, tx, index)
	return action.ContentStorage(payload)
}

func (s *storage) Exec_HashStorage(payload *storagetypes.HashOnlyNotaryStorage, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newStorageAction(s, tx, index)
	return action.HashStorage(payload)
}

func (s *storage) Exec_LinkStorage(payload *storagetypes.LinkNotaryStorage, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newStorageAction(s, tx, index)
	return action.LinkStorage(payload)
}

func (s *storage) Exec_EncryptStorage(payload *storagetypes.EncryptNotaryStorage, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newStorageAction(s, tx, index)
	return action.EncryptStorage(payload)
}

func (s *storage) Exec_EncryptShareStorage(payload *storagetypes.EncryptShareNotaryStorage, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newStorageAction(s, tx, index)
	return action.EncryptShareStorage(payload)
}

func (s *storage) Exec_EncryptAdd(payload *storagetypes.EncryptNotaryAdd, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newStorageAction(s, tx, index)
	return action.EncryptAdd(payload)
}

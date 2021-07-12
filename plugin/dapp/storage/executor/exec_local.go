package executor

import (
	"github.com/D-PlatformOperatingSystem/dpos/types"
	ety "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/storage/types"
)

/*
 *             ，
 *      ，    (localDB),       ，
 */

func (s *storage) ExecLocal_ContentStorage(payload *ety.ContentOnlyNotaryStorage, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	dbSet := &types.LocalDBSet{}
	cfg := s.GetAPI().GetConfig()
	if cfg.IsDappFork(s.GetHeight(), ety.StorageX, ety.ForkStorageLocalDB) {
		if receiptData.Ty == types.ExecOk {
			for _, log := range receiptData.Logs {
				switch log.Ty {
				case ety.TyContentStorageLog:
					storage := &ety.Storage{}
					if err := types.Decode(log.Log, storage); err != nil {
						return nil, err
					}
					kv := &types.KeyValue{Key: getLocalDBKey(storage.GetContentStorage().Key), Value: types.Encode(storage)}
					dbSet.KV = append(dbSet.KV, kv)
				}
			}
		}
	}
	return s.addAutoRollBack(tx, dbSet.KV), nil
}

func (s *storage) ExecLocal_HashStorage(payload *ety.HashOnlyNotaryStorage, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	dbSet := &types.LocalDBSet{}
	cfg := s.GetAPI().GetConfig()
	if cfg.IsDappFork(s.GetHeight(), ety.StorageX, ety.ForkStorageLocalDB) {
		if receiptData.Ty == types.ExecOk {
			for _, log := range receiptData.Logs {
				switch log.Ty {
				case ety.TyHashStorageLog:
					storage := &ety.Storage{}
					if err := types.Decode(log.Log, storage); err != nil {
						return nil, err
					}
					kv := &types.KeyValue{Key: getLocalDBKey(storage.GetHashStorage().Key), Value: types.Encode(storage)}
					dbSet.KV = append(dbSet.KV, kv)
				}
			}
		}
	}
	return s.addAutoRollBack(tx, dbSet.KV), nil
}

func (s *storage) ExecLocal_LinkStorage(payload *ety.LinkNotaryStorage, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	dbSet := &types.LocalDBSet{}
	cfg := s.GetAPI().GetConfig()
	if cfg.IsDappFork(s.GetHeight(), ety.StorageX, ety.ForkStorageLocalDB) {
		if receiptData.Ty == types.ExecOk {
			for _, log := range receiptData.Logs {
				switch log.Ty {
				case ety.TyLinkStorageLog:
					storage := &ety.Storage{}
					if err := types.Decode(log.Log, storage); err != nil {
						return nil, err
					}
					kv := &types.KeyValue{Key: getLocalDBKey(storage.GetLinkStorage().Key), Value: types.Encode(storage)}
					dbSet.KV = append(dbSet.KV, kv)
				}
			}
		}
	}
	return s.addAutoRollBack(tx, dbSet.KV), nil
}

func (s *storage) ExecLocal_EncryptStorage(payload *ety.EncryptNotaryStorage, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	dbSet := &types.LocalDBSet{}
	cfg := s.GetAPI().GetConfig()
	if cfg.IsDappFork(s.GetHeight(), ety.StorageX, ety.ForkStorageLocalDB) {
		if receiptData.Ty == types.ExecOk {
			for _, log := range receiptData.Logs {
				switch log.Ty {
				case ety.TyEncryptStorageLog:
					storage := &ety.Storage{}
					if err := types.Decode(log.Log, storage); err != nil {
						return nil, err
					}
					kv := &types.KeyValue{Key: getLocalDBKey(storage.GetEncryptStorage().Key), Value: types.Encode(storage)}
					dbSet.KV = append(dbSet.KV, kv)
				}
			}
		}
	}
	return s.addAutoRollBack(tx, dbSet.KV), nil
}

func (s *storage) ExecLocal_EncryptShareStorage(payload *ety.EncryptShareNotaryStorage, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	dbSet := &types.LocalDBSet{}
	cfg := s.GetAPI().GetConfig()
	if cfg.IsDappFork(s.GetHeight(), ety.StorageX, ety.ForkStorageLocalDB) {
		if receiptData.Ty == types.ExecOk {
			for _, log := range receiptData.Logs {
				switch log.Ty {
				case ety.TyEncryptShareStorageLog:
					storage := &ety.Storage{}
					if err := types.Decode(log.Log, storage); err != nil {
						return nil, err
					}
					kv := &types.KeyValue{Key: getLocalDBKey(storage.GetEncryptShareStorage().Key), Value: types.Encode(storage)}
					dbSet.KV = append(dbSet.KV, kv)
				}
			}
		}
	}
	return s.addAutoRollBack(tx, dbSet.KV), nil
}

func (s *storage) ExecLocal_EncryptAdd(payload *ety.EncryptNotaryAdd, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	dbSet := &types.LocalDBSet{}
	cfg := s.GetAPI().GetConfig()
	if cfg.IsDappFork(s.GetHeight(), ety.StorageX, ety.ForkStorageLocalDB) {
		if receiptData.Ty == types.ExecOk {
			for _, log := range receiptData.Logs {
				switch log.Ty {
				case ety.TyEncryptAddLog:
					storage := &ety.Storage{}
					if err := types.Decode(log.Log, storage); err != nil {
						return nil, err
					}
					kv := &types.KeyValue{Key: getLocalDBKey(storage.GetEncryptStorage().Key), Value: types.Encode(storage)}
					dbSet.KV = append(dbSet.KV, kv)
				}
			}
		}
	}
	return s.addAutoRollBack(tx, dbSet.KV), nil
}

//
func (s *storage) addAutoRollBack(tx *types.Transaction, kv []*types.KeyValue) *types.LocalDBSet {

	dbSet := &types.LocalDBSet{}
	dbSet.KV = s.AddRollbackKV(tx, tx.Execer, kv)
	return dbSet
}

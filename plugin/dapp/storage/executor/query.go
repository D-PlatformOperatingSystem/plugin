package executor

import (
	"github.com/D-PlatformOperatingSystem/dpos/types"
	storagetypes "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/storage/types"
)

// statedb       
func (s *storage) Query_QueryStorage(in *storagetypes.QueryStorage) (types.Message, error) {
	return QueryStorage(s.GetStateDB(), s.GetLocalDB(), in.TxHash)
}

//      ids
func (s *storage) Query_BatchQueryStorage(in *storagetypes.BatchQueryStorage) (types.Message, error) {
	return BatchQueryStorage(s.GetStateDB(), s.GetLocalDB(), in)
}

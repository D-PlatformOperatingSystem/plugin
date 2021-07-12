package executor

import (
	"github.com/D-PlatformOperatingSystem/dpos/types"
	et "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/accountmanager/types"
)

//Query_QueryAccountByID   ID      
func (a *Accountmanager) Query_QueryAccountByID(in *et.QueryAccountByID) (types.Message, error) {
	return findAccountByID(a.GetLocalDB(), in.AccountID)
}

//Query_QueryAccountByAddr   ID      
func (a *Accountmanager) Query_QueryAccountByAddr(in *et.QueryAccountByAddr) (types.Message, error) {
	return findAccountByAddr(a.GetLocalDB(), in.Addr)
}

//Query_QueryAccountsByStatus           ||       1   ， 2    , 3     4,    
func (a *Accountmanager) Query_QueryAccountsByStatus(in *et.QueryAccountsByStatus) (types.Message, error) {
	return findAccountListByStatus(a.GetLocalDB(), in.Status, in.Direction, in.PrimaryKey)
}

//Query_QueryExpiredAccounts            
func (a *Accountmanager) Query_QueryExpiredAccounts(in *et.QueryExpiredAccounts) (types.Message, error) {
	return findAccountListByIndex(a.GetLocalDB(), in.ExpiredTime, in.PrimaryKey)
}

//Query_QueryBalanceByID   ID      
func (a *Accountmanager) Query_QueryBalanceByID(in *et.QueryBalanceByID) (types.Message, error) {
	return queryBalanceByID(a.GetStateDB(), a.GetLocalDB(), a.GetAPI().GetConfig(), a.GetName(), in)
}

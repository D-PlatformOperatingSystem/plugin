package test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/D-PlatformOperatingSystem/dpos/common"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/exchange/executor"
	et "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/exchange/types"
	tt "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/token/types"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
)

//GRPCCli ...
type GRPCCli struct {
	client types.DplatformOSClient
}

//NewGRPCCli ...
func NewGRPCCli(grpcAddr string) *GRPCCli {
	conn, err := grpc.Dial(grpcAddr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	client := types.NewDplatformOSClient(conn)
	cfg := types.NewDplatformOSConfig(types.GetDefaultCfgstring())
	cfg.SetTitleOnlyForTest("dplatformos")
	executor.Init(et.ExchangeX, cfg, nil)
	return &GRPCCli{
		client: client,
	}
}

//Send ...
func (c *GRPCCli) Send(tx *types.Transaction, hexKey string) ([]*types.ReceiptLog, error) {
	txHash, logs, err := c.sendAndWaitReceipt(tx, hexKey)
	if txHash != nil {
		fmt.Println("txHash: ", common.ToHex(txHash))
	}
	if err != nil {
		return nil, parseError(err)
	}
	for _, l := range logs {
		if l.Ty == types.TyLogErr {
			return nil, errors.New(string(l.Log))
		}
	}
	return logs, nil
}

//Query ...
func (c *GRPCCli) Query(fn string, msg proto.Message) ([]byte, error) {
	ss := strings.Split(fn, ".")
	var in types.ChainExecutor
	if len(ss) == 2 {
		in.Driver = ss[0]
		in.FuncName = ss[1]
	} else {
		in.Driver = et.ExchangeX
		in.FuncName = fn
	}
	in.Param = types.Encode(msg)

	r, err := c.client.QueryChain(context.Background(), &in)
	if err != nil {
		return nil, err
	}
	if !r.IsOk {
		return nil, errors.New(string(r.Msg))
	}
	return r.Msg, nil
}

//GetExecAccount ...
func (c *GRPCCli) GetExecAccount(addr string, exec string, symbol string) (*types.Account, error) {
	if exec == "coins" {
		// dpos
		var addrs []string
		addrs = append(addrs, addr)
		params := &types.ReqBalance{
			Addresses: addrs,
			Execer:    et.ExchangeX,
		}

		accs, err := c.client.GetBalance(context.Background(), params)
		if err != nil {
			return nil, err
		}
		return accs.Acc[0], nil
	}

	// token: ccny
	param := &tt.ReqAccountTokenAssets{
		Address: addr,
		Execer:  et.ExchangeX,
	}
	msg, err := c.Query("token.GetAccountTokenAssets", param)
	if err != nil {
		return nil, err
	}

	var resp tt.ReplyAccountTokenAssets
	err = types.Decode(msg, &resp)
	if err != nil {
		return nil, err
	}

	for _, v := range resp.TokenAssets {
		if v.Symbol == symbol {
			return v.Account, nil
		}
	}

	return nil, types.ErrNotFound
}

//
//       ，
//         ，       、
func (c *GRPCCli) sendAndWaitReceipt(tx *types.Transaction, hexKey string) (txHash []byte, logs []*types.ReceiptLog, err error) {
	r, err := c.SendTx(tx, hexKey)
	if err != nil {
		// rpc error: code = Unknown desc = ErrNotBank
		return nil, nil, err
	}
	if !r.IsOk {
		return nil, nil, errors.New(string(r.Msg))
	}
	time.Sleep(time.Second)
	d, _ := c.client.QueryTransaction(context.Background(), &types.ReqHash{Hash: r.Msg})
	return r.Msg, d.Receipt.Logs, nil
}

//SendTx ...
func (c *GRPCCli) SendTx(tx *types.Transaction, hexKey string) (reply *types.Reply, err error) {
	cfg := types.NewDplatformOSConfig(types.GetDefaultCfgstring())
	cfg.SetTitleOnlyForTest("dplatformos")
	tx, err = types.FormatTx(cfg, et.ExchangeX, tx)
	if err != nil {
		return nil, err
	}
	tx, err = signTx(tx, hexKey)
	if err != nil {
		return nil, err
	}

	return c.client.SendTransaction(context.Background(), tx)
}

func parseError(err error) error {
	// rpc error: code = Unknown desc = ErrNotBank
	str := err.Error()
	sep := "desc = "
	i := strings.Index(str, sep)
	if i != -1 {
		return errors.New(str[i+len(sep):])
	}
	return err
}

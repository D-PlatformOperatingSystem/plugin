package rpc

import (
	"context"

	rpctypes "github.com/D-PlatformOperatingSystem/dpos/rpc/types"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	echotypes "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/echo/types/echo"
)

// Jrpc        RPC
type Jrpc struct {
	cli *channelClient
}

// RPC
type channelClient struct {
	rpctypes.ChannelClient
}

// Init    rpc
func Init(name string, s rpctypes.RPCServer) {
	cli := &channelClient{}
	//       ，     Jrpc，    grpc
	cli.Init(name, s, &Jrpc{cli: cli}, nil)
}

// QueryPing                Query  ，      rpc Query
//        ，           ，
func (c *Jrpc) QueryPing(param *echotypes.Query, result *interface{}) error {
	if param == nil {
		return types.ErrInvalidParam
	}
	//
	reply, err := c.cli.QueryPing(context.Background(), param)
	if err != nil {
		return err
	}
	*result = reply
	return nil
}

// QueryPing
func (c *channelClient) QueryPing(ctx context.Context, queryParam *echotypes.Query) (types.Message, error) {
	return c.Query(echotypes.EchoX, "GetPing", queryParam)
}

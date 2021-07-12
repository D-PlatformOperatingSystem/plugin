package rpc

import (
	rpctypes "github.com/D-PlatformOperatingSystem/dpos/rpc/types"
	exchangetypes "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/exchange/types"
)

/*
 * rpc
 */

//   grpc service
type channelClient struct {
	rpctypes.ChannelClient
}

// Jrpc   json rpc
type Jrpc struct {
	cli *channelClient
}

// Grpc grpc
type Grpc struct {
	*channelClient
}

// Init init rpc
func Init(name string, s rpctypes.RPCServer) {
	cli := &channelClient{}
	grpc := &Grpc{channelClient: cli}
	cli.Init(name, s, &Jrpc{cli: cli}, grpc)
	//  grpc service   grpc server，       pb.go
	exchangetypes.RegisterExchangeServer(s.GRPC(), grpc)
}

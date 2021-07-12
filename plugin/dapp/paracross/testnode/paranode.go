package testnode

import (
	"github.com/D-PlatformOperatingSystem/dpos/types"
	"github.com/D-PlatformOperatingSystem/dpos/util/testnode"
)

/*
1. solo   ，
2.          ：  ，       ，
*/

//ParaNode
type ParaNode struct {
	Main *testnode.DplatformOSMock
	Para *testnode.DplatformOSMock
}

//NewParaNode
func NewParaNode(main *testnode.DplatformOSMock, para *testnode.DplatformOSMock) *ParaNode {
	if main == nil {
		main = testnode.New("", nil)
		main.Listen()
	}
	if para == nil {
		cfg := types.NewDplatformOSConfig(DefaultConfig)
		testnode.ModifyParaClient(cfg, main.GetCfg().RPC.GrpcBindAddr)
		para = testnode.NewWithConfig(cfg, nil)
		para.Listen()
	}
	return &ParaNode{Main: main, Para: para}
}

//Close
func (node *ParaNode) Close() {
	node.Para.Close()
	node.Main.Close()
}

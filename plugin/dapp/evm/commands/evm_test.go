package commands

import (
	"testing"

	"github.com/D-PlatformOperatingSystem/dpos/rpc/jsonclient"
	rpctypes "github.com/D-PlatformOperatingSystem/dpos/rpc/types"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	"github.com/D-PlatformOperatingSystem/dpos/util/testnode"
	"github.com/stretchr/testify/assert"

	//          ，            ，        ，          
	evm "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/evm/executor"
	evmtypes "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/evm/types"

	//           ，         
	"github.com/D-PlatformOperatingSystem/dpos/client/mocks"
	_ "github.com/D-PlatformOperatingSystem/dpos/system"
	"github.com/stretchr/testify/mock"
)

// TestQueryDebug        rpc  
func TestQueryDebug(t *testing.T) {
	var cfg = types.NewDplatformOSConfig(types.GetDefaultCfgstring())
	evm.Init(evmtypes.ExecutorName, cfg, nil)
	var debugReq = evmtypes.EvmDebugReq{Optype: 1}
	js, err := types.PBToJSON(&debugReq)
	assert.Nil(t, err)
	in := &rpctypes.Query4Jrpc{
		Execer:   "evm",
		FuncName: "EvmDebug",
		Payload:  js,
	}

	var mockResp = evmtypes.EvmDebugResp{DebugStatus: "on"}

	mockapi := &mocks.QueueProtocolAPI{}
	//      mock     ,Close    ，        
	mockapi.On("Close").Return()
	mockapi.On("Query", "evm", "EvmDebug", &debugReq).Return(&mockResp, nil)
	mockapi.On("GetConfig", mock.Anything).Return(cfg, nil)

	mockDOM := testnode.New("", mockapi)
	defer mockDOM.Close()
	rpcCfg := mockDOM.GetCfg().RPC
	//           ，       
	rpcCfg.JrpcBindAddr = "127.0.0.1:8899"
	mockDOM.GetRPC().Listen()

	jsonClient, err := jsonclient.NewJSONClient("http://" + rpcCfg.JrpcBindAddr + "/")
	assert.Nil(t, err)
	assert.NotNil(t, jsonClient)

	var debugResp evmtypes.EvmDebugResp
	err = jsonClient.Call("DplatformOS.Query", in, &debugResp)
	assert.Nil(t, err)
	assert.Equal(t, "on", debugResp.DebugStatus)
}

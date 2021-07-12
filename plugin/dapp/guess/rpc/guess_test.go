package rpc

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/D-PlatformOperatingSystem/dpos/blockchain"
	"github.com/D-PlatformOperatingSystem/dpos/common/address"
	"github.com/D-PlatformOperatingSystem/dpos/common/crypto"
	"github.com/D-PlatformOperatingSystem/dpos/common/limits"
	"github.com/D-PlatformOperatingSystem/dpos/common/log"
	"github.com/D-PlatformOperatingSystem/dpos/executor"
	"github.com/D-PlatformOperatingSystem/dpos/mempool"
	"github.com/D-PlatformOperatingSystem/dpos/p2p"
	"github.com/D-PlatformOperatingSystem/dpos/queue"
	"github.com/D-PlatformOperatingSystem/dpos/rpc"
	"github.com/D-PlatformOperatingSystem/dpos/store"
	"github.com/D-PlatformOperatingSystem/dpos/system/consensus/solo"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	"github.com/D-PlatformOperatingSystem/dpos/util"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	_ "github.com/D-PlatformOperatingSystem/dpos/system"
	_ "github.com/D-PlatformOperatingSystem/plugin/plugin/store/init"

	jsonrpc "github.com/D-PlatformOperatingSystem/dpos/rpc/jsonclient"
	rpctypes "github.com/D-PlatformOperatingSystem/dpos/rpc/types"
	cty "github.com/D-PlatformOperatingSystem/dpos/system/dapp/coins/types"
	gty "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/guess/types"
	pty "github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/norm/types"
)

var (
	secp crypto.Crypto

	config = `# Title local，                。                  ，        solo  。
Title="local"
TestNet=true
FixTime=false
[log]
#     ，  debug(dbug)/info/warn/error(eror)/crit
loglevel = "info"
logConsoleLevel = "info"
#      ，    ，                
logFile = "logs/dplatformos.log"
#           （  ： ）
maxFileSize = 300
#              
maxBackups = 100
#            （  ： ）
maxAge = 28
#              （    UTC  ）
localTime = true
#           （     gz）
compress = true
#             
callerFile = false
#         
callerFunction = false
[blockchain]
#        
defCacheSize=128
#                   
maxFetchBlockNum=128
#                 
timeoutSeconds=5
#         
driver="leveldb"
#        
dbPath="datadir"
#        
dbCache=64
#       
singleMode=true
#            ，         ，             false，     
batchsync=false
#                ，         ，          ，     true
isRecordBlockSequence=true
#         
isParaChain=false
#             
enableTxQuickIndex=false
[p2p]
types=["dht"]
enable=true
#         
driver="leveldb"
#        
dbPath="datadir/addrbook"
#        
dbCache=4
# GRPC      
grpcLogFile="grpc.log"

[rpc]
# jrpc    
jrpcBindAddr="localhost:9801"
# grpc    
grpcBindAddr="localhost:9802"
#      ，     IP  ，   “*”，    IP  
whitelist=["127.0.0.1"]
# jrpc       ，   “*”，      RPC  
jrpcFuncWhitelist=["*"]
# jrpc       ，           rpc  ，          ，    
# jrpcFuncBlacklist=["xxxx"]
# grpc       ，   “*”，      RPC  
grpcFuncWhitelist=["*"]
# grpc       ，           rpc  ，          ，    
# grpcFuncBlacklist=["xxx"]
#     https
enableTLS=false
#     ，          cli    
certFile="cert.pem"
#     
keyFile="key.pem"
[mempool]
# mempool    ，  ，timeline，score，price
name="timeline"
# mempool      ，  10240
poolCacheSize=10240
#          ，       ，  ，   100000
minTxFeeRate=100000
#      mempool        ，  100
maxTxNumPerAccount=10000
# timeline               
[mempool.sub.timeline]
# mempool      ，  10240
poolCacheSize=10240
#          ，       ，  ，   100000
minTxFeeRate=100000
#      mempool        ，  100
maxTxNumPerAccount=10000
# score       (  =  a*   /     -  b*  *  c，     ，    ，  a，b   c   )，      
[mempool.sub.score]
# mempool      ，  10240
poolCacheSize=10240
#          ，       ，  ，   100000
minTxFeeRate=100000
#      mempool        ，  100
maxTxNumPerAccount=10000
#        
timeParam=1
#                 ，   unix       ，       1e-5   ~= 1s   
priceConstant=1544
#     
pricePower=1
# price       (  =   /     ，      ，        )
[mempool.sub.price]
# mempool      ，  10240
poolCacheSize=10240
#          ，       ，  ，   100000
minTxFeeRate=100000
#      mempool        ，  100
maxTxNumPerAccount=10000
[consensus]
#   ,    solo,ticket,raft,tendermint,para
name="solo"
#      ,          
minerstart=true
#      (UTC  )
genesisBlockTime=1514533394
#      
genesis="16ERTbYtKKQ64wMthAY9J4La4nAiidG45A"
[mver.consensus]
#      
fundKeyAddr = "1CQXE6TxaYCG5mADtWij4AxhZCUTpoABb3"
#    
coinReward = 18
#      
coinDevFund = 12
#ticket  
ticketPrice = 10000
#    
powLimitBits = "0x1f00ffff"
#            ，      4   ，    (1/4 - 4)，       4      ，            1/4 ，    ，                
retargetAdjustmentFactor = 4
#               16s ，             。
futureBlockTime = 16
#ticket    
ticketFrozenTime = 5    #5s only for test
ticketWithdrawTime = 10 #10s only for test
ticketMinerWaitTime = 2 #2s only for test
#         
maxTxNumber = 1600      #160
#         ，(ps:            ，     targetTimespan / targetTimePerBlock      )
targetTimespan = 2304
#           
targetTimePerBlock = 16
#       ，  consensus         
[consensus.sub.solo]
#      
genesis="16ABEbYtKKQm5wMuySK9J4La5nAiidGuyt"
#      (UTC  )
genesisBlockTime=1514533394
#        ,    
waitTxMs=10
[store]
#         ，    mavl,kvdb,kvmvcc,mpt
name="mavl"
#         ，    leveldb,goleveldb,memdb,gobadgerdb,ssdb,pegasus
driver="leveldb"
#         
dbPath="datadir/mavltree"
# Cache  
dbCache=128
#      
localdbVersion="1.0.0"
[store.sub.mavl]
#     mavl   
enableMavlPrefix=false
#     MVCC,  mavl enableMVCC true     true
enableMVCC=false
#     mavl    
enableMavlPrune=false
#       
pruneHeight=10000
[wallet]
#          ，  0.00000001DOM(1e-8),  100000， 0.001DOM
minFee=100000
# walletdb   ，  leveldb/memdb/gobadgerdb/ssdb/pegasus
driver="leveldb"
# walletdb  
dbPath="wallet"
# walletdb    
dbCache=16
#           
signType="secp256k1"
[wallet.sub.ticket]
#     ticket    ，  false
minerdisable=false
#     ticket        ，    “*”，        
minerwhitelist=["*"]
[exec]
#    stat  
enableStat=false
#    MVCC  
enableMVCC=false
alias=["token1:token","token2:token","token3:token"]
[exec.sub.token]
#    token    
saveTokenTxList=true
#token     
tokenApprs = [
    "1Bsg9j6gW83sShoee1fZAt9TkUjcrCgA9S",
    "1Q8hGLfoGe63efeWa8fJ4Pnukhkngt6poK",
    "1LY8GFia5EiyoTodMLfkB5PHNNpXRqxhyB",
    "1GCzJDS6HbgTQ2emade7mEJGGWFfA15pS9",
    "1JYB8sxi4He5pZWHCd3Zi2nypQ4JMB6AxN",
    "12oupcayRT7LvaC4qW4avxsTE7U41cKSio",
]
[exec.sub.cert]
#            
enable=false
#       
cryptoPath="authdir/crypto"
#        ，  "auth_ecdsa", "auth_sm2"
signType="auth_ecdsa"
[exec.sub.relay]
#relay     BTC       
genesis="16ABEbYtKKQm5wMuySK9J4La5nAiidGuyt"
[exec.sub.manage]
#manage          
superManager=[
    "1Bsg9j6gW83sShoee1fZAt9TkUjcrCgA9S", 
    "12oupcayRT7LvaC4qW4avxsTE7U41cKSio", 
    "1Q8hGLfoGe63efeWa8fJ4Pnukhkngt6poK"
]
`
)

var (
	random *rand.Rand

	loopCount = 1
	conn      *grpc.ClientConn
	c         types.DplatformOSClient
	adminPriv = "CC38546E9E659D15E6B4893F0AB32A06D103931A8230B0BDE71459D2B27D6944"
	adminAddr = "16ABEbYtKKQm5wMuySK9J4La5nAiidGuyt"

	//userAPubkey = "03EF0E1D3112CF571743A3318125EDE2E52A4EB904BCBAA4B1F75020C2846A7EB4"
	userAAddr = "15LsTP6tkYGZcN7tc1Xo2iYifQfowxot3b"
	userAPriv = "5A6A14DA6F5A42835E529D75D87CC8904544F59EEE5387A37D87EEAD194D7EB2"

	//userBPubkey = "027848E7FA630B759DB406940B5506B666A344B1060794BBF314EB459D40881BB3"
	userBAddr = "14w5JaGDrXTZwF5Wv51UAtuGgAupenLAok"
	userBPriv = "754F53FCEA0CB1F528918726A49B3551B7F1284D802A1D6AAF4522E8A8DA1B5B"
)

const fee = 1e6

func init() {
	err := limits.SetLimits()
	if err != nil {
		panic(err)
	}
	log.SetLogLevel("info")
	random = rand.New(rand.NewSource(types.Now().UnixNano()))

	cr2, err := crypto.New(types.GetSignName("", types.SECP256K1))
	if err != nil {
		fmt.Println("crypto.New failed for types.ED25519")
		return
	}
	secp = cr2
}

func Init() {
	fmt.Println("=======Init Data1!=======")
	os.RemoveAll("datadir")
	os.RemoveAll("wallet")
	os.Remove("dplatformos.test.toml")

	ioutil.WriteFile("dplatformos.test.toml", []byte(config), 0664)
}

func clearTestData() {
	fmt.Println("=======start clear test data!=======")

	os.Remove("dplatformos.test.toml")
	os.RemoveAll("wallet")
	err := os.RemoveAll("datadir")
	if err != nil {
		fmt.Println("delete datadir have a err:", err.Error())
	}

	fmt.Println("test data clear successfully!")
}

func TestGuess(t *testing.T) {
	Init()
	testGuessImp(t)
	fmt.Println("=======start clear test data!=======")
	clearTestData()
}

func testGuessImp(t *testing.T) {
	fmt.Println("=======start guess test!=======")
	q, chain, s, mem, exec, cs, p2p := initEnvGuess()
	cfg := q.GetConfig()
	defer chain.Close()
	defer mem.Close()
	defer exec.Close()
	defer s.Close()
	defer q.Close()
	defer cs.Close()
	defer p2p.Close()
	err := createConn()
	for err != nil {
		err = createConn()
	}
	time.Sleep(2 * time.Second)
	fmt.Println("=======start NormPut!=======")

	for i := 0; i < loopCount; i++ {
		NormPut()
		time.Sleep(time.Second)
	}

	fmt.Println("=======start sendTransferTx!=======")
	//          A
	sendTransferTx(cfg, adminPriv, userAAddr, 2000000000000)
	sendTransferTx(cfg, adminPriv, userBAddr, 2000000000000)

	time.Sleep(2 * time.Second)
	in := &types.ReqBalance{}
	in.Addresses = append(in.Addresses, userAAddr)
	acct, err1 := c.GetBalance(context.Background(), in)
	if err1 != nil || len(acct.Acc) == 0 {
		fmt.Println("no balance for ", userAAddr)
	} else {
		fmt.Println(userAAddr, " balance:", acct.Acc[0].Balance, "frozen:", acct.Acc[0].Frozen)
	}
	assert.Equal(t, true, acct.Acc[0].Balance == 2000000000000)

	in2 := &types.ReqBalance{}
	in2.Addresses = append(in.Addresses, userBAddr)
	acct2, err2 := c.GetBalance(context.Background(), in2)
	if err2 != nil || len(acct2.Acc) == 0 {
		fmt.Println("no balance for ", userBAddr)
	} else {
		fmt.Println(userBAddr, " balance:", acct2.Acc[0].Balance, "frozen:", acct2.Acc[0].Frozen)
	}
	assert.Equal(t, true, acct2.Acc[0].Balance == 2000000000000)

	fmt.Println("=======start sendTransferToExecTx!=======")
	//      dos
	sendTransferToExecTx(cfg, userAPriv, "guess", 1000000000000)
	sendTransferToExecTx(cfg, userBPriv, "guess", 1000000000000)
	time.Sleep(2 * time.Second)

	fmt.Println("=======start GetBalance!=======")

	in3 := &types.ReqBalance{}
	in3.Addresses = append(in3.Addresses, userAAddr)
	acct3, err3 := c.GetBalance(context.Background(), in3)
	if err3 != nil || len(acct3.Acc) == 0 {
		fmt.Println("no balance for ", userAAddr)
	} else {
		fmt.Println(userAAddr, " balance:", acct3.Acc[0].Balance, "frozen:", acct3.Acc[0].Frozen)
	}
	assert.Equal(t, true, acct3.Acc[0].Balance == 1000000000000)

	in4 := &types.ReqBalance{}
	in4.Addresses = append(in4.Addresses, userBAddr)
	acct4, err4 := c.GetBalance(context.Background(), in4)
	if err4 != nil || len(acct4.Acc) == 0 {
		fmt.Println("no balance for ", userBAddr)
	} else {
		fmt.Println(userBAddr, " balance:", acct4.Acc[0].Balance, "frozen:", acct4.Acc[0].Frozen)
	}
	assert.Equal(t, true, acct4.Acc[0].Balance == 1000000000000)

	fmt.Println("=======start sendGuessStartTx!=======")
	ok, gameid := sendGuessStartTx(cfg, "WorldCup Final", "A:France;B:Claodia", "football", adminPriv)
	if !ok {
		panic("Guess start failed.")
	} else {
		fmt.Println("txid: ", hex.EncodeToString(gameid))
	}
	time.Sleep(2 * time.Second)

	strGameID1 := "0x" + hex.EncodeToString(gameid)

	reply := queryGuessByIds(strGameID1)
	assert.Equal(t, true, reply.Games[0].Status == gty.GuessGameStatusStart)

	fmt.Println("=======start sendGuessBetTx!=======")
	ok, txid := sendGuessBetTx(cfg, strGameID1, "A", 5e8, userAPriv)
	if !ok {
		panic("Guess A bet failed.")
	} else {
		fmt.Println("Guess A bet txid: ", hex.EncodeToString(txid))
	}
	time.Sleep(2 * time.Second)
	reply = queryGuessByIds(strGameID1)
	assert.Equal(t, true, reply.Games[0].Status == gty.GuessGameStatusBet && reply.Games[0].BetStat.TotalBetTimes == 1)

	ok, txid = sendGuessBetTx(cfg, strGameID1, "B", 5e8, userBPriv)
	if !ok {
		panic("Guess B bet failed.")
	} else {
		fmt.Println("Guess B bet txid: ", hex.EncodeToString(txid))
	}
	time.Sleep(2 * time.Second)
	reply = queryGuessByIds(strGameID1)
	assert.Equal(t, true, reply.Games[0].Status == gty.GuessGameStatusBet && reply.Games[0].BetStat.TotalBetTimes == 2)

	fmt.Println("=======start sendGuessStopTx failed!=======")
	ok, txid = sendGuessStopTx(cfg, strGameID1, userBPriv)
	if !ok {
		panic("Guess stop failed,only admin can stop.")
	} else {
		fmt.Println("Guess stop txid: ", hex.EncodeToString(txid))
	}
	time.Sleep(2 * time.Second)
	reply = queryGuessByIds(strGameID1)
	assert.Equal(t, true, reply.Games[0].Status == gty.GuessGameStatusBet && reply.Games[0].BetStat.TotalBetTimes == 2)

	fmt.Println("=======start sendGuessStopTx!=======")
	ok, txid = sendGuessStopTx(cfg, strGameID1, adminPriv)
	if !ok {
		panic("Guess stop failed.")
	} else {
		fmt.Println("Guess stop txid: ", hex.EncodeToString(txid))
	}
	time.Sleep(2 * time.Second)
	reply = queryGuessByIds(strGameID1)
	assert.Equal(t, true, reply.Games[0].Status == gty.GuessGameStatusStopBet && reply.Games[0].BetStat.TotalBetTimes == 2)

	fmt.Println("=======start sendGuessBetTx failed!=======")
	ok, txid = sendGuessBetTx(cfg, strGameID1, "A", 5e8, userAPriv)
	if !ok {
		fmt.Println("Guess stopped, bet failed.")
	} else {
		fmt.Println(fmt.Sprintf("Guess A bet txid: %s", hex.EncodeToString(txid)))
	}
	time.Sleep(2 * time.Second)
	reply = queryGuessByIds(strGameID1)
	assert.Equal(t, true, reply.Games[0].Status == gty.GuessGameStatusStopBet && reply.Games[0].BetStat.TotalBetTimes == 2)

	fmt.Println("=======start sendGuessPublishTx failed!=======")
	ok, txid = sendGuessPublishTx(cfg, strGameID1, "A", userAPriv)
	if !ok {
		fmt.Println("sendGuessPublishTx failed,only admin can publish.")
	} else {
		fmt.Println(fmt.Sprintf("publish ok, but it's not correct, txid: %s", hex.EncodeToString(txid)))
	}
	time.Sleep(2 * time.Second)
	reply = queryGuessByIds(strGameID1)
	assert.Equal(t, true, reply.Games[0].Status == gty.GuessGameStatusStopBet && reply.Games[0].BetStat.TotalBetTimes == 2)

	fmt.Println("=======start sendGuessPublishTx!=======")
	ok, txid = sendGuessPublishTx(cfg, strGameID1, "A", adminPriv)
	if !ok {
		fmt.Println("sendGuessPublishTx failed.")
	} else {
		fmt.Println(fmt.Sprintf("publish ok, txid: %s", hex.EncodeToString(txid)))
	}
	time.Sleep(2 * time.Second)
	reply = queryGuessByIds(strGameID1)
	assert.Equal(t, true, reply.Games[0].Status == gty.GuessGameStatusPublish && reply.Games[0].BetStat.TotalBetTimes == 2)

	fmt.Println("=======start sendGuessAbortTx!=======")
	ok, txid = sendGuessAbortTx(cfg, strGameID1, adminPriv)
	if !ok {
		fmt.Println("Guess abort failed, already published.")
	} else {
		fmt.Println(fmt.Sprintf("Guess abort txid: %s", hex.EncodeToString(txid)))
	}
	time.Sleep(2 * time.Second)
	reply = queryGuessByIds(strGameID1)
	assert.Equal(t, true, reply.Games[0].Status == gty.GuessGameStatusPublish && reply.Games[0].BetStat.TotalBetTimes == 2)

	//    ，      :start->abort->stop
	fmt.Println("=======start sendGuessStartTx!=======")
	ok, gameid = sendGuessStartTx(cfg, "WorldCup Final", "A:France;B:Claodia", "football", adminPriv)
	if !ok {
		fmt.Println("Guess start failed.")
	} else {
		fmt.Println("txid: ", hex.EncodeToString(gameid))
	}
	time.Sleep(2 * time.Second)

	strGameID2 := "0x" + hex.EncodeToString(gameid)

	reply = queryGuessByIds(strGameID2)
	assert.Equal(t, true, reply.Games[0].Status == gty.GuessGameStatusStart)

	fmt.Println("=======start sendGuessAbortTx!=======")
	ok, txid = sendGuessAbortTx(cfg, strGameID2, adminPriv)
	if !ok {
		fmt.Println("Guess abort failed.")
	} else {
		fmt.Println(fmt.Sprintf("Guess abort txid: %s", hex.EncodeToString(txid)))
	}
	time.Sleep(2 * time.Second)
	reply = queryGuessByIds(strGameID2)
	assert.Equal(t, true, reply.Games[0].Status == gty.GuessGameStatusAbort)

	fmt.Println("=======start sendGuessStopTx failed!=======")
	ok, txid = sendGuessStopTx(cfg, strGameID2, adminPriv)
	if !ok {
		fmt.Println("Guess stop failed,it's already aborted.")
	} else {
		fmt.Println(fmt.Sprintf("Guess stop txid: %s", hex.EncodeToString(txid)))
	}
	time.Sleep(2 * time.Second)
	reply = queryGuessByIds(strGameID2)
	assert.Equal(t, true, reply.Games[0].Status == gty.GuessGameStatusAbort)

	//    ，    :start->stop->abort
	fmt.Println("=======start sendGuessStartTx!=======")
	ok, gameid = sendGuessStartTx(cfg, "WorldCup Final", "A:France;B:Claodia", "football", adminPriv)
	if !ok {
		fmt.Println("Guess start failed.")
	} else {
		fmt.Println("txid: ", hex.EncodeToString(gameid))
	}
	time.Sleep(2 * time.Second)

	strGameID3 := "0x" + hex.EncodeToString(gameid)
	reply = queryGuessByIds(strGameID3)
	assert.Equal(t, true, reply.Games[0].Status == gty.GuessGameStatusStart)

	fmt.Println("=======start sendGuessStopTx!=======")
	ok, txid = sendGuessStopTx(cfg, strGameID3, adminPriv)
	if !ok {
		fmt.Println("Guess stop failed")
	} else {
		fmt.Println("Guess stop txid: ", hex.EncodeToString(txid))
	}
	time.Sleep(2 * time.Second)
	reply = queryGuessByIds(strGameID3)
	assert.Equal(t, true, reply.Games[0].Status == gty.GuessGameStatusStopBet)

	fmt.Println("=======start sendGuessAbortTx!=======")
	ok, txid = sendGuessAbortTx(cfg, strGameID3, adminPriv)
	if !ok {
		fmt.Println("Guess abort failed.")
	} else {
		fmt.Println(fmt.Sprintf("Guess abort txid: %s", hex.EncodeToString(txid)))
	}
	time.Sleep(2 * time.Second)

	//
	fmt.Println("=======start queryGuessByIds!=======")
	reply = queryGuessByIds(strGameID3)
	assert.Equal(t, true, reply.Games[0].Status == gty.GuessGameStatusAbort)

	fmt.Println("=======start queryGuessByID!=======")
	reply2 := queryGuessByID(strGameID1)
	assert.Equal(t, true, reply2.Game.Status == gty.GuessGameStatusPublish && reply2.Game.BetStat.TotalBetTimes == 2)

	fmt.Println("=======start queryGuessByAddr!=======")
	record := queryGuessByAddr(userAAddr)
	assert.Equal(t, true, record.Records[0].GameID == strGameID1)

	fmt.Println("=======start queryGuessByStatus!=======")

	record = queryGuessByStatus(gty.GuessGameStatusPublish)
	assert.Equal(t, true, record.Records[0].GameID == strGameID1)

	record = queryGuessByStatus(gty.GuessGameStatusAbort)
	assert.Equal(t, true, len(record.Records) == 2)

	fmt.Println("=======start queryGuessByAdminAddr!=======")
	record = queryGuessByAdminAddr(adminAddr)
	assert.Equal(t, true, len(record.Records) == 3)

	fmt.Println("=======start queryGuessByAddrStatus!=======")
	record = queryGuessByAddrStatus(userBAddr, gty.GuessGameStatusPublish)
	assert.Equal(t, true, len(record.Records) == 1 && record.Records[0].GameID == strGameID1)
	record = queryGuessByAddrStatus(userBAddr, 10)
	assert.Equal(t, true, len(record.Records) == 0)

	fmt.Println("=======start queryGuessByAdminAddrStatus!=======")
	record = queryGuessByAdminAddrStatus(adminAddr, gty.GuessGameStatusAbort)
	assert.Equal(t, true, len(record.Records) == 2)

	record = queryGuessByAdminAddrStatus(adminAddr, gty.GuessGameStatusPublish)
	assert.Equal(t, true, len(record.Records) == 1)

	fmt.Println("=======start queryGuessByCategoryStatus!=======")
	record = queryGuessByCategoryStatus("football", gty.GuessGameStatusPublish)
	assert.Equal(t, true, len(record.Records) == 1)

	record = queryGuessByCategoryStatus("football", gty.GuessGameStatusAbort)
	assert.Equal(t, true, len(record.Records) == 2)

	time.Sleep(2 * time.Second)
}

func initEnvGuess() (queue.Queue, *blockchain.BlockChain, queue.Module, queue.Module, *executor.Executor, queue.Module, queue.Module) {
	flag.Parse()
	dplatformosCfg := types.NewDplatformOSConfig(types.ReadFile("dplatformos.test.toml"))
	var q = queue.New("channel")
	q.SetConfig(dplatformosCfg)
	cfg := dplatformosCfg.GetModuleConfig()
	sub := dplatformosCfg.GetSubConfig()
	chain := blockchain.New(dplatformosCfg)
	chain.SetQueueClient(q.Client())

	exec := executor.New(dplatformosCfg)
	exec.SetQueueClient(q.Client())
	dplatformosCfg.SetMinFee(0)
	s := store.New(dplatformosCfg)
	s.SetQueueClient(q.Client())

	cs := solo.New(cfg.Consensus, sub.Consensus["solo"])
	cs.SetQueueClient(q.Client())

	mem := mempool.New(dplatformosCfg)
	mem.SetQueueClient(q.Client())
	network := p2p.NewP2PMgr(dplatformosCfg)

	network.SetQueueClient(q.Client())

	rpc.InitCfg(cfg.RPC)
	gapi := rpc.NewGRpcServer(q.Client(), nil)
	go gapi.Listen()

	japi := rpc.NewJSONRPCServer(q.Client(), nil)
	go japi.Listen()

	return q, chain, s, mem, exec, cs, network
}

func createConn() error {
	var err error
	url := "127.0.0.1:9802"
	fmt.Println("grpc url:", url)
	conn, err = grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	c = types.NewDplatformOSClient(conn)
	return nil
}

func generateKey(i, valI int) string {
	key := make([]byte, valI)
	binary.PutUvarint(key[:10], uint64(valI))
	binary.PutUvarint(key[12:24], uint64(i))
	if _, err := rand.Read(key[24:]); err != nil {
		os.Exit(1)
	}
	return string(key)
}

func generateValue(i, valI int) string {
	value := make([]byte, valI)
	binary.PutUvarint(value[:16], uint64(i))
	binary.PutUvarint(value[32:128], uint64(i))
	if _, err := rand.Read(value[128:]); err != nil {
		os.Exit(1)
	}
	return string(value)
}

func getprivkey(key string) crypto.PrivKey {
	bkey, err := hex.DecodeString(key)
	if err != nil {
		panic(err)
	}
	priv, err := secp.PrivKeyFromBytes(bkey)
	if err != nil {
		panic(err)
	}
	return priv
}

func prepareTxList() *types.Transaction {
	var key string
	var value string
	var i int

	key = generateKey(i, 32)
	value = generateValue(i, 180)

	nput := &pty.NormAction_Nput{Nput: &pty.NormPut{Key: []byte(key), Value: []byte(value)}}
	action := &pty.NormAction{Value: nput, Ty: pty.NormActionPut}
	tx := &types.Transaction{Execer: []byte("norm"), Payload: types.Encode(action), Fee: fee}
	tx.To = address.ExecAddress("norm")
	tx.Nonce = random.Int63()
	tx.Sign(types.SECP256K1, getprivkey("CC38546E9E659D15E6B4893F0AB32A06D103931A8230B0BDE71459D2B27D6944"))
	return tx
}

func NormPut() {
	tx := prepareTxList()

	reply, err := c.SendTransaction(context.Background(), tx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	if !reply.IsOk {
		fmt.Fprintln(os.Stderr, errors.New(string(reply.GetMsg())))
		return
	}
}

func sendTransferTx(cfg *types.DplatformOSConfig, fromKey, to string, amount int64) bool {
	signer := util.HexToPrivkey(fromKey)
	var tx *types.Transaction
	transfer := &cty.CoinsAction{}
	v := &cty.CoinsAction_Transfer{Transfer: &types.AssetsTransfer{Amount: amount, Note: []byte(""), To: to}}
	transfer.Value = v
	transfer.Ty = cty.CoinsActionTransfer
	execer := []byte("coins")
	tx = &types.Transaction{Execer: execer, Payload: types.Encode(transfer), To: to, Fee: fee}
	tx, err := types.FormatTx(cfg, string(execer), tx)
	if err != nil {
		fmt.Println("in sendTransferTx formatTx failed")
		return false
	}

	tx.Sign(types.SECP256K1, signer)
	reply, err := c.SendTransaction(context.Background(), tx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Println("in sendTransferTx SendTransaction failed")

		return false
	}
	if !reply.IsOk {
		fmt.Fprintln(os.Stderr, errors.New(string(reply.GetMsg())))
		fmt.Println("in sendTransferTx SendTransaction failed,reply not ok.")

		return false
	}
	fmt.Println("sendTransferTx ok")

	return true
}

func sendTransferToExecTx(cfg *types.DplatformOSConfig, fromKey, execName string, amount int64) bool {
	signer := util.HexToPrivkey(fromKey)
	var tx *types.Transaction
	transfer := &cty.CoinsAction{}
	execAddr := address.ExecAddress(execName)
	v := &cty.CoinsAction_TransferToExec{TransferToExec: &types.AssetsTransferToExec{Amount: amount, Note: []byte(""), ExecName: execName, To: execAddr}}
	transfer.Value = v
	transfer.Ty = cty.CoinsActionTransferToExec
	execer := []byte("coins")
	tx = &types.Transaction{Execer: execer, Payload: types.Encode(transfer), To: address.ExecAddress("guess"), Fee: fee}
	tx, err := types.FormatTx(cfg, string(execer), tx)
	if err != nil {
		fmt.Println("sendTransferToExecTx formatTx failed.")

		return false
	}

	tx.Sign(types.SECP256K1, signer)
	reply, err := c.SendTransaction(context.Background(), tx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Println("in sendTransferToExecTx SendTransaction failed")

		return false
	}
	if !reply.IsOk {
		fmt.Fprintln(os.Stderr, errors.New(string(reply.GetMsg())))
		fmt.Println("in sendTransferToExecTx SendTransaction failed,reply not ok.")

		return false
	}

	fmt.Println("sendTransferToExecTx ok")

	return true
}

func sendGuessStartTx(cfg *types.DplatformOSConfig, topic, option, category, privKey string) (bool, []byte) {
	signer := util.HexToPrivkey(privKey)
	var tx *types.Transaction
	action := &gty.GuessGameAction{}

	v := &gty.GuessGameAction_Start{
		Start: &gty.GuessGameStart{
			Topic:          topic,
			Options:        option,
			Category:       category,
			MaxBetsOneTime: 100e8,
			MaxBetsNumber:  1000e8,
			DevFeeFactor:   5,
			DevFeeAddr:     "1D6RFZNp2rh6QdbcZ1d7RWuBUz61We6SD7",
			PlatFeeFactor:  5,
			PlatFeeAddr:    "1PHtChNt3UcfssR7v7trKSk3WJtAWjKjjX",
		},
	}

	action.Value = v
	action.Ty = gty.GuessGameActionStart
	execer := []byte("guess")
	tx = &types.Transaction{Execer: execer, Payload: types.Encode(action), To: address.ExecAddress(string(execer)), Fee: fee}
	tx, err := types.FormatTx(cfg, string(execer), tx)
	if err != nil {
		fmt.Println("sendGuessStartTx formatTx failed.")

		return false, nil
	}

	tx.Sign(types.SECP256K1, signer)
	reply, err := c.SendTransaction(context.Background(), tx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Println("in sendGuessStartTx SendTransaction failed")

		return false, nil
	}

	if !reply.IsOk {
		fmt.Fprintln(os.Stderr, errors.New(string(reply.GetMsg())))
		fmt.Println("in sendGuessStartTx SendTransaction failed,reply not ok.")

		return false, nil
	}

	fmt.Println("sendGuessStartTx ok")

	return true, reply.Msg
}

func sendGuessBetTx(cfg *types.DplatformOSConfig, gameID, option string, betsNum int64, privKey string) (bool, []byte) {
	signer := util.HexToPrivkey(privKey)
	var tx *types.Transaction
	action := &gty.GuessGameAction{}

	v := &gty.GuessGameAction_Bet{
		Bet: &gty.GuessGameBet{
			GameID:  gameID,
			Option:  option,
			BetsNum: betsNum,
		},
	}

	action.Value = v
	action.Ty = gty.GuessGameActionBet
	execer := []byte("guess")
	tx = &types.Transaction{Execer: execer, Payload: types.Encode(action), To: address.ExecAddress(string(execer)), Fee: fee}
	tx, err := types.FormatTx(cfg, string(execer), tx)
	if err != nil {
		fmt.Println("sendGuessBetTx formatTx failed.")

		return false, nil
	}

	tx.Sign(types.SECP256K1, signer)
	reply, err := c.SendTransaction(context.Background(), tx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Println("in sendGuessBetTx SendTransaction failed")

		return false, nil
	}

	if !reply.IsOk {
		fmt.Fprintln(os.Stderr, errors.New(string(reply.GetMsg())))
		fmt.Println("in sendGuessBetTx SendTransaction failed,reply not ok.")

		return false, nil
	}

	fmt.Println("sendGuessBetTx ok")

	return true, reply.Msg
}

func sendGuessStopTx(cfg *types.DplatformOSConfig, gameID, privKey string) (bool, []byte) {
	signer := util.HexToPrivkey(privKey)
	var tx *types.Transaction
	action := &gty.GuessGameAction{}

	v := &gty.GuessGameAction_StopBet{
		StopBet: &gty.GuessGameStopBet{
			GameID: gameID,
		},
	}

	action.Value = v
	action.Ty = gty.GuessGameActionStopBet
	execer := []byte("guess")
	tx = &types.Transaction{Execer: execer, Payload: types.Encode(action), To: address.ExecAddress(string(execer)), Fee: fee}
	tx, err := types.FormatTx(cfg, string(execer), tx)
	if err != nil {
		fmt.Println("sendGuessStopTx formatTx failed.")

		return false, nil
	}

	tx.Sign(types.SECP256K1, signer)
	reply, err := c.SendTransaction(context.Background(), tx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Println("in sendGuessStopTx SendTransaction failed")

		return false, nil
	}

	if !reply.IsOk {
		fmt.Fprintln(os.Stderr, errors.New(string(reply.GetMsg())))
		fmt.Println("in sendGuessStopTx SendTransaction failed,reply not ok.")

		return false, nil
	}

	fmt.Println("sendGuessStopTx ok")

	return true, reply.Msg
}

func sendGuessAbortTx(cfg *types.DplatformOSConfig, gameID, privKey string) (bool, []byte) {
	signer := util.HexToPrivkey(privKey)
	var tx *types.Transaction
	action := &gty.GuessGameAction{}

	v := &gty.GuessGameAction_Abort{
		Abort: &gty.GuessGameAbort{
			GameID: gameID,
		},
	}

	action.Value = v
	action.Ty = gty.GuessGameActionAbort
	execer := []byte("guess")
	tx = &types.Transaction{Execer: execer, Payload: types.Encode(action), To: address.ExecAddress(string(execer)), Fee: fee}
	tx, err := types.FormatTx(cfg, string(execer), tx)
	if err != nil {
		fmt.Println("sendGuessAbortTx formatTx failed.")

		return false, nil
	}

	tx.Sign(types.SECP256K1, signer)
	reply, err := c.SendTransaction(context.Background(), tx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Println("in sendGuessAbortTx SendTransaction failed")

		return false, nil
	}

	if !reply.IsOk {
		fmt.Fprintln(os.Stderr, errors.New(string(reply.GetMsg())))
		fmt.Println("in sendGuessAbortTx SendTransaction failed,reply not ok.")

		return false, nil
	}

	fmt.Println("sendGuessAbortTx ok")

	return true, reply.Msg
}

func sendGuessPublishTx(cfg *types.DplatformOSConfig, gameID, result, privKey string) (bool, []byte) {
	signer := util.HexToPrivkey(privKey)
	var tx *types.Transaction
	action := &gty.GuessGameAction{}

	v := &gty.GuessGameAction_Publish{
		Publish: &gty.GuessGamePublish{
			GameID: gameID,
			Result: result,
		},
	}

	action.Value = v
	action.Ty = gty.GuessGameActionPublish
	execer := []byte("guess")
	tx = &types.Transaction{Execer: execer, Payload: types.Encode(action), To: address.ExecAddress(string(execer)), Fee: fee}
	tx, err := types.FormatTx(cfg, string(execer), tx)
	if err != nil {
		fmt.Println("sendGuessPublishTx formatTx failed.")

		return false, nil
	}

	tx.Sign(types.SECP256K1, signer)
	reply, err := c.SendTransaction(context.Background(), tx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Println("in sendGuessPublishTx SendTransaction failed")

		return false, nil
	}

	if !reply.IsOk {
		fmt.Fprintln(os.Stderr, errors.New(string(reply.GetMsg())))
		fmt.Println("in sendGuessPublishTx SendTransaction failed,reply not ok.")

		return false, nil
	}

	fmt.Println("sendGuessPublishTx ok")

	return true, reply.Msg
}

func queryGuessByIds(gameIDs string) *gty.ReplyGuessGameInfos {
	var params rpctypes.Query4Jrpc
	params.Execer = gty.GuessX

	gameIds := strings.Split(gameIDs, ";")
	req := &gty.QueryGuessGameInfos{
		GameIDs: gameIds,
	}
	params.FuncName = gty.FuncNameQueryGamesByIDs
	params.Payload = types.MustPBToJSON(req)
	var res gty.ReplyGuessGameInfos
	//ctx := jsonrpc.NewRPCCtx("http://"+types.Conf("config.rpc").GStr("jrpcBindAddr"), "DplatformOS.Query", params, &res)
	ctx := jsonrpc.NewRPCCtx("http://127.0.0.1:9801", "DplatformOS.Query", params, &res)
	ctx.Run()
	return &res
}

func queryGuessByID(gameID string) *gty.ReplyGuessGameInfo {
	var params rpctypes.Query4Jrpc
	params.Execer = gty.GuessX

	req := &gty.QueryGuessGameInfo{
		GameID: gameID,
	}
	params.FuncName = gty.FuncNameQueryGameByID
	params.Payload = types.MustPBToJSON(req)
	var res gty.ReplyGuessGameInfo
	//ctx := jsonrpc.NewRPCCtx("http://"+types.Conf("config.rpc").GStr("jrpcBindAddr"), "DplatformOS.Query", params, &res)
	ctx := jsonrpc.NewRPCCtx("http://127.0.0.1:9801", "DplatformOS.Query", params, &res)
	ctx.Run()
	return &res
}

func queryGuessByAddr(addr string) *gty.GuessGameRecords {
	var params rpctypes.Query4Jrpc
	params.Execer = gty.GuessX

	req := &gty.QueryGuessGameInfo{
		Addr: addr,
	}
	params.FuncName = gty.FuncNameQueryGameByAddr
	params.Payload = types.MustPBToJSON(req)
	var res gty.GuessGameRecords
	//ctx := jsonrpc.NewRPCCtx("http://"+types.Conf("config.rpc").GStr("jrpcBindAddr"), "DplatformOS.Query", params, &res)
	ctx := jsonrpc.NewRPCCtx("http://127.0.0.1:9801", "DplatformOS.Query", params, &res)
	ctx.Run()
	return &res
}

func queryGuessByStatus(status int32) *gty.GuessGameRecords {
	var params rpctypes.Query4Jrpc
	params.Execer = gty.GuessX

	req := &gty.QueryGuessGameInfo{
		Status: status,
	}
	params.FuncName = gty.FuncNameQueryGameByStatus
	params.Payload = types.MustPBToJSON(req)
	var res gty.GuessGameRecords
	//ctx := jsonrpc.NewRPCCtx("http://"+types.Conf("config.rpc").GStr("jrpcBindAddr"), "DplatformOS.Query", params, &res)
	ctx := jsonrpc.NewRPCCtx("http://127.0.0.1:9801", "DplatformOS.Query", params, &res)
	ctx.Run()
	return &res
}

func queryGuessByAdminAddr(addr string) *gty.GuessGameRecords {
	var params rpctypes.Query4Jrpc
	params.Execer = gty.GuessX

	req := &gty.QueryGuessGameInfo{
		AdminAddr: addr,
	}
	params.FuncName = gty.FuncNameQueryGameByAdminAddr
	params.Payload = types.MustPBToJSON(req)
	var res gty.GuessGameRecords
	//ctx := jsonrpc.NewRPCCtx("http://"+types.Conf("config.rpc").GStr("jrpcBindAddr"), "DplatformOS.Query", params, &res)
	ctx := jsonrpc.NewRPCCtx("http://127.0.0.1:9801", "DplatformOS.Query", params, &res)
	ctx.Run()
	return &res
}

func queryGuessByAddrStatus(addr string, status int32) *gty.GuessGameRecords {
	var params rpctypes.Query4Jrpc
	params.Execer = gty.GuessX

	req := &gty.QueryGuessGameInfo{
		Addr:   addr,
		Status: status,
	}
	params.FuncName = gty.FuncNameQueryGameByAddrStatus
	params.Payload = types.MustPBToJSON(req)
	var res gty.GuessGameRecords
	//ctx := jsonrpc.NewRPCCtx("http://"+types.Conf("config.rpc").GStr("jrpcBindAddr"), "DplatformOS.Query", params, &res)
	ctx := jsonrpc.NewRPCCtx("http://127.0.0.1:9801", "DplatformOS.Query", params, &res)
	ctx.Run()
	return &res
}

func queryGuessByAdminAddrStatus(addr string, status int32) *gty.GuessGameRecords {
	var params rpctypes.Query4Jrpc
	params.Execer = gty.GuessX

	req := &gty.QueryGuessGameInfo{
		AdminAddr: addr,
		Status:    status,
	}
	params.FuncName = gty.FuncNameQueryGameByAdminStatus
	params.Payload = types.MustPBToJSON(req)
	var res gty.GuessGameRecords
	//ctx := jsonrpc.NewRPCCtx("http://"+types.Conf("config.rpc").GStr("jrpcBindAddr"), "DplatformOS.Query", params, &res)
	ctx := jsonrpc.NewRPCCtx("http://127.0.0.1:9801", "DplatformOS.Query", params, &res)
	ctx.Run()
	return &res
}

func queryGuessByCategoryStatus(category string, status int32) *gty.GuessGameRecords {
	var params rpctypes.Query4Jrpc
	params.Execer = gty.GuessX

	req := &gty.QueryGuessGameInfo{
		Category: category,
		Status:   status,
	}
	params.FuncName = gty.FuncNameQueryGameByCategoryStatus
	params.Payload = types.MustPBToJSON(req)
	var res gty.GuessGameRecords
	//ctx := jsonrpc.NewRPCCtx("http://"+types.Conf("config.rpc").GStr("jrpcBindAddr"), "DplatformOS.Query", params, &res)
	ctx := jsonrpc.NewRPCCtx("http://127.0.0.1:9801", "DplatformOS.Query", params, &res)
	ctx.Run()
	return &res
}

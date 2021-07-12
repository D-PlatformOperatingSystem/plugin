package testnode

//DefaultConfig default config for testnode
var DefaultConfig = `
Title="user.p.test."
CoinSymbol="DOM"
# TestNet=true

[log]
#     ，  debug(dbug)/info/warn/error(eror)/crit
loglevel = "debug"
logConsoleLevel = "info"
#      ，    ，                
logFile = "logs/dplatformos.para.log"
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
defCacheSize=128
maxFetchBlockNum=128
timeoutSeconds=5
batchBlockNum=128
driver="leveldb"
dbPath="paradatadir"
dbCache=64
isStrongConsistency=true
singleMode=true
batchsync=false
isRecordBlockSequence=false
isParaChain = true
enableTxQuickIndex=false

[p2p]
enable=false
msgCacheSize=10240
driver="leveldb"
dbPath="paradatadir/addrbook"
dbCache=4
grpcLogFile="grpc.log"


[rpc]
#          
jrpcBindAddr="localhost:8901"
grpcBindAddr="localhost:8902"
whitelist=["127.0.0.1"]
jrpcFuncWhitelist=["*"]
grpcFuncWhitelist=["*"]


[mempool]
name="timeline"
poolCacheSize=10240
minTxFeeRate=100000
maxTxNumPerAccount=10000

[mempool.sub.para]
poolCacheSize=102400

[consensus]
name="para"
genesisBlockTime=1514533390
genesis="16ERTbYtKKQ64wMthAY9J4La4nAiidG45A"
minerExecs=["paracross"]

[mver.consensus]
fundKeyAddr = "1CQXE6TxaYCG5mADtWij4AxhZCUTpoABb3"
powLimitBits = "0x1f00ffff"
maxTxNumber = 1600      #160


[mver.consensus.ticket]
coinReward = 18
coinDevFund = 12
ticketPrice = 10000
retargetAdjustmentFactor = 4
futureBlockTime = 16
ticketFrozenTime = 5    #5s only for test
ticketWithdrawTime = 10 #10s only for test
ticketMinerWaitTime = 2 #2s only for test
targetTimespan = 2304
targetTimePerBlock = 16

[mver.consensus.paracross]
coinReward = 18
coinDevFund = 12


[consensus.sub.para]
#     grpc   ip，       ip    ， “101.37.227.226:28804,39.97.20.242:28804,47.107.15.126:28804,jiedian2.dplatform.io”
ParaRemoteGrpcClient=""
#             
startHeight=1
#      ，   
writeBlockSeconds=2
#    ，             ，          ，       
authAccount="1EbDHAXpoiewjPLX9uqoz38HsKqMXayZrF"
#                    ，         ，   2
waitBlocks4CommitMsg=2
#         ，          block，            blockhash   
searchHashMatchedBlockDepth=10000
#      
genesisAmount=100000000
mainBlockHashForkHeight=1
mainForkParacrossCommitTx=1
mainLoopCheckCommitTxDoneForkHeight=11
selfConsensEnablePreContract=["0-1000"]
emptyBlockInterval=["0:2"]


[store]
name="mavl"
driver="leveldb"
dbPath="paradatadir/mavltree"
dbCache=128
enableMavlPrefix=false
enableMVCC=false
enableMavlPrune=false
pruneHeight=10000

[wallet]
minFee=100000
driver="leveldb"
dbPath="parawallet"
dbCache=16
signType="secp256k1"
minerdisable=true

[exec]
enableStat=false

[exec.sub.relay]
genesis="16ERTbYtKKQ64wMthAY9J4La4nAiidG45A"

[exec.sub.manage]
superManager=[
    "1Bsg9j6gW83sShoee1fZAt9TkUjcrCgA9S",
    "12oupcayRT7LvaC4qW4avxsTE7U41cKSio",
    "1Q8hGLfoGe63efeWa8fJ4Pnukhkngt6poK"
]

[exec.sub.token]
saveTokenTxList=true
tokenApprs = [
	"1Bsg9j6gW83sShoee1fZAt9TkUjcrCgA9S",
	"1Q8hGLfoGe63efeWa8fJ4Pnukhkngt6poK",
	"1LY8GFia5EiyoTodMLfkB5PHNNpXRqxhyB",
	"1GCzJDS6HbgTQ2emade7mEJGGWFfA15pS9",
	"1JYB8sxi4He5pZWHCd3Zi2nypQ4JMB6AxN",
	"12oupcayRT7LvaC4qW4avxsTE7U41cKSio",
]


[pprof]
listenAddr = "localhost:6062"
`

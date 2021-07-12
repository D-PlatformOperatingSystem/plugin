// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package gossip    gossip
package gossip

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/D-PlatformOperatingSystem/dpos/p2p"

	"github.com/D-PlatformOperatingSystem/dpos/client"
	l "github.com/D-PlatformOperatingSystem/dpos/common/log/log15"
	"github.com/D-PlatformOperatingSystem/dpos/queue"
	"github.com/D-PlatformOperatingSystem/dpos/types"

	_ "google.golang.org/grpc/encoding/gzip" // register gzip
)

// P2PTypeName p2p plugin name for gossip
const P2PTypeName = "gossip"

func init() {
	p2p.RegisterP2PCreate(P2PTypeName, New)
}

var (
	log = l.New("module", "p2p")
)

// p2p
type subConfig struct {
	// P2P
	Port int32 `protobuf:"varint,1,opt,name=port" json:"port,omitempty"`
	//     ，   ip:port，         ， seeds=
	Seeds []string `protobuf:"bytes,2,rep,name=seeds" json:"seeds,omitempty"`
	//
	IsSeed bool `protobuf:"varint,3,opt,name=isSeed" json:"isSeed,omitempty"`
	//      ，      seeds
	FixedSeed bool `protobuf:"varint,4,opt,name=fixedSeed" json:"fixedSeed,omitempty"`
	//
	InnerSeedEnable bool `protobuf:"varint,5,opt,name=innerSeedEnable" json:"innerSeedEnable,omitempty"`
	//     Github
	UseGithub bool `protobuf:"varint,6,opt,name=useGithub" json:"useGithub,omitempty"`
	//        ，
	ServerStart bool `protobuf:"varint,7,opt,name=serverStart" json:"serverStart,omitempty"`
	//
	InnerBounds int32 `protobuf:"varint,8,opt,name=innerBounds" json:"innerBounds,omitempty"`
	//           ttl
	LightTxTTL int32 `protobuf:"varint,9,opt,name=lightTxTTL" json:"lightTxTTL,omitempty"`
	//     ttl, ttl
	MaxTTL int32 `protobuf:"varint,10,opt,name=maxTTL" json:"maxTTL,omitempty"`
	// p2p    ,      /   /
	Channel int32 `protobuf:"varint,11,opt,name=channel" json:"channel,omitempty"`
	//           , KB
	MinLtBlockSize int32 `protobuf:"varint,12,opt,name=minLtBlockSize" json:"minLtBlockSize,omitempty"`
	//  p2p  ,   gossip, dht
}

// P2p interface
type P2p struct {
	api          client.QueueProtocolAPI
	client       queue.Client
	node         *Node
	p2pCli       EventInterface
	txCapcity    int32
	txFactory    chan struct{}
	otherFactory chan struct{}
	waitRestart  chan struct{}
	taskGroup    *sync.WaitGroup

	closed  int32
	restart int32
	p2pCfg  *types.P2P
	subCfg  *subConfig
	mgr     *p2p.Manager
	subChan chan interface{}
}

// New produce a p2p object
func New(mgr *p2p.Manager, subCfg []byte) p2p.IP2P {
	cfg := mgr.ChainCfg
	p2pCfg := cfg.GetModuleConfig().P2P
	mcfg := &subConfig{}
	types.MustDecode(subCfg, mcfg)
	//   channel    0,
	if cfg.IsTestNet() && mcfg.Channel == 0 {
		mcfg.Channel = defaultTestNetChannel
	}
	//ttl    2
	if mcfg.LightTxTTL <= 1 {
		mcfg.LightTxTTL = DefaultLtTxBroadCastTTL
	}
	if mcfg.MaxTTL <= 0 {
		mcfg.MaxTTL = DefaultMaxTxBroadCastTTL
	}

	if mcfg.MinLtBlockSize <= 0 {
		mcfg.MinLtBlockSize = defaultMinLtBlockSize
	}

	log.Info("p2p", "Channel", mcfg.Channel, "Version", VERSION, "IsTest", cfg.IsTestNet())
	if mcfg.InnerBounds == 0 {
		mcfg.InnerBounds = 500
	}
	log.Info("p2p", "InnerBounds", mcfg.InnerBounds)

	node, err := NewNode(mgr, mcfg)
	if err != nil {
		log.Error(err.Error())
		return nil
	}
	p2p := new(P2p)
	p2p.node = node
	p2p.p2pCli = NewP2PCli(p2p)
	p2p.txFactory = make(chan struct{}, 1000)    // 1000 task
	p2p.otherFactory = make(chan struct{}, 1000) //other task 1000
	p2p.waitRestart = make(chan struct{}, 1)
	p2p.txCapcity = 1000
	p2p.p2pCfg = p2pCfg
	p2p.subCfg = mcfg
	p2p.client = mgr.Client
	p2p.mgr = mgr
	p2p.api = mgr.SysAPI
	p2p.taskGroup = &sync.WaitGroup{}
	// p2p manger  pub
	p2p.subChan = p2p.mgr.PubSub.Sub(P2PTypeName)
	return p2p
}

//Wait wait for ready
func (network *P2p) Wait() {}

func (network *P2p) isClose() bool {
	return atomic.LoadInt32(&network.closed) == 1
}

func (network *P2p) isRestart() bool {
	return atomic.LoadInt32(&network.restart) == 1
}

//CloseP2P Close network client
func (network *P2p) CloseP2P() {
	log.Info("p2p network start shutdown")
	atomic.StoreInt32(&network.closed, 1)
	//
	network.waitTaskDone()
	network.node.Close()
	network.mgr.PubSub.Unsub(network.subChan)
}

// StartP2P set the queue
func (network *P2p) StartP2P() {
	network.node.SetQueueClient(network.client)

	go func(p2p *P2p) {

		if p2p.isRestart() {
			p2p.node.Start()
			atomic.StoreInt32(&p2p.restart, 0)
			//
			network.waitRestart <- struct{}{}
			return
		}

		p2p.subP2pMsg()
		key, pub := p2p.node.nodeInfo.addrBook.GetPrivPubKey()
		log.Debug("key pub:", pub, "")
		if key == "" {
			if p2p.p2pCfg.WaitPid { //key  ，      ，    ，           ，
				if p2p.genAirDropKeyFromWallet() != nil {
					return
				}
			} else {
				//    Pid,     node award ,airdropaddr
				p2p.node.nodeInfo.addrBook.ResetPeerkey(key, pub)
				go p2p.genAirDropKeyFromWallet()
			}

		} else {
			//key      ，      key,     seed key,
			go p2p.genAirDropKeyFromWallet()

		}
		p2p.node.Start()
		log.Debug("SetQueueClient gorountine ret")

	}(network)
}

func (network *P2p) loadP2PPrivKeyToWallet() error {
	var parm types.ReqWalletImportPrivkey
	parm.Privkey, _ = network.node.nodeInfo.addrBook.GetPrivPubKey()
	parm.Label = "node award"

ReTry:
	resp, err := network.api.ExecWalletFunc("wallet", "WalletImportPrivkey", &parm)
	if err != nil {
		if err == types.ErrPrivkeyExist {
			return nil
		}
		if err == types.ErrLabelHasUsed {
			//    lable
			parm.Label = fmt.Sprintf("node award %v", P2pComm.RandStr(3))
			time.Sleep(time.Second)
			goto ReTry
		}
		log.Error("loadP2PPrivKeyToWallet", "err", err.Error())
		return err
	}

	log.Debug("loadP2PPrivKeyToWallet", "resp", resp.(*types.WalletAccount))
	return nil
}

func (network *P2p) showTaskCapcity() {
	ticker := time.NewTicker(time.Second * 5)
	log.Info("ShowTaskCapcity", "Capcity", atomic.LoadInt32(&network.txCapcity))
	defer ticker.Stop()
	for {
		if network.isClose() {
			log.Debug("ShowTaskCapcity", "loop", "done")
			return
		}
		<-ticker.C
		log.Debug("ShowTaskCapcity", "Capcity", atomic.LoadInt32(&network.txCapcity))
	}
}

func (network *P2p) genAirDropKeyFromWallet() error {
	_, savePub := network.node.nodeInfo.addrBook.GetPrivPubKey()
	for {
		if network.isClose() {
			log.Error("genAirDropKeyFromWallet", "p2p closed", "")
			return fmt.Errorf("p2p closed")
		}

		resp, err := network.api.ExecWalletFunc("wallet", "GetWalletStatus", &types.ReqNil{})
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		if resp.(*types.WalletStatus).GetIsWalletLock() { //
			if savePub == "" {
				log.Warn("P2P Stuck ! Wallet must be unlock and save with mnemonics")

			}
			time.Sleep(time.Second)
			continue
		}

		if !resp.(*types.WalletStatus).GetIsHasSeed() { //
			if savePub == "" {
				log.Warn("P2P Stuck ! Wallet must be imported with mnemonics")

			}
			time.Sleep(time.Second * 5)
			continue
		}

		break
	}

	r := rand.New(rand.NewSource(types.Now().Unix()))
	var minIndex int32 = 100000000
	randIndex := minIndex + r.Int31n(1000000)
	reqIndex := &types.Int32{Data: randIndex}
	msg, err := network.api.ExecWalletFunc("wallet", "NewAccountByIndex", reqIndex)
	if err != nil {
		log.Error("genAirDropKeyFromWallet", "err", err)
		return err
	}
	var hexPrivkey string
	if reply, ok := msg.(*types.ReplyString); !ok {
		log.Error("genAirDropKeyFromWallet", "wrong format data", "")
		panic(err)

	} else {
		hexPrivkey = reply.GetData()
	}
	if hexPrivkey[:2] == "0x" {
		hexPrivkey = hexPrivkey[2:]
	}

	hexPubkey, err := P2pComm.Pubkey(hexPrivkey)
	if err != nil {
		log.Error("genAirDropKeyFromWallet", "gen pub error", err)
		panic(err)
	}

	log.Info("genAirDropKeyFromWallet", "pubkey", hexPubkey)

	if savePub == hexPubkey {
		return nil
	}

	if savePub != "" {
		//priv,pub       ，     ，
		err = network.loadP2PPrivKeyToWallet()
		if err != nil {
			log.Error("genAirDropKeyFromWallet", "loadP2PPrivKeyToWallet error", err)
			panic(err)
		}
		network.node.nodeInfo.addrBook.ResetPeerkey(hexPrivkey, hexPubkey)
		//  p2p
		log.Info("genAirDropKeyFromWallet", "p2p will Restart....")
		network.ReStart()
		return nil
	}
	//  addrbook
	network.node.nodeInfo.addrBook.ResetPeerkey(hexPrivkey, hexPubkey)

	return nil
}

//ReStart p2p
func (network *P2p) ReStart() {
	//
	if !atomic.CompareAndSwapInt32(&network.restart, 0, 1) {
		return
	}
	log.Info("p2p restart, wait p2p task done")
	network.waitTaskDone()
	network.node.Close()
	node, err := NewNode(network.mgr, network.subCfg) //    node
	if err != nil {
		panic(err.Error())
	}
	network.node = node
	network.StartP2P()

}

func (network *P2p) subP2pMsg() {
	if network.client == nil {
		return
	}

	go network.showTaskCapcity()
	go func() {

		var taskIndex int64
		for data := range network.subChan {

			msg, ok := data.(*queue.Message)
			if !ok {
				log.Debug("subP2pMsg", "assetMsg", ok)
				continue
			}
			if network.isClose() {
				log.Debug("subP2pMsg", "loop", "done")
				close(network.otherFactory)
				close(network.txFactory)
				return
			}
			taskIndex++
			log.Debug("p2p recv", "msg", types.GetEventName(int(msg.Ty)), "msg type", msg.Ty, "taskIndex", taskIndex)
			if msg.Ty == types.EventTxBroadcast {
				network.txFactory <- struct{}{} //allocal task
				atomic.AddInt32(&network.txCapcity, -1)
			} else {
				if msg.Ty != types.EventPeerInfo {
					network.otherFactory <- struct{}{}
				}
			}
			switch msg.Ty {

			case types.EventTxBroadcast: //  tx
				network.processEvent(msg, taskIndex, network.p2pCli.BroadCastTx)
			case types.EventBlockBroadcast: //  block
				network.processEvent(msg, taskIndex, network.p2pCli.BlockBroadcast)
			case types.EventFetchBlocks:
				network.processEvent(msg, taskIndex, network.p2pCli.GetBlocks)
			case types.EventGetMempool:
				network.processEvent(msg, taskIndex, network.p2pCli.GetMemPool)
			case types.EventPeerInfo:
				network.processEvent(msg, taskIndex, network.p2pCli.GetPeerInfo)
			case types.EventFetchBlockHeaders:
				network.processEvent(msg, taskIndex, network.p2pCli.GetHeaders)
			case types.EventGetNetInfo:
				network.processEvent(msg, taskIndex, network.p2pCli.GetNetInfo)
			default:
				log.Warn("unknown msgtype", "msg", msg)
				msg.Reply(network.client.NewMessage("", msg.Ty, types.Reply{Msg: []byte("unknown msgtype")}))
				<-network.otherFactory
				continue
			}
		}
		log.Info("subP2pMsg", "loop", "close")

	}()

}

func (network *P2p) processEvent(msg *queue.Message, taskIdx int64, eventFunc p2pEventFunc) {

	//      ，      ，
	if network.isRestart() {
		log.Info("wait for p2p restart....")
		<-network.waitRestart
		log.Info("p2p restart ok....")
	}
	network.taskGroup.Add(1)
	go func() {
		defer network.taskGroup.Done()
		eventFunc(msg, taskIdx)
	}()
}

func (network *P2p) waitTaskDone() {

	waitDone := make(chan struct{})
	go func() {
		defer close(waitDone)
		network.taskGroup.Wait()
	}()
	select {
	case <-waitDone:
	case <-time.After(time.Second * 20):
		log.Error("P2pWaitTaskDone", "err", "20s timeout")
	}
}

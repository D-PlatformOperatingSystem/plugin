// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gossip

import (
	"fmt"
	"math/rand"

	"github.com/D-PlatformOperatingSystem/dpos/p2p"

	//"strings"
	"sync/atomic"

	"sync"
	"time"

	"github.com/D-PlatformOperatingSystem/dpos/common/pubsub"
	"github.com/D-PlatformOperatingSystem/dpos/queue"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	"github.com/D-PlatformOperatingSystem/plugin/plugin/p2p/gossip/nat"
)

//   Node
// 1.    GRPC Server
// 2.
// 3.
// 4.      ，

// Start Node server
func (n *Node) Start() {
	if n.server != nil {
		n.server.Start()
	}
	n.detectNodeAddr()
	n.monitor()
	atomic.StoreInt32(&n.closed, 0)
	go n.doNat()

}

// Close node server
func (n *Node) Close() {
	//
	if !atomic.CompareAndSwapInt32(&n.closed, 0, 1) {
		return
	}
	if n.server != nil {
		n.server.Close()
	}
	log.Debug("stop", "listen", "closed")
	n.nodeInfo.addrBook.Close()
	n.nodeInfo.monitorChan <- nil
	log.Debug("stop", "addrBook", "closed")
	n.removeAll()
	if peerAddrFilter != nil {
		peerAddrFilter.Close()
	}
	n.deleteNatMapPort()
	n.pubsub.Shutdown()
	log.Info("stop", "PeerRemoeAll", "closed")

}

func (n *Node) isClose() bool {
	return atomic.LoadInt32(&n.closed) == 1
}

// Node attribute
type Node struct {
	omtx       sync.Mutex
	nodeInfo   *NodeInfo
	cmtx       sync.Mutex
	cacheBound map[string]*Peer
	outBound   map[string]*Peer
	server     *listener
	listenPort int
	innerSeeds sync.Map
	cfgSeeds   sync.Map
	closed     int32
	pubsub     *pubsub.PubSub
	chainCfg   *types.DplatformOSConfig
	p2pMgr     *p2p.Manager
}

// SetQueueClient return client for nodeinfo
func (n *Node) SetQueueClient(client queue.Client) {
	n.nodeInfo.client = client
}

// NewNode produce a node object
func NewNode(mgr *p2p.Manager, mcfg *subConfig) (*Node, error) {

	cfg := mgr.ChainCfg
	node := &Node{
		outBound:   make(map[string]*Peer),
		cacheBound: make(map[string]*Peer),
		pubsub:     pubsub.NewPubSub(10200),
		p2pMgr:     mgr,
	}
	node.listenPort = 28805
	if mcfg.Port != 0 && mcfg.Port <= 65535 && mcfg.Port > 1024 {
		node.listenPort = int(mcfg.Port)

	}

	if mcfg.InnerSeedEnable {
		seeds := MainNetSeeds
		if cfg.IsTestNet() {
			seeds = TestNetSeeds
		}

		for _, seed := range seeds {
			node.innerSeeds.Store(seed, "inner")
		}
	}

	for _, seed := range mcfg.Seeds {
		node.cfgSeeds.Store(seed, "cfg")
	}
	node.nodeInfo = NewNodeInfo(cfg.GetModuleConfig().P2P, mcfg)
	if mcfg.ServerStart {
		node.server = newListener(protocol, node)
	}
	node.chainCfg = cfg
	return node, nil
}

func (n *Node) flushNodePort(localport, export uint16) {

	if exaddr, err := NewNetAddressString(fmt.Sprintf("%v:%v", n.nodeInfo.GetExternalAddr().IP.String(), export)); err == nil {
		n.nodeInfo.SetExternalAddr(exaddr)
		n.nodeInfo.addrBook.AddOurAddress(exaddr)
	}

	if listenAddr, err := NewNetAddressString(fmt.Sprintf("%v:%v", n.nodeInfo.GetListenAddr().IP.String(), localport)); err == nil {
		n.nodeInfo.SetListenAddr(listenAddr)
		n.nodeInfo.addrBook.AddOurAddress(listenAddr)
	}

}

func (n *Node) natOk() bool {
	n.nodeInfo.natNoticeChain <- struct{}{}
	ok := <-n.nodeInfo.natResultChain
	return ok
}

func (n *Node) doNat() {
	//       ，        0   ，  28805
	for {
		if n.Size() > 0 {
			break
		}
		time.Sleep(time.Second)
	}
	testExaddr := fmt.Sprintf("%v:%v", n.nodeInfo.GetExternalAddr().IP.String(), n.listenPort)
	log.Info("TestNetAddr", "testExaddr", testExaddr)
	if len(P2pComm.AddrRouteble([]string{testExaddr}, n.nodeInfo.channelVersion)) != 0 {
		log.Info("node outside")
		n.nodeInfo.SetNetSide(true)
		if netexaddr, err := NewNetAddressString(testExaddr); err == nil {
			n.nodeInfo.SetExternalAddr(netexaddr)
			n.nodeInfo.addrBook.AddOurAddress(netexaddr)
		}
		return
	}
	log.Info("node inside")
	//   ，       ，
	if !n.nodeInfo.OutSide() && !n.nodeInfo.cfg.IsSeed && n.nodeInfo.cfg.ServerStart {

		go n.natMapPort()
		if !n.natOk() {
			log.Info("doNat", "Nat", "Faild")
		} else {
			//       ，
			for {
				if n.Size() > 0 {
					break
				}
				time.Sleep(time.Millisecond * 100)
			}

			p2pcli := NewNormalP2PCli()
			//                +
			if p2pcli.CheckPeerNatOk(n.nodeInfo.GetExternalAddr().String(), n.nodeInfo) ||
				p2pcli.CheckPeerNatOk(fmt.Sprintf("%v:%v", n.nodeInfo.GetExternalAddr().IP.String(), n.listenPort), n.nodeInfo) {
				n.nodeInfo.SetServiceTy(Service)
				log.Info("doNat", "NatOk", "Support Service")
			} else {
				n.nodeInfo.SetServiceTy(Service - nodeNetwork)
				log.Info("doNat", "NatOk", "No Support Service")
			}

		}

	}

	//n.nodeInfo.SetNatDone()
	n.nodeInfo.addrBook.AddOurAddress(n.nodeInfo.GetExternalAddr())
	n.nodeInfo.addrBook.AddOurAddress(n.nodeInfo.GetListenAddr())
	if selefNet, err := NewNetAddressString(fmt.Sprintf("127.0.0.1:%v", n.nodeInfo.GetListenAddr().Port)); err == nil {
		n.nodeInfo.addrBook.AddOurAddress(selefNet)
	}
}

func (n *Node) addPeer(pr *Peer) {
	n.omtx.Lock()
	defer n.omtx.Unlock()
	if peer, ok := n.outBound[pr.Addr()]; ok {
		log.Info("AddPeer", "delete peer", pr.Addr())
		n.nodeInfo.addrBook.RemoveAddr(peer.Addr())
		delete(n.outBound, pr.Addr())
		peer.Close()

	}
	log.Debug("AddPeer", "peer", pr.Addr())
	n.outBound[pr.Addr()] = pr
	pr.Start()
}

// AddCachePeer  add cacheBound map by addr
func (n *Node) AddCachePeer(pr *Peer) {
	n.cmtx.Lock()
	defer n.cmtx.Unlock()
	n.cacheBound[pr.Addr()] = pr
}

// RemoveCachePeer remove cacheBound by addr
func (n *Node) RemoveCachePeer(addr string) {
	n.cmtx.Lock()
	defer n.cmtx.Unlock()
	delete(n.cacheBound, addr)
}

// HasCacheBound peer whether exists according to address
func (n *Node) HasCacheBound(addr string) bool {
	n.cmtx.Lock()
	defer n.cmtx.Unlock()
	_, ok := n.cacheBound[addr]
	return ok

}

// CacheBoundsSize return node cachebount size
func (n *Node) CacheBoundsSize() int {
	n.cmtx.Lock()
	defer n.cmtx.Unlock()
	return len(n.cacheBound)
}

// GetCacheBounds get node cachebounds
func (n *Node) GetCacheBounds() []*Peer {
	n.cmtx.Lock()
	defer n.cmtx.Unlock()
	var peers []*Peer
	if len(n.cacheBound) == 0 {
		return peers
	}
	for _, peer := range n.cacheBound {
		peers = append(peers, peer)

	}
	return peers
}

// Size return size for peersize
func (n *Node) Size() int {

	return n.nodeInfo.peerInfos.PeerSize()
}

// Has peer whether exists according to address
func (n *Node) Has(paddr string) bool {
	n.omtx.Lock()
	defer n.omtx.Unlock()
	_, ok := n.outBound[paddr]
	return ok
}

// GetRegisterPeer return one peer according to paddr
func (n *Node) GetRegisterPeer(paddr string) *Peer {
	n.omtx.Lock()
	defer n.omtx.Unlock()
	if peer, ok := n.outBound[paddr]; ok {
		return peer
	}
	return nil
}

// GetRegisterPeers return peers
func (n *Node) GetRegisterPeers() []*Peer {
	n.omtx.Lock()
	defer n.omtx.Unlock()
	var peers []*Peer
	if len(n.outBound) == 0 {
		return peers
	}
	for _, peer := range n.outBound {
		peers = append(peers, peer)

	}
	return peers
}

// GetActivePeers return activities of the peers and infos
func (n *Node) GetActivePeers() (map[string]*Peer, map[string]*types.Peer) {
	regPeers := n.GetRegisterPeers()
	infos := n.nodeInfo.peerInfos.GetPeerInfos()

	var peers = make(map[string]*Peer)
	for _, peer := range regPeers {
		name := peer.GetPeerName()
		if _, ok := infos[name]; ok {

			peers[name] = peer
		}
	}
	return peers, infos
}
func (n *Node) remove(peerAddr string) {

	n.omtx.Lock()
	defer n.omtx.Unlock()
	peer, ok := n.outBound[peerAddr]
	if ok {
		delete(n.outBound, peerAddr)
		peer.Close()
	}
}

func (n *Node) removeAll() {
	n.omtx.Lock()
	defer n.omtx.Unlock()
	for addr, peer := range n.outBound {
		delete(n.outBound, addr)
		peer.Close()
	}
}

func (n *Node) monitor() {
	go n.monitorErrPeer()
	//    ,    seeds
	go n.monitorCfgSeeds()
	if !n.nodeInfo.cfg.FixedSeed {
		go n.getAddrFromOnline()
		go n.getAddrFromAddrBook()
	}
	go n.monitorPeerInfo()
	go n.monitorDialPeers()
	go n.monitorBlackList()
	go n.monitorFilter()
	go n.monitorPeers()
	go n.nodeReBalance()
}

func (n *Node) needMore() bool {
	outBoundNum := n.Size()
	return !(outBoundNum >= maxOutBoundNum)
}

func (n *Node) detectNodeAddr() {

	var externalIP string
	for {
		cfg := n.nodeInfo.cfg
		laddr := P2pComm.GetLocalAddr()
		//LocalAddr = laddr
		log.Info("DetectNodeAddr", "addr:", laddr)
		if laddr == "" {
			log.Error("DetectNodeAddr", "NetWork Disable p2p Disable", "Retry until Network enable")
			time.Sleep(time.Second * 5)
			continue
		}
		log.Info("detectNodeAddr", "LocalAddr", laddr)
		if cfg.IsSeed {
			log.Info("DetectNodeAddr", "ExIp", laddr)
			externalIP = laddr
			n.nodeInfo.SetNetSide(true)
			//goto SET_ADDR
		}

		//  nat,getSelfExternalAddr            ，  localaddr
		if len(externalIP) == 0 {
			externalIP = laddr
		}

		var externaladdr string
		var externalPort int

		if cfg.IsSeed {
			externalPort = n.listenPort
		} else {
			exportBytes, err := n.nodeInfo.addrBook.bookDb.Get([]byte(externalPortTag))
			if len(exportBytes) != 0 {
				externalPort = int(P2pComm.BytesToInt32(exportBytes))
			} else {
				externalPort = n.listenPort
			}
			if err != nil {
				log.Error("bookDb Get", "nodePort", n.listenPort, "externalPortTag fail err:", err)
			}
		}

		externaladdr = fmt.Sprintf("%v:%v", externalIP, externalPort)
		log.Debug("DetectionNodeAddr", "AddBlackList", externaladdr)
		n.nodeInfo.blacklist.Add(externaladdr, 0) //                ，    self
		if exaddr, err := NewNetAddressString(externaladdr); err == nil {
			n.nodeInfo.SetExternalAddr(exaddr)
			n.nodeInfo.addrBook.AddOurAddress(exaddr)

		} else {
			log.Error("DetectionNodeAddr", "error", err.Error())
		}

		if listaddr, err := NewNetAddressString(fmt.Sprintf("%v:%v", laddr, n.listenPort)); err == nil {
			n.nodeInfo.SetListenAddr(listaddr)
			n.nodeInfo.addrBook.AddOurAddress(listaddr)
		}

		break
	}
}

func (n *Node) natMapPort() {

	n.natNotice()
	for {
		if n.Size() > 0 {
			break
		}
		time.Sleep(time.Second)
	}
	var err error
	if len(P2pComm.AddrRouteble([]string{n.nodeInfo.GetExternalAddr().String()}, n.nodeInfo.channelVersion)) != 0 { //
		log.Info("natMapPort", "addr", "routeble")
		p2pcli := NewNormalP2PCli() //      IP
		ok := p2pcli.CheckSelf(n.nodeInfo.GetExternalAddr().String(), n.nodeInfo)
		if !ok {
			log.Info("natMapPort", "port is used", n.nodeInfo.GetExternalAddr().String())
			n.flushNodePort(uint16(n.listenPort), uint16(rand.Intn(64512)+1023))
		}

	}
	_, nodename := n.nodeInfo.addrBook.GetPrivPubKey()
	log.Info("natMapPort", "netport", n.nodeInfo.GetExternalAddr().Port)
	for i := 0; i < tryMapPortTimes; i++ {
		//       48
		err = nat.Any().AddMapping("TCP", int(n.nodeInfo.GetExternalAddr().Port), n.listenPort, nodename[:8], time.Hour*48)
		if err != nil {
			if i > tryMapPortTimes/2 { //
				log.Warn("TryNatMapPortFailed", "tryTimes", i, "err", err.Error())
				n.flushNodePort(uint16(n.listenPort), uint16(rand.Intn(64512)+1023))

			}
			log.Info("NatMapPort", "External Port", n.nodeInfo.GetExternalAddr().Port)
			continue
		}

		break
	}

	if err != nil {
		//
		log.Warn("NatMapPort", "Nat", "Faild")
		n.flushNodePort(uint16(n.listenPort), uint16(n.listenPort))
		n.nodeInfo.natResultChain <- false
		return
	}

	err = n.nodeInfo.addrBook.bookDb.Set([]byte(externalPortTag),
		P2pComm.Int32ToBytes(int32(n.nodeInfo.GetExternalAddr().Port))) //            db
	if err != nil {
		log.Error("NatMapPort", "dbErr", err)
		return
	}
	log.Info("natMapPort", "export insert into db", n.nodeInfo.GetExternalAddr().Port)
	n.nodeInfo.natResultChain <- true
	refresh := time.NewTimer(mapUpdateInterval)
	defer refresh.Stop()
	for {
		<-refresh.C
		log.Info("NatWorkRefresh")
		for {
			if err := nat.Any().AddMapping("TCP", int(n.nodeInfo.GetExternalAddr().Port), n.listenPort, nodename[:8], time.Hour*48); err != nil {
				log.Error("NatMapPort update", "err", err.Error())
				time.Sleep(time.Second)
				continue
			}
			break
		}
		refresh.Reset(mapUpdateInterval)

	}
}
func (n *Node) deleteNatMapPort() {

	if n.nodeInfo.OutSide() {
		return
	}

	err := nat.Any().DeleteMapping("TCP", int(n.nodeInfo.GetExternalAddr().Port), n.listenPort)
	if err != nil {
		log.Error("deleteNatMapPort", "DeleteMapping err", err.Error())
	}

}

func (n *Node) natNotice() {
	<-n.nodeInfo.natNoticeChain
}

func (n *Node) verifyP2PChannel(channel int32) bool {
	return channel == n.nodeInfo.cfg.Channel
}

//                ,           ,
func (n *Node) isInBoundPeer(peerName string) (bool, *innerpeer) {

	if n.server == nil || n.server.p2pserver == nil {
		return false, nil
	}
	//
	info := n.server.p2pserver.getInBoundPeerInfo(peerName)
	return info != nil, info
}

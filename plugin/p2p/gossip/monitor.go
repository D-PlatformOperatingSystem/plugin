// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gossip

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/D-PlatformOperatingSystem/dpos/p2p/utils"

	"github.com/D-PlatformOperatingSystem/dpos/types"
)

var (
	peerAddrFilter = utils.NewFilter(PeerAddrCacheNum)
)

func (n *Node) destroyPeer(peer *Peer) {
	log.Debug("deleteErrPeer", "Delete peer", peer.Addr(), "running", peer.GetRunning(),
		"version support", peer.version.IsSupport())

	n.nodeInfo.addrBook.RemoveAddr(peer.Addr())
	n.remove(peer.Addr())

}

func (n *Node) monitorErrPeer() {
	for {
		peer := <-n.nodeInfo.monitorChan
		if peer == nil {
			log.Info("monitorChan close")
			return
		}
		if !peer.version.IsSupport() {
			//       ,
			log.Info("VersoinMonitor", "NotSupport,addr", peer.Addr())
			n.destroyPeer(peer)
			//     12
			n.nodeInfo.blacklist.Add(peer.Addr(), int64(3600*12))
			continue
		}
		if peer.IsMaxInbouds {
			n.destroyPeer(peer)
		}
		if !peer.GetRunning() {
			n.destroyPeer(peer)
			continue
		}

		pstat, ok := n.nodeInfo.addrBook.setAddrStat(peer.Addr(), peer.peerStat.IsOk())
		if ok {
			if pstat.GetAttempts() > maxAttemps {
				log.Debug("monitorErrPeer", "over maxattamps", pstat.GetAttempts())
				n.destroyPeer(peer)
			}
		}
	}
}

func (n *Node) getAddrFromGithub() {
	if !n.nodeInfo.cfg.UseGithub {
		return
	}
	// github
	res, err := http.Get("https://raw.githubusercontent.com/chainseed/seeds/master/dpos.txt")
	if err != nil {
		log.Error("getAddrFromGithub", "http.Get", err.Error())
		return
	}

	bf := new(bytes.Buffer)
	_, err = io.Copy(bf, res.Body)
	if err != nil {
		log.Error("getAddrFromGithub", "io.Copy", err.Error())
		return
	}

	fileContent := bf.String()
	st := strings.TrimSpace(fileContent)
	strs := strings.Split(st, "\n")
	log.Info("getAddrFromGithub", "download file", fileContent)
	for _, linestr := range strs {
		pidaddr := strings.Split(linestr, "@")
		if len(pidaddr) == 2 {
			addr := pidaddr[1]
			if n.Has(addr) || n.nodeInfo.blacklist.Has(addr) {
				return
			}
			n.pubsub.FIFOPub(addr, "addr")

		}
	}
}

// getAddrFromOnline gets the address list from the online node
func (n *Node) getAddrFromOnline() {
	ticker := time.NewTicker(GetAddrFromOnlineInterval)
	defer ticker.Stop()
	pcli := NewNormalP2PCli()

	var ticktimes int
	for {

		<-ticker.C

		seedsMap := make(map[string]bool)
		//    seed
		seedArr := n.nodeInfo.cfg.Seeds
		for _, seed := range seedArr {
			seedsMap[seed] = true
		}

		ticktimes++
		if n.isClose() {
			log.Debug("GetAddrFromOnLine", "loop", "done")
			return
		}

		if n.Size() == 0 && ticktimes > 2 {
			//   Seed
			var rangeCount int
			for addr := range seedsMap {
				//  seed
				rangeCount++
				if rangeCount < maxOutBoundNum {
					n.pubsub.FIFOPub(addr, "addr")
				}

			}

			if rangeCount < maxOutBoundNum {
				// innerSeeds
				n.innerSeeds.Range(func(k, v interface{}) bool {
					rangeCount++
					if rangeCount < maxOutBoundNum {
						n.pubsub.FIFOPub(k.(string), "addr")
						return true
					}
					return false

				})
			}

			continue
		}

		peers, _ := n.GetActivePeers()
		for _, peer := range peers { //         ，

			var addrlist []string
			var addrlistMap map[string]int64

			var err error

			addrlistMap, err = pcli.GetAddrList(peer)

			P2pComm.CollectPeerStat(err, peer)
			if err != nil {
				log.Error("getAddrFromOnline", "ERROR", err.Error())
				continue
			}

			for addr := range addrlistMap {
				addrlist = append(addrlist, addr)
			}

			for _, addr := range addrlist {

				if !n.needMore() {

					//      25   ，
					localBlockHeight, err := pcli.GetBlockHeight(n.nodeInfo)
					if err != nil {
						continue
					}
					//       ，          ,          ，
					if peerHeight, ok := addrlistMap[addr]; ok {

						if localBlockHeight-peerHeight < 1024 {
							if _, ok := seedsMap[addr]; ok {
								continue
							}

							//

							n.innerSeeds.Range(func(k, v interface{}) bool {
								if n.Has(k.(string)) {
									//     cfgseed
									if _, ok := n.cfgSeeds.Load(k.(string)); ok {
										return true
									}
									n.remove(k.(string))
									n.nodeInfo.addrBook.RemoveAddr(k.(string))
									return false
								}
								return true
							})
						}
					}
					time.Sleep(MonitorPeerInfoInterval)
					continue
				}

				if !n.nodeInfo.blacklist.Has(addr) || !peerAddrFilter.Contains(addr) {
					if ticktimes < 10 {
						//         ，
						if _, ok := n.innerSeeds.Load(addr); !ok {
							//  seed
							n.pubsub.FIFOPub(addr, "addr")

						}
					} else {
						n.pubsub.FIFOPub(addr, "addr")
					}

				}
			}

		}

	}
}

func (n *Node) getAddrFromAddrBook() {
	ticker := time.NewTicker(GetAddrFromAddrBookInterval)
	defer ticker.Stop()
	var tickerTimes int64

	for {
		<-ticker.C
		tickerTimes++
		if n.isClose() {
			log.Debug("GetAddrFromOnLine", "loop", "done")
			return
		}
		//12    ，   github
		if tickerTimes > 12 && n.Size() == 0 {
			n.getAddrFromGithub()
			tickerTimes = 0
		}

		log.Debug("OUTBOUND NUM", "NUM", n.Size(), "start getaddr from peer,peernum", len(n.nodeInfo.addrBook.GetPeers()))

		addrNetArr := n.nodeInfo.addrBook.GetPeers()

		for _, addr := range addrNetArr {
			if !n.Has(addr.String()) && !n.nodeInfo.blacklist.Has(addr.String()) {
				log.Debug("GetAddrFromOffline", "Add addr", addr.String())

				if n.needMore() || n.CacheBoundsSize() < maxOutBoundNum {
					n.pubsub.FIFOPub(addr.String(), "addr")

				}
			}
		}

		log.Debug("Node Monitor process", "outbound num", n.Size())
	}

}

func (n *Node) nodeReBalance() {
	p2pcli := NewNormalP2PCli()
	ticker := time.NewTicker(MonitorReBalanceInterval)
	defer ticker.Stop()

	for {
		if n.isClose() {
			log.Debug("nodeReBalance", "loop", "done")
			return
		}

		<-ticker.C
		log.Info("nodeReBalance", "cacheSize", n.CacheBoundsSize())
		if n.CacheBoundsSize() == 0 {
			continue
		}
		peers, _ := n.GetActivePeers()
		//          ，
		var MaxInBounds int32
		var MaxInBoundPeer *Peer
		for _, peer := range peers {
			if peer.GetInBouns() > MaxInBounds {
				MaxInBounds = peer.GetInBouns()
				MaxInBoundPeer = peer
			}
		}
		if MaxInBoundPeer == nil {
			continue
		}

		//
		cachePeers := n.GetCacheBounds()
		var MinCacheInBounds int32 = 1000
		var MinCacheInBoundPeer *Peer
		var MaxCacheInBounds int32
		var MaxCacheInBoundPeer *Peer
		for _, peer := range cachePeers {
			inbounds, err := p2pcli.GetInPeersNum(peer)
			if err != nil {
				n.RemoveCachePeer(peer.Addr())
				peer.Close()
				continue
			}
			//
			if int32(inbounds) < MinCacheInBounds {
				MinCacheInBounds = int32(inbounds)
				MinCacheInBoundPeer = peer
			}

			//
			if int32(inbounds) > MaxCacheInBounds {
				MaxCacheInBounds = int32(inbounds)
				MaxCacheInBoundPeer = peer
			}
		}

		if MinCacheInBoundPeer == nil || MaxCacheInBoundPeer == nil {
			continue
		}

		//
		if MaxInBounds < MaxCacheInBounds {
			n.RemoveCachePeer(MaxCacheInBoundPeer.Addr())
			MaxCacheInBoundPeer.Close()
		}
		//                  ，
		if MaxInBounds < MinCacheInBounds {
			cachePeers := n.GetCacheBounds()
			for _, peer := range cachePeers {
				n.RemoveCachePeer(peer.Addr())
				peer.Close()
			}

			continue
		}
		log.Info("nodeReBalance", "MaxInBounds", MaxInBounds, "MixCacheInBounds", MinCacheInBounds)
		if MaxInBounds-MinCacheInBounds < 50 {
			continue
		}

		if MinCacheInBoundPeer != nil {
			info, err := MinCacheInBoundPeer.GetPeerInfo()
			if err != nil {
				n.RemoveCachePeer(MinCacheInBoundPeer.Addr())
				MinCacheInBoundPeer.Close()
				continue
			}
			localBlockHeight, err := p2pcli.GetBlockHeight(n.nodeInfo)
			if err != nil {
				continue
			}
			peerBlockHeight := info.GetHeader().GetHeight()
			if localBlockHeight-peerBlockHeight < 2048 {
				log.Info("noReBalance", "Repalce node new node", MinCacheInBoundPeer.Addr(), "old node", MaxInBoundPeer.Addr())
				n.addPeer(MinCacheInBoundPeer)
				n.nodeInfo.addrBook.AddAddress(MinCacheInBoundPeer.peerAddr, nil)

				n.remove(MaxInBoundPeer.Addr())
				n.RemoveCachePeer(MinCacheInBoundPeer.Addr())
			}
		}
	}
}

func (n *Node) monitorPeers() {

	p2pcli := NewNormalP2PCli()

	ticker := time.NewTicker(MonitorPeerNumInterval)
	defer ticker.Stop()
	_, selfName := n.nodeInfo.addrBook.GetPrivPubKey()
	for {
		if n.isClose() {
			log.Debug("monitorPeers", "loop", "done")
			return
		}
		<-ticker.C
		localBlockHeight, err := p2pcli.GetBlockHeight(n.nodeInfo)
		if err != nil {
			continue
		}

		peers, infos := n.GetActivePeers()
		for name, pinfo := range infos {
			peerheight := pinfo.GetHeader().GetHeight()
			paddr := pinfo.GetAddr()
			if name == selfName && !pinfo.GetSelf() { //       ，
				//
				n.remove(pinfo.GetAddr())
				n.nodeInfo.addrBook.RemoveAddr(paddr)
				n.nodeInfo.blacklist.Add(paddr, 0)
			}

			if localBlockHeight-peerheight > 2048 {
				//
				if addrMap, err := p2pcli.GetAddrList(peers[paddr]); err == nil {

					for addr := range addrMap {
						if !n.Has(addr) && !n.nodeInfo.blacklist.Has(addr) {
							n.pubsub.FIFOPub(addr, "addr")
						}
					}

				}

				if n.Size() <= stableBoundNum {
					continue
				}
				//       ，
				if _, ok := n.cfgSeeds.Load(paddr); ok {
					continue
				}
				//
				n.remove(paddr)
				n.nodeInfo.addrBook.RemoveAddr(paddr)
			}

		}

	}

}

func (n *Node) monitorPeerInfo() {

	go func() {
		n.nodeInfo.FetchPeerInfo(n)
		ticker := time.NewTicker(MonitorPeerInfoInterval)
		defer ticker.Stop()
		for {
			if n.isClose() {
				return
			}

			<-ticker.C
			n.nodeInfo.FetchPeerInfo(n)
		}
	}()
}

// monitorDialPeers connect the node address concurrently
func (n *Node) monitorDialPeers() {
	var dialCount int
	addrChan := n.pubsub.Sub("addr")
	p2pcli := NewNormalP2PCli()
	for addr := range addrChan {

		if n.isClose() {
			log.Info("monitorDialPeers", "loop", "done")
			return
		}
		if peerAddrFilter.Contains(addr.(string)) {
			//          ，
			continue
		}

		netAddr, err := NewNetAddressString(addr.(string))
		if err != nil {
			continue
		}

		if n.nodeInfo.addrBook.ISOurAddress(netAddr) {
			continue
		}

		//                      TODO:     ,               (             ,     )
		if n.Has(netAddr.String()) || n.nodeInfo.blacklist.Has(netAddr.String()) || n.HasCacheBound(netAddr.String()) {
			log.Debug("DialPeers", "find hash", netAddr.String())
			continue
		}

		//
		if !n.needMore() && n.CacheBoundsSize() >= maxOutBoundNum {
			n.pubsub.FIFOPub(addr, "addr")
			time.Sleep(time.Second * 10)
			continue
		}

		log.Debug("DialPeers", "peer", netAddr.String())
		//      ，
		if dialCount >= maxOutBoundNum*2 {
			n.pubsub.FIFOPub(addr, "addr")
			time.Sleep(time.Second * 10)
			dialCount = len(n.GetRegisterPeers()) + n.CacheBoundsSize()
			continue
		}
		dialCount++
		//
		peerAddrFilter.Add(addr.(string), types.Now().Unix())
		log.Debug("monitorDialPeer", "dialCount", dialCount)
		go func(netAddr *NetAddress) {
			defer peerAddrFilter.Remove(netAddr.String())
			peer, err := P2pComm.dialPeer(netAddr, n)
			if err != nil {
				//
				n.nodeInfo.addrBook.RemoveAddr(netAddr.String())
				log.Error("monitorDialPeers", "peerAddr", netAddr.str, "err", err.Error())
				if err == types.ErrVersion { //     ，     12
					peer.version.SetSupport(false)
					P2pComm.CollectPeerStat(err, peer)
					return
				}
				//    ，   10
				if peer != nil {
					peer.Close()
				}
				if _, ok := n.cfgSeeds.Load(netAddr.String()); !ok {
					n.nodeInfo.blacklist.Add(netAddr.String(), int64(60*10))
				}
				return
			}
			//
			inbounds, err := p2pcli.GetInPeersNum(peer)
			if err != nil {
				peer.Close()
				return
			}
			//      ,      90%，    ，
			if int32(inbounds*100)/n.nodeInfo.cfg.InnerBounds > 90 {
				peer.Close()
				return
			}
			//
			if len(n.GetRegisterPeers()) >= maxOutBoundNum {
				if n.CacheBoundsSize() < maxOutBoundNum {
					n.AddCachePeer(peer)
				} else {
					peer.Close()
				}
				return
			}

			n.addPeer(peer)
			n.nodeInfo.addrBook.AddAddress(netAddr, nil)

		}(netAddr)

	}

}

func (n *Node) monitorBlackList() {
	ticker := time.NewTicker(CheckBlackListInterVal)
	defer ticker.Stop()
	for {
		if n.isClose() {
			log.Info("monitorBlackList", "loop", "done")
			return
		}

		<-ticker.C
		badPeers := n.nodeInfo.blacklist.GetBadPeers()
		now := types.Now().Unix()
		for badPeer, intime := range badPeers {
			if n.nodeInfo.addrBook.IsOurStringAddress(badPeer) {
				continue
			}
			//0
			if 0 == intime {
				continue
			}
			if now > intime {
				n.nodeInfo.blacklist.Delete(badPeer)
			}
		}
	}
}

func (n *Node) monitorFilter() {
	tickTime := time.Second * 30
	peerAddrFilter.ManageRecvFilter(tickTime)
}

//  goroutine

func (n *Node) monitorCfgSeeds() {

	ticker := time.NewTicker(CheckCfgSeedsInterVal)
	defer ticker.Stop()

	for {
		if n.isClose() {
			log.Info("monitorCfgSeeds", "loop", "done")
			return
		}

		<-ticker.C
		n.cfgSeeds.Range(func(k, v interface{}) bool {

			if !n.Has(k.(string)) {
				//
				if n.needMore() { //
					n.pubsub.FIFOPub(k.(string), "addr")
				} else {
					//
					peers, _ := n.GetActivePeers()
					//          ，
					var MaxInBounds int32
					MaxInBoundPeer := &Peer{}
					for _, peer := range peers {
						if peer.GetInBouns() > MaxInBounds {
							MaxInBounds = peer.GetInBouns()
							MaxInBoundPeer = peer
						}
					}

					n.remove(MaxInBoundPeer.Addr())
					n.pubsub.FIFOPub(k.(string), "addr")

				}

			}
			return true
		})
	}

}

package score

import (
	"log"
	"testing"
	"time"

	"github.com/D-PlatformOperatingSystem/dpos/common"
	"github.com/D-PlatformOperatingSystem/dpos/common/address"
	"github.com/D-PlatformOperatingSystem/dpos/common/crypto"
	cty "github.com/D-PlatformOperatingSystem/dpos/system/dapp/coins/types"
	drivers "github.com/D-PlatformOperatingSystem/dpos/system/mempool"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	"github.com/D-PlatformOperatingSystem/dpos/util"
	"github.com/D-PlatformOperatingSystem/dpos/util/testnode"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"

	_ "github.com/D-PlatformOperatingSystem/dpos/system"
)

var (
	c, _       = crypto.New(types.GetSignName("", types.SECP256K1))
	hex        = "CC38546E9E659D15E6B4893F0AB32A06D103931A8230B0BDE71459D2B27D6944"
	a, _       = common.FromHex(hex)
	privKey, _ = c.PrivKeyFromBytes(a)
	toAddr     = address.PubKeyToAddress(privKey.PubKey().Bytes()).String()
	amount     = int64(1e8)
	v          = &cty.CoinsAction_Transfer{Transfer: &types.AssetsTransfer{Amount: amount}}
	transfer   = &cty.CoinsAction{Value: v, Ty: cty.CoinsActionTransfer}
	tx1        = &types.Transaction{Execer: []byte("coins"), Payload: types.Encode(transfer), Fee: 1000000, Expire: 1, To: toAddr}
	tx2        = &types.Transaction{Execer: []byte("coins"), Payload: types.Encode(transfer), Fee: 1000000, Expire: 2, To: toAddr}
	tx3        = &types.Transaction{Execer: []byte("coins"), Payload: types.Encode(transfer), Fee: 1000000, Expire: 3, To: toAddr}
	tx4        = &types.Transaction{Execer: []byte("coins"), Payload: types.Encode(transfer), Fee: 2000000, Expire: 4, To: toAddr}
	tx5        = &types.Transaction{Execer: []byte("coins"), Payload: types.Encode(transfer), Fee: 1000000, Expire: 5, To: toAddr}
	item1      = &drivers.Item{Value: tx1, Priority: tx1.Fee, EnterTime: types.Now().Unix()}
	item2      = &drivers.Item{Value: tx2, Priority: tx2.Fee, EnterTime: types.Now().Unix()}
	item3      = &drivers.Item{Value: tx3, Priority: tx3.Fee, EnterTime: types.Now().Unix() - 1000}
	item4      = &drivers.Item{Value: tx4, Priority: tx4.Fee, EnterTime: types.Now().Unix() - 1000}
	item5      = &drivers.Item{Value: tx5, Priority: tx5.Fee, EnterTime: types.Now().Unix() - 1000}
)

func initEnv(size int64) *Queue {
	if size == 0 {
		size = 100
	}
	_, sub := types.InitCfg("dplatformos.test.toml")
	var subcfg subConfig
	types.MustDecode(sub.Mempool["score"], &subcfg)
	subcfg.PoolCacheSize = size
	cache := NewQueue(subcfg)
	return cache
}

func TestMemFull(t *testing.T) {
	cache := initEnv(1)
	hash := string(tx1.Hash())
	err := cache.Push(item1)
	assert.Nil(t, err)
	assert.Equal(t, true, cache.Exist(hash))
	it, err := cache.GetItem(hash)
	assert.Nil(t, err)
	assert.Equal(t, item1, it)

	_, err = cache.GetItem(hash + ":")
	assert.Equal(t, types.ErrNotFound, err)

	err = cache.Push(item1)
	assert.Equal(t, types.ErrTxExist, err)

	err = cache.Push(item2)
	assert.Equal(t, types.ErrMemFull, err)

	cache.Remove(hash)
	assert.Equal(t, 0, cache.Size())
}

func TestWalk(t *testing.T) {
	//push to item
	cache := initEnv(2)
	cache.Push(item1)
	cache.Push(item2)
	assert.Equal(t, 2, cache.Size())
	var data [2]*drivers.Item
	i := 0
	cache.Walk(1, func(value *drivers.Item) bool {
		data[i] = value
		i++
		return true
	})
	assert.Equal(t, 1, i)
	assert.Equal(t, data[0], item1)

	i = 0
	cache.Walk(2, func(value *drivers.Item) bool {
		data[i] = value
		i++
		return true
	})
	assert.Equal(t, 2, i)
	assert.Equal(t, data[0], item1)
	assert.Equal(t, data[1], item2)

	i = 0
	cache.Walk(2, func(value *drivers.Item) bool {
		data[i] = value
		i++
		return false
	})
	assert.Equal(t, 1, i)
}

func TestTimeCompetition(t *testing.T) {
	cache := initEnv(1)
	cache.Push(item1)
	cache.Push(item3)
	assert.Equal(t, false, cache.Exist(string(item1.Value.Hash())))
	assert.Equal(t, true, cache.Exist(string(item3.Value.Hash())))
	assert.Equal(t, int64(item3.Value.Size()), cache.GetCacheBytes())
}

func TestPriceCompetition(t *testing.T) {
	cache := initEnv(1)
	cache.Push(item3)
	cache.Push(item4)
	assert.Equal(t, false, cache.Exist(string(item3.Value.Hash())))
	assert.Equal(t, true, cache.Exist(string(item4.Value.Hash())))
	assert.Equal(t, int64(item4.Value.Size()), cache.GetCacheBytes())
}

func TestAddDuplicateItem(t *testing.T) {
	cache := initEnv(1)
	cache.Push(item1)
	err := cache.Push(item1)
	assert.Equal(t, types.ErrTxExist, err)
}

func TestQueueDirection(t *testing.T) {
	cache := initEnv(0)
	cache.Push(item1)
	cache.Push(item2)
	cache.Push(item3)
	cache.Push(item4)
	cache.Push(item5)
	i := 0
	lastScore := cache.First().GetScore()
	var tmpScore int64
	cache.Walk(5, func(value *drivers.Item) bool {
		tmpScore = cache.CreateSkipValue(&scoreScore{Item: value, subConfig: cache.subConfig}).Score
		if lastScore < tmpScore {
			return false
		}
		lastScore = tmpScore
		i++
		return true
	})
	assert.Equal(t, 5, i)
	assert.Equal(t, true, lastScore == cache.Last().GetScore())
}

func TestRealNodeMempool(t *testing.T) {
	mockDOM := testnode.New("dplatformos.test.toml", nil)
	cfg := mockDOM.GetClient().GetConfig()
	defer mockDOM.Close()
	mockDOM.Listen()
	mockDOM.WaitHeight(0)
	mockDOM.SendHot()
	mockDOM.WaitHeight(1)
	n := 10
	done := make(chan struct{}, n)
	keys := make([]crypto.PrivKey, n)
	for i := 0; i < n; i++ {
		addr, priv := util.Genaddress()
		tx := util.CreateCoinsTx(cfg, mockDOM.GetHotKey(), addr, 10*types.Coin)
		mockDOM.SendTx(tx)
		keys[i] = priv
	}
	mockDOM.Wait()
	for i := 0; i < n; i++ {
		go func(priv crypto.PrivKey) {
			for i := 0; i < 30; i++ {
				tx := util.CreateCoinsTx(cfg, priv, mockDOM.GetGenesisAddress(), types.Coin/1000)
				reply, err := mockDOM.GetAPI().SendTx(tx)
				if err != nil {
					log.Println(err)
					continue
				}
				//
				tx1 := util.CreateCoinsTx(cfg, priv, mockDOM.GetGenesisAddress(), types.Coin/1000)
				tx2 := util.CreateCoinsTx(cfg, priv, mockDOM.GetGenesisAddress(), types.Coin/1000)
				txgroup, err := types.CreateTxGroup([]*types.Transaction{tx1, tx2}, cfg.GetMinTxFeeRate())
				if err != nil {
					log.Println(err)
					continue
				}
				for i := 0; i < len(txgroup.GetTxs()); i++ {
					err = txgroup.SignN(i, types.SECP256K1, priv)
					if err != nil {
						t.Error(err)
						return
					}
				}
				reply, err = mockDOM.GetAPI().SendTx(txgroup.Tx())
				if err != nil {
					log.Println(err)
					continue
				}
				mockDOM.SetLastSend(reply.GetMsg())
			}
			done <- struct{}{}
		}(keys[i])
	}
	for i := 0; i < n; i++ {
		<-done
	}
	for {
		txs, err := mockDOM.GetAPI().GetMempool(&types.ReqGetMempool{})
		assert.Nil(t, err)
		println("len", len(txs.GetTxs()))
		if len(txs.GetTxs()) > 0 {
			mockDOM.Wait()
			continue
		}
		break
	}
	peer, err := mockDOM.GetAPI().PeerInfo(nil)
	assert.Nil(t, err)
	assert.Equal(t, len(peer.Peers), 0)
	//assert.Equal(t, peer.Peers[0].MempoolSize, int32(0))
}

func TestGetProperFee(t *testing.T) {
	cache := initEnv(0)
	assert.Equal(t, cache.subConfig.ProperFee, cache.GetProperFee())
	cache.Push(item3)
	cache.Push(item4)
	size3 := proto.Size(item3.Value)
	size4 := proto.Size(item3.Value)
	score3 := cache.subConfig.PriceConstant*cache.subConfig.PricePower*(item3.Value.Fee/int64(size3)) -
		item3.EnterTime*cache.subConfig.TimeParam
	score4 := cache.subConfig.PriceConstant*cache.subConfig.PricePower*(item4.Value.Fee/int64(size4)) -
		item4.EnterTime*cache.subConfig.TimeParam
	properFee := ((score3+score4)/2 + time.Now().Unix()*cache.subConfig.TimeParam) * int64(100) /
		(cache.subConfig.PriceConstant * cache.subConfig.PricePower)
	assert.Equal(t, int64(1), properFee/cache.GetProperFee())
}

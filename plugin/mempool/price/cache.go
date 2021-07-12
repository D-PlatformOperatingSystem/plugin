package price

import (
	"github.com/D-PlatformOperatingSystem/dpos/common/skiplist"
	"github.com/D-PlatformOperatingSystem/dpos/system/mempool"
	"github.com/golang/protobuf/proto"
)

// Queue       (  =   /     ,      ,        )
type Queue struct {
	*skiplist.Queue
	subConfig subConfig
}

type priceScore struct {
	*mempool.Item
}

func (item *priceScore) GetScore() int64 {
	txSize := proto.Size(item.Value)
	return item.Value.Fee / int64(txSize)
}

func (item *priceScore) Hash() []byte {
	return item.Value.Hash()
}

func (item *priceScore) Compare(cmp skiplist.Scorer) int {
	it := cmp.(*priceScore)
	//    ，
	if item.EnterTime < it.EnterTime {
		return skiplist.Big
	}
	if item.EnterTime == it.EnterTime {
		return skiplist.Equal
	}
	return skiplist.Small
}

func (item *priceScore) ByteSize() int64 {
	return int64(proto.Size(item.Value))
}

// NewQueue
func NewQueue(subcfg subConfig) *Queue {
	return &Queue{
		Queue:     skiplist.NewQueue(subcfg.PoolCacheSize),
		subConfig: subcfg,
	}
}

//GetItem        key
func (cache *Queue) GetItem(hash string) (*mempool.Item, error) {
	item, err := cache.Queue.GetItem(hash)
	if err != nil {
		return nil, err
	}
	return item.(*priceScore).Item, nil
}

//Push
func (cache *Queue) Push(item *mempool.Item) error {
	return cache.Queue.Push(&priceScore{Item: item})
}

//Walk        key
func (cache *Queue) Walk(count int, cb func(tx *mempool.Item) bool) {
	cache.Queue.Walk(count, func(item skiplist.Scorer) bool {
		return cb(item.(*priceScore).Item)
	})
}

// GetProperFee          ,  100
func (cache *Queue) GetProperFee() int64 {
	var sumFeeRate int64
	var properFeeRate int64
	if cache.Size() < 100 {
		return cache.subConfig.ProperFee
	}
	i := 0
	var feeRate int64
	cache.Walk(100, func(item *mempool.Item) bool {
		//        ,       txsize/1000 + 1
		unitFeeNum := proto.Size(item.Value)/1000 + 1
		//
		if count := item.Value.GetGroupCount(); count > 0 {
			unitFeeNum = int(count)
			txs, err := item.Value.GetTxGroup()
			if err != nil {
				for _, tx := range txs.GetTxs() {
					unitFeeNum += proto.Size(tx) / 1000
				}
			}
		}
		feeRate = item.Value.Fee / int64(unitFeeNum)
		sumFeeRate += feeRate
		i++
		return true
	})
	properFeeRate = sumFeeRate / int64(i)
	return properFeeRate
}

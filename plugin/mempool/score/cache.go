package score

import (
	"time"

	"github.com/D-PlatformOperatingSystem/dpos/common/skiplist"
	"github.com/D-PlatformOperatingSystem/dpos/system/mempool"
	"github.com/golang/protobuf/proto"
)

// Queue       (  =  a*  b*   /     -  c*  ,     ,    ,  a   b,c   )
type Queue struct {
	*skiplist.Queue
	subConfig subConfig
}

type scoreScore struct {
	*mempool.Item
	subConfig subConfig
}

func (item *scoreScore) GetScore() int64 {
	size := proto.Size(item.Value)
	score := item.subConfig.PriceConstant*(item.Value.Fee/int64(size))*
		item.subConfig.PricePower - item.subConfig.TimeParam*item.EnterTime
	return score
}

func (item *scoreScore) Hash() []byte {
	return item.Value.Hash()
}

func (item *scoreScore) Compare(cmp skiplist.Scorer) int {
	it := cmp.(*scoreScore)
	//    ，
	if item.EnterTime < it.EnterTime {
		return skiplist.Big
	}
	if item.EnterTime == it.EnterTime {
		return skiplist.Equal
	}
	return skiplist.Small
}

func (item *scoreScore) ByteSize() int64 {
	return int64(proto.Size(item.Value))
}

// NewQueue
func NewQueue(subcfg subConfig) *Queue {
	return &Queue{
		Queue:     skiplist.NewQueue(subcfg.PoolCacheSize),
		subConfig: subcfg,
	}
}

//func (cache *Queue) newSkipValue(item *mempool.Item) (*skiplist.SkipValue, error) {
//	size := proto.Size(item.Value)
//	return &skiplist.SkipValue{Score: cache.subConfig.PriceConstant*(item.Value.Fee/int64(size))*
//		cache.subConfig.PricePower - cache.subConfig.TimeParam*item.EnterTime, Value: item}, nil
//}

//GetItem        key
func (cache *Queue) GetItem(hash string) (*mempool.Item, error) {
	item, err := cache.Queue.GetItem(hash)
	if err != nil {
		return nil, err
	}
	return item.(*scoreScore).Item, nil
}

// Push    tx   Queue；  tx    Queue  Mempool       error
func (cache *Queue) Push(item *mempool.Item) error {
	return cache.Queue.Push(&scoreScore{Item: item, subConfig: cache.subConfig})
}

// Walk
func (cache *Queue) Walk(count int, cb func(value *mempool.Item) bool) {
	cache.Queue.Walk(count, func(item skiplist.Scorer) bool {
		return cb(item.(*scoreScore).Item)
	})
}

// GetProperFee
func (cache *Queue) GetProperFee() int64 {
	var sumScore int64
	var properFeerate int64
	if cache.Size() == 0 {
		return cache.subConfig.ProperFee
	}
	i := 0
	cache.Queue.Walk(0, func(score skiplist.Scorer) bool {
		if i == 100 {
			return false
		}
		sumScore += score.GetScore()
		i++
		return true
	})
	//   int64(100)
	properFeerate = (sumScore/int64(i) + cache.subConfig.TimeParam*time.Now().Unix()) * int64(100) /
		(cache.subConfig.PriceConstant * cache.subConfig.PricePower)
	return properFeerate
}

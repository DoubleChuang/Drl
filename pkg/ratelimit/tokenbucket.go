package ratelimit

import (
	"sync"
	"time"
)

type Bucket struct {
	availableTokens int64
	capacity        int64
	fillInterval    time.Duration
	latestTick      int64
	lock            sync.Mutex
	startTime       time.Time
}

func NewBucket(fillInterval time.Duration, capacity int64) *Bucket {

	if fillInterval <= 0 {
		panic("token bucket fill interval is not > 0")
	}
	if capacity <= 0 {
		panic("token bucket capacity is not > 0")
	}
	return &Bucket{
		availableTokens: capacity,
		capacity:        capacity,
		fillInterval:    fillInterval,
		latestTick:      0,
		startTime:       time.Now(),
	}
}

//TakeOnce 拿取一個令牌, 並回傳剩餘可用數量
//	成功: 返回剩餘可用的數量
//	失敗: 則返回零
func (tb *Bucket) TakeOnce() int64 {
	tb.lock.Lock()
	defer tb.lock.Unlock()
	now := time.Now()

	tb.computeAvailableTokens(tb.currentTick(now))
	if tb.availableTokens < 0 {
		return 0
	}
	if tb.availableTokens == 0 {
		tb.availableTokens = 1
	}
	tb.availableTokens--
	return tb.availableTokens
}

//Get 獲取當前可用令牌數量
//	成功: 返回剩餘可用的數量
//	失敗: 則返回零
func (tb *Bucket) Get() int64 {
	tb.lock.Lock()
	defer tb.lock.Unlock()

	return tb.availableTokens
}

//currentTick 計算當前時間到起始時間需要補充幾次
func (tb *Bucket) currentTick(now time.Time) int64 {
	return int64(now.Sub(tb.startTime) / tb.fillInterval)
}

//computeAvailableTokens 計算令牌當下的數量
func (b *Bucket) computeAvailableTokens(tick int64) {
	lastTick := b.latestTick
	b.latestTick = tick
	if b.availableTokens >= b.capacity {
		return
	}
	b.availableTokens += (tick - lastTick) * b.capacity
	if b.availableTokens > b.capacity {
		b.availableTokens = b.capacity
	}
	return
}

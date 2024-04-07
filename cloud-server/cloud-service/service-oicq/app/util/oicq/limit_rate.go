package oicq

import (
	"sync"
	"time"
)

type LimitRate struct {
	rate       float64
	capacity   float64
	tokens     float64
	lastFilled time.Time
	mutex      sync.Mutex
}

func NewLimitRate(rate float64, capacity float64) *LimitRate {
	return &LimitRate{
		rate:       rate,
		capacity:   capacity,
		tokens:     capacity,
		lastFilled: time.Now(),
	}
}

func (tb *LimitRate) fillTokens() {
	now := time.Now()
	delta := now.Sub(tb.lastFilled).Seconds()
	tb.tokens = tb.tokens + tb.rate*delta
	if tb.tokens > tb.capacity {
		tb.tokens = tb.capacity
	}
	tb.lastFilled = now
}

func (tb *LimitRate) TryConsume() bool {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()

	tb.fillTokens()

	if tb.tokens >= 1 {
		tb.tokens = tb.tokens - 1
		return true
	}
	return false
}

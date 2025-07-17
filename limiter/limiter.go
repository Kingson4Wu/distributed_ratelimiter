package limiter

import (
	"sync"

	"golang.org/x/time/rate"
)

type Limiter struct {
	mu   sync.Mutex
	rl   *rate.Limiter
	rate float64
}

func NewLimiter(r float64) *Limiter {
	return &Limiter{
		rl:   rate.NewLimiter(rate.Limit(r), int(r)),
		rate: r,
	}
}

func (l *Limiter) Allow() bool {
	return l.rl.Allow()
}

func (l *Limiter) UpdateRate(newRate float64) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.rate = newRate
	l.rl.SetLimit(rate.Limit(newRate))
	l.rl.SetBurst(int(newRate))
}

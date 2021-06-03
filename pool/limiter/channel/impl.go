package channel

import "example/pool/limiter"

type channelLimiter struct {
	limiter chan struct{}
}

func New(limit int) limiter.Limiter {
	c := channelLimiter{limiter: make(chan struct{}, limit)}
	for i := 0; i < limit; i++ {
		c.limiter <- struct{}{}
	}
	return &c
}

func (c channelLimiter) Next() chan struct{} {
	return c.limiter
}

func (c channelLimiter) Ready() {
	c.limiter <- struct{}{}
}

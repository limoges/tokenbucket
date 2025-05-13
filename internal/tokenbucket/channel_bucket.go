package tokenbucket

import (
	"context"
	"time"
)

type ChannelBucket struct {
	tokens chan struct{}
	rate   int
	stop   chan struct{}
}

func NewChannelBucket(initial, max, rate int) *ChannelBucket {
	b := &ChannelBucket{}
	b.tokens = make(chan struct{}, max)
	for range initial {
		b.tokens <- struct{}{}
	}
	b.rate = rate
	b.stop = make(chan struct{})
	go b.run()
	return b
}

func (b *ChannelBucket) Take(ctx context.Context) error {
	select {
	case <-b.tokens:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (b *ChannelBucket) TryTake(ctx context.Context) error {
	select {
	case <-b.tokens:
		return nil
	default:
		return ErrRateLimited
	}
}

func (b *ChannelBucket) run() {
	ticker := time.Tick(1 * time.Second)
	for {
		select {
		case <-ticker:
			for range b.rate {
				select {
				case b.tokens <- struct{}{}:
				case <-b.stop:
					return
				}
			}
		}
	}
}

func (b *ChannelBucket) Stop() {
	b.tokens <- struct{}{}
}

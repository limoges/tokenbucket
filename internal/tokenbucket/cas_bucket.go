package tokenbucket

import (
	"context"
	"errors"
	"sync/atomic"
	"time"
)

type CASBucket struct {
	maxCapacity int32
	tokenCount  atomic.Int32
}

var ErrRateLimited = errors.New("rate limited")

func NewCASBucket(tokenCount, maxCapacity, tokenRate int32) *CASBucket {
	b := &CASBucket{}
	b.maxCapacity = maxCapacity
	b.tokenCount.Store(tokenCount)
	go b.start(1*time.Second, tokenRate)
	return b
}

func (b *CASBucket) Take(ctx context.Context) error {

	for {
		old := b.tokenCount.Load()
		newCount := old - 1

		if newCount < 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				continue
			}
		}

		if b.tokenCount.CompareAndSwap(old, newCount) {
			return nil
		}
	}
}

func (b *CASBucket) TryTake(ctx context.Context) error {
	old := b.tokenCount.Load()
	newCount := old - 1

	if newCount < 0 {
		return ErrRateLimited
	}
	return nil
}

func (b *CASBucket) start(duration time.Duration, count int32) {
	c := time.Tick(duration)
	for {
		select {
		case <-c:
			for {
				if b.replenish(count) {
					break
				}
			}
		}
	}
}

func (b *CASBucket) replenish(count int32) (succeeded bool) {
	old := b.tokenCount.Load()
	newCount := old + count

	if newCount > b.maxCapacity {
		newCount = b.maxCapacity
	}

	return b.tokenCount.CompareAndSwap(old, newCount)
}

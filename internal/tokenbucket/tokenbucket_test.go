package tokenbucket_test

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/limoges/tokenbucket/internal/tokenbucket"
)

func TestTryTake(t *testing.T) {
	b := tokenbucket.NewChannelBucket(0, 1, 0)
	err := b.TryTake(context.TODO())
	if err == nil {
		t.Fail()
	}
}

func TestTakeTwice(t *testing.T) {

	b := tokenbucket.NewChannelBucket(1, 1, 1)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(1 * time.Second)
		cancel()
	}()

	var count atomic.Int32
	var wg sync.WaitGroup

	wg.Add(2)

	take := func(ctx context.Context, b tokenbucket.RateLimit) {
		err := b.Take(ctx)
		if err == nil {
			count.Add(1)
		}
		wg.Done()
	}

	go take(ctx, b)
	go take(ctx, b)

	wg.Wait()

	if count.Load() != 1 {
		t.Log("Count is not 1")
		t.Fail()
	}
}

func TestTakeBasicSuccess(t *testing.T) {

	b := tokenbucket.NewChannelBucket(1, 1, 1)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(1 * time.Second)
		cancel()
	}()

	err := b.Take(ctx)
	if err != nil {
		t.Fail()
	}
}

func TestTakeBasicFailure(t *testing.T) {
	b := tokenbucket.NewChannelBucket(0, 1, 0)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(1 * time.Second)
		cancel()
	}()

	if err := b.Take(ctx); err != nil {
		if !strings.Contains(err.Error(), "context canceled") {
			t.Fail()
		}
	} else {
		t.Log("Should fail.")
		t.Fail()
	}
}

func TestTakeUnderLoad(t *testing.T) {

	var (
		wg     sync.WaitGroup
		count  atomic.Int32
		bucket = tokenbucket.NewChannelBucket(5, 5, 5)
		n      = 100
	)
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	wg.Add(n)

	take := func(b tokenbucket.Bucket, ctx context.Context) {
		for {
			if err := b.Take(ctx); err == nil {
				count.Add(1)
				fmt.Println("took 1")
			}
			select {
			case <-ctx.Done():
				wg.Done()
				return
			default:
				continue
			}
		}
	}

	for range n {
		go take(bucket, ctx)
	}

	time.Sleep(10 * time.Second)
	cancel()
	wg.Wait()
	c := count.Load()
	t.Log("count is", c)
	if !(c <= 55 || c >= 45) {
		t.Fail()
	}
}

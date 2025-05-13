package main

import (
	"context"
	"fmt"
	"os"
	"runtime/pprof"
	"sync"
	"sync/atomic"
	"time"

	"github.com/limoges/tokenbucket/internal/tokenbucket"
	cli "github.com/urfave/cli/v3"
)

func main() {

	cmd := &cli.Command{
		Name:   "tokenbucket",
		Usage:  "Run a simulation of a token bucket implementation",
		Action: run,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "profile",
				Usage: "Write cpu profile to file",
				Value: "profile.txt",
			},
			&cli.StringFlag{
				Name:  "type",
				Usage: `Which bucket implementation to use. One of {"cas", "channel"}`,
				Value: "channel",
			},
		},
	}

	cmd.Run(context.Background(), os.Args)
}

func run(cctx context.Context, cmd *cli.Command) error {

	f, err := os.Create(cmd.String("profile"))
	if err != nil {
		return err
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	var (
		concurrency = 1000
		init        = 0
		max         = 5
		rate        = 5 // per second
		bucketType  = cmd.String("type")
		bucket      tokenbucket.Bucket
		count       atomic.Int64
		wg          sync.WaitGroup
		duration    = 20 * time.Second
	)

	switch bucketType {
	case "cas":
		bucket = tokenbucket.NewCASBucket(int32(init), int32(max), int32(rate))
	case "channel":
		bucket = tokenbucket.NewChannelBucket(init, max, rate)
	}

	ctx, cancel := context.WithTimeout(cctx, duration)
	defer cancel()

	wg.Add(concurrency)
	for range concurrency {
		go func() {
			consume(ctx, bucket, &count)
			wg.Done()
		}()
	}
	wg.Wait()

	return nil
}

func consume(ctx context.Context, bucket tokenbucket.Bucket, count *atomic.Int64) {
	for {
		if err := bucket.Take(ctx); err == nil {
			newCount := count.Add(1)
			fmt.Println(time.Now().UnixNano(), "take successful", "newCount", newCount)
		}
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(100 * time.Millisecond)
			continue
		}
	}
}

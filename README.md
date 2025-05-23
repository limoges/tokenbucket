# tokenbucket

This is a toy project to showcase the difference between [busy waiting] and using
the proper system calls which can block and wake upon lock acquisition.

## Comparison

Here is the difference between two (2) simple implementation of a tokenbucket:

### CAS-based (busy waiting)
![cas-based](./cas-cpu.svg)
Top:
```
File: tokenbucket
Type: cpu
Time: 2025-05-13 18:40:22 CEST
Duration: 20.70s, Total samples = 178.40s (861.76%)
Showing nodes accounting for 177.05s, 99.24% of 178.40s total
Dropped 94 nodes (cum <= 0.89s)
      flat  flat%   sum%        cum   cum%
    48.67s 27.28% 27.28%    177.05s 99.24%  github.com/limoges/tokenbucket/internal/tokenbucket.(*CASBucket).Take
    39.75s 22.28% 49.56%     67.76s 37.98%  runtime.chanrecv
    30.98s 17.37% 66.93%     50.78s 28.46%  context.(*cancelCtx).Done
    27.99s 15.69% 82.62%     27.99s 15.69%  runtime.empty
    19.20s 10.76% 93.38%     19.38s 10.86%  sync/atomic.(*Value).Load (inline)
     9.45s  5.30% 98.68%     77.19s 43.27%  runtime.selectnbrecv
     1.01s  0.57% 99.24%      1.02s  0.57%  runtime.asyncPreempt
         0     0% 99.24%    177.19s 99.32%  main.consume
         0     0% 99.24%    177.19s 99.32%  main.run.func1
         0     0% 99.24%      0.99s  0.55%  runtime.gopreempt_m (inline)
         0     0% 99.24%      0.99s  0.55%  runtime.goschedImpl
         0     0% 99.24%         1s  0.56%  runtime.morestack
         0     0% 99.24%      1.02s  0.57%  runtime.newstack
```

### Channel-based (uses signals under the hood)

![channel-based](./channel-cpu.svg)
Top:
```
File: tokenbucket
Type: cpu
Time: 2025-05-13 18:40:44 CEST
Duration: 20.18s, Total samples = 440ms ( 2.18%)
Showing nodes accounting for 440ms, 100% of 440ms total
      flat  flat%   sum%        cum   cum%
     290ms 65.91% 65.91%      290ms 65.91%  runtime.usleep
     140ms 31.82% 97.73%      140ms 31.82%  runtime.pthread_cond_wait
      10ms  2.27%   100%       10ms  2.27%  runtime.stackpoolalloc
         0     0%   100%      430ms 97.73%  github.com/limoges/tokenbucket/internal/tokenbucket.(*ChannelBucket).Take
         0     0%   100%       10ms  2.27%  github.com/urfave/cli/v3.(*Command).Run (inline)
         0     0%   100%       10ms  2.27%  github.com/urfave/cli/v3.(*Command).run
         0     0%   100%      430ms 97.73%  main.consume
         0     0%   100%       10ms  2.27%  main.main
         0     0%   100%       10ms  2.27%  main.run
         0     0%   100%      430ms 97.73%  main.run.func1
         0     0%   100%      430ms 97.73%  runtime.lock (inline)
         0     0%   100%      430ms 97.73%  runtime.lock2
         0     0%   100%      430ms 97.73%  runtime.lockWithRank (inline)
         0     0%   100%       10ms  2.27%  runtime.main
         0     0%   100%       10ms  2.27%  runtime.malg
         0     0%   100%       10ms  2.27%  runtime.malg.func1
         0     0%   100%       10ms  2.27%  runtime.newproc
         0     0%   100%       10ms  2.27%  runtime.newproc.func1
         0     0%   100%       10ms  2.27%  runtime.newproc1
         0     0%   100%      290ms 65.91%  runtime.osyield (inline)
         0     0%   100%      430ms 97.73%  runtime.selectgo
         0     0%   100%      430ms 97.73%  runtime.sellock
         0     0%   100%      140ms 31.82%  runtime.semasleep
         0     0%   100%       10ms  2.27%  runtime.stackalloc
         0     0%   100%       10ms  2.27%  runtime.stackcacherefill
         0     0%   100%       10ms  2.27%  runtime.systemstack
```


## Inspecting the profiles

You'll need to install `pprof`. Then, use:
```
./tokenbucket --type cas --profile cas.pprof
pprof -http=0.0.0.0:8080 tokenbucket cas.pprof
```

[busy waiting]: https://en.wikipedia.org/wiki/Busy_waiting

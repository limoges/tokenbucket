# tokenbucket

This is a toy project to showcase the difference between [busy waiting] and using
the proper system calls which can block and wake upon lock acquisition.

## Inspecting the profiles

You'll need to install `pprof`. Then, simply use:
```
./tokenbucket --type cas --profile cas.pprof
pprof -http=0.0.0.0:8080 tokenbucket cas.pprof
```

[busy waiting]: https://en.wikipedia.org/wiki/Busy_waiting

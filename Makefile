
generate: channel.pprof cas.pprof

tokenbucket: main.go internal/tokenbucket/
	go build

cas.pprof: tokenbucket
	./tokenbucket --type cas --profile cas.pprof

channel.pprof: tokenbucket
	./tokenbucket --type channel --profile channel.pprof

view:
	pprof -http=0.0.0.0:8080 tokenbucket cas.pprof channel.pprof


clean:
	$(RM) tokenbucket

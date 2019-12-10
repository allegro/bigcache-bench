# bigcache-bench

Benchmarks for BigCache project

### GC pause time
```
go version
go version go1.13 linux/amd64

Number of entries:  20000000
Number of repeats:  50
GC pause for startup:  20.811µs
GC pause for warmup:  3.715446ms
GC pause for freecache:  38.607127ms
GC pause for bigcache:  159.066277ms
GC pause for map:  54.543926ms
```

### Writes and reads
```
# go version
go version go1.13.5 linux/amd64
# go test -bench=. -benchmem -benchtime=4s ./... -timeout 30m
goos: linux
goarch: amd64
pkg: github.com/allegro/bigcache-bench
BenchmarkMapSet-8                     	12373646	       392 ns/op	     201 B/op	       3 allocs/op
BenchmarkConcurrentMapSet-8           	 4108234	      1389 ns/op	     344 B/op	       8 allocs/op
BenchmarkFreeCacheSet-8               	 8555763	       783 ns/op	     342 B/op	       2 allocs/op
BenchmarkBigCacheSet-8                	 8892351	       554 ns/op	     308 B/op	       2 allocs/op
BenchmarkMapGet-8                     	15402984	       329 ns/op	      24 B/op	       1 allocs/op
BenchmarkConcurrentMapGet-8           	12573116	       496 ns/op	      24 B/op	       2 allocs/op
BenchmarkFreeCacheGet-8               	 8183350	       781 ns/op	     136 B/op	       2 allocs/op
BenchmarkBigCacheGet-8                	10523967	       478 ns/op	     152 B/op	       4 allocs/op
BenchmarkBigCacheSetParallel-8        	31305084	       152 ns/op	     319 B/op	       3 allocs/op
BenchmarkFreeCacheSetParallel-8       	20131963	       266 ns/op	     341 B/op	       3 allocs/op
BenchmarkConcurrentMapSetParallel-8   	18498877	       280 ns/op	     200 B/op	       6 allocs/op
BenchmarkBigCacheGetParallel-8        	52064035	       101 ns/op	     152 B/op	       4 allocs/op
BenchmarkFreeCacheGetParallel-8       	44531870	       180 ns/op	     136 B/op	       3 allocs/op
BenchmarkConcurrentMapGetParallel-8   	25893908	       197 ns/op	      24 B/op	       2 allocs/op
PASS
ok  	github.com/allegro/bigcache-bench	254.953s
```

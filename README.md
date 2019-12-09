# bigcache-bench

Benchmarks for BigCache project

### GC pause time
```
go version
go version go1.13 linux/amd64

go run caches_gc_overhead_comparison.go

Number of entries:  20000000
GC pause for bigcache:  1.506077ms
GC pause for freecache:  5.594416ms
GC pause for map:  9.347015ms
```

### Writes and reads
```
# go version
go version go1.13.5 linux/amd64
# go test -bench=. -benchmem -benchtime=4s ./... -timeout 30m

goos: linux
goarch: amd64
pkg: github.com/allegro/bigcache-bench
BenchmarkMapSet-8                     	12408804	       395 ns/op	     201 B/op	       3 allocs/op
BenchmarkConcurrentMapSet-8           	 4251478	      1337 ns/op	     340 B/op	       8 allocs/op
BenchmarkFreeCacheSet-8               	10797914	       708 ns/op	     329 B/op	       2 allocs/op
BenchmarkBigCacheSet/<nil>-8          	 9989348	       490 ns/op	     330 B/op	       2 allocs/op
BenchmarkBigCacheSet/*main.xxHasher-8 	10094407	       488 ns/op	     329 B/op	       2 allocs/op
BenchmarkMapGet-8                     	16063958	       326 ns/op	      24 B/op	       1 allocs/op
BenchmarkConcurrentMapGet-8           	13032650	       409 ns/op	      24 B/op	       2 allocs/op
BenchmarkFreeCacheGet-8               	 9862021	       699 ns/op	     135 B/op	       2 allocs/op
BenchmarkBigCacheGet/<nil>-8          	11149669	       461 ns/op	     152 B/op	       4 allocs/op
BenchmarkBigCacheGet/*main.xxHasher-8 	11311851	       456 ns/op	     152 B/op	       4 allocs/op
BenchmarkBigCacheSetParallel/<nil>-8  	33519153	       158 ns/op	     347 B/op	       3 allocs/op
BenchmarkBigCacheSetParallel/*main.xxHasher-8         	31663240	       144 ns/op	     351 B/op	       3 allocs/op
BenchmarkFreeCacheSetParallel-8                       	35883649	       276 ns/op	     347 B/op	       3 allocs/op
BenchmarkConcurrentMapSetParallel-8                   	19331784	       248 ns/op	     200 B/op	       6 allocs/op
BenchmarkBigCacheGetParallel/<nil>-8                  	58548170	        90.4 ns/op	     152 B/op	       4 allocs/op
BenchmarkBigCacheGetParallel/*main.xxHasher-8         	56315374	       110 ns/op	     152 B/op	       4 allocs/op
BenchmarkFreeCacheGetParallel-8                       	48692251	       122 ns/op	     136 B/op	       3 allocs/op
BenchmarkConcurrentMapGetParallel-8                   	26139178	       191 ns/op	      24 B/op	       2 allocs/op
PASS
ok  	github.com/allegro/bigcache-bench	314.404s
```

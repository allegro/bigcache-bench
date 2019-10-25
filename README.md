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
go version
go version go1.13 linux/amd64

go test -bench=. -benchmem -benchtime=4s ./... -timeout 30m
goos: linux
goarch: amd64
pkg: github.com/allegro/bigcache/v2/caches_bench
BenchmarkMapSet-8                     	12999889	       376 ns/op	     199 B/op	       3 allocs/op
BenchmarkConcurrentMapSet-8           	 4355726	      1275 ns/op	     337 B/op	       8 allocs/op
BenchmarkFreeCacheSet-8               	11068976	       703 ns/op	     328 B/op	       2 allocs/op
BenchmarkBigCacheSet-8                	10183717	       478 ns/op	     304 B/op	       2 allocs/op
BenchmarkMapGet-8                     	16536015	       324 ns/op	      23 B/op	       1 allocs/op
BenchmarkConcurrentMapGet-8           	13165708	       401 ns/op	      24 B/op	       2 allocs/op
BenchmarkFreeCacheGet-8               	10137682	       690 ns/op	     136 B/op	       2 allocs/op
BenchmarkBigCacheGet-8                	11423854	       450 ns/op	     152 B/op	       4 allocs/op
BenchmarkBigCacheSetParallel-8        	34233472	       148 ns/op	     317 B/op	       3 allocs/op
BenchmarkFreeCacheSetParallel-8       	34222654	       268 ns/op	     350 B/op	       3 allocs/op
BenchmarkConcurrentMapSetParallel-8   	19635688	       240 ns/op	     200 B/op	       6 allocs/op
BenchmarkBigCacheGetParallel-8        	60547064	        86.1 ns/op	     152 B/op	       4 allocs/op
BenchmarkFreeCacheGetParallel-8       	50701280	       147 ns/op	     136 B/op	       3 allocs/op
BenchmarkConcurrentMapGetParallel-8   	27353288	       175 ns/op	      24 B/op	       2 allocs/op
PASS
ok  	github.com/allegro/bigcache/v2/caches_bench	256.257s
```

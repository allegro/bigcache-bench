# bigcache-bench

Benchmarks for BigCache project

```
go version
go version go1.13 linux/amd64

go run caches_gc_overhead_comparison.go

Number of entries:  20000000
GC pause for bigcache:  1.506077ms
GC pause for freecache:  5.594416ms
GC pause for map:  9.347015ms
```

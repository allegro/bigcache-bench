# bigcache-bench

Benchmarks for BigCache project

### GC pause time
```
# go version
go version go1.20.2 linux/amd64

# go run caches_gc_overhead_comparison.go -cache bigcache
Cache:              bigcache
Number of entries:  20000000
Number of repeats:  50
Value size:         100
GC pause for startup:  3.591222ms
GC pause for bigcache: 273.027743ms
# go run caches_gc_overhead_comparison.go -cache stdmap
Cache:              stdmap
Number of entries:  20000000
Number of repeats:  50
Value size:         100
GC pause for startup:  2.893986ms
GC pause for stdmap: 111.621318ms

```

### Writes and reads
```
# go version
go version go1.20.2 linux/amd64
# go test -bench=. -benchmem -benchtime=4s ./... -timeout 30m
goos: linux
goarch: amd64
pkg: github.com/allegro/bigcache-bench
cpu: Intel(R) Core(TM) i7-6700K CPU @ 4.00GHz
BenchmarkMapSetForStruct-8                   	    2523	   1887529 ns/op	  696961 B/op	   19746 allocs/op
BenchmarkSyncMapSetForStruct-8               	     934	   4957356 ns/op	 1851193 B/op	   69774 allocs/op
BenchmarkOracamanMapSetForStruct-8           	    1513	   2870550 ns/op	 1226702 B/op	   20195 allocs/op
BenchmarkFreeCacheSetForStruct-8             	     688	   6652622 ns/op	 6983119 B/op	   40287 allocs/op
BenchmarkBigCacheSetForStruct-8              	     760	   5973545 ns/op	 3772362 B/op	   42204 allocs/op
BenchmarkMapSetForBytes-8                    	    1708	   2826723 ns/op	 2096269 B/op	   29749 allocs/op
BenchmarkSyncMapSetForBytes-8                	     814	   5821909 ns/op	 3133981 B/op	   80034 allocs/op
BenchmarkOracamanMapSetForBytes-8            	    1372	   3579288 ns/op	 2976597 B/op	   30276 allocs/op
BenchmarkFreeCacheSetForBytes-8              	     852	   5575751 ns/op	 7864114 B/op	   30538 allocs/op
BenchmarkBigCacheSetForBytes-8               	     906	   5206221 ns/op	 4653427 B/op	   32457 allocs/op
BenchmarkMapGetForStruct-8                   	23807222	       201.1 ns/op	      23 B/op	       1 allocs/op
BenchmarkSyncMapGetForStruct-8               	17040469	       238.3 ns/op	      23 B/op	       1 allocs/op
BenchmarkOracamanMapGetForStruct-8           	18912418	       238.2 ns/op	      23 B/op	       1 allocs/op
BenchmarkFreeCacheGetForStruct-8             	 5201450	       935.8 ns/op	     295 B/op	       9 allocs/op
BenchmarkBigCacheGetForStruct-8              	 5397729	       875.0 ns/op	     287 B/op	       9 allocs/op
BenchmarkMapGetForBytes-8                    	23130621	       207.2 ns/op	      23 B/op	       1 allocs/op
BenchmarkSyncMapGetForBytes-8                	18846595	       237.4 ns/op	      23 B/op	       1 allocs/op
BenchmarkOracamanMapGetForBytes-8            	19344432	       246.4 ns/op	      23 B/op	       1 allocs/op
BenchmarkFreeCacheGetForBytes-8              	11845938	       395.3 ns/op	     159 B/op	       3 allocs/op
BenchmarkBigCacheGetForBytes-8               	12868870	       346.8 ns/op	     151 B/op	       3 allocs/op
BenchmarkSyncMapSetParallelForStruct-8       	 6717847	       694.8 ns/op	      70 B/op	       5 allocs/op
BenchmarkOracamanMapSetParallelForStruct-8   	21306747	       218.4 ns/op	      31 B/op	       2 allocs/op
BenchmarkFreeCacheSetParallelForStruct-8     	18417992	       266.4 ns/op	      54 B/op	       4 allocs/op
BenchmarkBigCacheSetParallelForStruct-8      	17324250	       290.5 ns/op	     203 B/op	       4 allocs/op
BenchmarkSyncMapSetParallelForBytes-8        	 6361482	       773.0 ns/op	     199 B/op	       6 allocs/op
BenchmarkOracamanMapSetParallelForBytes-8    	20362628	       234.1 ns/op	     139 B/op	       3 allocs/op
BenchmarkFreeCacheSetParallelForBytes-8      	19983202	       237.1 ns/op	     143 B/op	       3 allocs/op
BenchmarkBigCacheSetParallelForBytes-8       	17371614	       280.4 ns/op	     443 B/op	       3 allocs/op
BenchmarkSyncMapGetParallelForStruct-8       	24168936	       193.5 ns/op	      23 B/op	       1 allocs/op
BenchmarkOracamanMapGetParallelForStruct-8   	21529862	       186.5 ns/op	      23 B/op	       1 allocs/op
BenchmarkFreeCacheGetParallelForStruct-8     	14119288	       350.2 ns/op	     295 B/op	       9 allocs/op
BenchmarkBigCacheGetParallelForStruct-8      	13292924	       368.0 ns/op	     287 B/op	       9 allocs/op
BenchmarkSyncMapGetParallelForBytes-8        	25061154	       188.6 ns/op	      23 B/op	       1 allocs/op
BenchmarkOracamanMapGetParallelForBytes-8    	25262446	       190.0 ns/op	      23 B/op	       1 allocs/op
BenchmarkFreeCacheGetParallelForBytes-8      	20815525	       232.5 ns/op	     159 B/op	       3 allocs/op
BenchmarkBigCacheGetParallelForBytes-8       	21059220	       220.7 ns/op	     151 B/op	       3 allocs/op
PASS
ok  	github.com/allegro/bigcache-bench	183.467s
```

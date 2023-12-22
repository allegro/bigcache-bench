package main

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/allegro/bigcache/v2"
	"github.com/coocood/freecache"
	cmap "github.com/orcaman/concurrent-map/v2"
)

const maxEntrySize = 256
const maxEntryCount = 10000

type myStruct struct {
	Id int
}

type constructor[T any] interface {
	Get(int) T
	Parse([]byte) T
	ToBytes(T) []byte
}

type byteConstructor []byte

func (bc byteConstructor) Get(n int) []byte {
	return value()
}

func (bc byteConstructor) Parse(data []byte) []byte {
	return data
}

func (bc byteConstructor) ToBytes(v []byte) []byte {
	return v
}

type structConstructor struct {
}

func (sc structConstructor) Get(n int) myStruct {
	return myStruct{Id: n}
}

func (sc structConstructor) Parse(data []byte) myStruct {
	return myStruct{Id: int(binary.BigEndian.Uint64(data))}
}

func (sc structConstructor) ToBytes(v myStruct) []byte {
	b := [8]byte{}
	binary.BigEndian.PutUint64(b[:], uint64(v.Id))
	return b[:]
}

func MapSet[T any](cs constructor[T], b *testing.B) {
	m := make(map[string]T, maxEntryCount)

	id := rand.Intn(maxEntryCount)
	for i := 0; i < b.N; i++ {
		if id >= maxEntryCount {
			id = 0
		}
		m[keys[id]] = cs.Get(id)
	}
}

func SyncMapSet[T any](cs constructor[T], b *testing.B) {
	var m sync.Map

	id := rand.Intn(maxEntryCount)
	for i := 0; i < b.N; i++ {
		if id >= maxEntryCount {
			id = 0
		}
		m.Store(keys[id], cs.Get(id))
	}
}

func OracamanMapSet[T any](cs constructor[T], b *testing.B) {
	m := cmap.New[T]()

	id := rand.Intn(maxEntryCount)
	for i := 0; i < b.N; i++ {
		if id >= maxEntryCount {
			id = 0
		}
		m.Set(keys[id], cs.Get(id))
	}
}

func FreeCacheSet[T any](cs constructor[T], b *testing.B) {
	cache := freecache.NewCache(maxEntryCount * maxEntrySize)

	id := rand.Intn(maxEntryCount)
	for i := 0; i < b.N; i++ {
		if id >= maxEntryCount {
			id = 0
		}
		data := cs.ToBytes(cs.Get(id))
		cache.Set([]byte(keys[id]), data, 0)
	}
}

func BigCacheSet[T any](cs constructor[T], b *testing.B) {
	cache := initBigCache(maxEntryCount)

	id := rand.Intn(maxEntryCount)
	for i := 0; i < b.N; i++ {
		if id >= maxEntryCount {
			id = 0
		}
		data := cs.ToBytes(cs.Get(id))
		cache.Set(keys[id], data)
	}
}

func BenchmarkMapSetForStruct(b *testing.B) {
	MapSet[myStruct](structConstructor{}, b)
}

func BenchmarkSyncMapSetForStruct(b *testing.B) {
	SyncMapSet[myStruct](structConstructor{}, b)
}

func BenchmarkOracamanMapSetForStruct(b *testing.B) {
	OracamanMapSet[myStruct](structConstructor{}, b)
}

func BenchmarkFreeCacheSetForStruct(b *testing.B) {
	FreeCacheSet[myStruct](structConstructor{}, b)
}

func BenchmarkBigCacheSetForStruct(b *testing.B) {
	BigCacheSet[myStruct](structConstructor{}, b)
}

func BenchmarkMapSetForBytes(b *testing.B) {
	MapSet[[]byte](byteConstructor{}, b)
}

func BenchmarkSyncMapSetForBytes(b *testing.B) {
	SyncMapSet[[]byte](byteConstructor{}, b)
}

func BenchmarkOracamanMapSetForBytes(b *testing.B) {
	OracamanMapSet[[]byte](byteConstructor{}, b)
}

func BenchmarkFreeCacheSetForBytes(b *testing.B) {
	FreeCacheSet[[]byte](byteConstructor{}, b)
}

func BenchmarkBigCacheSetForBytes(b *testing.B) {
	BigCacheSet[[]byte](byteConstructor{}, b)
}

func MapGet[T any](cs constructor[T], b *testing.B) {
	b.StopTimer()
	m := make(map[string]T)
	for n := 0; n < maxEntryCount; n++ {
		m[keys[n]] = cs.Get(n)
	}
	b.StartTimer()

	hitCount := 0
	id := rand.Intn(maxEntryCount)
	for i := 0; i < b.N; i++ {
		if id >= maxEntryCount {
			id = 0
		}
		if e, ok := m[keys[id]]; ok {
			_ = (T)(e)
			hitCount++
		}
	}
}

func SyncMapGet[T any](cs constructor[T], b *testing.B) {
	b.StopTimer()
	var m sync.Map
	for n := 0; n < maxEntryCount; n++ {
		m.Store(keys[n], cs.Get(n))
	}
	b.StartTimer()

	hitCounter := 0
	id := rand.Intn(maxEntryCount)
	for i := 0; i < b.N; i++ {
		if id >= maxEntryCount {
			id = 0
		}
		e, ok := m.Load(keys[id])
		if ok {
			_ = (T)(e.(T))
			hitCounter++
		}
	}
}

func OracamanMapGet[T any](cs constructor[T], b *testing.B) {
	b.StopTimer()
	m := cmap.New[T]()
	for n := 0; n < maxEntryCount; n++ {
		m.Set(keys[n], cs.Get(n))
	}
	b.StartTimer()

	hitCounter := 0
	id := rand.Intn(maxEntryCount)
	for i := 0; i < b.N; i++ {
		if id >= maxEntryCount {
			id = 0
		}
		e, ok := m.Get(keys[id])
		if ok {
			_ = (T)(e)
			hitCounter++
		}
	}
}

func FreeCacheGet[T any](cs constructor[T], b *testing.B) {
	b.StopTimer()
	cache := freecache.NewCache(maxEntryCount * maxEntrySize)
	for n := 0; n < maxEntryCount; n++ {
		data := cs.ToBytes(cs.Get(n))
		cache.Set([]byte(keys[n]), data, 0)
	}
	b.StartTimer()

	hitCounter := 0
	id := rand.Intn(maxEntryCount)
	for i := 0; i < b.N; i++ {
		if id >= maxEntryCount {
			id = 0
		}
		data, _ := cache.Get([]byte(keys[id]))
		v := cs.Parse(data)
		_ = (T)(v)
		hitCounter++
	}
}

func BigCacheGet[T any](cs constructor[T], b *testing.B) {
	b.StopTimer()
	cache := initBigCache(maxEntryCount)
	for n := 0; n < maxEntryCount; n++ {
		data := cs.ToBytes(cs.Get(n))
		cache.Set(keys[n], data)
	}
	b.StartTimer()

	hitCount := 0
	id := rand.Intn(maxEntryCount)
	for i := 0; i < b.N; i++ {
		if id >= maxEntryCount {
			id = 0
		}
		data, _ := cache.Get(keys[id])
		v := cs.Parse(data)
		_ = (T)(v)
		hitCount++
	}
}

func BenchmarkMapGetForStruct(b *testing.B) {
	MapGet[myStruct](structConstructor{}, b)
}

func BenchmarkSyncMapGetForStruct(b *testing.B) {
	SyncMapGet[myStruct](structConstructor{}, b)
}

func BenchmarkOracamanMapGetForStruct(b *testing.B) {
	OracamanMapGet[myStruct](structConstructor{}, b)
}

func BenchmarkFreeCacheGetForStruct(b *testing.B) {
	FreeCacheGet[myStruct](structConstructor{}, b)
}

func BenchmarkBigCacheGetForStruct(b *testing.B) {
	BigCacheGet[myStruct](structConstructor{}, b)
}

func BenchmarkMapGetForBytes(b *testing.B) {
	MapGet[[]byte](byteConstructor{}, b)
}

func BenchmarkSyncMapGetForBytes(b *testing.B) {
	SyncMapGet[[]byte](byteConstructor{}, b)
}

func BenchmarkOracamanMapGetForBytes(b *testing.B) {
	OracamanMapGet[[]byte](byteConstructor{}, b)
}

func BenchmarkFreeCacheGetForBytes(b *testing.B) {
	FreeCacheGet[[]byte](byteConstructor{}, b)
}

func BenchmarkBigCacheGetForBytes(b *testing.B) {
	BigCacheGet[[]byte](byteConstructor{}, b)
}

func SyncMapSetParallel[T any](cs constructor[T], b *testing.B) {
	var m sync.Map

	var threadIDCount atomic.Int32

	b.RunParallel(func(pb *testing.PB) {
		threadID := int(threadIDCount.Add(1)) - 1
		for id := rand.Intn(maxEntryCount); pb.Next(); id++ {
			if id >= maxEntryCount {
				id = 0
			}
			m.Store(parallelKey(threadID, id), cs.Get(id))
		}
	})
}

func OracamanMapSetParallel[T any](cs constructor[T], b *testing.B) {
	m := cmap.New[T]()

	var threadIDCount atomic.Int32

	b.RunParallel(func(pb *testing.PB) {
		threadID := int(threadIDCount.Add(1)) - 1
		for id := rand.Intn(maxEntryCount); pb.Next(); id++ {
			if id >= maxEntryCount {
				id = 0
			}
			m.Set(parallelKey(threadID, id), cs.Get(id))
		}
	})
}

func FreeCacheSetParallel[T any](cs constructor[T], b *testing.B) {
	cache := freecache.NewCache(maxEntryCount * maxEntrySize)

	var threadIDCount atomic.Int32

	b.RunParallel(func(pb *testing.PB) {
		threadID := int(threadIDCount.Add(1)) - 1
		for id := rand.Intn(maxEntryCount); pb.Next(); id++ {
			if id >= maxEntryCount {
				id = 0
			}
			data := cs.ToBytes(cs.Get(id))
			cache.Set([]byte(parallelKey(threadID, id)), data, 0)
		}
	})
}

func BigCacheSetParallel[T any](cs constructor[T], b *testing.B) {
	cache := initBigCache(maxEntryCount)

	var threadIDCount atomic.Int32

	b.RunParallel(func(pb *testing.PB) {
		threadID := int(threadIDCount.Add(1)) - 1
		for id := rand.Intn(maxEntryCount); pb.Next(); id++ {
			if id >= maxEntryCount {
				id = 0
			}
			data := cs.ToBytes(cs.Get(id))
			cache.Set(parallelKey(threadID, id), data)
		}
	})
}

func BenchmarkSyncMapSetParallelForStruct(b *testing.B) {
	SyncMapSetParallel[myStruct](structConstructor{}, b)
}

func BenchmarkOracamanMapSetParallelForStruct(b *testing.B) {
	OracamanMapSetParallel[myStruct](structConstructor{}, b)
}

func BenchmarkFreeCacheSetParallelForStruct(b *testing.B) {
	FreeCacheSetParallel[myStruct](structConstructor{}, b)
}

func BenchmarkBigCacheSetParallelForStruct(b *testing.B) {
	BigCacheSetParallel[myStruct](structConstructor{}, b)
}

func BenchmarkSyncMapSetParallelForBytes(b *testing.B) {
	SyncMapSetParallel[[]byte](byteConstructor{}, b)
}

func BenchmarkOracamanMapSetParallelForBytes(b *testing.B) {
	OracamanMapSetParallel[[]byte](byteConstructor{}, b)
}

func BenchmarkFreeCacheSetParallelForBytes(b *testing.B) {
	FreeCacheSetParallel[[]byte](byteConstructor{}, b)
}

func BenchmarkBigCacheSetParallelForBytes(b *testing.B) {
	BigCacheSetParallel[[]byte](byteConstructor{}, b)
}

func SyncMapGetParallel[T any](cs constructor[T], b *testing.B) {
	b.StopTimer()
	var m sync.Map
	for i := 0; i < maxEntryCount; i++ {
		m.Store(keys[i], cs.Get(i))
	}
	b.StartTimer()

	b.RunParallel(func(pb *testing.PB) {
		for id := rand.Intn(maxEntryCount); pb.Next(); id++ {
			if id >= maxEntryCount {
				id = 0
			}
			e, ok := m.Load(keys[id])
			if ok {
				_ = (T)(e.(T))
			}
		}
	})
}

func OracamanMapGetParallel[T any](cs constructor[T], b *testing.B) {
	b.StopTimer()
	m := cmap.New[T]()
	for i := 0; i < maxEntryCount; i++ {
		m.Set(keys[i], cs.Get(i))
	}
	b.StartTimer()

	b.RunParallel(func(pb *testing.PB) {
		for id := rand.Intn(maxEntryCount); pb.Next(); id++ {
			if id >= maxEntryCount {
				id = 0
			}
			e, _ := m.Get(keys[id])
			_ = (T)(e)
		}
	})
}

func FreeCacheGetParallel[T any](cs constructor[T], b *testing.B) {
	b.StopTimer()
	cache := freecache.NewCache(maxEntryCount * maxEntrySize)
	for i := 0; i < maxEntryCount; i++ {
		data := cs.ToBytes(cs.Get(i))
		cache.Set([]byte(keys[i]), data, 0)
	}
	b.StartTimer()

	b.RunParallel(func(pb *testing.PB) {
		for id := rand.Intn(maxEntryCount); pb.Next(); id++ {
			if id >= maxEntryCount {
				id = 0
			}
			data, _ := cache.Get([]byte(keys[id]))
			v := cs.Parse(data)
			_ = (T)(v)
		}
	})
}

func BigCacheGetParallel[T any](cs constructor[T], b *testing.B) {
	b.StopTimer()
	cache := initBigCache(maxEntryCount)
	for i := 0; i < maxEntryCount; i++ {
		data := cs.ToBytes(cs.Get(i))
		cache.Set(keys[i], data)
	}
	b.StartTimer()

	b.RunParallel(func(pb *testing.PB) {
		for id := rand.Intn(maxEntryCount); pb.Next(); id++ {
			if id >= maxEntryCount {
				id = 0
			}
			data, _ := cache.Get(keys[id])
			v := cs.Parse(data)
			_ = (T)(v)
		}
	})
}

func BenchmarkSyncMapGetParallelForStruct(b *testing.B) {
	SyncMapGetParallel[myStruct](structConstructor{}, b)
}

func BenchmarkOracamanMapGetParallelForStruct(b *testing.B) {
	OracamanMapGetParallel[myStruct](structConstructor{}, b)
}

func BenchmarkFreeCacheGetParallelForStruct(b *testing.B) {
	FreeCacheGetParallel[myStruct](structConstructor{}, b)
}

func BenchmarkBigCacheGetParallelForStruct(b *testing.B) {
	BigCacheGetParallel[myStruct](structConstructor{}, b)
}

func BenchmarkSyncMapGetParallelForBytes(b *testing.B) {
	SyncMapGetParallel[[]byte](byteConstructor{}, b)
}

func BenchmarkOracamanMapGetParallelForBytes(b *testing.B) {
	OracamanMapGetParallel[[]byte](byteConstructor{}, b)
}

func BenchmarkFreeCacheGetParallelForBytes(b *testing.B) {
	FreeCacheGetParallel[[]byte](byteConstructor{}, b)
}

func BenchmarkBigCacheGetParallelForBytes(b *testing.B) {
	BigCacheGetParallel[[]byte](byteConstructor{}, b)
}

var (
	keys         = make([]string, maxEntryCount)
	parallelKeys [][]string
)

func init() {
	nThreads := runtime.GOMAXPROCS(0)
	parallelKeys = make([][]string, nThreads)

	for threadID := 0; threadID < nThreads; threadID++ {
		parallelKeys[threadID] = make([]string, maxEntryCount)
		for i := 0; i < maxEntryCount; i++ {
			parallelKeys[threadID][i] = fmt.Sprintf("key-%04d-%06d", threadID, rand.Uint64())
		}
	}

	for i := 0; i < maxEntryCount; i++ {
		keys[i] = fmt.Sprintf("key-%010d", rand.Uint64())
	}
}

func value() []byte {
	return make([]byte, 100)
}

func parallelKey(threadID int, i int) string {
	return parallelKeys[threadID][i]
}

func initBigCache(entriesInWindow int) *bigcache.BigCache {
	cache, _ := bigcache.NewBigCache(bigcache.Config{
		Shards:             256,
		LifeWindow:         10 * time.Minute,
		MaxEntriesInWindow: entriesInWindow,
		MaxEntrySize:       maxEntrySize,
		Verbose:            false,
	})

	return cache
}

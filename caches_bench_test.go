package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/coocood/freecache"
	cmap "github.com/orcaman/concurrent-map/v2"
)

const maxEntrySize = 256
const maxEntryCount = 10000

type myStruct struct {
	Id int `json:"id"`
}

type constructor[T any] interface {
	Get(int) T
	Parse([]byte) (T, error)
	ToBytes(T) ([]byte, error)
}

type byteConstructor []byte

func (bc byteConstructor) Get(n int) []byte {
	return value()
}

func (bc byteConstructor) Parse(data []byte) ([]byte, error) {
	return data, nil
}

func (bc byteConstructor) ToBytes(v []byte) ([]byte, error) {
	return v, nil
}

type structConstructor struct {
}

func (sc structConstructor) Get(n int) myStruct {
	return myStruct{Id: n}
}

func (sc structConstructor) Parse(data []byte) (myStruct, error) {
	var s myStruct
	err := json.Unmarshal(data, &s)
	return s, err
}

func (sc structConstructor) ToBytes(v myStruct) ([]byte, error) {
	return json.Marshal(v)
}

func MapSet[T any](cs constructor[T], b *testing.B) {
	for i := 0; i < b.N; i++ {
		m := make(map[string]T, maxEntryCount)
		for n := 0; n < maxEntryCount; n++ {
			m[key(n)] = cs.Get(n)
		}
	}
}

func SyncMapSet[T any](cs constructor[T], b *testing.B) {
	for i := 0; i < b.N; i++ {
		var m sync.Map
		for n := 0; n < maxEntryCount; n++ {
			m.Store(key(n), cs.Get(n))
		}
	}
}

func OracamanMapSet[T any](cs constructor[T], b *testing.B) {
	for i := 0; i < b.N; i++ {
		m := cmap.New[T]()
		for n := 0; n < maxEntryCount; n++ {
			m.Set(key(n), cs.Get(n))
		}
	}
}

func FreeCacheSet[T any](cs constructor[T], b *testing.B) {
	for i := 0; i < b.N; i++ {
		cache := freecache.NewCache(maxEntryCount * maxEntrySize)
		for n := 0; n < maxEntryCount; n++ {
			data, _ := cs.ToBytes(cs.Get(n))
			cache.Set([]byte(key(n)), data, 0)
		}
	}
}

func BigCacheSet[T any](cs constructor[T], b *testing.B) {
	for i := 0; i < b.N; i++ {
		cache := initBigCache(maxEntryCount)
		for n := 0; n < maxEntryCount; n++ {
			data, _ := cs.ToBytes(cs.Get(n))
			cache.Set(key(n), data)
		}
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
		m[key(n)] = cs.Get(n)
	}
	b.StartTimer()

	hitCount := 0
	for i := 0; i < b.N; i++ {
		id := rand.Intn(maxEntryCount)
		if e, ok := m[key(id)]; ok {
			_ = (T)(e)
			hitCount++
		}
	}
}

func SyncMapGet[T any](cs constructor[T], b *testing.B) {
	b.StopTimer()
	var m sync.Map
	for n := 0; n < maxEntryCount; n++ {
		m.Store(key(n), cs.Get(n))
	}
	b.StartTimer()

	hitCounter := 0
	for i := 0; i < b.N; i++ {
		id := rand.Intn(maxEntryCount)
		e, ok := m.Load(key(id))
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
		m.Set(key(n), cs.Get(n))
	}
	b.StartTimer()

	hitCounter := 0
	for i := 0; i < b.N; i++ {
		id := rand.Intn(maxEntryCount)
		e, ok := m.Get(key(id))
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
		data, _ := cs.ToBytes(cs.Get(n))
		cache.Set([]byte(key(n)), data, 0)
	}
	b.StartTimer()

	hitCounter := 0
	for i := 0; i < b.N; i++ {
		id := rand.Intn(maxEntryCount)
		data, _ := cache.Get([]byte(key(id)))
		v, _ := cs.Parse(data)
		_ = (T)(v)
		hitCounter++
	}
}

func BigCacheGet[T any](cs constructor[T], b *testing.B) {
	b.StopTimer()
	cache := initBigCache(maxEntryCount)
	for n := 0; n < maxEntryCount; n++ {
		data, _ := cs.ToBytes(cs.Get(n))
		cache.Set(key(n), data)
	}
	b.StartTimer()

	hitCount := 0
	for i := 0; i < b.N; i++ {
		id := rand.Intn(maxEntryCount)
		data, _ := cache.Get(key(id))
		v, _ := cs.Parse(data)
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
	b.RunParallel(func(pb *testing.PB) {
		thread := rand.Intn(1000)
		for pb.Next() {
			id := rand.Intn(maxEntryCount)
			m.Store(parallelKey(thread, id), cs.Get(id))
		}
	})
}

func OracamanMapSetParallel[T any](cs constructor[T], b *testing.B) {
	m := cmap.New[T]()

	b.RunParallel(func(pb *testing.PB) {
		thread := rand.Intn(1000)
		for pb.Next() {
			id := rand.Intn(maxEntryCount)
			m.Set(parallelKey(thread, id), cs.Get(id))
		}
	})
}

func FreeCacheSetParallel[T any](cs constructor[T], b *testing.B) {
	cache := freecache.NewCache(maxEntryCount * maxEntrySize)

	b.RunParallel(func(pb *testing.PB) {
		thread := rand.Intn(1000)
		for pb.Next() {
			id := rand.Intn(maxEntryCount)
			data, _ := cs.ToBytes(cs.Get(id))
			cache.Set([]byte(parallelKey(thread, id)), data, 0)
		}
	})
}

func BigCacheSetParallel[T any](cs constructor[T], b *testing.B) {
	cache := initBigCache(maxEntryCount)

	b.RunParallel(func(pb *testing.PB) {
		thread := rand.Intn(1000)
		for pb.Next() {
			id := rand.Intn(maxEntryCount)
			data, _ := cs.ToBytes(cs.Get(id))
			cache.Set(parallelKey(thread, id), data)
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
		m.Store(key(i), cs.Get(i))
	}
	b.StartTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			id := rand.Intn(maxEntryCount)
			e, ok := m.Load(key(id))
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
		m.Set(key(i), cs.Get(i))
	}
	b.StartTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			id := rand.Intn(maxEntryCount)
			e, _ := m.Get(key(id))
			_ = (T)(e)
		}
	})
}

func FreeCacheGetParallel[T any](cs constructor[T], b *testing.B) {
	b.StopTimer()
	cache := freecache.NewCache(maxEntryCount * maxEntrySize)
	for i := 0; i < maxEntryCount; i++ {
		data, _ := cs.ToBytes(cs.Get(i))
		cache.Set([]byte(key(i)), data, 0)
	}
	b.StartTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			id := rand.Intn(maxEntryCount)
			data, _ := cache.Get([]byte(key(id)))
			v, _ := cs.Parse(data)
			_ = (T)(v)
		}
	})
}

func BigCacheGetParallel[T any](cs constructor[T], b *testing.B) {
	b.StopTimer()
	cache := initBigCache(maxEntryCount)
	for i := 0; i < maxEntryCount; i++ {
		data, _ := cs.ToBytes(cs.Get(i))
		cache.Set(key(i), data)
	}
	b.StartTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			id := rand.Intn(maxEntryCount)
			data, _ := cache.Get(key(id))
			v, _ := cs.Parse(data)
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

func key(i int) string {
	return fmt.Sprintf("key-%010d", i)
}

func value() []byte {
	return make([]byte, 100)
}

func parallelKey(threadID int, counter int) string {
	return fmt.Sprintf("key-%04d-%06d", threadID, counter)
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

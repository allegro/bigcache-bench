package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/allegro/bigcache/v2"
	"github.com/coocood/freecache"
)

var previousPause time.Duration

func gcPause() time.Duration {
	runtime.GC()
	var stats debug.GCStats
	debug.ReadGCStats(&stats)
	pause := stats.PauseTotal - previousPause
	previousPause = stats.PauseTotal
	return pause
}

func main() {

	c := ""
	entries := 0
	repeat := 0
	valueSize := 0
	flag.StringVar(&c, "cache", "bigcache", "cache to bench.")
	flag.IntVar(&entries, "entries", 20000000, "number of entries to test")
	flag.IntVar(&repeat, "repeat", 50, "number of repetitions")
	flag.IntVar(&valueSize, "value-size", 100, "size of single entry value in bytes")
	flag.Parse()

	debug.SetGCPercent(10)
	fmt.Println("Cache:             ", c)
	fmt.Println("Number of entries: ", entries)
	fmt.Println("Number of repeats: ", repeat)
	fmt.Println("Value size:        ", valueSize)

	var benchFunc func(kv *keyValueStore)

	switch c {
	case "freecache":
		benchFunc = freeCache
	case "bigcache":
		benchFunc = bigCache
	case "stdmap":
		benchFunc = stdMap
	default:
		fmt.Printf("unknown cache: %s", c)
		os.Exit(1)
	}

	kv := newKeyValueStore(entries, valueSize)

	benchFunc(kv)
	fmt.Println("GC pause for startup: ", gcPause())
	for i := 0; i < repeat; i++ {
		benchFunc(kv)
	}

	fmt.Printf("GC pause for %s: %s\n", c, gcPause())
}

func stdMap(kv *keyValueStore) {
	mapCache := make(map[string][]byte)
	for i := 0; i < kv.Size(); i++ {
		mapCache[kv.Key(i)] = kv.Value(i)
	}
}

func freeCache(kv *keyValueStore) {
	freeCache := freecache.NewCache(kv.Size() * 200) //allocate entries * 200 bytes
	for i := 0; i < kv.Size(); i++ {
		if err := freeCache.Set([]byte(kv.Key(i)), kv.Value(i), 0); err != nil {
			fmt.Println("Error in set: ", err.Error())
		}
	}

	v, err := freeCache.Get([]byte(kv.Key(1)))
	checkFirstElement(kv.Value(1), v, err)

	if freeCache.OverwriteCount() != 0 {
		fmt.Println("Overwritten: ", freeCache.OverwriteCount())
	}
}

func bigCache(kv *keyValueStore) {
	config := bigcache.Config{
		Shards:             256,
		LifeWindow:         100 * time.Minute,
		MaxEntriesInWindow: kv.Size(),
		MaxEntrySize:       200,
		Verbose:            true,
	}

	bigcache, _ := bigcache.NewBigCache(config)
	for i := 0; i < kv.Size(); i++ {
		bigcache.Set(kv.Key(i), kv.Value(i))
	}

	v, err := bigcache.Get(kv.Key(1))
	checkFirstElement(kv.Value(1), v, err)
}

func checkFirstElement(expectedVal []byte, val []byte, err error) {
	if err != nil {
		fmt.Println("Error in get: ", err.Error())
	} else if string(val) != string(expectedVal) {
		fmt.Println("Wrong first element: ", string(val))
	}
}

type keyValueStore struct {
	valueSize int
	keys      []string
	values    []byte
}

func newKeyValueStore(entries int, valueSize int) *keyValueStore {
	keys := make([]string, entries)
	values := make([]byte, entries*valueSize)

	for i := 0; i < entries; i++ {
		keys[i] = fmt.Sprintf("key-%010d", i)
		// Reuse the underlying data of key to generate the value without allocating more memory.
		value := (([]byte)(keys[i]))[4:]
		copy(values[(i+1)*valueSize-len(value):], value)
	}

	return &keyValueStore{
		valueSize: valueSize,
		keys:      keys,
		values:    values,
	}
}

func (store *keyValueStore) Size() int {
	return len(store.keys)
}

func (store *keyValueStore) Key(index int) string {
	return store.keys[index]
}

func (store *keyValueStore) Value(index int) []byte {
	return store.values[index*store.valueSize : (index+1)*store.valueSize]
}

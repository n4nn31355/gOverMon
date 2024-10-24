package series

import (
	"reflect"
	"strings"
	"sync"

	"n4/gui-test/pkg/tickstore"
)

type EntryData = tickstore.TickData[float64]

func NewEntryData(size int) *EntryData {
	return tickstore.NewTickData[float64](size)
}

type Subscriber struct {
	// NOTE: all empty structs have the same address, so we use byte
	byte //nolint:unused
}

type SubscribersMap map[*Subscriber]struct{}

type IEntry interface {
	Subscribe(sub *Subscriber) *EntryData
	Unsubscribe(sub *Subscriber)
}

var _ IEntry = (*Entry)(nil)

type Entry struct {
	lock sync.Mutex

	size int
	data *EntryData

	subscribers SubscribersMap
}

func NewEntry(size int) *Entry {
	if size < 1 {
		panic("size must be greater than zero")
	}
	entry := Entry{size: size, subscribers: make(SubscribersMap)}
	return &entry
}

func (e *Entry) GetData() *EntryData {
	e.lock.Lock()
	defer e.lock.Unlock()
	return e.data
}

func (e *Entry) Subscribe(sub *Subscriber) *EntryData {
	e.lock.Lock()
	defer e.lock.Unlock()

	_, present := e.subscribers[sub]
	if present {
		panic("already subscribed with the same subscriber")
	}
	e.subscribers[sub] = struct{}{}
	if len(e.subscribers) == 1 {
		e.data = NewEntryData(e.size)
	}
	return e.data
}

func (e *Entry) Unsubscribe(sub *Subscriber) {
	e.lock.Lock()
	defer e.lock.Unlock()

	_, present := e.subscribers[sub]
	if !present {
		panic("not subscribed")
	}
	delete(e.subscribers, sub)
	if len(e.subscribers) == 0 {
		e.data = NewEntryData(e.size)
	}
}

func (e *Entry) IsActive() bool {
	e.lock.Lock()
	defer e.lock.Unlock()

	return len(e.subscribers) > 0
}

// TODO: Benchmark GetEntries
func GetEntries[T any](structPtr *T, prefixes ...string) (entries map[string]*Entry) {
	prefix := strings.Join(prefixes, "_")
	entries = make(map[string]*Entry)
	entryType := reflect.TypeOf(Entry{})
	colType := reflect.TypeOf(structPtr).Elem()
	colVal := reflect.ValueOf(structPtr).Elem()
	for x := 0; x < colType.NumField(); x++ {
		field := colType.Field(x)
		if field.Type != entryType {
			continue
		}

		name := field.Tag.Get("json")
		// TODO: validate name(spaces, uniq, etc)
		if name == "" {
			panic("SeriesEntry must have 'json' tag")
		}

		if prefix != "" {
			name = prefix + "_" + name
		}

		_, present := entries[name]
		if present {
			panic("SeriesEntry 'json' tag value must be unique: " + name)
		}

		entries[name] = colVal.Field(x).Addr().Interface().(*Entry)
	}
	return entries
}

// TODO: Benchmark GetEntryNames
// func GetEntryNames[T any](structPtr *T, prefixes ...string) []string {
// 	entries := GetEntries(structPtr, prefixes...)
// 	return slices.Sorted(maps.Keys(entries))
// }

// TODO: Benchmark HasActiveEntries
func HasActiveEntries[T any](structPtr *T) bool {
	entries := GetEntries(structPtr)
	for _, entry := range entries {
		if entry.IsActive() {
			return true
		}
	}
	return false
}

type valToEntry struct {
	val   float64
	entry *Entry
}

func mapValues(
	valToEntryMapping []valToEntry,
	valComp func(val float64, entry *Entry) float64,
) {
	for _, mapping := range valToEntryMapping {
		if !mapping.entry.IsActive() {
			continue
		}
		if mapping.entry.data == nil {
			mapping.entry.data = NewEntryData(mapping.entry.size)
		}
		mapping.entry.data.AddValues(valComp(mapping.val, mapping.entry))
	}
}

func MapValues(valToEntryMapping []valToEntry) {
	mapValues(valToEntryMapping, func(val float64, _ *Entry) float64 {
		return val
	})
}

func MapValuesWithDiff(valToEntryMapping []valToEntry) {
	mapValues(valToEntryMapping, func(val float64, entry *Entry) float64 {
		return val - entry.data.GetFirstValue()
	})
}

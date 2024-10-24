package app

import (
	"os"
	"reflect"
	"runtime"
	"time"

	"n4/gui-test/pkg/series"
	"n4/gui-test/pkg/tickstore"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/shirou/gopsutil/v4/process"
)

type TickDataFloat64 = tickstore.TickData[float64]

func NewTickData(size int) *TickDataFloat64 {
	return tickstore.NewTickData[float64](size)
}

type Stats struct {
	size    int
	Updated time.Time

	StatsUpdate   series.Entry `json:"stats_update"`
	SelfFramerate series.Entry `json:"self_framerate"`

	SelfCPU *series.ProcessCPUCollector `json:"self_cpu"`
	SelfMem *series.ProcessMemCollector `json:"self_mem"`

	RuntimeMemAlloc series.Entry `json:"runtime_mem_alloc"`
	RuntimeMemSys   series.Entry `json:"runtime_mem_sys"`

	SysMem   *series.SysMemCollector   `json:"sys_mem"`
	SysMemEx *series.SysMemExCollector `json:"sys_mem_ex"`

	Disks *series.DiskCollector `json:"disk"`
}

func NewStats(size int) *Stats {
	if size < 1 {
		panic("size must be greater than zero")
	}
	proc, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		panic(err) // FIXME: do not panic?
	}

	stats := Stats{
		size: size,

		StatsUpdate:   *series.NewEntry(size),
		SelfFramerate: *series.NewEntry(size),

		SelfCPU: series.NewProcessCPUCollector(proc, size),
		SelfMem: series.NewProcessMemCollector(proc, size),

		RuntimeMemAlloc: *series.NewEntry(size),
		RuntimeMemSys:   *series.NewEntry(size),

		SysMem:   series.NewSysMemCollector(size),
		SysMemEx: series.NewSysMemExCollector(size),

		Disks: series.NewDiskCollector(size),
	}

	sub := &series.Subscriber{}
	stats.StatsUpdate.Subscribe(sub)
	stats.SelfFramerate.Subscribe(sub)
	stats.RuntimeMemAlloc.Subscribe(sub)
	stats.RuntimeMemSys.Subscribe(sub)

	return &stats
}

func (s *Stats) getFieldPrefix(name string) string {
	colType := reflect.TypeOf(s).Elem()
	field, found := colType.FieldByName(name)
	if !found {
		panic("field not found")
	}
	prefix := field.Tag.Get("json")
	// TODO: validate name(spaces, uniq, etc)
	if prefix == "" {
		panic("SeriesEntry must have 'json' tag")
	}
	return prefix
}

// FIXME: Handle errors
func (s *Stats) Update() {
	perfStart := time.Now()
	defer func() {
		s.StatsUpdate.GetData().AddValues(time.Since(perfStart).Seconds())
	}()

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	s.RuntimeMemAlloc.GetData().AddValues(float64(memStats.Alloc))
	s.RuntimeMemSys.GetData().AddValues(float64(memStats.Sys))
	s.SelfFramerate.GetData().AddValues(ebiten.ActualTPS())

	err := s.SelfCPU.Collect()
	if err != nil {
		panic(err) // FIXME: do not panic?
	}
	err = s.SelfMem.Collect()
	if err != nil {
		panic(err) // FIXME: do not panic?
	}
	err = s.SysMem.Collect()
	if err != nil {
		panic(err) // FIXME: do not panic?
	}
	// FIXME: handle non-windows
	err = s.SysMemEx.Collect()
	if err != nil {
		panic(err) // FIXME: do not panic?
	}

	err = s.Disks.Collect()
	if err != nil {
		panic(err) // FIXME: do not panic?
	}

	s.Updated = time.Now()
}

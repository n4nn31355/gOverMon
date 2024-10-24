package series

import (
	"fmt"

	"github.com/shirou/gopsutil/v4/process"
)

type ProcessMemCollector struct {
	Collector

	proc *process.Process

	RSS    Entry `json:"rss"`
	VMS    Entry `json:"vms"`
	HWM    Entry `json:"hwm"`
	Data   Entry `json:"data"`
	Stack  Entry `json:"stack"`
	Locked Entry `json:"locked"`
	Swap   Entry `json:"swap"`
}

func NewProcessMemCollector(proc *process.Process, size int) *ProcessMemCollector {
	if size < 1 {
		panic("size must be greater than zero")
	}
	collector := ProcessMemCollector{
		Collector: Collector{size: size},

		proc: proc,

		RSS:    *NewEntry(size),
		VMS:    *NewEntry(size),
		HWM:    *NewEntry(size),
		Data:   *NewEntry(size),
		Stack:  *NewEntry(size),
		Locked: *NewEntry(size),
		Swap:   *NewEntry(size),
	}
	return &collector
}

func (c *ProcessMemCollector) Collect() error {
	if !HasActiveEntries(c) {
		return nil
	}

	mem, err := c.proc.MemoryInfo()
	if err != nil {
		return fmt.Errorf("failed to get memory stats: %w", err)
	}

	MapValues([]valToEntry{
		{float64(mem.RSS), &c.RSS},
		{float64(mem.VMS), &c.VMS},
		{float64(mem.HWM), &c.HWM},
		{float64(mem.Data), &c.Data},
		{float64(mem.Stack), &c.Stack},
		{float64(mem.Locked), &c.Locked},
		{float64(mem.Swap), &c.Swap},
	})

	return nil
}

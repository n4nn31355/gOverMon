package series

import (
	"fmt"

	"github.com/shirou/gopsutil/v4/mem"
)

type SysMemCollector struct {
	Collector

	Total       Entry `json:"total"`
	Available   Entry `json:"available"`
	Used        Entry `json:"used"`
	UsedPercent Entry `json:"used_percent"`
}

func NewSysMemCollector(size int) *SysMemCollector {
	if size < 1 {
		panic("size must be greater than zero")
	}
	collector := SysMemCollector{
		Collector: Collector{size: size},

		Total:       *NewEntry(size),
		Available:   *NewEntry(size),
		Used:        *NewEntry(size),
		UsedPercent: *NewEntry(size),
	}
	return &collector
}

func (c *SysMemCollector) Collect() error {
	if !HasActiveEntries(c) {
		return nil
	}

	sysMem, err := mem.VirtualMemory()
	if err != nil {
		return fmt.Errorf("failed to get memory stats: %w", err)
	}

	MapValues([]valToEntry{
		{float64(sysMem.Total), &c.Total},
		{float64(sysMem.Available), &c.Available},
		{float64(sysMem.Used), &c.Used},
		{float64(sysMem.UsedPercent), &c.UsedPercent},
	})

	return nil
}

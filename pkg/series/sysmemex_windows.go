package series

import (
	"fmt"

	"github.com/shirou/gopsutil/v4/mem"
)

type SysMemExCollector struct {
	Collector

	CommitLimit Entry `json:"commit_limit"`
	CommitTotal Entry `json:"commit_total"`
}

func NewSysMemExCollector(size int) *SysMemExCollector {
	if size < 1 {
		panic("size must be greater than zero")
	}
	collector := SysMemExCollector{
		Collector: Collector{size: size},

		CommitLimit: *NewEntry(size),
		CommitTotal: *NewEntry(size),
	}
	return &collector
}

func (c *SysMemExCollector) Collect() error {
	if !HasActiveEntries(c) {
		return nil
	}

	sysMem, err := mem.NewExWindows().VirtualMemory()
	if err != nil {
		return fmt.Errorf("failed to get memory stats: %w", err)
	}

	MapValues([]valToEntry{
		{float64(sysMem.CommitLimit), &c.CommitLimit},
		{float64(sysMem.CommitTotal), &c.CommitTotal},
	})

	return nil
}

package series

import (
	"fmt"

	"github.com/shirou/gopsutil/v4/process"
)

type ProcessCPUCollector struct {
	Collector

	proc *process.Process

	Sys       Entry `json:"sys"`
	User      Entry `json:"user"`
	Idle      Entry `json:"idle"`
	Iowait    Entry `json:"iowait"`
	Steal     Entry `json:"steal"`
	Irq       Entry `json:"irq"`
	Guest     Entry `json:"guest"`
	Softirq   Entry `json:"softirq"`
	Nice      Entry `json:"nice"`
	GuestNice Entry `json:"guestnice"`

	Perc Entry `json:"perc"`
}

func NewProcessCPUCollector(proc *process.Process, size int) *ProcessCPUCollector {
	if size < 1 {
		panic("size must be greater than zero")
	}
	collector := ProcessCPUCollector{
		Collector: Collector{size: size},

		proc: proc,

		Sys:       *NewEntry(size),
		User:      *NewEntry(size),
		Idle:      *NewEntry(size),
		Iowait:    *NewEntry(size),
		Steal:     *NewEntry(size),
		Irq:       *NewEntry(size),
		Guest:     *NewEntry(size),
		Softirq:   *NewEntry(size),
		Nice:      *NewEntry(size),
		GuestNice: *NewEntry(size),

		Perc: *NewEntry(size),
	}
	return &collector
}

func (c *ProcessCPUCollector) Collect() error {
	if !HasActiveEntries(c) {
		return nil
	}

	times, err := c.proc.Times()
	if err != nil {
		return fmt.Errorf("failed to get cpu times: %w", err)
	}

	MapValuesWithDiff([]valToEntry{
		{times.System, &c.Sys},
		{times.User, &c.User},
		{times.Idle, &c.Idle},
		{times.Iowait, &c.Iowait},
		{times.Steal, &c.Steal},
		{times.Irq, &c.Irq},
		{times.Guest, &c.Guest},
		{times.Softirq, &c.Softirq},
		{times.Nice, &c.Nice},
		{times.GuestNice, &c.GuestNice},
	})

	if c.Perc.IsActive() {
		perc, err := c.proc.CPUPercent()
		if err != nil {
			return fmt.Errorf("failed to get cpu percent: %w", err)
		}
		c.Perc.data.AddValues(perc)
	}

	return nil
}

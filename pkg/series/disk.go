package series

import (
	"fmt"

	"github.com/shirou/gopsutil/v4/disk"
)

type DiskStats struct {
	Name string

	mountpoint string

	Total       Entry `json:"total"`
	Free        Entry `json:"free"`
	Used        Entry `json:"used"`
	UsedPercent Entry `json:"used_percent"`

	InodesTotal       Entry `json:"inodes_total"`
	InodesUsed        Entry `json:"inodes_used"`
	InodesFree        Entry `json:"inodes_free"`
	InodesUsedPercent Entry `json:"inodes_used_percent"`
}

func NewDiskStats(size int, mountpoint string) *DiskStats {
	if size < 1 {
		panic("size must be greater than zero")
	}
	ret := &DiskStats{
		Name: mountpoint,

		mountpoint: mountpoint,

		Total:       *NewEntry(size),
		Free:        *NewEntry(size),
		Used:        *NewEntry(size),
		UsedPercent: *NewEntry(size),

		InodesTotal:       *NewEntry(size),
		InodesUsed:        *NewEntry(size),
		InodesFree:        *NewEntry(size),
		InodesUsedPercent: *NewEntry(size),
	}
	return ret
}

type DiskCollector struct {
	Collector

	Disks map[string]*DiskStats
}

func NewDiskCollector(size int) *DiskCollector {
	if size < 1 {
		panic("size must be greater than zero")
	}
	ret := &DiskCollector{
		Collector: Collector{size: size},

		Disks: map[string]*DiskStats{},
	}

	return ret
}

// TODO: handle new devices?
// TODO: handle disconnected devices?
// TODO: filter from config
func (c *DiskCollector) Discover() error {
	partitions, err := disk.Partitions(true)
	if err != nil {
		return fmt.Errorf("failed to get disk stats: %w", err)
	}

	for _, part := range partitions {
		_, present := c.Disks[part.Mountpoint]
		if !present {
			c.Disks[part.Mountpoint] = NewDiskStats(c.size, part.Mountpoint)
		}
	}

	return nil
}

func (c *DiskCollector) Collect() error {
	for _, dStats := range c.Disks {
		if !HasActiveEntries(dStats) {
			continue
		}

		usage, err := disk.Usage(dStats.mountpoint)
		if err != nil {
			return fmt.Errorf("failed to get disk stats: %w", err)
		}

		MapValues([]valToEntry{
			{float64(usage.Total), &dStats.Total},
			{float64(usage.Free), &dStats.Free},
			{float64(usage.Used), &dStats.Used},
			{float64(usage.UsedPercent), &dStats.UsedPercent},

			{float64(usage.InodesTotal), &dStats.InodesTotal},
			{float64(usage.InodesUsed), &dStats.InodesUsed},
			{float64(usage.InodesFree), &dStats.InodesFree},
			{float64(usage.InodesUsedPercent), &dStats.InodesUsedPercent},
		})
	}

	return nil
}

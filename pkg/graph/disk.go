package graph

import (
	"maps"
	"slices"

	"n4/gui-test/pkg/app"
	"n4/gui-test/pkg/series"
)

func Disk(disk *series.DiskStats) *Graph {
	usedSeries := []*series.Entry{
		&disk.Free,
		&disk.Total,
	}
	sub := &series.Subscriber{}
	total := disk.Total.Subscribe(sub)
	data := disk.Free.Subscribe(sub)

	getMax := func() float64 { return total.GetFirstValue() }

	setts := NewSettings("Free "+disk.Name, fmtCBMem)
	setts.AutoMinMaxPadding = 0

	gr := newGraph(setts, data, usedSeries, sub)
	gr.updateFunc = func(g *Graph) { g.Limits.Max = getMax() }

	return gr
}

// TODO: filter from config
func Disks(stats *app.Stats) []*Graph {
	disks := make([]*Graph, len(stats.Disks.Disks))
	for x, key := range slices.Sorted(maps.Keys(stats.Disks.Disks)) {
		disk := stats.Disks.Disks[key]
		disks[x] = Disk(disk)
	}

	return disks
}

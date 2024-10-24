package graph

import (
	"n4/gui-test/pkg/app"
	"n4/gui-test/pkg/series"
)

func SysMemAvailable(stats *app.Stats) *Graph {
	usedSeries := []*series.Entry{
		&stats.SysMem.Available,
		&stats.SysMem.Total,
	}
	sub := &series.Subscriber{}
	limit := stats.SysMem.Total.Subscribe(sub)
	data := stats.SysMem.Available.Subscribe(sub)

	setts := NewSettings("MemAvail", fmtCBMem)
	setts.configName = "sys_mem_available"
	setts.AutoMinMaxPadding = 0
	setts.Description = "System memory available"

	gr := newGraph(setts, data, usedSeries, sub)
	gr.updateFunc = func(g *Graph) { g.Limits.Max = limit.GetFirstValue() }

	return gr
}

func SysMemUsed(stats *app.Stats) *Graph {
	usedSeries := []*series.Entry{
		&stats.SysMem.Used,
		&stats.SysMem.Total,
	}
	sub := &series.Subscriber{}
	limit := stats.SysMem.Total.Subscribe(sub)
	data := stats.SysMem.Used.Subscribe(sub)

	setts := NewSettings("MemUsed", fmtCBMem)
	setts.configName = "sys_mem_used"
	setts.AutoMinMaxPadding = 0
	setts.Description = "System memory used"

	gr := newGraph(setts, data, usedSeries, sub)
	gr.updateFunc = func(g *Graph) { g.Limits.Max = limit.GetFirstValue() }

	return gr
}

func SysMemUsedPercent(stats *app.Stats) *Graph {
	usedSeries := []*series.Entry{
		&stats.SysMem.UsedPercent,
	}
	sub := &series.Subscriber{}
	data := stats.SysMem.UsedPercent.Subscribe(sub)

	setts := NewSettings("MemUsed%", fmtCBFloatMaker(2))
	setts.configName = "sys_mem_used_percent"
	setts.Limits = Limits{0, 100}
	setts.AutoMinMaxPadding = 0
	setts.Description = "System memory used(percent)"

	gr := newGraph(setts, data, usedSeries, sub)

	return gr
}

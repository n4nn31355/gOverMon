package graph

import (
	"n4/gui-test/pkg/app"
	"n4/gui-test/pkg/series"

	"github.com/dustin/go-humanize"
)

func newSelfMem(name string, label string, entry *series.Entry) *Graph {
	usedSeries := []*series.Entry{entry}
	sub := &series.Subscriber{}
	data := entry.Subscribe(sub)

	setts := NewSettings(label, fmtCBMem)
	setts.configName = "self_mem_" + name
	setts.Limits = Limits{0, humanize.MiByte}
	setts.Description = "Overlay process memory usage"

	gr := newGraph(setts, data, usedSeries, sub)

	return gr
}

func SelfMemRSS(stats *app.Stats) *Graph {
	return newSelfMem("rss", "RSS", &stats.SelfMem.RSS)
}

func SelfMemVMS(stats *app.Stats) *Graph {
	return newSelfMem("vms", "VMS", &stats.SelfMem.VMS)
}

func SelfMemHWM(stats *app.Stats) *Graph {
	return newSelfMem("hwm", "HWM", &stats.SelfMem.HWM)
}

func SelfMemData(stats *app.Stats) *Graph {
	return newSelfMem("data", "Data", &stats.SelfMem.Data)
}

func SelfMemStack(stats *app.Stats) *Graph {
	return newSelfMem("stack", "Stack", &stats.SelfMem.Stack)
}

func SelfMemLocked(stats *app.Stats) *Graph {
	return newSelfMem("locked", "Locked", &stats.SelfMem.Locked)
}

func SelfMemSwap(stats *app.Stats) *Graph {
	return newSelfMem("swap", "Swap", &stats.SelfMem.Swap)
}

func SelfRuntimeMemAlloc(stats *app.Stats) *Graph {
	return newSelfMem("runtime_alloc", "rMemAlloc", &stats.RuntimeMemAlloc)
}

func SelfRuntimeMemSys(stats *app.Stats) *Graph {
	return newSelfMem("runtime_sys", "rMemSys", &stats.RuntimeMemSys)
}

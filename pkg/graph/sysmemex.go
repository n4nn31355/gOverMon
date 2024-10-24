package graph

import (
	"n4/gui-test/pkg/app"
	"n4/gui-test/pkg/series"
)

func SysMemCommit(stats *app.Stats) *Graph {
	usedSeries := []*series.Entry{
		&stats.SysMemEx.CommitTotal,
		&stats.SysMemEx.CommitLimit,
	}
	sub := &series.Subscriber{}
	limit := stats.SysMemEx.CommitLimit.Subscribe(sub)
	data := stats.SysMemEx.CommitTotal.Subscribe(sub)

	setts := NewSettings("MemCommit", fmtCBMem)
	setts.configName = "sys_mem_commit"
	setts.AutoMinMaxPadding = 0
	setts.Description = "System memory commited"

	gr := newGraph(setts, data, usedSeries, sub)
	gr.updateFunc = func(g *Graph) { g.Limits.Max = limit.GetFirstValue() }

	return gr
}

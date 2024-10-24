package graph

import (
	"n4/gui-test/pkg/app"
	"n4/gui-test/pkg/series"
)

func SelfUpdate(stats *app.Stats) *Graph {
	entry := &stats.StatsUpdate
	usedSeries := []*series.Entry{entry}
	sub := &series.Subscriber{}
	data := entry.Subscribe(sub)

	setts := NewSettings("Update", fmtCBFloatMaker(4))
	setts.configName = "self_update"
	setts.Limits = Limits{0, 0.05}
	setts.Description = "Time in seconds spent to update stats"

	gr := newGraph(setts, data, usedSeries, sub)

	return gr
}

func SelfFramerate(stats *app.Stats) *Graph {
	entry := &stats.SelfFramerate
	usedSeries := []*series.Entry{entry}
	sub := &series.Subscriber{}
	data := entry.Subscribe(sub)

	setts := NewSettings("FPS", fmtCBFloatMaker(1))
	setts.configName = "self_framerate"
	setts.Limits = Limits{0, 10}
	setts.Description = "GUI Ticks per Seconds"

	gr := newGraph(setts, data, usedSeries, sub)

	return gr
}

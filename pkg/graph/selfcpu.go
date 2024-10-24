package graph

import (
	"n4/gui-test/pkg/app"
	"n4/gui-test/pkg/series"
)

func selfCPU(name string, label string, entry *series.Entry) *Graph {
	usedSeries := []*series.Entry{entry}
	sub := &series.Subscriber{}
	data := entry.Subscribe(sub)

	setts := NewSettings(label, fmtCBCPU)
	setts.configName = "self_cpu_" + name
	setts.Description = "Overlay process CPU usage"

	gr := newGraph(setts, data, usedSeries, sub)

	return gr
}

func SelfCPUSys(stats *app.Stats) *Graph {
	return selfCPU("sys", "CPU Sys", &stats.SelfCPU.Sys)
}

func SelfCPUUser(stats *app.Stats) *Graph {
	return selfCPU("user", "CPU User", &stats.SelfCPU.User)
}

func SelfCPUIdle(stats *app.Stats) *Graph {
	return selfCPU("idle", "CPU Idle", &stats.SelfCPU.Idle)
}

func SelfCPUIowait(stats *app.Stats) *Graph {
	return selfCPU("iowait", "CPU Iowait", &stats.SelfCPU.Iowait)
}

func SelfCPUSteal(stats *app.Stats) *Graph {
	return selfCPU("steal", "CPU Steal", &stats.SelfCPU.Steal)
}

func SelfCPUIrq(stats *app.Stats) *Graph {
	return selfCPU("irq", "CPU Irq", &stats.SelfCPU.Irq)
}

func SelfCPUGuest(stats *app.Stats) *Graph {
	return selfCPU("guest", "CPU Guest", &stats.SelfCPU.Guest)
}

func SelfCPUSoftirq(stats *app.Stats) *Graph {
	return selfCPU("softirq", "CPU Softirq", &stats.SelfCPU.Softirq)
}

func SelfCPUNice(stats *app.Stats) *Graph {
	return selfCPU("nice", "CPU Nice", &stats.SelfCPU.Nice)
}

func SelfCPUGuestNice(stats *app.Stats) *Graph {
	return selfCPU("guestnice", "CPU GuestNice", &stats.SelfCPU.GuestNice)
}

func SelfCPUPerc(stats *app.Stats) *Graph {
	entry := &stats.SelfCPU.Perc
	usedSeries := []*series.Entry{entry}
	sub := &series.Subscriber{}
	data := entry.Subscribe(sub)

	setts := NewSettings("CPU %", fmtCBCPU)
	setts.configName = "self_cpu_perc"
	setts.Limits = Limits{0, 100}
	setts.AutoMinMaxPadding = 0
	setts.Description = "Overlay process CPU usage percent"

	gr := newGraph(setts, data, usedSeries, sub)

	return gr
}

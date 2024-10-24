package graph

import (
	"n4/gui-test/pkg/series"
)

type Graph struct {
	*Settings

	data       *series.EntryData
	series     []*series.Entry
	subscriber *series.Subscriber

	// TODO: I'm not happy with this solution
	updateFunc func(g *Graph)
}

func newGraph(
	settings *Settings,
	data *series.EntryData,
	series []*series.Entry,
	subscriber *series.Subscriber,
) *Graph {
	return &Graph{
		Settings:   settings,
		data:       data,
		series:     series,
		subscriber: subscriber,
	}
}

// FIXME: register/unregister series usage
func (g *Graph) SetActive(active bool) {
	g.active = active
}

func (g *Graph) GetData() *series.EntryData {
	return g.data
}

func (g *Graph) Update() {
	if g.updateFunc != nil {
		g.updateFunc(g)
	}
}

type Collection []*Graph

func (gl Collection) ActiveNum() int {
	count := 0
	for _, g := range gl {
		if g.IsActive() {
			count++
		}
	}
	return count
}

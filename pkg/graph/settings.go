package graph

import (
	"n4/gui-test/pkg/plot"
)

type Limits struct {
	Min float64
	Max float64
}

type Settings struct {
	// TODO: use go generate to generate configName
	configName string

	active bool

	NameLabel   string
	Description string

	ValueLabelFormatCb plot.FormatCallback

	Limits            Limits
	AutoMinMaxPadding float64

	Flags []plot.Flag
}

func NewSettings(nameLabel string, formatCb plot.FormatCallback) *Settings {
	return &Settings{
		active: true,

		NameLabel: nameLabel,

		ValueLabelFormatCb: formatCb,

		Limits: Limits{0, 1},

		AutoMinMaxPadding: 0.1,
	}
}

func (s *Settings) IsActive() bool {
	return s.active
}

func (s *Settings) GetName() string {
	return s.configName
}

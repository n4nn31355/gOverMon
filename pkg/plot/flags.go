package plot

import "n4/gui-test/pkg/bitflags"

type Flag bitflags.BitFlag

const (
	FlagsNone Flag = 0

	FlagsAutoMin Flag = 1 << (iota - 1)
	FlagsAutoMax

	// Always keep auto min/max outside of initial min/max boundaries
	FlagsAutoKeepMinMax

	FlagsBorderTop
	FlagsBorderBottom
	FlagsBorderLeft
	FlagsBorderRight

	FlagsLabelsAll

	// Draw latest values from right side
	FlagsReverseOrder

	// Draw plot outside of set widget size
	FlagsDebugIgnoreCanvasBounds
)

const (
	FlagsAutoMinMax = FlagsAutoMin | FlagsAutoMax

	FlagsBorderTopBottom = FlagsBorderTop | FlagsBorderBottom
	FlagsBorderLeftRight = FlagsBorderLeft | FlagsBorderRight
	FlagsBorderAll       = FlagsBorderTopBottom | FlagsBorderLeftRight
)

package graph

import (
	"strconv"

	"github.com/dustin/go-humanize"
)

func fmtCBFloatMaker(precision int) func(value float64) string {
	return func(value float64) string {
		return strconv.FormatFloat(value, 'f', precision, 64)
	}
}

var fmtCBCPU = fmtCBFloatMaker(2)

func fmtCBMem(value float64) string {
	return humanize.IBytes(uint64(value))
}

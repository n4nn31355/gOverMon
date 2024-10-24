package plot

import (
	"fmt"
	"image"
	"math"

	"n4/gui-test/internal/utils"
	"n4/gui-test/pkg/bitflags"
)

type FormatCallback func(value float64) string

type WidgetData []float64

type Widget struct {
	Label string

	data WidgetData

	// TODO: handle zero Width, Height
	Width, Height     int
	xMin, xMax        float64
	autoMinMaxPadding float64

	// TODO: Implement threshold line and color
	// thresholds []Threshold

	barWidth, barSpacing int

	LabelPadding image.Point

	FormatCallback FormatCallback

	Flags Flag
}

func NewWidget(label string, data WidgetData) *Widget {
	return &Widget{
		Label:  label,
		data:   data,
		Width:  200,
		Height: 100,

		xMin: 0.0,
		xMax: 1.0,

		autoMinMaxPadding: 0.1,

		barWidth:   4,
		barSpacing: 2,

		LabelPadding: image.Pt(4, 4),

		FormatCallback: func(value float64) string {
			return fmt.Sprint(value)
		},

		Flags: FlagsAutoMinMax |
			FlagsBorderAll |
			FlagsLabelsAll |
			FlagsReverseOrder,
	}
}

func (w *Widget) SetSize(width, height int) *Widget {
	w.Width = width
	w.Height = height
	return w
}

func (w *Widget) SetBarSize(width, spacing int) *Widget {
	w.barWidth = width
	w.barSpacing = spacing
	return w
}

func (w *Widget) SetLimits(xMin, xMax float64) *Widget {
	w.xMin = xMin
	w.xMax = xMax
	return w
}

func (w *Widget) SetFlags(flags Flag, resetRest bool) *Widget {
	if resetRest {
		w.Flags = flags
	} else {
		w.Flags = bitflags.Set(w.Flags, flags)
	}
	return w
}

func (w *Widget) ClearFlags(flags Flag) *Widget {
	w.Flags = bitflags.Clear(w.Flags, flags)
	return w
}

func (w *Widget) SetFormatCallback(cb FormatCallback) *Widget {
	w.FormatCallback = cb
	return w
}

func (w *Widget) SetAutoHeightPadding(ratio float64) *Widget {
	w.autoMinMaxPadding = ratio
	return w
}

// TODO: remove GetData()?
func (w *Widget) GetData() WidgetData {
	return w.data
}

func (w *Widget) GetSanitizedMinMax() (xMin, xMax float64) {
	xMin, xMax = w.xMin, w.xMax

	if bitflags.Has(w.Flags, FlagsAutoMinMax) {
		xMin, xMax = utils.SliceMinMax(w.data)

		if xMin < 0.0 {
			xMin *= 1 + w.autoMinMaxPadding
		}
		if xMax > 0.0 {
			xMax *= 1 + w.autoMinMaxPadding
		}
		if bitflags.Has(w.Flags, FlagsAutoKeepMinMax) {
			xMin = min(xMin, w.xMin)
			xMax = max(xMax, w.xMax)
		}
		if xMin == xMax {
			xMax += 0.1
		}
	}

	return xMin, xMax
}

// Returns array {left, right, top, bottom}. Value may be nil
func (w *Widget) GetBorders() (borders [4]*image.Rectangle) {
	if bitflags.Has(w.Flags, FlagsBorderTop) {
		borders[0] = &image.Rectangle{
			image.Pt(0, 0),
			image.Pt(w.Width, 1),
		}
	}
	if bitflags.Has(w.Flags, FlagsBorderBottom) {
		borders[1] = &image.Rectangle{
			image.Pt(0, w.Height-1),
			image.Pt(w.Width, w.Height),
		}
	}
	if bitflags.Has(w.Flags, FlagsBorderLeft) {
		borders[2] = &image.Rectangle{
			image.Pt(0, 0),
			image.Pt(1, w.Height),
		}
	}
	if bitflags.Has(w.Flags, FlagsBorderRight) {
		borders[3] = &image.Rectangle{
			image.Pt(w.Width-1, 0),
			image.Pt(w.Width, w.Height),
		}
	}
	return borders
}

func (w *Widget) GetWidgetRect() image.Rectangle {
	return image.Rect(0, 0, w.Width, w.Height)
}

func (w *Widget) GetPlotMidLine() image.Rectangle {
	xMin, xMax := w.GetSanitizedMinMax()
	barRange := xMax - xMin
	fracSize := barRange / float64(w.Height-1)
	midPoint := float64(w.Height-1) - (barRange-xMax)/fracSize
	midPointInt := int(math.Round(midPoint))

	midline := image.Rectangle{
		Min: image.Pt(0, midPointInt),
		Max: image.Pt(w.Width, midPointInt+1),
	}

	// TODO: Flag to avoid drawing midline on border?
	if !bitflags.Has(w.Flags, FlagsDebugIgnoreCanvasBounds) {
		midline = utils.LimitRectangleToRectangle(midline, w.GetWidgetRect())
	}

	return midline
}

// TODO: Handle Inf/NaN xMin/xMax
// TODO: widget auto size
func (w *Widget) GetValueRect(x int) (rect image.Rectangle) {
	var val float64
	// TODO: make sure that flag FlagsReverseOrder handled everywhere where value accessed
	if bitflags.Has(w.Flags, FlagsReverseOrder) {
		val = w.data[len(w.data)-1-x]
	} else {
		val = w.data[x]
	}
	if math.IsNaN(val) || math.IsInf(val, 0) {
		return rect
	}
	xMin, xMax := w.GetSanitizedMinMax()
	barRange := xMax - xMin
	fracSize := barRange / float64(w.Height-1)
	midPoint := float64(w.Height-1) - (barRange-xMax)/fracSize
	midPointInt := int(math.Round(midPoint))

	barHeight := int(math.Round(val / fracSize))
	// TODO: Flag for zero value bar appearance? Also draw it only on midline
	// if barHeight == 0 {
	// 	barHeight = 1
	// 	barColor = c.barStyle.zeroColor
	// }
	pOffset := x * (w.barWidth + w.barSpacing)
	midPointOffset := 0
	// TODO: Flag for drawing on top of midline?
	if barHeight < 0.0 {
		midPointOffset = 1
	}
	rect.Min.X = pOffset
	rect.Max.X = pOffset + w.barWidth
	rect.Min.Y = midPointInt + midPointOffset
	rect.Max.Y = midPointInt - barHeight + midPointOffset

	if !bitflags.Has(w.Flags, FlagsDebugIgnoreCanvasBounds) {
		rect = utils.LimitRectangleToRectangle(rect.Canon(), w.GetWidgetRect())
	}

	return rect
}

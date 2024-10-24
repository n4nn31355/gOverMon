package ebiten

import (
	"bytes"
	"image"
	"image/color"
	"log"

	"n4/gui-test/pkg/bitflags"
	"n4/gui-test/pkg/plot"
	"n4/gui-test/resources/fonts"

	"github.com/ebitengine/microui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var fontFace *text.GoTextFace

func init() {
	// TODO: Implement setting of font face and size in config
	// source, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.SonoPropBold))
	// source, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.SonoPropMedium))
	// source, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.SonoPropRegular))
	// source, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.SonoPropSemiBold))
	source, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.UbuntuSansMedium))
	// source, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.UbuntuSansRegular))
	if err != nil {
		log.Fatal(err)
	}
	fontFace = &text.GoTextFace{
		Source: source,
		Size:   10,
	}
}

func textWidth(str string) int {
	return int(text.Advance(str, fontFace))
}

func lineHeight() int {
	return int(fontFace.Metrics().HAscent + fontFace.Metrics().HDescent + fontFace.Metrics().HLineGap)
}

func (g *Game) DrawPlot(widget *plot.Widget) {
	g.ctx.Control(0, 0, func(r image.Rectangle) microui.Res {
		widget.SetSize(r.Dx(), r.Dy())

		for _, border := range widget.GetBorders() {
			if border == nil {
				continue
			}
			rectGlobal := image.Rectangle{
				r.Min.Add(border.Min),
				r.Min.Add(border.Max),
			}
			g.ctx.DrawControl(func(screen *ebiten.Image) {
				vector.DrawFilledRect(
					screen,
					float32(rectGlobal.Min.X),
					float32(rectGlobal.Min.Y),
					float32(rectGlobal.Dx()),
					float32(rectGlobal.Dy()),
					g.cfg.App.Theme.Plot.Border,
					false)
			})
		}

		{
			midLineRect := widget.GetPlotMidLine()
			if !midLineRect.Empty() {
				rectGlobal := image.Rectangle{
					r.Min.Add(midLineRect.Min),
					r.Min.Add(midLineRect.Max),
				}

				g.ctx.DrawControl(func(screen *ebiten.Image) {
					vector.DrawFilledRect(
						screen,
						float32(rectGlobal.Min.X),
						float32(rectGlobal.Min.Y),
						float32(rectGlobal.Dx()),
						float32(rectGlobal.Dy()),
						g.cfg.App.Theme.Plot.Midline,
						false)
				})
			}
		}

		for x := range widget.GetData() {
			barRect := widget.GetValueRect(x).Canon()
			if !barRect.Empty() {
				rectGlobal := barRect.Add(r.Min)
				g.ctx.DrawControl(func(screen *ebiten.Image) {
					vector.DrawFilledRect(
						screen,
						float32(rectGlobal.Min.X),
						float32(rectGlobal.Min.Y),
						float32(rectGlobal.Dx()),
						float32(rectGlobal.Dy()),
						g.cfg.App.Theme.Plot.Bar,
						false)
				})
			}
		}

		if bitflags.Has(widget.Flags, plot.FlagsLabelsAll) {
			xMin, xMax := widget.GetSanitizedMinMax()

			labels := []struct {
				text string
				pos  image.Point
			}{
				{
					text: widget.Label,
					pos:  r.Min.Add(widget.LabelPadding),
				},
				{
					text: widget.FormatCallback(widget.GetData()[0]),
					pos: r.Min.Add(image.Pt(
						widget.LabelPadding.X,
						widget.LabelPadding.Y+lineHeight(),
					)),
				},
				{
					text: widget.FormatCallback(xMax),
					pos: r.Min.Add(image.Pt(
						widget.Width-textWidth(widget.FormatCallback(xMax))-
							widget.LabelPadding.X,
						widget.LabelPadding.Y,
					)),
				},
				{
					text: widget.FormatCallback(xMin),
					pos: r.Min.Add(image.Pt(
						widget.Width-textWidth(widget.FormatCallback(xMin))-widget.LabelPadding.X,
						widget.Height-lineHeight()-widget.LabelPadding.Y,
					)),
				},
			}

			for _, label := range labels {
				labelRect := image.Rectangle{
					label.pos, image.Point{textWidth(label.text), lineHeight()},
				}
				g.ctx.DrawControl(func(screen *ebiten.Image) {
					op := &text.DrawOptions{}
					op.GeoM.Translate(float64(label.pos.X), float64(label.pos.Y))
					op.ColorScale.ScaleWithColor(g.cfg.App.Theme.Plot.LabelText)
					text.Draw(screen, label.text, fontFace, op)

					drawFilledRect(screen, labelRect, g.cfg.App.Theme.Plot.LabelBackground)
				})
			}
		}

		return 0
	})
}

func drawFilledRect(screen *ebiten.Image, rect image.Rectangle, color color.RGBA) {
	vector.DrawFilledRect(
		screen,
		float32(rect.Min.X),
		float32(rect.Min.Y),
		float32(rect.Max.X),
		float32(rect.Max.Y),
		color,
		false)
}

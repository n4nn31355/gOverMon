package ebiten

import (
	"image"
	"slices"

	"n4/gui-test/pkg/plot"

	"github.com/ebitengine/microui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func (g *Game) intSlider(fvalue *float64, value *int, low, high int) microui.Res {
	*fvalue = float64(*value)
	res := g.ctx.SliderEx(fvalue, float64(low), float64(high), 0, "%.0f", microui.OptAlignCenter)
	*value = int(*fvalue)
	return res
}

func (g *Game) byteSlider(fvalue *float64, value *byte, low, high byte) microui.Res {
	*fvalue = float64(*value)
	res := g.ctx.SliderEx(fvalue, float64(low), float64(high), 0, "%.0f", microui.OptAlignCenter)
	*value = byte(*fvalue)
	return res
}

var (
	fcolors = [14]struct {
		R, G, B, A float64
	}{}
	colors = []struct {
		Label   string
		ColorID int
	}{
		{"text:", microui.ColorText},
		{"border:", microui.ColorBorder},
		{"windowbg:", microui.ColorWindowBG},
		{"titlebg:", microui.ColorTitleBG},
		{"titletext:", microui.ColorTitleText},
		{"panelbg:", microui.ColorPanelBG},
		{"button:", microui.ColorButton},
		{"buttonhover:", microui.ColorButtonHover},
		{"buttonfocus:", microui.ColorButtonFocus},
		{"base:", microui.ColorBase},
		{"basehover:", microui.ColorBaseHover},
		{"basefocus:", microui.ColorBaseFocus},
		{"scrollbase:", microui.ColorScrollBase},
		{"scrollthumb:", microui.ColorScrollThumb},
	}

	barWidth   float64
	barSpacing float64
	plotHeight float64
)

func (g *Game) styleWindow(size image.Rectangle) {
	g.ctx.Window("Style Editor", size, func(res microui.Res) {
		sw := int(float64(g.ctx.CurrentContainer().Body.Dx()) * 0.14)
		g.ctx.SetLayoutRow([]int{80, sw, sw, sw, sw, -1}, 0)
		for _, c := range colors {
			g.ctx.Label(c.Label)
			g.byteSlider(&fcolors[c.ColorID].R, &g.ctx.Style.Colors[c.ColorID].R, 0, 255)
			g.byteSlider(&fcolors[c.ColorID].G, &g.ctx.Style.Colors[c.ColorID].G, 0, 255)
			g.byteSlider(&fcolors[c.ColorID].B, &g.ctx.Style.Colors[c.ColorID].B, 0, 255)
			g.byteSlider(&fcolors[c.ColorID].A, &g.ctx.Style.Colors[c.ColorID].A, 0, 255)
			g.ctx.Control(0, 0, func(r image.Rectangle) microui.Res {
				clr := g.ctx.Style.Colors[c.ColorID]
				g.ctx.DrawControl(func(target *ebiten.Image) {
					vector.DrawFilledRect(
						target,
						float32(r.Min.X),
						float32(r.Min.Y),
						float32(r.Dx()),
						float32(r.Dy()),
						clr,
						false)
				})
				return 0
			})
		}
	})
}

func (g *Game) drawSettings() {
	g.ctx.LayoutColumn(func() {
		g.ctx.SetLayoutRow([]int{settingsBtnWidth, settingsDescriptionWidth}, 0)

		g.ctx.Label("Bar Width")
		res := g.intSlider(&barWidth, &g.cfg.App.BarWidth, 1, 6)
		if res == microui.ResChange {
			g.cfg.Save()
			g.updateSize()
		}
		g.ctx.Label("Bar Spacing")
		res = g.intSlider(&barSpacing, &g.cfg.App.BarSpacing, 0, 6)
		if res == microui.ResChange {
			g.cfg.Save()
			g.updateSize()
		}
		g.ctx.Label("Plot Height")
		res = g.intSlider(&plotHeight, &g.cfg.App.PlotHeight, 30, 200)
		if res == microui.ResChange {
			g.cfg.Save()
			g.updateSize()
		}

		g.ctx.SetLayoutRow(slices.Repeat(
			[]int{settingsBtnWidth, settingsDescriptionWidth},
			settingsBtnNumInRow,
		), 0)

		for _, gr := range g.graphs {
			cfgName := gr.GetName()
			grCfg, present := g.cfg.App.GraphSettings[cfgName]
			if !present {
				continue
			}

			active := gr.IsActive()

			var btnText string
			if active {
				btnText = "on"
			} else {
				btnText = "off"
			}

			btnRes := g.ctx.ButtonEx(
				gr.NameLabel+": "+btnText, 0, microui.OptAutoSize,
			)
			if btnRes != 0 {
				grCfg.Enabled = !active
				gr.SetActive(!active)
				g.cfg.Save()
				g.updateSize()
			}

			g.ctx.Label(gr.Description)
		}
	})
}

func (g *Game) settingsWindow() {
	flags := microui.OptNoResize |
		microui.OptNoTitle

	rect := image.Rect(
		settingsWindowPadding,
		settingsWindowPadding,
		g.width-g.ctx.Style.Padding,
		settingsWindowHeight,
	)
	g.ctx.WindowEx("Settings", rect, flags, func(_ microui.Res) {
		// g.ctx.BringToFront(g.ctx.CurrentContainer())
		g.drawSettings()
	})
}

func (g *Game) mainWindow() {
	flags := microui.OptNoTitle |
		microui.OptNoFrame |
		microui.OptAutoSize |
		microui.OptNoScroll |
		microui.OptNoResize |
		microui.OptNoInteract
	rect := image.Rect(0, 0, g.width, g.height)
	g.ctx.WindowEx("Main", rect, flags, func(_ microui.Res) {
		g.ctx.SetLayoutRow([]int{0, -1}, 0)

		g.ctx.LayoutColumn(func() {
			if !ebiten.IsWindowMousePassthrough() {
				g.ctx.SetLayoutRow([]int{-1, 22}, titleHeight)
				g.ctx.Label("Stats")

				if g.ctx.Button("cfg") != 0 {
					g.toggleSettings()
				}
			}

			tRange := g.cfg.App.TimeRangeSeconds
			g.ctx.SetLayoutRow([]int{
				tRange*g.cfg.App.BarWidth + (tRange-1)*g.cfg.App.BarSpacing,
			}, g.cfg.App.PlotHeight)

			for _, graph := range g.graphs {
				if !graph.IsActive() {
					continue
				}
				graph.Update()
				data := graph.GetData().GetValues()
				plotWidget := plot.NewWidget(graph.NameLabel, data).
					// SetSize(r.Dx(), r.Dy()).
					SetLimits(graph.Limits.Min, graph.Limits.Max).
					SetAutoHeightPadding(graph.AutoMinMaxPadding).
					SetBarSize(g.cfg.App.BarWidth, g.cfg.App.BarSpacing).
					SetFlags(
						plot.FlagsDebugIgnoreCanvasBounds|
							plot.FlagsAutoKeepMinMax,
						false).
					SetFormatCallback(graph.ValueLabelFormatCb)
				g.DrawPlot(plotWidget)
			}
		})
	})
}

func (g *Game) ProcessFrame() {
	g.ctx.Update(func() {
		g.mainWindow()
		if g.isSettingsShown() {
			g.settingsWindow()
		}
	})
}

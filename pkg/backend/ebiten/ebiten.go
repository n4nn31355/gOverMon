package ebiten

import (
	"log"
	"time"

	"n4/gui-test/pkg/app"
	"n4/gui-test/pkg/config"
	"n4/gui-test/pkg/graph"

	"github.com/ebitengine/microui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"go.uber.org/zap"
)

const (
	TPSHigh = ebiten.DefaultTPS // Suitable for interaction
	TPSLow  = 1                 // Suitable for power saving

	dragButton = ebiten.MouseButtonRight
	// dragButton = ebiten.MouseButtonLeft

	// FIXME: remove hardcode
	titleHeight = 20

	settingsWindowPadding = 32
	settingsWindowHeight  = 250
	settingsBtnWidth      = 120
	settingsBtnNumInRow   = 1

	settingsDescriptionWidth = 240
)

type Game struct {
	ctx *microui.Context
	cfg *config.Config

	width, height int

	// TODO: uncouple Game from app.Stats?
	stats        *app.Stats
	statsUpdated time.Time

	graphs graph.Collection

	input      bool
	introShown bool
	close      bool

	showSettings bool

	dragging         bool
	dragStartWindowX int
	dragStartWindowY int
	dragStartCursorX int
	dragStartCursorY int

	cursorToWindowX float64
	cursorToWindowY float64
}

func (g *Game) handleDrag() {
	if inpututil.IsMouseButtonJustReleased(dragButton) {
		g.dragging = false
		winX, winY := ebiten.WindowPosition()
		g.cfg.App.Position.X = winX
		g.cfg.App.Position.Y = winY
		err := g.cfg.Save()
		if err != nil {
			// TODO: better config save error handling
			log.Println(err)
		}
	}
	if !g.dragging && inpututil.IsMouseButtonJustPressed(dragButton) {
		g.dragging = true
		g.dragStartWindowX, g.dragStartWindowY = ebiten.WindowPosition()
		g.dragStartCursorX, g.dragStartCursorY = ebiten.CursorPosition()
	}
	if g.dragging {
		curX, curY := ebiten.CursorPosition()
		winX, winY := ebiten.WindowPosition()
		distX := int(float64(curX-g.dragStartCursorX) * g.cursorToWindowX)
		distY := int(float64(curY-g.dragStartCursorY) * g.cursorToWindowY)
		// TODO: handle screen bounds
		// TODO: handle resolution change
		ebiten.SetWindowPosition(winX+distX, winY+distY)
	}
}

// TODO: Don't depend on (g *Game)?
func (g *Game) getSettingsWindowWidth() int {
	return settingsWindowPadding +
		g.ctx.Style.Padding*3 +
		g.ctx.Style.ScrollbarSize +
		settingsBtnWidth + g.ctx.Style.Padding + settingsDescriptionWidth
}

// TODO: Don't depend on (g *Game)?
func (g *Game) getSettingsWindowHeight() int {
	return settingsWindowPadding + settingsWindowHeight
}

// TODO: Don't depend on (g *Game)?
func (g *Game) getPlotColumnWidth() int {
	timeRange := g.cfg.App.TimeRangeSeconds
	return g.ctx.Style.Padding*2 +
		timeRange*g.cfg.App.BarWidth +
		// TODO: probably should be ticks - 1
		timeRange*g.cfg.App.BarSpacing
}

// TODO: Don't depend on (g *Game)?
func (g *Game) getPlotColumnHeight() int {
	plotsCount := g.graphs.ActiveNum()

	titleOffset := 0
	if !ebiten.IsWindowMousePassthrough() {
		titleOffset += g.ctx.Style.Spacing + titleHeight
	}

	return g.ctx.Style.Spacing +
		titleOffset +
		(g.ctx.Style.Spacing+g.cfg.App.PlotHeight)*plotsCount
}

func (g *Game) setPassthrough(enable bool) {
	ebiten.SetWindowMousePassthrough(enable)
	if enable {
		ebiten.SetTPS(TPSLow)
	} else {
		ebiten.SetTPS(TPSHigh)
	}
	g.updateSize()
}

func (g *Game) togglePassthrough() {
	g.setPassthrough(!ebiten.IsWindowMousePassthrough())
}

func (g *Game) toggleSettings() {
	g.showSettings = !g.showSettings
	g.updateSize()
}

func (g *Game) isSettingsShown() bool {
	return g.showSettings && !ebiten.IsWindowMousePassthrough()
}

func (g *Game) updateSize() {
	if g.isSettingsShown() {
		g.width = max(g.getPlotColumnWidth(), g.getSettingsWindowWidth())
		g.height = max(g.getPlotColumnHeight(), g.getSettingsWindowHeight())
	} else {
		g.width = g.getPlotColumnWidth()
		g.height = g.getPlotColumnHeight()
	}
	ebiten.SetWindowSize(g.width, g.height)
}

func (g *Game) Update() error {
	if g.close {
		return ebiten.Termination
	}

	g.input = len(inpututil.AppendPressedKeys(nil)) > 0 ||
		ebiten.IsMouseButtonPressed(dragButton) ||
		// NOTE: MouseButtonLeft helps to avoid flickering when closing settings
		ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)

	g.handleDrag()

	g.ProcessFrame()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	draw := g.input ||
		g.isSettingsShown() ||
		!g.introShown ||
		g.stats.Updated.After(g.statsUpdated)
	if !draw {
		return
	}

	screen.Clear()

	if ebiten.IsWindowMousePassthrough() {
		screen.Fill(g.cfg.App.Theme.Window.Passthrough.Background)
	} else {
		screen.Fill(g.cfg.App.Theme.Window.Active.Background)
	}

	g.ctx.Draw(screen)

	g.introShown = true
	g.statsUpdated = g.stats.Updated
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	// The cursor position is in a "logical" coordinate, which is determined by
	// the game width and height.
	// Calculate the factors to convert a cursor position to a window position.
	g.cursorToWindowX = float64(outsideWidth) / float64(g.width)
	g.cursorToWindowY = float64(outsideHeight) / float64(g.height)
	return g.width, g.height
}

func (g *Game) Close() {
	g.close = true
}

func Window(
	cfg *config.Config,
	graphs graph.Collection,
	stats *app.Stats,
	togglePassthrough <-chan struct{},
	exit <-chan struct{},
) {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

	ebiten.SetWindowTitle("gOverMon")
	ebiten.SetWindowDecorated(false)
	ebiten.SetWindowFloating(true)
	ebiten.SetScreenClearedEveryFrame(false)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)
	ebiten.SetWindowPosition(cfg.App.Position.X, cfg.App.Position.Y)
	ebiten.SetVsyncEnabled(false)

	game := &Game{
		ctx: microui.NewContext(),
		cfg: cfg,

		stats:  stats,
		graphs: graphs,
	}
	game.ctx.Style.Padding = 2
	game.ctx.Style.Spacing = 2

	game.setPassthrough(false)

	game.width = game.getPlotColumnWidth()
	game.height = game.getPlotColumnHeight()

	ebiten.SetWindowSize(game.width, game.height)

	go func() {
		<-exit
		logger.Info("We're exitting...")
		game.Close()
	}()
	go func() {
		for {
			<-togglePassthrough
			game.togglePassthrough()
		}
	}()

	opts := &ebiten.RunGameOptions{
		InitUnfocused:     true,
		ScreenTransparent: true,
		SkipTaskbar:       true,
	}

	err = ebiten.RunGameWithOptions(game, opts)
	if err != nil {
		logger.Fatal("err: ", zap.Error(err))
	}
}

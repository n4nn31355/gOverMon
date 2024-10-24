package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"slices"
	"syscall"
	"time"

	"n4/gui-test/pkg/app"
	"n4/gui-test/pkg/backend/ebiten"
	"n4/gui-test/pkg/config"
	"n4/gui-test/pkg/graph"

	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

const (
	appName = "govermon"
)

var (
	exit              = make(chan struct{})
	togglePassthrough = make(chan struct{})
	stopUpdates       = make(chan struct{})
)

type flagVars struct {
	cfgPath         string
	useDebugHotkeys bool
}

func handleFlags() (fVars flagVars) {
	flagSet := pflag.NewFlagSet("app", pflag.ExitOnError)
	flagSet.Usage = func() {
		fmt.Println(flagSet.FlagUsages())
		os.Exit(0)
	}
	fVars.cfgPath = ""
	defaultPath, _ := config.GetDefaultPath()
	// NOTE: actual default is "" because we create dirs only for path == ""
	flagSet.StringVarP(
		&fVars.cfgPath, "config", "c", "",
		fmt.Sprintf("path to .yaml config file (default %s)", defaultPath),
	)
	flagSet.BoolVar(
		&fVars.useDebugHotkeys, "debug-hotkeys", false,
		"debug: use debug set of hotkeys",
	)
	flagSet.Parse(os.Args[1:])
	return fVars
}

func loadConfig(logger *zap.Logger, cfgPath string) *config.Config {
	var err error
	if cfgPath == "" {
		cfgPath, err = config.InitDefaultFile()
		if err != nil {
			logger.Fatal("failed to init config", zap.Error(err))
		}
	}
	logger.Info("Loading config file", zap.String("path", cfgPath))

	cfg := config.NewConfig(cfgPath)
	err = cfg.Load()
	if err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
	}
	return cfg
}

func handleSignals(logger *zap.Logger) {
	sigInt := make(chan os.Signal, 1)
	signal.Notify(sigInt, os.Interrupt, syscall.SIGTERM)
	go func() {
		s := <-sigInt
		logger.Info("Received signal, exitting", zap.String("signal", s.String()))
		close(exit)
	}()
}

func handleExit() {
	go func() {
		<-exit
		close(stopUpdates)
	}()
}

func handleUpdates(logger *zap.Logger, rateSeconds int, cb func()) {
	go func() {
		for {
			ctx, cancel := context.WithTimeout(
				context.Background(), time.Duration(rateSeconds)*time.Second,
			)
			defer cancel()

			select {
			case <-stopUpdates:
				logger.Info("We're done here")
				return
			// case <-time.After(rateSeconds * time.Second):
			case <-ctx.Done():
				cb()
				cancel()
			}
		}
	}()
}

func main() {
	fVars := handleFlags()

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

	cfg := loadConfig(logger, fVars.cfgPath)

	if cfg.App.Debug {
		initPyroscope()

		// mux := http.NewServeMux()
		// statsviz.Register(mux)

		go func() {
			// log.Println(http.ListenAndServe("localhost:8080", mux))
			logger.Info(
				"ListenAndServe",
				zap.Error(http.ListenAndServe("localhost:8080", nil)),
			)
		}()
	}

	go registerHotkeys(fVars.useDebugHotkeys)

	stats := app.NewStats(cfg.App.TimeRangeSeconds)
	// TODO: Is there a better solution for collectors of dynamic instances
	stats.Disks.Discover() // NOTE: Prefetch available disks to use it in graph init

	handleExit()
	handleSignals(logger)
	handleUpdates(logger, cfg.App.UpdateRateSeconds, stats.Update)

	graphs := graph.Collection{
		graph.SelfUpdate(stats),
		graph.SelfFramerate(stats),
		graph.SelfRuntimeMemAlloc(stats),
		graph.SelfRuntimeMemSys(stats),

		graph.SelfMemRSS(stats),
		graph.SelfMemVMS(stats),
		graph.SelfMemHWM(stats),
		graph.SelfMemData(stats),
		graph.SelfMemStack(stats),
		graph.SelfMemLocked(stats),
		graph.SelfMemSwap(stats),

		graph.SelfCPUSys(stats),
		graph.SelfCPUUser(stats),
		graph.SelfCPUIdle(stats),
		graph.SelfCPUIowait(stats),
		graph.SelfCPUSteal(stats),
		graph.SelfCPUIrq(stats),
		graph.SelfCPUGuest(stats),
		graph.SelfCPUSoftirq(stats),
		graph.SelfCPUNice(stats),
		graph.SelfCPUGuestNice(stats),

		graph.SelfCPUPerc(stats),

		graph.SysMemCommit(stats),
		graph.SysMemAvailable(stats),
		graph.SysMemUsed(stats),
		graph.SysMemUsedPercent(stats),
	}

	for _, graph := range graphs {
		settingName := graph.GetName()
		if settingName == "" {
			logger.Fatal(
				"graph settings name is empty",
				zap.String("graph", graph.NameLabel),
			)
		}
		settings, present := cfg.App.GraphSettings[settingName]
		if present {
			graph.SetActive(settings.Enabled)
		} else {
			cfg.App.GraphSettings[settingName] = &config.GraphSettings{
				Enabled: graph.IsActive(),
			}
		}
	}
	cfg.Save()

	graphs = slices.Concat(graphs, graph.Disks(stats))

	ebiten.Window(cfg, graphs, stats, togglePassthrough, exit)

	logger.Info("Aloha!")
}

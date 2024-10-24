package main

import (
	"log"

	"go.uber.org/zap"
	"golang.design/x/hotkey"
)

type hotkeyDefinition struct {
	mods []hotkey.Modifier
	key  hotkey.Key
}

type hotkeyData struct {
	name    string
	hotkey  *hotkey.Hotkey
	handler func(hk *hotkey.Hotkey, logger *zap.Logger)
}

func registerHotkeys(useDebugSet bool) {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

	var hkPassthru *hotkey.Hotkey
	var hkExit *hotkey.Hotkey
	if useDebugSet {
		hkPassthru = hotkey.New(hotkeyTogglePassthroughDebug.mods, hotkeyTogglePassthroughDebug.key)
		hkExit = hotkey.New(hotkeyExitDebug.mods, hotkeyExitDebug.key)
	} else {
		hkPassthru = hotkey.New(hotkeyTogglePassthrough.mods, hotkeyTogglePassthrough.key)
		hkExit = hotkey.New(hotkeyExit.mods, hotkeyExit.key)
	}

	hkeys := []hotkeyData{
		{
			"TogglePassthrough", hkPassthru,
			func(hk *hotkey.Hotkey, logger *zap.Logger) {
				for range hk.Keyup() {
					logger.Info("hotkey event", zap.String("event_name", "Keyup"))
					togglePassthrough <- struct{}{}
				}
			},
		},
		{
			"Exit", hkExit,
			func(hk *hotkey.Hotkey, logger *zap.Logger) {
				for range hk.Keyup() {
					logger.Info("hotkey event", zap.String("event_name", "Keyup"))
					close(exit)
				}
			},
		},
	}

	for _, hkData := range hkeys {
		hkLogger := logger.With(
			zap.String("name", hkData.name),
			zap.String("hk", hkData.hotkey.String()),
		)
		err = hkData.hotkey.Register()
		if err != nil {
			hkLogger.Fatal("failed to register hotkey", zap.Error(err))
			return
		}

		hkLogger.Info("hotkey is registered")
		go hkData.handler(hkData.hotkey, hkLogger)
	}
}

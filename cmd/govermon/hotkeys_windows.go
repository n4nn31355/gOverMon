package main

import "golang.design/x/hotkey"

// TODO: use config for hotkeys

var hotkeyTogglePassthrough = hotkeyDefinition{
	[]hotkey.Modifier{hotkey.ModWin, hotkey.ModShift}, hotkey.KeyO,
}

var hotkeyTogglePassthroughDebug = hotkeyDefinition{
	[]hotkey.Modifier{hotkey.ModWin, hotkey.ModShift, hotkey.ModAlt}, hotkey.KeyO,
}

var hotkeyExit = hotkeyDefinition{
	[]hotkey.Modifier{hotkey.ModWin, hotkey.ModShift}, hotkey.KeyI,
}

var hotkeyExitDebug = hotkeyDefinition{
	[]hotkey.Modifier{hotkey.ModWin, hotkey.ModShift, hotkey.ModAlt}, hotkey.KeyI,
}

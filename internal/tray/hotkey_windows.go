//go:build windows

package tray

import "golang.design/x/hotkey"

// Ctrl+Shift+B on Windows
var hotkeyModifiers = []hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}

const hotkeyLabel = "Ctrl+Shift+B"

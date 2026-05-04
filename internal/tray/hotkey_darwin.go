//go:build darwin

package tray

import "golang.design/x/hotkey"

// Cmd+Shift+B on macOS
var hotkeyModifiers = []hotkey.Modifier{hotkey.ModCmd, hotkey.ModShift}

const hotkeyLabel = "⌘⇧B"

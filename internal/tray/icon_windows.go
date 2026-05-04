//go:build windows

package tray

import "github.com/getlantern/systray"

func setIcon(icon []byte) {
	// Windows requires ICO format; fall back to title until an ICO asset is provided.
	// To use a proper icon: convert tray_icon.png to tray_icon.ico and call
	// systray.SetIcon(icon) with the ICO bytes instead.
	systray.SetTitle("ZB")
}

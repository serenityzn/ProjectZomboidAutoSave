package main

import (
	_ "embed"

	"github.com/getlantern/systray"
	"zomboidautobackup/internal/tray"
)

//go:embed assets/tray_icon.png
var trayIcon []byte

func main() {
	t := tray.New()
	systray.Run(func() { t.Setup(trayIcon) }, func() {})
}

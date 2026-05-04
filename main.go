package main

import (
	"github.com/getlantern/systray"
	"zomboidautobackup/internal/tray"
)

func main() {
	t := tray.New()
	systray.Run(t.Setup, func() {})
}

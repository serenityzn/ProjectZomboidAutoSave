//go:build darwin

package tray

import "github.com/getlantern/systray"

func setIcon(icon []byte) {
	systray.SetTemplateIcon(icon, icon)
}

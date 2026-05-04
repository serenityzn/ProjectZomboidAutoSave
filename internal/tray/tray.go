package tray

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/getlantern/systray"
	"zomboidautobackup/internal/backup"
	"zomboidautobackup/internal/config"
)

var intervalPresets = []int{5, 10, 15, 20, 25, 30, 35, 40, 45, 50, 55, 60}
var maxBackupPresets = []int{5, 10, 15}

const maxBackupCap = 20

type Tray struct {
	cfg *config.Config

	autoBackupItem *systray.MenuItem

	zomboidDisplay *systray.MenuItem
	zomboidChange  *systray.MenuItem
	backupDisplay  *systray.MenuItem
	backupChange   *systray.MenuItem

	intervalParent *systray.MenuItem
	intervalItems  []*systray.MenuItem
	intervalSet    *systray.MenuItem

	maxBackupParent *systray.MenuItem
	maxBackupItems  []*systray.MenuItem
	maxBackupSet    *systray.MenuItem

	backupBeforeItem *systray.MenuItem

	// Pre-allocated restore slots (Show/Hide at runtime)
	restoreManualEmpty *systray.MenuItem
	restoreAutoEmpty   *systray.MenuItem
	restoreManualItems []*systray.MenuItem
	restoreAutoItems   []*systray.MenuItem
	restoreManualSnaps []string
	restoreAutoSnaps   []string
}

func New() *Tray {
	cfg, _ := config.Load()
	cfg.EnsureBackupFolder() //nolint:errcheck
	cfg.Save()               //nolint:errcheck
	return &Tray{cfg: cfg}
}

func (t *Tray) Setup() {
	systray.SetTitle("ZB")
	systray.SetTooltip("Project Zomboid Auto Backup")

	t.buildSettingsSubmenu()

	systray.AddSeparator()

	autoLabel := "Auto Backup: OFF"
	if t.cfg.AutoBackup {
		autoLabel = "Auto Backup: ON"
	}
	t.autoBackupItem = systray.AddMenuItemCheckbox(autoLabel, "Toggle automatic backups", t.cfg.AutoBackup)
	manualItem := systray.AddMenuItem("Manual Backup", "Trigger a backup now")

	t.buildRestoreSubmenu()

	systray.AddSeparator()

	quitItem := systray.AddMenuItem("Quit", "Quit ZomboidAutoBackup")

	go t.handleEvents(manualItem, quitItem)
	go t.autoBackupLoop()
}

func (t *Tray) buildSettingsSubmenu() {
	settingsItem := systray.AddMenuItem("Settings", "Configure backup settings")

	t.zomboidDisplay = settingsItem.AddSubMenuItem(
		fmt.Sprintf("Zomboid Folder: %s", t.cfg.ZomboidFolder), "")
	t.zomboidDisplay.Disable()
	t.zomboidChange = settingsItem.AddSubMenuItem("Change Zomboid Folder...", "")

	sep1 := settingsItem.AddSubMenuItem("───────────────────", "")
	sep1.Disable()

	t.backupDisplay = settingsItem.AddSubMenuItem(
		fmt.Sprintf("Backup Folder: %s", t.cfg.BackupFolder), "")
	t.backupDisplay.Disable()
	t.backupChange = settingsItem.AddSubMenuItem("Change Backup Folder...", "")

	sep2 := settingsItem.AddSubMenuItem("───────────────────", "")
	sep2.Disable()

	t.intervalParent = settingsItem.AddSubMenuItem(
		fmt.Sprintf("Backup Every: %d min", t.cfg.BackupInterval), "")
	for _, v := range intervalPresets {
		item := t.intervalParent.AddSubMenuItem(fmt.Sprintf("%d min", v), "")
		if v == t.cfg.BackupInterval {
			item.Check()
		}
		t.intervalItems = append(t.intervalItems, item)
	}
	t.intervalSet = t.intervalParent.AddSubMenuItem("Set...", "Enter a custom interval in minutes")

	sep3 := settingsItem.AddSubMenuItem("───────────────────", "")
	sep3.Disable()

	t.maxBackupParent = settingsItem.AddSubMenuItem(
		fmt.Sprintf("Max Backups: %d", t.cfg.MaxBackupFiles), "")
	for _, v := range maxBackupPresets {
		item := t.maxBackupParent.AddSubMenuItem(fmt.Sprintf("%d", v), "")
		if v == t.cfg.MaxBackupFiles {
			item.Check()
		}
		t.maxBackupItems = append(t.maxBackupItems, item)
	}
	t.maxBackupSet = t.maxBackupParent.AddSubMenuItem("Set...  (max 20)", "Enter a custom max backup count")
}

func (t *Tray) buildRestoreSubmenu() {
	restoreItem := systray.AddMenuItem("Restore", "Restore a saved backup")

	t.backupBeforeItem = restoreItem.AddSubMenuItemCheckbox(
		"Backup before restore", "Create a safety backup of current Saves before restoring",
		t.cfg.BackupBeforeRestore)

	sep := restoreItem.AddSubMenuItem("───────────────────", "")
	sep.Disable()

	// Manual restore sub-submenu
	manualParent := restoreItem.AddSubMenuItem("Manual", "Restore from a manual backup")
	t.restoreManualEmpty = manualParent.AddSubMenuItem("No manual backups yet", "")
	t.restoreManualEmpty.Disable()
	for i := 0; i < maxBackupCap; i++ {
		item := manualParent.AddSubMenuItem("", "")
		item.Hide()
		t.restoreManualItems = append(t.restoreManualItems, item)
	}

	// Auto restore sub-submenu
	autoParent := restoreItem.AddSubMenuItem("Auto", "Restore from an auto backup")
	t.restoreAutoEmpty = autoParent.AddSubMenuItem("No auto backups yet", "")
	t.restoreAutoEmpty.Disable()
	for i := 0; i < maxBackupCap; i++ {
		item := autoParent.AddSubMenuItem("", "")
		item.Hide()
		t.restoreAutoItems = append(t.restoreAutoItems, item)
	}

	t.refreshRestoreMenus()
}

// refreshRestoreMenus re-reads backup dirs and updates the pre-allocated slots.
func (t *Tray) refreshRestoreMenus() {
	t.refreshRestoreDir("manual", t.restoreManualItems, t.restoreManualEmpty, &t.restoreManualSnaps)
	t.refreshRestoreDir("auto", t.restoreAutoItems, t.restoreAutoEmpty, &t.restoreAutoSnaps)
}

func (t *Tray) refreshRestoreDir(subdir string, items []*systray.MenuItem, emptyItem *systray.MenuItem, snaps *[]string) {
	dir := filepath.Join(t.cfg.BackupFolder, subdir)
	newSnaps, _ := backup.ListSnapshots(dir) // newest first; nil on missing dir = empty list
	*snaps = newSnaps

	if len(newSnaps) == 0 {
		emptyItem.Show()
	} else {
		emptyItem.Hide()
	}

	for i, item := range items {
		if i < len(newSnaps) {
			item.SetTitle(newSnaps[i])
			item.Show()
		} else {
			item.Hide()
		}
	}
}

func (t *Tray) handleEvents(manualItem, quitItem *systray.MenuItem) {
	intervalCh := make(chan int, 1)
	for i, item := range t.intervalItems {
		i, item := i, item
		go func() {
			for range item.ClickedCh {
				intervalCh <- i
			}
		}()
	}

	maxBackupCh := make(chan int, 1)
	for i, item := range t.maxBackupItems {
		i, item := i, item
		go func() {
			for range item.ClickedCh {
				maxBackupCh <- i
			}
		}()
	}

	restoreManualCh := make(chan int, 1)
	for i, item := range t.restoreManualItems {
		i, item := i, item
		go func() {
			for range item.ClickedCh {
				restoreManualCh <- i
			}
		}()
	}

	restoreAutoCh := make(chan int, 1)
	for i, item := range t.restoreAutoItems {
		i, item := i, item
		go func() {
			for range item.ClickedCh {
				restoreAutoCh <- i
			}
		}()
	}

	for {
		select {
		case <-t.autoBackupItem.ClickedCh:
			t.toggleAutoBackup()
		case <-manualItem.ClickedCh:
			go func() {
				backup.Manual(t.cfg.ZomboidFolder, t.cfg.BackupFolder, t.cfg.MaxBackupFiles)
				t.refreshRestoreMenus()
			}()
		case <-t.zomboidChange.ClickedCh:
			t.changeZomboidFolder()
		case <-t.backupChange.ClickedCh:
			t.changeBackupFolder()
		case idx := <-intervalCh:
			t.selectInterval(idx)
		case <-t.intervalSet.ClickedCh:
			t.promptInterval()
		case idx := <-maxBackupCh:
			t.selectMaxBackup(idx)
		case <-t.maxBackupSet.ClickedCh:
			t.promptMaxBackups()
		case <-t.backupBeforeItem.ClickedCh:
			t.toggleBackupBefore()
		case idx := <-restoreManualCh:
			if idx < len(t.restoreManualSnaps) {
				go t.doRestore("manual", t.restoreManualSnaps[idx])
			}
		case idx := <-restoreAutoCh:
			if idx < len(t.restoreAutoSnaps) {
				go t.doRestore("auto", t.restoreAutoSnaps[idx])
			}
		case <-quitItem.ClickedCh:
			systray.Quit()
			return
		}
	}
}

func (t *Tray) doRestore(subdir, snapName string) {
	confirmed := osascriptConfirm(fmt.Sprintf(
		"Restore \"%s\"?\n\nThis will replace your current Saves folder.", snapName))
	if !confirmed {
		return
	}
	backupPath := filepath.Join(t.cfg.BackupFolder, subdir, snapName)
	backup.Restore(backupPath, t.cfg.ZomboidFolder, t.cfg.BackupFolder, t.cfg.BackupBeforeRestore)
}

func (t *Tray) toggleAutoBackup() {
	t.cfg.AutoBackup = !t.cfg.AutoBackup
	if t.cfg.AutoBackup {
		t.autoBackupItem.SetTitle("Auto Backup: ON")
		t.autoBackupItem.Check()
	} else {
		t.autoBackupItem.SetTitle("Auto Backup: OFF")
		t.autoBackupItem.Uncheck()
	}
	t.cfg.Save() //nolint:errcheck
}

func (t *Tray) toggleBackupBefore() {
	t.cfg.BackupBeforeRestore = !t.cfg.BackupBeforeRestore
	if t.cfg.BackupBeforeRestore {
		t.backupBeforeItem.Check()
	} else {
		t.backupBeforeItem.Uncheck()
	}
	t.cfg.Save() //nolint:errcheck
}

func (t *Tray) selectInterval(idx int) {
	for i, item := range t.intervalItems {
		if i == idx {
			item.Check()
		} else {
			item.Uncheck()
		}
	}
	t.cfg.BackupInterval = intervalPresets[idx]
	t.intervalParent.SetTitle(fmt.Sprintf("Backup Every: %d min", t.cfg.BackupInterval))
	t.cfg.Save() //nolint:errcheck
}

func (t *Tray) promptInterval() {
	result, ok := osascriptPrompt("Backup interval (minutes, 1–1440):", strconv.Itoa(t.cfg.BackupInterval))
	if !ok {
		return
	}
	val, err := strconv.Atoi(strings.TrimSpace(result))
	if err != nil || val < 1 || val > 1440 {
		return
	}
	for _, item := range t.intervalItems {
		item.Uncheck()
	}
	for i, v := range intervalPresets {
		if v == val {
			t.intervalItems[i].Check()
			break
		}
	}
	t.cfg.BackupInterval = val
	t.intervalParent.SetTitle(fmt.Sprintf("Backup Every: %d min", val))
	t.cfg.Save() //nolint:errcheck
}

func (t *Tray) selectMaxBackup(idx int) {
	for i, item := range t.maxBackupItems {
		if i == idx {
			item.Check()
		} else {
			item.Uncheck()
		}
	}
	t.cfg.MaxBackupFiles = maxBackupPresets[idx]
	t.maxBackupParent.SetTitle(fmt.Sprintf("Max Backups: %d", t.cfg.MaxBackupFiles))
	t.cfg.Save() //nolint:errcheck
}

func (t *Tray) promptMaxBackups() {
	result, ok := osascriptPrompt("Max number of backup files (1–20):", strconv.Itoa(t.cfg.MaxBackupFiles))
	if !ok {
		return
	}
	val, err := strconv.Atoi(strings.TrimSpace(result))
	if err != nil || val < 1 {
		return
	}
	if val > maxBackupCap {
		val = maxBackupCap
	}
	for _, item := range t.maxBackupItems {
		item.Uncheck()
	}
	for i, v := range maxBackupPresets {
		if v == val {
			t.maxBackupItems[i].Check()
			break
		}
	}
	t.cfg.MaxBackupFiles = val
	t.maxBackupParent.SetTitle(fmt.Sprintf("Max Backups: %d", val))
	t.cfg.Save() //nolint:errcheck
}

func (t *Tray) changeBackupFolder() {
	path, ok := osascriptChooseFolder("Select a folder to store backups:", t.cfg.BackupFolder)
	if !ok {
		return
	}
	t.cfg.BackupFolder = path
	t.backupDisplay.SetTitle(fmt.Sprintf("Backup Folder: %s", path))
	t.cfg.Save() //nolint:errcheck
	t.refreshRestoreMenus()
}

func (t *Tray) changeZomboidFolder() {
	path, ok := osascriptChooseFolder("Select your Zomboid game folder:", t.cfg.ZomboidFolder)
	if !ok {
		return
	}
	t.cfg.ZomboidFolder = path
	t.zomboidDisplay.SetTitle(fmt.Sprintf("Zomboid Folder: %s", path))
	t.cfg.Save() //nolint:errcheck
}

func (t *Tray) autoBackupLoop() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	var lastRun time.Time
	for range ticker.C {
		if !t.cfg.AutoBackup {
			continue
		}
		interval := time.Duration(t.cfg.BackupInterval) * time.Minute
		if time.Since(lastRun) >= interval {
			lastRun = time.Now()
			go func() {
				backup.Auto(t.cfg.ZomboidFolder, t.cfg.BackupFolder, t.cfg.MaxBackupFiles)
				t.refreshRestoreMenus()
			}()
		}
	}
}

func osascriptChooseFolder(prompt, defaultPath string) (string, bool) {
	script := fmt.Sprintf(
		`POSIX path of (choose folder with prompt %q default location %q)`,
		prompt, defaultPath,
	)
	out, err := exec.Command("osascript", "-e", script).Output()
	if err != nil {
		return "", false
	}
	return strings.TrimSpace(string(out)), true
}

func osascriptConfirm(message string) bool {
	script := fmt.Sprintf(
		`button returned of (display dialog %q buttons {"Cancel", "Restore"} default button "Cancel" with icon caution)`,
		message,
	)
	out, err := exec.Command("osascript", "-e", script).Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(out)) == "Restore"
}

func osascriptPrompt(message, defaultValue string) (string, bool) {
	script := fmt.Sprintf(
		`display dialog %q default answer %q buttons {"Cancel", "OK"} default button "OK"`,
		message, defaultValue,
	)
	out, err := exec.Command("osascript", "-e", script).Output()
	if err != nil {
		return "", false
	}
	parts := strings.SplitN(string(out), "text returned:", 2)
	if len(parts) < 2 {
		return "", false
	}
	return strings.TrimSpace(parts[1]), true
}

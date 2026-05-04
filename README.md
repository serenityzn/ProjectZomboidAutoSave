# 🧟 ZomboidAutoBackup

A lightweight **macOS and Windows** system tray app for automatically and manually backing up your **Project Zomboid** save files.

---

## Features

- **System tray app** — lives in your menu bar / taskbar, no windows, no clutter
- **Auto backup** — backs up your saves on a configurable schedule
- **Manual backup** — trigger a backup instantly with one click or global hotkey (`⌘⇧B` / `Ctrl+Shift+B`)
- **Restore** — restore any previous backup directly from the menu
- **Safety backup before restore** — automatically preserves your current saves before any restore
- **Efficient compression** — archives use `.tar.zst` (Zstandard), faster and smaller than zip or gzip
- **Backup rotation** — automatically deletes oldest backups when the limit is reached
- **Persistent settings** — all configuration is saved and restored between app launches

---

## Menu Structure

```
ZB (tray icon)
├── Settings
│   ├── Zomboid Folder: ~/Zomboid
│   ├── Change Zomboid Folder...
│   ├── Backup Folder: ~/ZomboidAutoBackup
│   ├── Change Backup Folder...
│   ├── Backup Every: 15 min  ▶  (5–60 min presets + custom)
│   └── Max Backups: 10       ▶  (5 / 10 / 15 presets + custom, max 20)
├── ─────────────────────────────
├── Auto Backup: OFF           (click to toggle)
├── Manual Backup  ⌘⇧B        (click or hotkey to backup now)
├── Restore
│   ├── ✓ Backup before restore
│   ├── Manual ▶  [list of manual snapshots, newest first]
│   └── Auto   ▶  [list of auto snapshots, newest first]
└── Quit
```

---

## Backup Layout

All backups are stored inside your configured **Backup Folder** (default `~/ZomboidAutoBackup`):

```
~/ZomboidAutoBackup/
├── .state               ← saved settings (JSON)
├── manual/
│   ├── m-2026-05-04_11-00-00.tar.zst
│   └── m-2026-05-04_10-30-00.tar.zst
├── auto/
│   ├── a-2026-05-04_11-15-00.tar.zst
│   └── ...
└── before-restore/
    └── r-2026-05-04_11-20-00.tar.zst   ← safety backup created before each restore
```

Each snapshot contains the full `Saves/` directory tree from your Zomboid folder.

---

## Settings

| Setting | Default (macOS) | Default (Windows) | Description |
|---|---|---|---|
| Zomboid Folder | `~/Zomboid` | `%USERPROFILE%\Zomboid` | Path to your Project Zomboid data folder |
| Backup Folder | `~/ZomboidAutoBackup` | `%USERPROFILE%\ZomboidAutoBackup` | Where backups are stored |
| Backup Every | 15 min | 15 min | Auto backup interval (1–1440 min) |
| Max Backups | 10 | 10 | Max snapshots per folder (manual/auto independently, capped at 20) |
| Backup before restore | ON | ON | Archives current saves before any restore |

Settings are persisted to `<BackupFolder>/.state` and loaded on every launch.

---

## Requirements

| | macOS | Windows |
|---|---|---|
| OS | macOS 10.13+ | Windows 10+ |
| Build from source | [Go 1.21+](https://go.dev/dl/) + Xcode CLT | [Go 1.21+](https://go.dev/dl/) + GCC (MinGW-w64) |

---

## Installation

**macOS** — extract and run:
```bash
tar xzf ZomboidAutoBackup_*_darwin_*.tar.gz
./ZomboidAutoBackup
```
If macOS blocks it on first launch: right-click → Open, or run:
```bash
xattr -cr ZomboidAutoBackup
```
To run at login: **System Settings → General → Login Items** → add the binary.

**Windows** — extract the `.zip` and run `ZomboidAutoBackup.exe`.

---

## Build from Source

**macOS:**
```bash
git clone https://github.com/serenityzn/ProjectZomboidAutoSave.git
cd ProjectZomboidAutoSave
go build -o ZomboidAutoBackup .
./ZomboidAutoBackup
```

**Windows (cross-compile from macOS):**
```bash
brew install mingw-w64
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc \
  go build -ldflags "-H=windowsgui" -o ZomboidAutoBackup.exe .
```

---

## Dependencies

| Package | Purpose |
|---|---|
| [`github.com/getlantern/systray`](https://github.com/getlantern/systray) | Cross-platform system tray |
| [`github.com/klauspost/compress`](https://github.com/klauspost/compress) | Zstandard (zstd) compression |
| [`github.com/ncruces/zenity`](https://github.com/ncruces/zenity) | Native OS dialogs (folder picker, prompts) |
| [`golang.design/x/hotkey`](https://pkg.go.dev/golang.design/x/hotkey) | Global hotkey registration |

---

## License

MIT © serenityzn

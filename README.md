# 🧟 ZomboidAutoBackup

A lightweight macOS menu bar app for automatically and manually backing up your **Project Zomboid** save files.

---

## Features

- **System tray app** — lives in your menu bar, no windows, no clutter
- **Auto backup** — backs up your saves on a configurable schedule
- **Manual backup** — trigger a backup instantly with one click
- **Restore** — restore any previous backup directly from the menu
- **Safety backup before restore** — automatically preserves your current saves before any restore
- **Efficient compression** — archives use `.tar.zst` (Zstandard), faster and smaller than zip or gzip
- **Backup rotation** — automatically deletes oldest backups when the limit is reached
- **Persistent settings** — all configuration is saved and restored between app launches

---

## Menu Structure

```
ZB (menu bar)
├── Settings
│   ├── Zomboid Folder: ~/Zomboid
│   ├── Change Zomboid Folder...
│   ├── Backup Folder: ~/ZomboidAutoBackup
│   ├── Change Backup Folder...
│   ├── Backup Every: 15 min  ▶  (5–60 min presets + custom)
│   └── Max Backups: 10       ▶  (5 / 10 / 15 presets + custom, max 20)
├── ─────────────────────────────
├── Auto Backup: OFF           (click to toggle)
├── Manual Backup              (click to backup now)
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
│   ├── snap-2026-05-04_11-00-00.tar.zst
│   └── snap-2026-05-04_10-30-00.tar.zst
├── auto/
│   ├── snap-2026-05-04_11-15-00.tar.zst
│   └── ...
└── before-restore/
    └── snap-2026-05-04_11-20-00.tar.zst   ← safety backup created before each restore
```

Each snapshot contains the full `Saves/` directory tree from your Zomboid folder.

---

## Settings

| Setting | Default | Description |
|---|---|---|
| Zomboid Folder | `~/Zomboid` | Path to your Project Zomboid data folder |
| Backup Folder | `~/ZomboidAutoBackup` | Where backups are stored |
| Backup Every | 15 min | Auto backup interval (1–1440 min) |
| Max Backups | 10 | Max snapshots per folder (manual/auto independently, capped at 20) |
| Backup before restore | ON | Archives current saves before any restore |

Settings are persisted to `~/ZomboidAutoBackup/.state` and loaded on every launch.

---

## Requirements

- macOS 10.13+
- [Go 1.21+](https://go.dev/dl/) (to build from source)

---

## Build & Run

```bash
git clone https://github.com/serenityzn/ProjectZomboidAutoSave.git
cd ProjectZomboidAutoSave
go build -o ZomboidAutoBackup .
./ZomboidAutoBackup
```

To run at login, add the built binary to **System Settings → General → Login Items**.

---

## Dependencies

| Package | Purpose |
|---|---|
| [`github.com/getlantern/systray`](https://github.com/getlantern/systray) | macOS menu bar integration |
| [`github.com/klauspost/compress`](https://github.com/klauspost/compress) | Zstandard (zstd) compression |

---

## License

MIT © serenityzn

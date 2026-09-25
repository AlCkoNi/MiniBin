<div align="center">

# MiniBin

**A compact Windows system-tray utility for managing the Recycle Bin.**

English · [Русский](README.md)

![Version](https://img.shields.io/badge/version-1.1-blue)
![Latest release](https://img.shields.io/github/v/release/AlCkoNi/MiniBin?label=release)
![Platform](https://img.shields.io/badge/platform-Windows-0078D6)
![Architecture](https://img.shields.io/badge/architecture-x86-lightgrey)
![Go](https://img.shields.io/badge/Go-1.23-00ADD8)
![License](https://img.shields.io/badge/license-MIT-green)

[**Download the latest release**](https://github.com/AlCkoNi/MiniBin/releases/latest)

</div>

MiniBin lives in the Windows system tray, shows the current Recycle Bin fill state using one of five icons, and provides quick access to opening or emptying the Recycle Bin. Manual and automatic cleanup also attempt to clear the current user's temporary directory and `%WINDIR%\Temp`.

## Screenshots

<div align="center">

<img src="docs/screenshots/tray.png" alt="MiniBin tray icon" width="180">

<br><br>

<img src="docs/screenshots/menu.png" alt="MiniBin tray menu" width="250">

</div>

## Features

- `Open` opens the standard Windows Recycle Bin.
- `Empty` empties the Recycle Bin without confirmation, sound, or progress UI.
- `Empty` is automatically disabled while the Recycle Bin is empty.
- Double-clicking the tray icon runs `Empty` when items are present.
- `Auto Clean`: `Off`, `30 min`, `1 hour`, `3 hour`.
- `My Function` is an intentionally empty extension submenu for custom features.
- Five fixed tray icon states: `Empty`, `25%`, `50%`, `75%`, `Full`.
- Cleans `%TEMP%` and `%WINDIR%\Temp` during manual and automatic cleanup.
- Configuration is stored next to the executable in `minibin.ini`.
- Portable: no installer is required.
- The public baseline contains no personal absolute paths or machine-specific automation.

## Menu

```text
Open
Empty
────────────
Auto Clean >
    Off
    30 min
    1 hour
    3 hour
────────────
My Function >
    (empty)
────────────
About
Exit
```

## Quick start

1. Open the **Releases** section.
2. Download `MiniBin - 1.1.zip`.
3. Extract the archive.
4. Run `MiniBin.exe`.
5. Keep `empty.ico`, `25.ico`, `50.ico`, `75.ico`, and `full.ico` next to the executable.

## Configuration

```ini
[Configure]
AutoCleanMinutes=0

[Display]
MaxFillSizeMB=1024
```

`AutoCleanMinutes` accepts `0`, `30`, `60`, or `180`. `MaxFillSizeMB` is the reference size used to select the `25% / 50% / 75% / Full` tray icon state.

## Temp cleanup note

MiniBin removes only files that Windows allows it to remove. Locked and protected files can remain. Parts of the system temporary directory may require elevated permissions.

## Build from source

Go 1.23 or later is required. The project uses only the Go standard library.

```powershell
.\scripts\build_release.ps1
```

The Russian developer documentation is available in [`docs/DEVELOPMENT_RU.md`](docs/DEVELOPMENT_RU.md).

## Extending MiniBin

The project intentionally keeps `My Function` as a stable extension point. The public 1.1 build ships it empty so developers can add their own features without changing the top-level menu structure.

See:

- [`docs/MY_FUNCTION_RU.md`](docs/MY_FUNCTION_RU.md)
- [`docs/ADDING_FUNCTIONS_RU.md`](docs/ADDING_FUNCTIONS_RU.md)

## License

MIT. See [`LICENSE`](LICENSE).

## Author

**AlCkoNi**

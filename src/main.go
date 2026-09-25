//go:build windows

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"
	"syscall"
	"unsafe"
)

const windowClassName = "MiniBinHiddenWindow"

var appState struct {
	hwnd      uintptr
	menu      uintptr
	autoMenu  uintptr
	myMenu    uintptr
	nid       notifyIconData
	icons     [5]uintptr
	iconState int32
	config    Config
	cleaning  int32
}

func main() {
	appState.config = loadConfig()
	callback := syscall.NewCallback(windowProc)
	if !registerWindowClass(windowClassName, callback) {
		return
	}
	appState.hwnd = createHiddenWindow(windowClassName, AppName)
	if appState.hwnd == 0 {
		return
	}

	loadTrayIcons()
	buildMenus()
	addTrayIcon()
	refreshRecycleState()
	setTimer(appState.hwnd, timerRefresh, 5000)
	applyAutoCleanTimer()
	runMessageLoop()
	shutdown()
}

func executableDir() string {
	exe, err := os.Executable()
	if err != nil {
		return "."
	}
	return filepath.Dir(exe)
}

func loadTrayIcons() {
	names := []string{"empty.ico", "25.ico", "50.ico", "75.ico", "full.ico"}
	for i, name := range names {
		appState.icons[i] = loadIconFromFile(filepath.Join(executableDir(), name))
		if appState.icons[i] == 0 {
			appState.icons[i] = loadDefaultIcon()
		}
	}
	atomic.StoreInt32(&appState.iconState, -1)
}

func buildMenus() {
	appState.menu = createPopupMenu()
	appState.autoMenu = createPopupMenu()
	appState.myMenu = createPopupMenu()

	appendMenu(appState.menu, mfString, uintptr(idOpen), "Open")
	appendMenu(appState.menu, mfString, uintptr(idEmpty), "Empty")
	appendMenu(appState.menu, mfSeparator, 0, "")
	appendMenu(appState.autoMenu, mfString, uintptr(idAutoOff), "Off")
	appendMenu(appState.autoMenu, mfString, uintptr(idAuto30), "30 min")
	appendMenu(appState.autoMenu, mfString, uintptr(idAuto60), "1 hour")
	appendMenu(appState.autoMenu, mfString, uintptr(idAuto180), "3 hour")
	appendMenu(appState.menu, mfPopup, appState.autoMenu, "Auto Clean")
	appendMenu(appState.menu, mfSeparator, 0, "")

	// Reserved extension point for optional user/developer features.
	// The public build intentionally ships this submenu empty.
	appendMenu(appState.menu, mfPopup, appState.myMenu, "My Function")
	appendMenu(appState.menu, mfSeparator, 0, "")

	appendMenu(appState.menu, mfString, uintptr(idAbout), "About")
	appendMenu(appState.menu, mfString, uintptr(idExit), "Exit")
	updateAutoCleanChecks()
}

func addTrayIcon() {
	appState.nid = notifyIconData{
		CbSize:           uint32(unsafe.Sizeof(notifyIconData{})),
		HWnd:             appState.hwnd,
		UID:              1,
		UFlags:           nifMessage | nifIcon | nifTip,
		UCallbackMessage: wmTray,
		HIcon:            appState.icons[0],
	}
	setUTF16(appState.nid.SzTip[:], AppName+" "+AppVersion)
	shellNotifyIcon(nimAdd, &appState.nid)
}

func shutdown() {
	killTimer(appState.hwnd, timerRefresh)
	killTimer(appState.hwnd, timerAuto)
	shellNotifyIcon(nimDelete, &appState.nid)
	if appState.menu != 0 {
		destroyMenu(appState.menu)
	}
	// autoMenu and myMenu are owned by the parent menu after MF_POPUP; do not destroy twice.
	// Icon handles are intentionally left to process teardown. A fallback icon
	// may be a shared system handle returned by LoadIconW and must not be destroyed.
}

func windowProc(hwnd uintptr, message uint32, wParam, lParam uintptr) uintptr {
	switch message {
	case wmCommand:
		handleCommand(loword(wParam))
		return 0
	case wmTimer:
		switch wParam {
		case timerRefresh:
			refreshRecycleState()
		case timerAuto:
			startCleanup(true)
		}
		return 0
	case wmTray:
		switch uint32(lParam) {
		case wmRButtonUp:
			refreshRecycleState()
			updateAutoCleanChecks()
			showPopupMenu(hwnd, appState.menu)
		case wmLButtonDblClk:
			items, _, ok := queryRecycleBin()
			if ok && items > 0 {
				startCleanup(false)
			}
		}
		return 0
	case wmCleanDone:
		atomic.StoreInt32(&appState.cleaning, 0)
		refreshRecycleState()
		return 0
	case wmDestroy:
		postQuitMessage()
		return 0
	}
	return defWindowProc(hwnd, message, wParam, lParam)
}

func handleCommand(id uint32) {
	switch id {
	case idOpen:
		openRecycleBin(appState.hwnd)
	case idEmpty:
		items, _, ok := queryRecycleBin()
		if ok && items > 0 {
			startCleanup(false)
		}
	case idAutoOff:
		setAutoClean(0)
	case idAuto30:
		setAutoClean(30)
	case idAuto60:
		setAutoClean(60)
	case idAuto180:
		setAutoClean(180)
	case idAbout:
		showAbout()
	case idExit:
		postQuitMessage()
	}
}

func startCleanup(fromTimer bool) {
	if !atomic.CompareAndSwapInt32(&appState.cleaning, 0, 1) {
		return
	}
	if !fromTimer {
		enableMenuItem(appState.menu, idEmpty, false)
	}
	go func() {
		emptyRecycleBin(appState.hwnd)
		clearTemporaryFiles()
		postMessage(appState.hwnd, wmCleanDone)
	}()
}

func refreshRecycleState() {
	items, size, ok := queryRecycleBin()
	if !ok {
		return
	}
	enableMenuItem(appState.menu, idEmpty, items > 0 && atomic.LoadInt32(&appState.cleaning) == 0)
	state := iconStateFor(items, size)
	if atomic.LoadInt32(&appState.iconState) == int32(state) {
		return
	}
	appState.nid.HIcon = appState.icons[state]
	shellNotifyIcon(nimModify, &appState.nid)
	atomic.StoreInt32(&appState.iconState, int32(state))
}

func iconStateFor(items, size int64) int {
	if items <= 0 {
		return 0
	}
	maxBytes := appState.config.MaxFillSizeMB * 1024 * 1024
	if maxBytes <= 0 {
		maxBytes = 1024 * 1024 * 1024
	}
	ratio := float64(size) / float64(maxBytes)
	switch {
	case ratio < 0.25:
		return 1
	case ratio < 0.50:
		return 2
	case ratio < 0.75:
		return 3
	default:
		return 4
	}
}

func setAutoClean(minutes int) {
	if !validAutoCleanMinutes(minutes) {
		return
	}
	appState.config.AutoCleanMinutes = minutes
	_ = saveConfig(appState.config)
	applyAutoCleanTimer()
	updateAutoCleanChecks()
}

func applyAutoCleanTimer() {
	killTimer(appState.hwnd, timerAuto)
	if appState.config.AutoCleanMinutes > 0 {
		ms := uint32(appState.config.AutoCleanMinutes * 60 * 1000)
		setTimer(appState.hwnd, timerAuto, ms)
	}
}

func updateAutoCleanChecks() {
	checks := map[uint32]int{
		idAutoOff: 0,
		idAuto30:  30,
		idAuto60:  60,
		idAuto180: 180,
	}
	for id, minutes := range checks {
		checkMenuItem(appState.autoMenu, id, appState.config.AutoCleanMinutes == minutes)
	}
}

func showAbout() {
	sizeText := "unknown"
	if info, err := os.Stat(os.Args[0]); err == nil {
		sizeText = formatSize(info.Size())
	}
	text := fmt.Sprintf("%s %s\nSize: %s\n%s", AppName, AppVersion, sizeText, Author)
	messageBox(appState.hwnd, text, "About")
}

func formatSize(n int64) string {
	if n < 1024 {
		return fmt.Sprintf("%d B", n)
	}
	if n < 1024*1024 {
		return fmt.Sprintf("%.1f KB", float64(n)/1024)
	}
	return fmt.Sprintf("%.1f MB", float64(n)/(1024*1024))
}

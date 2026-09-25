//go:build windows

package main

import (
	"syscall"
	"unsafe"
)

const (
	wmDestroy       = 0x0002
	wmCommand       = 0x0111
	wmTimer         = 0x0113
	wmLButtonDblClk = 0x0203
	wmRButtonUp     = 0x0205
	wmTray          = 0x0400 + 1
	wmCleanDone     = 0x8000 + 1

	mfString    = 0x0000
	mfGray      = 0x0001
	mfChecked   = 0x0008
	mfPopup     = 0x0010
	mfSeparator = 0x0800
	mfByCommand = 0x0000

	nifMessage = 0x0001
	nifIcon    = 0x0002
	nifTip     = 0x0004
	nimAdd     = 0x00000000
	nimModify  = 0x00000001
	nimDelete  = 0x00000002

	imageIcon      = 1
	lrLoadFromFile = 0x0010

	shERBNoConfirmation = 0x00000001
	shERBNoProgressUI   = 0x00000002
	shERBNoSound        = 0x00000004

	idcArrow = 32512
	idiApp   = 32512

	mbOK            = 0x00000000
	mbIconInfo      = 0x00000040
	mbSetForeground = 0x00010000
)

const (
	idOpen uint32 = 1001 + iota
	idEmpty
	idAutoOff
	idAuto30
	idAuto60
	idAuto180
	idAbout
	idExit
)

const (
	timerRefresh uintptr = 1
	timerAuto    uintptr = 2
)

type point struct {
	X int32
	Y int32
}

type msg struct {
	HWnd    uintptr
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      point
}

type wndClassEx struct {
	CbSize        uint32
	Style         uint32
	LpfnWndProc   uintptr
	CbClsExtra    int32
	CbWndExtra    int32
	HInstance     uintptr
	HIcon         uintptr
	HCursor       uintptr
	HbrBackground uintptr
	LpszMenuName  *uint16
	LpszClassName *uint16
	HIconSm       uintptr
}

type notifyIconData struct {
	CbSize           uint32
	HWnd             uintptr
	UID              uint32
	UFlags           uint32
	UCallbackMessage uint32
	HIcon            uintptr
	SzTip            [128]uint16
	DwState          uint32
	DwStateMask      uint32
	SzInfo           [256]uint16
	UVersion         uint32
	SzInfoTitle      [64]uint16
	DwInfoFlags      uint32
	GuidItem         [16]byte
	HBalloonIcon     uintptr
}

type shQueryRBInfo struct {
	CbSize      uint32
	I64Size     int64
	I64NumItems int64
}

var (
	user32   = syscall.NewLazyDLL("user32.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	shell32  = syscall.NewLazyDLL("shell32.dll")

	procRegisterClassExW    = user32.NewProc("RegisterClassExW")
	procCreateWindowExW     = user32.NewProc("CreateWindowExW")
	procDefWindowProcW      = user32.NewProc("DefWindowProcW")
	procGetMessageW         = user32.NewProc("GetMessageW")
	procTranslateMessage    = user32.NewProc("TranslateMessage")
	procDispatchMessageW    = user32.NewProc("DispatchMessageW")
	procPostQuitMessage     = user32.NewProc("PostQuitMessage")
	procPostMessageW        = user32.NewProc("PostMessageW")
	procCreatePopupMenu     = user32.NewProc("CreatePopupMenu")
	procAppendMenuW         = user32.NewProc("AppendMenuW")
	procDestroyMenu         = user32.NewProc("DestroyMenu")
	procTrackPopupMenu      = user32.NewProc("TrackPopupMenu")
	procSetForegroundWindow = user32.NewProc("SetForegroundWindow")
	procGetCursorPos        = user32.NewProc("GetCursorPos")
	procEnableMenuItem      = user32.NewProc("EnableMenuItem")
	procCheckMenuItem       = user32.NewProc("CheckMenuItem")
	procSetTimer            = user32.NewProc("SetTimer")
	procKillTimer           = user32.NewProc("KillTimer")
	procLoadImageW          = user32.NewProc("LoadImageW")
	procLoadIconW           = user32.NewProc("LoadIconW")
	procDestroyIcon         = user32.NewProc("DestroyIcon")
	procLoadCursorW         = user32.NewProc("LoadCursorW")
	procMessageBoxW         = user32.NewProc("MessageBoxW")

	procGetModuleHandleW   = kernel32.NewProc("GetModuleHandleW")
	procShellNotifyIconW   = shell32.NewProc("Shell_NotifyIconW")
	procSHQueryRecycleBinW = shell32.NewProc("SHQueryRecycleBinW")
	procSHEmptyRecycleBinW = shell32.NewProc("SHEmptyRecycleBinW")
	procShellExecuteW      = shell32.NewProc("ShellExecuteW")
)

func utf16Ptr(s string) *uint16 {
	p, _ := syscall.UTF16PtrFromString(s)
	return p
}

func setUTF16(dst []uint16, s string) {
	u, _ := syscall.UTF16FromString(s)
	if len(u) > len(dst) {
		u = u[:len(dst)]
		u[len(u)-1] = 0
	}
	copy(dst, u)
}

func loword(v uintptr) uint32 { return uint32(v & 0xffff) }

func defWindowProc(hwnd uintptr, message uint32, wParam, lParam uintptr) uintptr {
	r, _, _ := procDefWindowProcW.Call(hwnd, uintptr(message), wParam, lParam)
	return r
}

func getModuleHandle() uintptr {
	r, _, _ := procGetModuleHandleW.Call(0)
	return r
}

func loadCursor() uintptr {
	r, _, _ := procLoadCursorW.Call(0, idcArrow)
	return r
}

func registerWindowClass(className string, callback uintptr) bool {
	name := utf16Ptr(className)
	wc := wndClassEx{
		CbSize:        uint32(unsafe.Sizeof(wndClassEx{})),
		LpfnWndProc:   callback,
		HInstance:     getModuleHandle(),
		HCursor:       loadCursor(),
		LpszClassName: name,
	}
	r, _, _ := procRegisterClassExW.Call(uintptr(unsafe.Pointer(&wc)))
	return r != 0
}

func createHiddenWindow(className, title string) uintptr {
	r, _, _ := procCreateWindowExW.Call(
		0,
		uintptr(unsafe.Pointer(utf16Ptr(className))),
		uintptr(unsafe.Pointer(utf16Ptr(title))),
		0,
		0, 0, 0, 0,
		0, 0,
		getModuleHandle(),
		0,
	)
	return r
}

func runMessageLoop() {
	var m msg
	for {
		r, _, _ := procGetMessageW.Call(uintptr(unsafe.Pointer(&m)), 0, 0, 0)
		if int32(r) <= 0 {
			return
		}
		procTranslateMessage.Call(uintptr(unsafe.Pointer(&m)))
		procDispatchMessageW.Call(uintptr(unsafe.Pointer(&m)))
	}
}

func postQuitMessage() { procPostQuitMessage.Call(0) }

func postMessage(hwnd uintptr, message uint32) {
	procPostMessageW.Call(hwnd, uintptr(message), 0, 0)
}

func createPopupMenu() uintptr {
	r, _, _ := procCreatePopupMenu.Call()
	return r
}

func appendMenu(menu uintptr, flags uint32, id uintptr, text string) {
	var p uintptr
	if text != "" {
		p = uintptr(unsafe.Pointer(utf16Ptr(text)))
	}
	procAppendMenuW.Call(menu, uintptr(flags), id, p)
}

func destroyMenu(menu uintptr) { procDestroyMenu.Call(menu) }

func enableMenuItem(menu uintptr, id uint32, enabled bool) {
	flags := uint32(mfByCommand)
	if !enabled {
		flags |= mfGray
	}
	procEnableMenuItem.Call(menu, uintptr(id), uintptr(flags))
}

func checkMenuItem(menu uintptr, id uint32, checked bool) {
	flags := uint32(mfByCommand)
	if checked {
		flags |= mfChecked
	}
	procCheckMenuItem.Call(menu, uintptr(id), uintptr(flags))
}

func showPopupMenu(hwnd, menu uintptr) {
	var pt point
	procGetCursorPos.Call(uintptr(unsafe.Pointer(&pt)))
	procSetForegroundWindow.Call(hwnd)
	procTrackPopupMenu.Call(menu, 0x0002, uintptr(pt.X), uintptr(pt.Y), 0, hwnd, 0)
	postMessage(hwnd, 0)
}

func setTimer(hwnd, id uintptr, milliseconds uint32) {
	procSetTimer.Call(hwnd, id, uintptr(milliseconds), 0)
}

func killTimer(hwnd, id uintptr) { procKillTimer.Call(hwnd, id) }

func shellNotifyIcon(action uint32, data *notifyIconData) bool {
	r, _, _ := procShellNotifyIconW.Call(uintptr(action), uintptr(unsafe.Pointer(data)))
	return r != 0
}

func loadIconFromFile(path string) uintptr {
	r, _, _ := procLoadImageW.Call(0, uintptr(unsafe.Pointer(utf16Ptr(path))), imageIcon, 16, 16, lrLoadFromFile)
	return r
}

func loadDefaultIcon() uintptr {
	r, _, _ := procLoadIconW.Call(0, idiApp)
	return r
}

func destroyIcon(icon uintptr) {
	if icon != 0 {
		procDestroyIcon.Call(icon)
	}
}

func queryRecycleBin() (items int64, size int64, ok bool) {
	info := shQueryRBInfo{CbSize: uint32(unsafe.Sizeof(shQueryRBInfo{}))}
	r, _, _ := procSHQueryRecycleBinW.Call(0, uintptr(unsafe.Pointer(&info)))
	if int32(r) != 0 {
		return 0, 0, false
	}
	return info.I64NumItems, info.I64Size, true
}

func emptyRecycleBin(hwnd uintptr) {
	flags := uintptr(shERBNoConfirmation | shERBNoProgressUI | shERBNoSound)
	procSHEmptyRecycleBinW.Call(hwnd, 0, flags)
}

func openRecycleBin(hwnd uintptr) {
	procShellExecuteW.Call(
		hwnd,
		uintptr(unsafe.Pointer(utf16Ptr("open"))),
		uintptr(unsafe.Pointer(utf16Ptr("shell:RecycleBinFolder"))),
		0, 0, 1,
	)
}

func messageBox(hwnd uintptr, text, caption string) {
	procMessageBoxW.Call(
		hwnd,
		uintptr(unsafe.Pointer(utf16Ptr(text))),
		uintptr(unsafe.Pointer(utf16Ptr(caption))),
		mbOK|mbIconInfo|mbSetForeground,
	)
}

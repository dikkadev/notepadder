package win

import (
	"fmt"
	"os/exec"
	"syscall"
	"time"
	"unsafe"
)

var (
	user32                    = syscall.MustLoadDLL("user32.dll")
	procFindWindowW           = user32.MustFindProc("FindWindowW")
	procSetForegroundWindow   = user32.MustFindProc("SetForegroundWindow")
	procShowWindow            = user32.MustFindProc("ShowWindow")
	procKeybdEvent            = user32.MustFindProc("keybd_event")
)

// overrideable functions for testing
var (
	findWindow   = findNotepadWindow
	launchNotepad = launchNotepadInternal
	activateWindow = setForegroundWindow
	sendCtrlT      = sendCtrlTInternal
)

// Run finds or launches Notepad, activates it, and optionally sends Ctrl+T
func Run(noNew bool) error {
	hwnd, err := findWindow()
	if err != nil {
		return err
	}
	if hwnd == 0 {
		if err := launchNotepad(); err != nil {
			return err
		}
		for i := 0; i < 50; i++ {
			time.Sleep(100 * time.Millisecond)
			hwnd, err = findWindow()
			if err != nil {
				return err
			}
			if hwnd != 0 {
				break
			}
		}
		if hwnd == 0 {
			return fmt.Errorf("could not find Notepad window after launch")
		}
	}
	if err := activateWindow(hwnd); err != nil {
		return err
	}
	if !noNew {
		if err := sendCtrlT(); err != nil {
			return err
		}
	}
	return nil
}

func findNotepadWindow() (syscall.Handle, error) {
	ptr, err := syscall.UTF16PtrFromString("Notepad")
	if err != nil {
		return 0, err
	}
	r1, _, e1 := procFindWindowW.Call(
		uintptr(unsafe.Pointer(ptr)),
		uintptr(0),
	)
	if r1 == 0 {
		return 0, e1
	}
	return syscall.Handle(r1), nil
}

func launchNotepadInternal() error {
	return exec.Command("notepad.exe").Start()
}

func setForegroundWindow(hwnd syscall.Handle) error {
	const SW_RESTORE = 9
	// restore if minimized
	procShowWindow.Call(uintptr(hwnd), uintptr(SW_RESTORE))
	r1, _, e1 := procSetForegroundWindow.Call(uintptr(hwnd))
	if r1 == 0 {
		return e1
	}
	return nil
}

func sendCtrlTInternal() error {
	const (
		KEYEVENTF_KEYUP = 0x0002
		VK_CONTROL     = 0x11
		VK_T           = 0x54
	)
	// press Control
	procKeybdEvent.Call(uintptr(VK_CONTROL), 0, 0, 0)
	// press T
	procKeybdEvent.Call(uintptr(VK_T), 0, 0, 0)
	// release T
	procKeybdEvent.Call(uintptr(VK_T), 0, uintptr(KEYEVENTF_KEYUP), 0)
	// release Control
	procKeybdEvent.Call(uintptr(VK_CONTROL), 0, uintptr(KEYEVENTF_KEYUP), 0)
	return nil
} 
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
	procIsIconic              = user32.MustFindProc("IsIconic")
)

// overrideable functions for testing
var (
	findWindow   = findNotepadWindow
	launchNotepad = launchNotepadInternal
	activateWindow = setForegroundWindow
	sendCtrlT      = sendCtrlTInternal
)

// Run finds or launches Notepad, activates it, and optionally sends Ctrl+T
func Run(noNew bool, debug bool) error {
	if debug { fmt.Println("Attempting to find Notepad window...") }
	hwnd, err := findWindow()
	if err != nil {
		if debug { fmt.Printf("Error finding window: %v\n", err) }
		return err
	}
	if debug { fmt.Printf("FindWindow result: hwnd=%v\n", hwnd) }

	if hwnd == 0 {
		if debug { fmt.Println("Notepad window not found. Launching...") }
		if err := launchNotepad(); err != nil {
			if debug { fmt.Printf("Error launching notepad: %v\n", err) }
			return err
		}
		if debug { fmt.Println("Notepad launched. Waiting for window...") }
		found := false
		for i := 0; i < 50; i++ {
			time.Sleep(50 * time.Millisecond)
			if debug { fmt.Printf("  Retry %d: Calling FindWindow...\n", i+1) }
			hwnd, err = findWindow()
			if err != nil {
				if debug { fmt.Printf("  Retry %d: Error finding window: %v\n", i+1, err) }
				// Don't return immediately, maybe transient error
			}
			if debug { fmt.Printf("  Retry %d: FindWindow result: hwnd=%v\n", i+1, hwnd) }
			if hwnd != 0 {
				if debug { fmt.Println("  Window found!") }
				found = true
				break
			}
		}
		if !found {
			 err := fmt.Errorf("could not find Notepad window after launch and 50 retries")
			 if debug { fmt.Printf("Error: %v\n", err)}
			 return err
		}
	}

	if debug { fmt.Printf("Activating window: hwnd=%v\n", hwnd) }
	if err := activateWindow(hwnd); err != nil {
		if debug { fmt.Printf("Error activating window: %v\n", err) }
		return err
	}
	if debug { fmt.Println("Window activated.") }

	if !noNew {
		if debug { fmt.Println("Sending Ctrl+T...") }
		// Add a small delay to allow Notepad to fully initialize after activation
		time.Sleep(100 * time.Millisecond)
		if err := sendCtrlT(); err != nil {
			if debug { fmt.Printf("Error sending Ctrl+T: %v\n", err) }
			return err
		}
		if debug { fmt.Println("Ctrl+T sent.") }
	} else {
		if debug { fmt.Println("Skipping new tab (--no-new specified).") }
	}

	// Add a small delay before exiting to allow activated window to settle
	time.Sleep(100 * time.Millisecond)
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
		if e1.(syscall.Errno) != 0 {
			return 0, e1
		}
		return 0, nil
	}
	return syscall.Handle(r1), nil
}

func launchNotepadInternal() error {
	return exec.Command("notepad.exe").Start()
}

func setForegroundWindow(hwnd syscall.Handle) error {
	// Check if minimized
	isMinimized, _, _ := procIsIconic.Call(uintptr(hwnd))
	if isMinimized != 0 {
		const SW_RESTORE = 9
		// restore if minimized
		procShowWindow.Call(uintptr(hwnd), uintptr(SW_RESTORE))
	}

	// Now set foreground
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
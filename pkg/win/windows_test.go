package win

import (
	"errors"
	"syscall"
	"testing"
)

func TestRun_WithNewTab(t *testing.T) {
	// Save original functions and restore after
	origFindWindow := findWindow
	origLaunch := launchNotepad
	origActivate := activateWindow
	origSend := sendCtrlT
	defer func() {
		findWindow = origFindWindow
		launchNotepad = origLaunch
		activateWindow = origActivate
		sendCtrlT = origSend
	}()

	var callOrder []string
	firstCall := true
	// simulate initial no window, then window found
	findWindow = func() (syscall.Handle, error) {
		if firstCall {
			callOrder = append(callOrder, "find")
			firstCall = false
			return 0, nil
		}
		callOrder = append(callOrder, "find2")
		return syscall.Handle(123), nil
	}
	launchNotepad = func() error {
		callOrder = append(callOrder, "launch")
		return nil
	}
	activateWindow = func(hwnd syscall.Handle) error {
		callOrder = append(callOrder, "activate")
		if hwnd != 123 {
			t.Errorf("expected handle 123, got %v", hwnd)
		}
		return nil
	}
	sendCtrlT = func() error {
		callOrder = append(callOrder, "sendCtrlT")
		return nil
	}

	if err := Run(false); err != nil {
		t.Fatalf("Run(false) error: %v", err)
	}

	expected := []string{"find", "launch", "find2", "activate", "sendCtrlT"}
	if len(callOrder) != len(expected) {
		t.Errorf("callOrder length %d, want %d, calls: %v", len(callOrder), len(expected), callOrder)
	}
	for i, want := range expected {
		if callOrder[i] != want {
			t.Errorf("callOrder[%d] = %q, want %q", i, callOrder[i], want)
		}
	}
}

func TestRun_NoNewTab(t *testing.T) {
	// Save original functions and restore after
	origFind := findWindow
	origActivate := activateWindow
	origSend := sendCtrlT
	defer func() {
		findWindow = origFind
		activateWindow = origActivate
		sendCtrlT = origSend
	}()

	var callOrder []string
	findWindow = func() (syscall.Handle, error) {
		callOrder = append(callOrder, "find")
		return syscall.Handle(456), nil
	}
	activateWindow = func(hwnd syscall.Handle) error {
		callOrder = append(callOrder, "activate")
		return nil
	}
	sendCtrlT = func() error {
		callOrder = append(callOrder, "sendCtrlT")
		return errors.New("should not be called")
	}

	if err := Run(true); err != nil {
		t.Fatalf("Run(true) error: %v", err)
	}

	expected := []string{"find", "activate"}
	if len(callOrder) != len(expected) {
		t.Errorf("callOrder length %d, want %d, calls: %v", len(callOrder), len(expected), callOrder)
	}
	for i, want := range expected {
		if callOrder[i] != want {
			t.Errorf("callOrder[%d] = %q, want %q", i, callOrder[i], want)
		}
	}
} 
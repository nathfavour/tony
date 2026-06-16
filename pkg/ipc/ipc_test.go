package ipc

import (
	"net"
	"os"
	"testing"
)

func TestIPCFD(t *testing.T) {
	socketPath := "/tmp/tony_test.sock"
	engine := NewEngine(socketPath)
	defer engine.CloseSocket()

	ln, err := engine.Listen()
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	defer ln.Close()

	// Client goroutine
	go func() {
		conn, err := engine.Dial()
		if err != nil {
			t.Errorf("failed to dial: %v", err)
			return
		}
		defer conn.Close()

		// Send an FD (stdout for test)
		err = SendFD(conn, int(os.Stdout.Fd()))
		if err != nil {
			t.Errorf("failed to send FD: %v", err)
		}
	}()

	conn, err := ln.Accept()
	if err != nil {
		t.Fatalf("failed to accept: %v", err)
	}
	defer conn.Close()

	unixConn := conn.(*net.UnixConn)
	fds, err := RecvFD(unixConn, 1)
	if err != nil {
		t.Fatalf("failed to receive FD: %v", err)
	}

	if len(fds) != 1 {
		t.Errorf("expected 1 FD, got %d", len(fds))
	} else {
		// Just close the received FD
		syscallClose(fds[0])
	}
}

// Minimal syscall close wrapper
func syscallClose(fd int) {
	// Import syscall or unix would be better, but we just want to close it.
	f := os.NewFile(uintptr(fd), "received_fd")
	if f != nil {
		f.Close()
	}
}

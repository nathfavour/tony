package ipc

import (
	"fmt"
	"net"
	"os"
	"syscall"
)

// Engine handles the low-level UDS communication and FD passing.
type Engine struct {
	SocketPath string
}

// NewEngine creates a new IPC engine.
func NewEngine(path string) *Engine {
	return &Engine{SocketPath: path}
}

// Listen creates a Unix domain socket listener.
func (e *Engine) Listen() (net.Listener, error) {
	if _, err := os.Stat(e.SocketPath); err == nil {
		if err := os.Remove(e.SocketPath); err != nil {
			return nil, err
		}
	}
	return net.Listen("unix", e.SocketPath)
}

// Dial connects to a Unix domain socket.
func (e *Engine) Dial() (*net.UnixConn, error) {
	addr, err := net.ResolveUnixAddr("unix", e.SocketPath)
	if err != nil {
		return nil, err
	}
	return net.DialUnix("unix", nil, addr)
}

// SendFD sends file descriptors over a Unix domain socket connection.
func SendFD(conn *net.UnixConn, fds ...int) error {
	rights := syscall.UnixRights(fds...)
	// We send a dummy byte because some systems require at least one byte of data.
	_, _, err := conn.WriteMsgUnix([]byte{0}, rights, nil)
	return err
}

// RecvFD receives file descriptors over a Unix domain socket connection.
func RecvFD(conn *net.UnixConn, maxFds int) ([]int, error) {
	// Calculate space for the control message.
	buf := make([]byte, syscall.CmsgSpace(maxFds*4))
	dummy := make([]byte, 1)
	
	_, oobn, _, _, err := conn.ReadMsgUnix(dummy, buf)
	if err != nil {
		return nil, err
	}
	
	msgs, err := syscall.ParseSocketControlMessage(buf[:oobn])
	if err != nil {
		return nil, err
	}
	
	var fds []int
	for _, msg := range msgs {
		fds2, err := syscall.ParseUnixRights(&msg)
		if err != nil {
			return nil, err
		}
		fds = append(fds, fds2...)
	}
	
	if len(fds) == 0 {
		return nil, fmt.Errorf("no file descriptors received")
	}
	
	return fds, nil
}

// CloseSocket removes the socket file.
func (e *Engine) CloseSocket() error {
	return os.Remove(e.SocketPath)
}

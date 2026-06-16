package runtime

import (
	"os"
	"os/exec"
	"syscall"
)

// Environment represents an isolated runtime for an unmanned agent.
type Environment struct {
	ID   string
	Root string // Path to the agent's isolated root
}

// NewEnvironment creates a new isolated environment.
func NewEnvironment(id, root string) *Environment {
	return &Environment{
		ID:   id,
		Root: root,
	}
}

// Spawn executes a command within the isolated environment.
// It uses Linux-specific namespaces for isolation.
func (e *Environment) Spawn(command string, args ...string) (*exec.Cmd, error) {
	cmd := exec.Command(command, args...)
	
	// Configure namespaces for isolation.
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUSER | 
		            syscall.CLONE_NEWNS | 
		            syscall.CLONE_NEWPID | 
		            syscall.CLONE_NEWUTS | 
		            syscall.CLONE_NEWNET,
	}

	// Only set mappings if we are running on a system that supports it
	// and if we have the necessary privileges or are in a user namespace.
	cmd.SysProcAttr.UidMappings = []syscall.SysProcIDMap{
		{
			ContainerID: 0,
			HostID:      os.Getuid(),
			Size:        1,
		},
	}
	cmd.SysProcAttr.GidMappings = []syscall.SysProcIDMap{
		{
			ContainerID: 0,
			HostID:      os.Getgid(),
			Size:        1,
		},
	}

	return cmd, nil
}

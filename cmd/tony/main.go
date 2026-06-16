package main

import (
	"fmt"
	"os"

	"github.com/nathfavour/tony/pkg/identity"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("tony: agentic kernel initialized")
		fmt.Println("usage: tony <command> [args]")
		return
	}

	cmd := os.Args[1]
	switch cmd {
	case "version":
		fmt.Println("tony v0.1.0")
	case "derive":
		if len(os.Args) < 4 {
			fmt.Println("usage: tony derive <master_seed_hex> <path>")
			return
		}
		// Minimal placeholder that uses the package
		_ = identity.NewManager([32]byte{})
		fmt.Println("derivation engine active")
	default:
		fmt.Printf("unknown command: %s\n", cmd)
	}
}

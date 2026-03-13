// main.go - kwtsms-cli entry point.
// Delegates immediately to cmd.Execute() following cobra conventions.
// Related files: cmd/root.go
package main

import (
	"fmt"
	"os"

	"github.com/boxlinknet/kwtsms-cli/cmd"
	"github.com/boxlinknet/kwtsms-cli/internal/update"
)

func main() {
	notice := update.CheckAsync(cmd.Version())
	cmd.Execute()
	if msg := <-notice; msg != "" {
		fmt.Fprintln(os.Stderr, msg)
	}
}

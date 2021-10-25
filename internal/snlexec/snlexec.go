package snlexec

import (
	"context"
	"fmt"
	"time"

	"github.com/kassybas/shannel/internal/snlvar"
)

// ShExec prepares the context and options for a shell execution
func ShExec(script string, vt *snlvar.VarTable, timeout time.Duration) (exitStatus int, err error) {

	var ctx context.Context
	var cancel context.CancelFunc
	if timeout != 0 {
		ctx, cancel = context.WithTimeout(context.Background(), timeout)
		defer cancel()
	} else {
		ctx = context.Background()
	}

	opts := Options{
		Silent:       false, // if false: streams to stdout and stderr
		IgnoreResult: false, // if false: returns stdout, stderr and status code
		ShieldEnv:    true,  // if false: exposes the current process' environment variables for the shell process

		ShellPath:       "/bin/bash",         // shell to be executed (default: sh)
		ShellExtraFlags: []string{"--posix"}, // extra options passed to the shell
	}
	_, _, exitStatus, err = ShellExec(ctx, script, vt.DumpShellFormat(), opts)
	if ctx.Err() == context.DeadlineExceeded {
		return exitStatus, fmt.Errorf("timeout of '%s' exceeded", timeout.String())
	}
	return
}

package snlexec

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
)

// Options defines the mode of shell execution
// TODO
type Options struct {
	Silent       bool // no output to stdout and stderr
	IgnoreResult bool // result of stdout and stderr is not returned (enable to save resources + you don't care about the values)
	ShieldEnv    bool // do not expose the current process' environment variables (ignore wath is in os.Environ())

	// ShellExec invocation scheme: "ShellPath ShellExtraFlags ShellCmdFlag Script"
	ShellPath       string   // path of the shell to be executed (defaults to sh)
	ShellCmdFlag    string   // flag that precedes the script string (default '-c')
	ShellExtraFlags []string // list of flags for the invocation of the shell (eg. --posix for bash)
}

const defaultShellCmdFlag = "-c"
const defaultShell = "sh"

func copyAndCapture(w io.Writer, r io.Reader, copyEnabled, capture bool) ([]byte, error) {
	var out []byte
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			d := buf[:n]
			if capture {
				out = append(out, d...)
			}
			if copyEnabled {
				_, err := w.Write(d)
				if err != nil {
					return out, err
				}
			}
		}
		if err != nil {
			// Read returns io.EOF at the end of file, which is not an error for us
			if err == io.EOF {
				err = nil
			}
			return out, err
		}
	}
}

func createCommand(ctx context.Context, script string, envVars []string, opts Options) *exec.Cmd {
	// Create Cmd with options
	if opts.ShellPath == "" {
		opts.ShellPath = defaultShell
	}
	// TODO: handle extra flags
	args := append(opts.ShellExtraFlags, defaultShellCmdFlag, script)
	cmd := exec.CommandContext(ctx, opts.ShellPath, args...)
	// Add environment variables
	if opts.ShieldEnv {
		cmd.Env = []string{}
	} else {
		cmd.Env = os.Environ()
	}
	cmd.Env = append(cmd.Env, envVars...)
	return cmd

}

// func getOutputFileDescriptors(silent bool) (*os.File, *os.File) {
// 	if silent {
// 		devNull := os.NewFile(0, os.DevNull)
// 		return devNull, devNull
// 	}
// 	return os.Stdout, os.Stderr
// }

// ShellExec executes the given shell command with the passed options and environment variables
// Returns the stdout, stderr, status code and execution errors
func ShellExec(ctx context.Context, script string, envVars []string, opts Options) (outStr, errStr string, statusCode int, err error) {
	var stdout, stderr []byte
	var errStdout, errStderr error

	cmd := createCommand(ctx, script, envVars, opts)

	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()

	err = cmd.Start()
	if err != nil {
		return
	}

	if !opts.Silent || !opts.IgnoreResult {
		var wg sync.WaitGroup
		// cmd.Wait() should be called only after we finish reading
		// from stdoutIn and stderrIn.
		// wg ensures that we finish
		wg.Add(1)
		go func() {
			stdout, errStdout = copyAndCapture(os.Stdout, stdoutIn, !opts.Silent, !opts.IgnoreResult)
			wg.Done()
		}()

		stderr, errStderr = copyAndCapture(os.Stderr, stderrIn, !opts.Silent, !opts.IgnoreResult)

		wg.Wait()
	}
	err = cmd.Wait()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			statusCode = exitError.ExitCode()
			err = nil
		}
	}
	if errStdout != nil {
		err = fmt.Errorf("failed to capture stdout: %w", errStdout)
	}
	if errStderr != nil {
		err = fmt.Errorf("failed to capture stderr: %w", errStderr)
	}
	outStr, errStr = strings.TrimSuffix(string(stdout), "\n"), strings.TrimSuffix(string(stderr), "\n")
	return
}

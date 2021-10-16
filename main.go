package main

import (
	"github.com/kassybas/shell-exec/exec"
)

func main() {
	script := `
		echo hello ${FOO}
	  `
	envVars := []string{
		"FOO=world",
	}
	runSh(script, envVars)

}

func runSh(script string, envVars []string) {
	opts := exec.Options{
		Silent:       false, // if false: streams to stdout and stderr
		IgnoreResult: false, // if false: returns stdout, stderr and status code
		ShieldEnv:    false, // if false: exposes the current process' environment variables for the shell process

		ShellPath:       "/bin/bash",         // shell to be executed (default: sh)
		ShellExtraFlags: []string{"--posix"}, // extra options passed to the shell
	}
	exec.ShellExec(script, envVars, opts)

}

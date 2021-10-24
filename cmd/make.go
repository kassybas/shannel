package cmd

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/kassybas/shannel/internal/snlargs"
	"github.com/kassybas/shannel/internal/snlexec"
	"github.com/kassybas/shannel/internal/snlloader"
	"github.com/kassybas/shannel/internal/snlvar"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func Make(c *cli.Context) error {
	loglevel, err := logrus.ParseLevel(c.String("loglevel"))
	if err != nil {
		logrus.Fatal("unknown log level")
	}
	logrus.SetLevel(loglevel)

	logrus.Trace("make command started")
	targetName := c.Args().First()
	path := c.String("file")
	makelog := logrus.WithField("target", targetName).WithField("path", path)

	makelog.Trace("loading shannel file...")
	snlfile, err := snlloader.Load(path)
	if err != nil {
		makelog.WithField("error", err).Fatalf("failed to load file")
	}

	if targetName == "" { // && no default target set
		makelog.Trace("target name empty, printing help")
		PrintHelpText(snlfile)
		return nil
	}

	targetCode, exists := snlfile.Target[targetName]
	if !exists {
		makelog.Fatalf("target does not exist")
	}

	vt := snlvar.NewVarTable()

	for k, v := range snlfile.Vars {
		if v != nil {
			vt.Set(k, *v)
		}
	}
	args, err := snlargs.LoadArgs(snlfile.Args, c.Args().Tail())
	if err != nil {
		PrintHelpText(snlfile)
		logrus.Fatal("argument error: ", err)
	}
	for k, v := range args {
		if _, exists := vt.Get(k); exists {
			logrus.WithField("name", k).Fatalf("arg already exists as variable")
		}
		vt.Set(k, *v)
	}

	timeout, ok := vt.Eval(targetCode.Timeout)
	if !ok {
		return fmt.Errorf("could not evaluate timeout expression: %s", timeout)
	}
	timeoutDuration, err := time.ParseDuration(timeout)
	if err != nil {
		return fmt.Errorf("could not parse timeout duration '%s': '%w", timeout, err)
	}

	addBuiltInTargetVars(vt, timeoutDuration, targetName)

	logrus.Info("starting execution of target: ", targetName)

	now := time.Now()

	exitStatus, err := snlexec.ShExec(targetCode.Sh, vt, timeoutDuration)
	if err != nil {
		logrus.WithField("target", targetName).Error("target execution failed: ", err)
	}
	if exitStatus != 0 {
		os.Exit(exitStatus)
	}

	logrus.Debug("ellapsed time:", time.Since(now))

	return nil
}

func addBuiltInTargetVars(vt *snlvar.VarTable, timeout time.Duration, targetName string) {
	vt.Set("SNL_TARGET_TIMEOUT_DURATION", timeout.String())
	vt.Set("SNL_TARGET_TIMEOUT_SEC", strconv.Itoa(int(timeout.Seconds())))
	vt.Set("SNL_TARGET_NAME", targetName)
}

package main

import (
	"os"

	"github.com/kassybas/shannel/cmd"
	"github.com/kassybas/shannel/internal/snlipc/snlrpc"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:          true,
		TimestampFormat:        "15:04:05",
		DisableLevelTruncation: true,
	})

	logrus.Trace("starting shannel")

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "file",
				Aliases: []string{"f"},
				Value:   "shannel.yaml",
				Usage:   "Use file as shannel file",
			},
			&cli.StringFlag{
				Name:    "loglevel",
				Aliases: []string{"l"},
				Value:   "warn",
				Usage:   "Logging level, one of error, warn, info, debug, trace",
			},
		},
		Commands: []*cli.Command{
			{
				Name:   "make",
				Usage:  "make a target",
				Action: cmd.Make,
			},
			{
				Name:  "publish",
				Usage: "publish a message to a channel",
				Action: func(c *cli.Context) error {
					channelName := c.Args().Get(1)
					msg := c.Args().Get(2)
					snlrpc.Send(channelName, msg, 1234)
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}

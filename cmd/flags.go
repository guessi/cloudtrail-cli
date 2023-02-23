package cmd

import (
	"time"

	"github.com/urfave/cli/v2"
)

var Flags = []cli.Flag{
	&cli.StringFlag{
		Name:     "profile",
		Aliases:  []string{"p"},
		Required: false,
	},
	&cli.StringFlag{
		Name:     "region",
		Aliases:  []string{"r"},
		Required: false,
	},
	&cli.TimestampFlag{
		Name:     "start-time",
		Aliases:  []string{"s"},
		Layout:   "2006-01-02T15:04:05",
		Timezone: time.UTC,
		Usage:    "Timestamp in 2023-01-01T00:00:00 format (UTC)",
		Required: false,
	},
	&cli.TimestampFlag{
		Name:     "end-time",
		Aliases:  []string{"e"},
		Layout:   "2006-01-02T15:04:05",
		Timezone: time.UTC,
		Usage:    "Timestamp in 2023-01-01T00:00:00 format (UTC)",
		Required: false,
	},
	&cli.StringFlag{
		Name:     "event-id",
		Usage:    "Filter events with event id",
		Required: false,
	},
	&cli.StringFlag{
		Name:     "event-name",
		Usage:    "Filter events with event name",
		Required: false,
	},
	&cli.StringFlag{
		Name:     "event-source",
		Usage:    "Filter events with event source",
		Required: false,
	},
	&cli.StringFlag{
		Name:     "user-name",
		Usage:    "Filter events with user name",
		Required: false,
	},
	&cli.BoolFlag{
		Name:     "read-only",
		Usage:    "Filter events with ReadOnly=true",
		Required: false,
	},
	&cli.BoolFlag{
		Name:     "no-read-only",
		Usage:    "Filter events with ReadOnly=false",
		Required: false,
	},
	&cli.StringFlag{
		Name:     "max-results",
		Aliases:  []string{"n"},
		Value:    "20",
		Required: false,
	},
	&cli.BoolFlag{
		Name:     "error-only",
		Usage:    "Filter events with errors",
		Required: false,
	},
}

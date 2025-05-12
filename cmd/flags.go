package cmd

import (
	"time"

	"github.com/urfave/cli/v3"
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
		Name:    "start-time",
		Aliases: []string{"s"},
		Config: cli.TimestampConfig{
			Layouts:  []string{time.RFC3339},
			Timezone: time.UTC,
		},
		Usage:    "Timestamp in RFC3339 format",
		Required: false,
	},
	&cli.TimestampFlag{
		Name:    "end-time",
		Aliases: []string{"e"},
		Config: cli.TimestampConfig{
			Layouts:  []string{time.RFC3339},
			Timezone: time.UTC,
		},
		Usage:    "Timestamp in RFC3339 format",
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
		Name:     "user-name",
		Usage:    "Filter events with user name",
		Required: false,
	},
	&cli.StringFlag{
		Name:     "resource-name",
		Usage:    "Filter events with resource name",
		Required: false,
	},
	&cli.StringFlag{
		Name:     "resource-type",
		Usage:    "Filter events with resource type",
		Required: false,
	},
	&cli.StringFlag{
		Name:     "event-source",
		Usage:    "Filter events with event source",
		Required: false,
	},
	&cli.StringFlag{
		Name:     "access-key-id",
		Usage:    "Filter events with access key id",
		Required: false,
	},
	&cli.BoolFlag{
		Name:     "read-only",
		Usage:    "Filter events with ReadOnly=true",
		Required: false,
	},
	&cli.IntFlag{
		Name:     "max-results",
		Aliases:  []string{"n"},
		Value:    20,
		Required: false,
	},
	&cli.BoolFlag{
		Name:     "error-only",
		Usage:    "Filter events with errors",
		Required: false,
	},
	&cli.BoolFlag{
		Name:     "truncate-user-name",
		Usage:    "Truncate user name string",
		Value:    false,
		Required: false,
	},
	&cli.BoolFlag{
		Name:     "truncate-user-agent",
		Usage:    "Truncate user agent string",
		Value:    false,
		Required: false,
	},
}

package cmd

import (
	"github.com/guessi/cloudtrail-cli/pkg/utils"
	"github.com/urfave/cli/v2"
)

func Wrapper(c *cli.Context) {
	utils.EventsHandler(
		c.String("profile"),
		c.String("region"),
		c.Timestamp("start-time"),
		c.Timestamp("end-time"),
		c.String("event-id"),
		c.String("event-name"),
		c.String("user-name"),
		c.String("resource-name"),
		c.String("resource-type"),
		c.String("event-source"),
		c.String("access-key-id"),
		c.Bool("read-only"),
		c.Bool("no-read-only"),
		c.Int("max-results"),
		c.Bool("error-only"),
		c.Bool("truncate-user-name"),
		c.Bool("truncate-user-agent"),
	)
}

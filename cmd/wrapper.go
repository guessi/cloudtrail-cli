package cmd

import (
	"github.com/guessi/cloudtrail-cli/pkg/utils"
	"github.com/urfave/cli/v2"
)

func QueryHandlerWrapper(c *cli.Context) {
	utils.EventsHandler(
		c.String("profile"),
		c.String("region"),
		c.Timestamp("start-time"),
		c.Timestamp("end-time"),
		c.String("event-id"),
		c.String("event-name"),
		c.String("event-source"),
		c.String("user-name"),
		c.Bool("read-only"),
		c.Bool("no-read-only"),
		c.Int("max-results"),
		c.Bool("error-only"),
	)
}

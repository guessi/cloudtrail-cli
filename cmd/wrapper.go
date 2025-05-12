package cmd

import (
	"github.com/guessi/cloudtrail-cli/pkg/types"
	"github.com/guessi/cloudtrail-cli/pkg/utils"
	"github.com/urfave/cli/v3"
)

func Wrapper(c *cli.Command) {
	var isReadOnlyFlagSet bool
	if c.IsSet("read-only") {
		isReadOnlyFlagSet = true
	}

	cloudTrailCliInput := types.CloudTrailCliInput{
		Profile:           c.String("profile"),
		Region:            c.String("region"),
		StartTime:         c.Timestamp("start-time"),
		EndTime:           c.Timestamp("end-time"),
		EventId:           c.String("event-id"),
		EventName:         c.String("event-name"),
		UserName:          c.String("user-name"),
		ResourceName:      c.String("resource-name"),
		ResourceType:      c.String("resource-type"),
		EventSource:       c.String("event-source"),
		AccessKeyId:       c.String("access-key-id"),
		IsReadOnlyFlagSet: isReadOnlyFlagSet,
		ReadOnly:          c.Bool("read-only"),
		MaxResults:        c.Int("max-results"),
		ErrorOnly:         c.Bool("error-only"),
		TruncateUserName:  c.Bool("truncate-user-name"),
		TruncateUserAgent: c.Bool("truncate-user-agent"),
	}

	utils.EventsHandler(cloudTrailCliInput)
}

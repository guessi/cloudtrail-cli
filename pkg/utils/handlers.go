package utils

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/guessi/cloudtrail-cli/pkg/types"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

func EventsHandler(i types.CloudTrailCliInput) {
	// do nothing if maxResults is invalid input
	if i.MaxResults <= 0 {
		log.Fatalln("Can not pass --max-results with a value lower or equal to 0.")
	}

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(i.Region),
		config.WithSharedConfigProfile(i.Profile),
	)
	if err != nil {
		log.Fatalf("Unable to load SDK config. Error: %s\n", err.Error())
	}

	svc := cloudtrail.NewFromConfig(cfg)

	lookupEventsInput := &cloudtrail.LookupEventsInput{
		StartTime: &i.StartTime,
		EndTime:   &i.EndTime,
		LookupAttributes: composeLookupAttributesInput(
			i.EventId,
			i.EventName,
			i.IsReadOnlyFlagSet,
			i.ReadOnly,
			i.UserName,
			i.ResourceName,
			i.ResourceType,
			i.EventSource,
			i.AccessKeyId,
		),
		MaxResults: getBatchSize(i.MaxResults),
	}

	events, err := LookupEvents(context.TODO(), svc, lookupEventsInput, i.MaxResults)
	if err != nil {
		log.Fatalf("Unable to ListTrails. Error: %s\n", err.Error())
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{
		"EventId",
		"EventName",
		"EventTime",
		"Username",
		"EventSource",
		"UserAgent",
		"SourceIPAddress",
		"AccessKeyId",
		"ErrorCode",
		"ReadOnly",
	})

	for _, event := range events {
		var c types.CloudTrailEvent
		if err := json.Unmarshal([]byte(*event.CloudTrailEvent), &c); err != nil {
			panic(err)
		}

		// early exit if errorOnly flag is set
		if i.ErrorOnly && len(c.ErrorCode) <= 0 {
			continue
		}

		username := getDisplayUserName(c.UserIdentity)

		t.AppendRow(table.Row{
			c.EventId,
			c.EventName,
			c.EventTime,
			truncateString(i.TruncateUserName, username, 24),
			c.EventSource,
			truncateString(i.TruncateUserAgent, c.UserAgent, 24),
			c.SourceIPAddress,
			c.UserIdentity.AccessKeyId,
			c.ErrorCode,
			c.ReadOnly,
		})
	}

	t.Style().Format.Header = text.FormatDefault
	t.Render()
}

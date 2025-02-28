package utils

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/guessi/cloudtrail-cli/pkg/types"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

func EventsHandler(profile, region string, startTime, endTime *time.Time, eventId, eventName, userName, resourceName, resourceType, eventSource, accessKeyId string, isReadOnlyFlagSet, readOnly bool, maxResults int, errorOnly, truncateUserName, truncateUserAgent bool) {
	// do nothing if maxResults is invalid input
	if maxResults <= 0 {
		log.Fatalf("Can not pass --max-results with a value lower or equal to 0.\n")
	}

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(region),
		config.WithSharedConfigProfile(profile),
	)
	if err != nil {
		log.Fatalf("Unable to load SDK config. Error: %s\n", err.Error())
	}

	svc := cloudtrail.NewFromConfig(cfg)

	input := &cloudtrail.LookupEventsInput{
		StartTime: startTime,
		EndTime:   endTime,
		LookupAttributes: composeLookupAttributesInput(
			eventId,
			eventName,
			isReadOnlyFlagSet,
			readOnly,
			userName,
			resourceName,
			resourceType,
			eventSource,
			accessKeyId,
		),
		MaxResults: getBatchSize(maxResults),
	}

	events, err := LookupEvents(context.TODO(), svc, input, maxResults)
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
		if errorOnly && len(c.ErrorCode) <= 0 {
			continue
		}

		username := getDisplayUserName(c.UserIdentity)

		t.AppendRow(table.Row{
			c.EventId,
			c.EventName,
			c.EventTime,
			truncateString(truncateUserName, username, 24),
			c.EventSource,
			truncateString(truncateUserAgent, c.UserAgent, 24),
			c.SourceIPAddress,
			c.UserIdentity.AccessKeyId,
			c.ErrorCode,
			c.ReadOnly,
		})
	}

	t.Style().Format.Header = text.FormatDefault
	t.Render()
}

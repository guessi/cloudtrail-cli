package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	ctypes "github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"
	"github.com/guessi/cloudtrail-cli/pkg/constants"
	"github.com/guessi/cloudtrail-cli/pkg/types"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// parseCloudTrailEvent safely parses CloudTrail event JSON with size validation
func parseCloudTrailEvent(event ctypes.Event) (*types.CloudTrailEvent, error) {
	if event.CloudTrailEvent == nil {
		return nil, fmt.Errorf("event data is nil")
	}

	if len(*event.CloudTrailEvent) > constants.MaxJSONPayloadSize {
		log.Printf("Skipping event with oversized JSON payload")
		return nil, fmt.Errorf("payload exceeds maximum size limit")
	}

	var cloudTrailEvent types.CloudTrailEvent
	if err := json.Unmarshal([]byte(*event.CloudTrailEvent), &cloudTrailEvent); err != nil {
		log.Printf("Failed to parse event JSON: skipping event")
		return nil, fmt.Errorf("failed to unmarshal event JSON: %w", err)
	}

	return &cloudTrailEvent, nil
}

// setDefaultTimeRange applies default time range (last 24 hours) when not specified
func setDefaultTimeRange(input *types.CloudTrailCliInput) {
	now := time.Now()
	if input.EndTime.IsZero() {
		input.EndTime = now
	}
	if input.StartTime.IsZero() {
		input.StartTime = input.EndTime.AddDate(0, 0, -1)
	}
}

// EventsHandler retrieves and displays CloudTrail events based on provided filters
// validateInput validates CloudTrail CLI input parameters
func validateInput(i types.CloudTrailCliInput) error {
	if i.MaxResults <= 0 {
		return fmt.Errorf("cannot pass --max-results with a value lower or equal to 0")
	}
	if i.MaxResults > constants.MaxCloudTrailResults {
		return fmt.Errorf("--max-results cannot exceed %d", constants.MaxCloudTrailResults)
	}
	return nil
}

// processEvents processes CloudTrail events and returns table rows
func processEvents(events []ctypes.Event, config types.CloudTrailCliInput) []table.Row {
	var rows []table.Row

	for _, event := range events {
		cloudTrailEvent, err := parseCloudTrailEvent(event)
		if err != nil {
			continue
		}

		// Apply error-only filter if requested
		if config.ErrorOnly && cloudTrailEvent.ErrorCode == "" {
			continue
		}

		username := getDisplayUserName(cloudTrailEvent.UserIdentity)
		rows = append(rows, table.Row{
			cloudTrailEvent.EventId,
			cloudTrailEvent.EventName,
			cloudTrailEvent.EventTime,
			truncateString(config.TruncateUserName, username, constants.DefaultTruncateLength),
			cloudTrailEvent.EventSource,
			truncateString(config.TruncateUserAgent, cloudTrailEvent.UserAgent, constants.DefaultTruncateLength),
			cloudTrailEvent.SourceIPAddress,
			cloudTrailEvent.UserIdentity.AccessKeyId,
			cloudTrailEvent.ErrorCode,
			cloudTrailEvent.ReadOnly,
		})
	}

	return rows
}

// createCloudTrailClient creates and configures AWS CloudTrail client
func createCloudTrailClient(ctx context.Context, region, profile string) (*cloudtrail.Client, error) {
	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(region),
		config.WithSharedConfigProfile(profile),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS configuration. Please check your credentials and region settings")
	}

	return cloudtrail.NewFromConfig(cfg), nil
}

// renderTable creates and renders the output table
func renderTable(rows []table.Row) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{
		"EventId", "EventName", "EventTime", "Username", "EventSource",
		"UserAgent", "SourceIPAddress", "AccessKeyId", "ErrorCode", "ReadOnly",
	})

	for _, row := range rows {
		t.AppendRow(row)
	}

	t.Style().Format.Header = text.FormatDefault
	t.Render()
}

// buildCloudTrailInput creates the CloudTrail API request input
func buildCloudTrailInput(i types.CloudTrailCliInput) *cloudtrail.LookupEventsInput {
	return &cloudtrail.LookupEventsInput{
		StartTime:        &i.StartTime,
		EndTime:          &i.EndTime,
		LookupAttributes: buildLookupAttributes(i),
		MaxResults:       getBatchSize(i.MaxResults),
	}
}

// LookupEventsFunc type for dependency injection
type LookupEventsFunc func(ctx context.Context, svc *cloudtrail.Client, input *cloudtrail.LookupEventsInput, maxResults int) ([]ctypes.Event, error)

// eventsHandlerWithLookup allows injection of LookupEvents function for testing
func eventsHandlerWithLookup(i types.CloudTrailCliInput, lookupFunc LookupEventsFunc) error {
	// Validate input parameters
	if err := validateInput(i); err != nil {
		return err
	}

	// Setup AWS client with timeout protection
	ctx, cancel := context.WithTimeout(context.Background(), constants.OperationTimeout)
	defer cancel()

	svc, err := createCloudTrailClient(ctx, i.Region, i.Profile)
	if err != nil {
		return err
	}

	// Configure time range and validate
	setDefaultTimeRange(&i)
	if i.StartTime.After(i.EndTime) {
		return fmt.Errorf("start time cannot be after end time")
	}

	// Build CloudTrail API request
	input := buildCloudTrailInput(i)

	// Retrieve events from CloudTrail
	events, err := lookupFunc(ctx, svc, input, i.MaxResults)
	if err != nil {
		return fmt.Errorf("unable to retrieve CloudTrail events. Please check your permissions and try again")
	}

	// Process and display events
	rows := processEvents(events, i)
	renderTable(rows)
	return nil
}

func EventsHandler(i types.CloudTrailCliInput) error {
	return eventsHandlerWithLookup(i, LookupEvents)
}

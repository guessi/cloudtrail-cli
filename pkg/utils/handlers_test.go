package utils

import (
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	ctypes "github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"
	"github.com/guessi/cloudtrail-cli/pkg/types"
)

func TestValidateInput(t *testing.T) {
	testCases := []struct {
		name        string
		input       types.CloudTrailCliInput
		expectError bool
	}{
		{
			"Valid input with positive max results",
			types.CloudTrailCliInput{
				MaxResults: 10,
			},
			false,
		},
		{
			"Invalid input with zero max results",
			types.CloudTrailCliInput{
				MaxResults: 0,
			},
			true,
		},
		{
			"Invalid input with negative max results",
			types.CloudTrailCliInput{
				MaxResults: -1,
			},
			true,
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			err := validateInput(tc.input)
			if (err != nil) != tc.expectError {
				t.Errorf("validateInput() error = %v, expectError %v", err, tc.expectError)
			}
		})
	}
}

func TestSetDefaultTimeRange(t *testing.T) {
	input := types.CloudTrailCliInput{}
	setDefaultTimeRange(&input)

	if input.StartTime.IsZero() || input.EndTime.IsZero() {
		t.Error("Default time range should set both start and end times")
	}

	if !input.StartTime.Before(input.EndTime) {
		t.Error("Start time should be before end time in default range")
	}
}

func TestBuildCloudTrailInput(t *testing.T) {
	startTime := time.Now().Add(-1 * time.Hour)
	endTime := time.Now()
	input := types.CloudTrailCliInput{
		MaxResults: 10,
		StartTime:  startTime,
		EndTime:    endTime,
	}
	result := buildCloudTrailInput(input)

	if result.MaxResults == nil || *result.MaxResults != int32(10) {
		t.Errorf("MaxResults = %v, want 10", result.MaxResults)
	}
}

func TestParseCloudTrailEvent(t *testing.T) {
	event := ctypes.Event{
		CloudTrailEvent: aws.String(`{
			"eventID": "test-event",
			"eventName": "TestEvent",
			"eventTime": "2023-01-01T12:00:00Z",
			"eventSource": "test.amazonaws.com"
		}`),
	}

	result, err := parseCloudTrailEvent(event)
	if err != nil {
		t.Fatalf("parseCloudTrailEvent() failed: %v", err)
	}

	if result.EventId != "test-event" {
		t.Errorf("EventId = %s, want test-event", result.EventId)
	}
}

func TestProcessEvents(t *testing.T) {
	events := []ctypes.Event{
		{
			CloudTrailEvent: aws.String(`{
				"eventID": "test-event",
				"eventName": "TestEvent",
				"eventTime": "2023-01-01T12:00:00Z",
				"eventSource": "test.amazonaws.com"
			}`),
		},
	}

	input := types.CloudTrailCliInput{
		TruncateUserName:  false,
		TruncateUserAgent: false,
	}

	rows := processEvents(events, input)
	if len(rows) == 0 {
		t.Error("Expected at least one table row from the event")
	}

	if len(rows) > 0 && len(rows[0]) > 0 {
		firstCell := rows[0][0].(string)
		if !strings.Contains(firstCell, "test-event") {
			t.Errorf("First table cell should contain event ID, got %v", firstCell)
		}
	}
}

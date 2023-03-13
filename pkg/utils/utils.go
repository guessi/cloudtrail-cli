package utils

import (
	"context"
	"strings"
	"unicode/utf8"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	ctypes "github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"
	"github.com/google/uuid"
	"github.com/guessi/cloudtrail-cli/pkg/types"
)

func isValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func truncateString(truncate bool, input string, maxLength int) string {
	if !truncate || len(input) < maxLength {
		return input
	}

	if utf8.ValidString(input[:maxLength]) {
		return input[:maxLength]
	}

	return input[:maxLength+1]
}

func getDisplayUserName(u types.UserIdentity) string {
	var username string
	switch t := u.Type; t {
	case "IAMUser":
		username = u.UserName
	case "WebIdentityUser":
		username = u.UserName
	case "AssumedRole":
		username = strings.Split(u.Arn, "/")[2]
	default:
		username = "-"
	}
	return username
}

func getBatchSize(i int) *int32 {
	const defaultBatchSize int = 50
	var r int32
	if i > 0 && i <= defaultBatchSize {
		r = int32(i)
	} else {
		r = int32(defaultBatchSize)
	}
	return &r
}

func LookupEvents(ctx context.Context, svc *cloudtrail.Client, input *cloudtrail.LookupEventsInput, maxResults int) ([]ctypes.Event, error) {
	var events []ctypes.Event
	var returnSize int

	paginator := cloudtrail.NewLookupEventsPaginator(svc, input)
	for paginator.HasMorePages() {
		out, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		events = append(events, out.Events...)
		if len(events) > maxResults {
			break
		}
	}

	returnSize = maxResults
	if len(events) < maxResults {
		returnSize = len(events)
	}

	return events[:returnSize], nil
}

func composeLookupAttributesInput(eventId, eventName string, readOnly, noReadOnly bool, userName, eventSource, accessKeyId string) []ctypes.LookupAttribute {
	lookupAttributesInput := []ctypes.LookupAttribute{}

	// LookupAttributeKeyEventId
	if isValidUUID(eventId) {
		attrEventId := ctypes.LookupAttribute{
			AttributeKey:   ctypes.LookupAttributeKeyEventId,
			AttributeValue: aws.String(eventId),
		}
		lookupAttributesInput = append(lookupAttributesInput, attrEventId)
	}

	// LookupAttributeKeyEventName
	if len(eventName) > 0 {
		attrEventName := ctypes.LookupAttribute{
			AttributeKey:   ctypes.LookupAttributeKeyEventName,
			AttributeValue: aws.String(eventName),
		}
		lookupAttributesInput = append(lookupAttributesInput, attrEventName)
	}

	// LookupAttributeKeyReadOnly
	var shouldPassReadonly bool
	var lookupAttributeKeyReadOnlyValue *string
	if readOnly != noReadOnly {
		shouldPassReadonly = true
		if readOnly {
			lookupAttributeKeyReadOnlyValue = aws.String("true")
		}
		if noReadOnly {
			lookupAttributeKeyReadOnlyValue = aws.String("false")
		}
	}
	if shouldPassReadonly {
		attrReadOnly := ctypes.LookupAttribute{
			AttributeKey:   ctypes.LookupAttributeKeyReadOnly,
			AttributeValue: lookupAttributeKeyReadOnlyValue,
		}
		lookupAttributesInput = append(lookupAttributesInput, attrReadOnly)
	}

	// LookupAttributeKeyUsername
	if len(userName) > 0 {
		attrUserName := ctypes.LookupAttribute{
			AttributeKey:   ctypes.LookupAttributeKeyUsername,
			AttributeValue: aws.String(userName),
		}
		lookupAttributesInput = append(lookupAttributesInput, attrUserName)
	}

	// LookupAttributeKeyEventSource
	const EVENT_SOURCE_SUFFIX = ".amazonaws.com"
	if len(eventSource) > len(EVENT_SOURCE_SUFFIX) && strings.HasSuffix(eventSource, EVENT_SOURCE_SUFFIX) {
		attrEventSource := ctypes.LookupAttribute{
			AttributeKey:   ctypes.LookupAttributeKeyEventSource,
			AttributeValue: aws.String(eventSource),
		}
		lookupAttributesInput = append(lookupAttributesInput, attrEventSource)
	}

	// LookupAttributeKeyAccessKeyId
	if len(accessKeyId) == 20 {
		attrAccessKeyId := ctypes.LookupAttribute{
			AttributeKey:   ctypes.LookupAttributeKeyAccessKeyId,
			AttributeValue: aws.String(accessKeyId),
		}
		lookupAttributesInput = append(lookupAttributesInput, attrAccessKeyId)
	}

	return lookupAttributesInput
}

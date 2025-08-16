package utils

import (
	"context"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	ctypes "github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"
	"github.com/google/uuid"
	"github.com/guessi/cloudtrail-cli/pkg/constants"
	"github.com/guessi/cloudtrail-cli/pkg/types"
)

// isValidEventSource checks if event source is a valid AWS service domain
func isValidEventSource(eventSource string) bool {
	return len(eventSource) > len(constants.AWSServiceSuffix) &&
		strings.HasSuffix(eventSource, constants.AWSServiceSuffix)
}

// isValidUUID validates UUID format using the google/uuid library
func isValidUUID(u string) bool {
	return u != "" && uuid.Validate(u) == nil
}

// isValidAccessKeyID validates AWS access key ID format (AKIA/ASIA prefixes only)
// Reference: https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_identifiers.html#identifiers-prefixes
func isValidAccessKeyID(accessKeyID string) bool {
	if len(accessKeyID) != 20 {
		return false
	}
	if !strings.HasPrefix(accessKeyID, "AKIA") && !strings.HasPrefix(accessKeyID, "ASIA") {
		return false
	}
	for _, char := range accessKeyID {
		if (char < 'A' || char > 'Z') && (char < '0' || char > '9') {
			return false
		}
	}
	return true
}

// truncateString safely truncates strings while preserving UTF-8 character boundaries
func truncateString(truncate bool, input string, maxLength int) string {
	if !truncate || len(input) <= maxLength {
		return input
	}

	// Find the largest valid UTF-8 boundary within maxLength
	for i := maxLength; i > 0; i-- {
		if utf8.ValidString(input[:i]) {
			return input[:i]
		}
	}

	return input[:maxLength] // fallback if no valid boundary found
}

// getDisplayUserName extracts a readable username from CloudTrail user identity
func getDisplayUserName(u types.UserIdentity) string {
	switch u.Type {
	case "IAMUser", "WebIdentityUser":
		return u.UserName
	case "AssumedRole":
		// Extract role name from ARN (format: arn:aws:sts::account:assumed-role/role-name/session-name)
		parts := strings.Split(u.Arn, "/")
		if len(parts) >= 3 {
			return parts[2]
		}
		return u.Arn // fallback to full ARN if parsing fails
	default:
		return "-"
	}
}

// getBatchSize returns the appropriate batch size for CloudTrail API pagination
func getBatchSize(requested int) *int32 {
	if requested > 0 && requested <= constants.DefaultBatchSize {
		return aws.Int32(int32(requested))
	}
	return aws.Int32(int32(constants.DefaultBatchSize))
}

// CloudTrailPaginator interface for dependency injection
type CloudTrailPaginator interface {
	HasMorePages() bool
	NextPage(ctx context.Context, optFns ...func(*cloudtrail.Options)) (*cloudtrail.LookupEventsOutput, error)
}

// LookupEvents retrieves CloudTrail events using AWS SDK pagination
func LookupEvents(ctx context.Context, svc *cloudtrail.Client, input *cloudtrail.LookupEventsInput, maxResults int) ([]ctypes.Event, error) {
	return LookupEventsWithPaginator(ctx, cloudtrail.NewLookupEventsPaginator(svc, input), maxResults)
}

// LookupEventsWithPaginator allows injection of paginator for testing
func LookupEventsWithPaginator(ctx context.Context, paginator CloudTrailPaginator, maxResults int) ([]ctypes.Event, error) {
	var events []ctypes.Event

	for paginator.HasMorePages() {
		out, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		events = append(events, out.Events...)
		if len(events) >= maxResults {
			break
		}
	}

	// Return only the requested number of events
	if len(events) > maxResults {
		return events[:maxResults], nil
	}

	return events, nil
}

// buildLookupAttributes constructs CloudTrail API filters from user input
func buildLookupAttributes(input types.CloudTrailCliInput) []ctypes.LookupAttribute {
	var attributes []ctypes.LookupAttribute

	// Add filters only for non-empty, valid inputs
	if isValidUUID(input.EventId) {
		attributes = append(attributes, ctypes.LookupAttribute{
			AttributeKey:   ctypes.LookupAttributeKeyEventId,
			AttributeValue: aws.String(input.EventId),
		})
	}

	if input.EventName != "" {
		attributes = append(attributes, ctypes.LookupAttribute{
			AttributeKey:   ctypes.LookupAttributeKeyEventName,
			AttributeValue: aws.String(input.EventName),
		})
	}

	if input.IsReadOnlyFlagSet {
		attributes = append(attributes, ctypes.LookupAttribute{
			AttributeKey:   ctypes.LookupAttributeKeyReadOnly,
			AttributeValue: aws.String(strconv.FormatBool(input.ReadOnly)),
		})
	}

	if input.UserName != "" {
		attributes = append(attributes, ctypes.LookupAttribute{
			AttributeKey:   ctypes.LookupAttributeKeyUsername,
			AttributeValue: aws.String(input.UserName),
		})
	}

	if input.ResourceName != "" {
		attributes = append(attributes, ctypes.LookupAttribute{
			AttributeKey:   ctypes.LookupAttributeKeyResourceName,
			AttributeValue: aws.String(input.ResourceName),
		})
	}

	if input.ResourceType != "" {
		attributes = append(attributes, ctypes.LookupAttribute{
			AttributeKey:   ctypes.LookupAttributeKeyResourceType,
			AttributeValue: aws.String(input.ResourceType),
		})
	}

	if isValidEventSource(input.EventSource) {
		attributes = append(attributes, ctypes.LookupAttribute{
			AttributeKey:   ctypes.LookupAttributeKeyEventSource,
			AttributeValue: aws.String(input.EventSource),
		})
	}

	if isValidAccessKeyID(input.AccessKeyId) {
		attributes = append(attributes, ctypes.LookupAttribute{
			AttributeKey:   ctypes.LookupAttributeKeyAccessKeyId,
			AttributeValue: aws.String(input.AccessKeyId),
		})
	}

	return attributes
}

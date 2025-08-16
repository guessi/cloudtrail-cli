package constants

import "time"

const (
	NAME  string = "cloudtrail-cli"
	USAGE string = "Blazing fast single purpose cli for CloudTrail log filtering"
)

const (
	// CloudTrail API limits
	MaxCloudTrailResults = 50000
	DefaultBatchSize     = 50

	// Memory and performance limits
	MaxJSONPayloadSize    = 1024 * 1024 // 1MB
	DefaultTruncateLength = 24
	OperationTimeout      = 5 * time.Minute

	// AWS service validation
	AWSServiceSuffix = ".amazonaws.com"
)

var (
	GitVersion string
	GoVersion  string
	BuildTime  string
)

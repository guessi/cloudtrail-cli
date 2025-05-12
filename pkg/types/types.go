package types

import "time"

type CloudTrailCliInput struct {
	Profile           string
	Region            string
	StartTime         time.Time
	EndTime           time.Time
	EventId           string
	EventName         string
	UserName          string
	ResourceName      string
	ResourceType      string
	EventSource       string
	AccessKeyId       string
	IsReadOnlyFlagSet bool
	ReadOnly          bool
	MaxResults        int
	ErrorOnly         bool
	TruncateUserName  bool
	TruncateUserAgent bool
}

// References:
// - https://docs.aws.amazon.com/awscloudtrail/latest/userguide/cloudtrail-event-reference.html
// - https://docs.aws.amazon.com/awscloudtrail/latest/userguide/cloudtrail-event-reference-record-contents.html
// - https://docs.aws.amazon.com/awscloudtrail/latest/userguide/cloudtrail-event-reference-user-identity.html

// AssumedRole
type SessionIssuer struct {
	Type        string `json:"type"`
	PrincipalId string `json:"principalId"`
	Arn         string `json:"arn"`
	AccountId   string `json:"accountId"`
	UserName    string `json:"userName"`
}

type WebIdFederationData struct {
	FederatedProvider string      `json:"federatedProvider"`
	Attributes        interface{} `json:"attributes"`
}

type Attributes struct {
	CreationDate     string `json:"creationDate"`
	MfaAuthenticated string `json:"mfaAuthenticated"`
}

type SessionContext struct {
	SessionIssuer       SessionIssuer       `json:"sessionIssuer"`
	WebIdFederationData WebIdFederationData `json:"webIdFederationData"`
	Attributes          Attributes          `json:"attributes"`
}

type UserIdentity struct {
	Type           string         `json:"type"`
	PrincipalId    string         `json:"principalId"`
	Arn            string         `json:"arn"`
	AccountId      string         `json:"accountId"`
	AccessKeyId    string         `json:"accessKeyId"`
	UserName       string         `json:"userName,omitempty"`
	SessionContext SessionContext `json:"sessionContext,omitempty"`
	InvokedBy      string         `json:"invokedBy"`
}

type CloudTrailEvent struct {
	EventVersion       string       `json:"eventVersion"`
	UserIdentity       UserIdentity `json:"userIdentity"`
	EventTime          string       `json:"eventTime"`
	EventSource        string       `json:"eventSource"`
	EventName          string       `json:"eventName"`
	AwsRegion          string       `json:"awsRegion"`
	SourceIPAddress    string       `json:"sourceIPAddress"`
	UserAgent          string       `json:"userAgent"`
	ErrorCode          string       `json:"errorCode,omitempty"`
	ErrorMessage       string       `json:"errorMessage,omitempty"`
	RequestParameters  interface{}  `json:"requestParameters"`
	ResponseElements   interface{}  `json:"responseElements"`
	RequestId          string       `json:"requestID"`
	EventId            string       `json:"eventID"`
	ReadOnly           bool         `json:"readOnly"`
	ManagementEvent    bool         `json:"managementEvent"`
	RecipientAccountId string       `json:"recipientAccountId"`
	EventCategory      string       `json:"eventCategory"`
}

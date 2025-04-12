package domain

import "github.com/aws/aws-sdk-go-v2/service/sesv2/types"

type SendEmailRequest struct {
	SenderEmail string
	ReplyTo     string
	Destination string
	ConfigSet   string
	EmailTags   []types.MessageTag
	Subject     string
	Body        string
}

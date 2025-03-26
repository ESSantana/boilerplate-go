package email

import (
	"context"
	"errors"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

type EmailManager struct {
	ses *sesv2.Client
}

func NewEmailManager() (*EmailManager, error) {
	awsRegion, found := os.LookupEnv("AWS_DEFAULT_REGION")
	if !found {
		return nil, errors.New("AWS_DEFAULT_REGION not found")
	}

	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(awsRegion),
	)

	if err != nil {
		return nil, err
	}

	ses := sesv2.NewFromConfig(cfg)
	return &EmailManager{ses: ses}, nil
}

type sendemail struct {
	SenderEmail string
	ReplyTo     string
	Destination string
	ConfigSet   string
	EmailTags   []types.MessageTag
	Subject     string
	Body        string
}

func (em *EmailManager) SendEmail(ctx context.Context, sendemail sendemail) (messageID string, err error) {
	sesInput := sesv2.SendEmailInput{
		FromEmailAddress: aws.String(sendemail.SenderEmail),
		ReplyToAddresses: []string{sendemail.ReplyTo},
		Destination: &types.Destination{
			ToAddresses: []string{sendemail.Destination},
		},
		ConfigurationSetName: aws.String(sendemail.ConfigSet),
		ListManagementOptions: &types.ListManagementOptions{
			ContactListName: aws.String("main"),
			TopicName:       aws.String("main"),
		},
		EmailTags: sendemail.EmailTags,
		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{
					Data: aws.String(sendemail.Subject),
				},
				Body: &types.Body{
					Html: &types.Content{
						Data: aws.String(sendemail.Body),
					},
				},
			},
		},
	}

	out, err := em.ses.SendEmail(ctx, &sesInput)
	if err != nil {
		return messageID, err
	}

	return *out.MessageId, nil
}

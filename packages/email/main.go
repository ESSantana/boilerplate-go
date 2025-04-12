package email

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/application-ellas/ella-backend/internal/domain/models"
	"github.com/application-ellas/ella-backend/packages/email/domain"
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

func (em *EmailManager) SendRecoverPasswordEmail(ctx context.Context, customer models.Customer) (err error) {
	if os.Getenv("ENV") == "local" {
		fmt.Println("Mock sending email to:", customer.Email)
		return nil
	}

	sendRequest := domain.SendEmailRequest{
		Destination: customer.Email,
	}

	messageID, err := em.sendEmail(ctx, sendRequest)
	if err != nil {
		return err
	}

	// TODO: save messageID to database or cache for tracking purposes
	fmt.Println("Message ID:", messageID)

	return nil
}

func (em *EmailManager) sendEmail(ctx context.Context, request domain.SendEmailRequest) (messageID string, err error) {
	sesInput := sesv2.SendEmailInput{
		FromEmailAddress: aws.String(request.SenderEmail),
		ReplyToAddresses: []string{request.ReplyTo},
		Destination: &types.Destination{
			ToAddresses: []string{request.Destination},
		},
		ConfigurationSetName: aws.String(request.ConfigSet),
		ListManagementOptions: &types.ListManagementOptions{
			ContactListName: aws.String("main"),
			TopicName:       aws.String("main"),
		},
		EmailTags: request.EmailTags,
		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{
					Data: aws.String(request.Subject),
				},
				Body: &types.Body{
					Html: &types.Content{
						Data: aws.String(request.Body),
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

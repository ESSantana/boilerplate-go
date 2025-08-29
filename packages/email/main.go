package email

import (
	"context"
	"fmt"
	"net/url"

	app_cfg "github.com/ESSantana/boilerplate-backend/internal/config"
	"github.com/ESSantana/boilerplate-backend/internal/domain/models"
	"github.com/ESSantana/boilerplate-backend/packages/email/domain"
	"github.com/ESSantana/boilerplate-backend/packages/email/templates"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

type EmailManager struct {
	cfg *app_cfg.Config
	ses *sesv2.Client
}

func NewEmailManager(appCfg *app_cfg.Config) (*EmailManager, error) {
	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(appCfg.AWS.DefaultRegion),
	)

	if err != nil {
		return nil, err
	}

	ses := sesv2.NewFromConfig(cfg)
	return &EmailManager{ses: ses, cfg: appCfg}, nil
}

func (em *EmailManager) SendRecoverPasswordEmail(ctx context.Context, customer models.Customer) (err error) {
	if em.cfg.Server.Environment == "development" {
		fmt.Println("Mock sending email to:", customer.Email)
		return nil
	}

	base, _ := url.Parse(em.cfg.Frontend.AuthRedirect)
	base.Path = "/recover-password"
	q := base.Query()
	q.Set("email", customer.Email)
	base.RawQuery = q.Encode()

	subject := "Redefinição de senha"
	htmlBody := templates.RecoverPasswordHTML("Boilerplate App", customer.Name, base.String())

	sendRequest := domain.SendEmailRequest{
		SenderEmail: em.cfg.AWS.SESSenderEmail,
		ReplyTo:     em.cfg.AWS.SESReplyTo,
		Destination: customer.Email,
		ConfigSet:   em.cfg.AWS.SESConfigSet,
		Subject:     subject,
		Body:        htmlBody,
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

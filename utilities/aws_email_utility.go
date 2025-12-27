package utilities

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

// AmazonSesUtility wraps the AWS SES client for sending emails.
type AmazonSesUtility struct {
	client *ses.Client
}

// NewAmazonSesUtility initializes an AWS SES client and returns a new AmazonSesUtility instance.
func NewAmazonSesUtility() (*AmazonSesUtility, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := ses.NewFromConfig(cfg)

	return &AmazonSesUtility{
		client: client,
	}, nil
}

// SendEmail sends an email to a recipient with the provided content using SendRawEmail.
// recipient: the email address of the recipient
// subject: the email subject line
// content: the email body content
func (e *AmazonSesUtility) SendEmail(ctx context.Context, recipient string, subject string, content string) error {
	if recipient == "" {
		return errors.New("recipient is required")
	}
	if subject == "" {
		return errors.New("subject is required")
	}
	if content == "" {
		return errors.New("content is required")
	}

	// Get sender email from environment variable
	sender := os.Getenv("AWS_SES_SENDER_EMAIL")
	if sender == "" {
		return errors.New("AWS_SES_SENDER_EMAIL environment variable is not set")
	}

	// Create raw email message
	rawMessage := e.createRawEmailMessage(sender, recipient, subject, content)

	// Create SendRawEmail input
	// RawMessage.Data expects the raw email bytes (not base64 encoded)
	input := &ses.SendRawEmailInput{
		RawMessage: &types.RawMessage{
			Data: []byte(rawMessage),
		},
	}

	// Send the email
	_, err := e.client.SendRawEmail(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// createRawEmailMessage creates a raw email message in MIME format.
func (e *AmazonSesUtility) createRawEmailMessage(sender, recipient, subject, content string) string {
	var buf bytes.Buffer

	// Email headers
	buf.WriteString(fmt.Sprintf("From: %s\n", sender))
	buf.WriteString(fmt.Sprintf("To: %s\n", recipient))
	buf.WriteString(fmt.Sprintf("Subject: %s\n", subject))
	buf.WriteString("MIME-Version: 1.0\n")
	buf.WriteString("Content-Type: text/plain; charset=UTF-8\n")
	buf.WriteString("Content-Transfer-Encoding: 7bit\n")
	buf.WriteString("\n")

	// Email body
	buf.WriteString(content)
	buf.WriteString("\r\n")

	return buf.String()
}

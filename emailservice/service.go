package service

import (
	"apidootoday/config"
	"log"

	sendgrid "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// EmailService :
type EmailService struct {
	Client *sendgrid.Client
}

// NewEmailService :
func NewEmailService() *EmailService {
	return &EmailService{
		Client: sendgrid.NewSendClient(config.SendgridKey),
	}
}

// SendEmail :
func (s *EmailService) SendWelcomeEmail(
	toEmail string,
	toName string,
	shortName string,
) error {
	from := mail.NewEmail("Doo.Today", "contact@doo.today")
	to := mail.NewEmail(toName, toEmail)
	sgMail := mail.NewV3Mail()
	sgMail.SetFrom(from)
	sgMail.SetReplyTo(from)
	p := mail.NewPersonalization()
	p.AddTos(to)
	p.SetDynamicTemplateData("name", shortName)
	p.SetDynamicTemplateData("subject", "Welcomet to Doo.Today")
	sgMail.AddPersonalizations(p)
	sgMail.SetTemplateID("d-d30ca85798ab4d1faf927c8240e79715")
	response, err := s.Client.Send(sgMail)

	if err != nil {
		log.Println(err)
	} else {
		log.Println(response.StatusCode)
		log.Println(response.Body)
		log.Println(response.Headers)
	}
	return err
}

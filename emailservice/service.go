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

// SendWelcomeEmail :
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

// SendTaskMoveEmail :
func (s *EmailService) SendTaskMoveEmail(
	toEmail string,
	toName string,
	shortName string,
	tasks []string,
) error {
	from := mail.NewEmail("Doo.Today", "contact@doo.today")
	to := mail.NewEmail(toName, toEmail)
	sgMail := mail.NewV3Mail()
	sgMail.SetFrom(from)
	sgMail.SetReplyTo(from)
	p := mail.NewPersonalization()
	p.AddTos(to)
	p.SetDynamicTemplateData("name", shortName)
	p.SetDynamicTemplateData("subject", "Moving undone tasks")
	p.SetDynamicTemplateData("tasks", tasks)
	sgMail.AddPersonalizations(p)
	sgMail.SetTemplateID("d-0398b256912942f983f3e14f4ba335de")
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

// SendEmptyListReminder :
func (s *EmailService) SendEmptyListReminder(
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
	p.SetDynamicTemplateData("subject", "What you gotta do today?")
	sgMail.AddPersonalizations(p)
	sgMail.SetTemplateID("d-5a3e1944da44407182d8a46e81b32acd")
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

// SendYouHaveTasks :
func (s *EmailService) SendYouHaveTasks(
	toEmail string,
	toName string,
	shortName string,
	tasks []string,
) error {
	from := mail.NewEmail("Doo.Today", "contact@doo.today")
	to := mail.NewEmail(toName, toEmail)
	sgMail := mail.NewV3Mail()
	sgMail.SetFrom(from)
	sgMail.SetReplyTo(from)
	p := mail.NewPersonalization()
	p.AddTos(to)
	p.SetDynamicTemplateData("name", shortName)
	p.SetDynamicTemplateData("subject", "Hey you gotta things to do today!")
	p.SetDynamicTemplateData("tasks", tasks)
	sgMail.AddPersonalizations(p)
	sgMail.SetTemplateID("d-98b006cbecbe4db795a020b285b70bff")
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

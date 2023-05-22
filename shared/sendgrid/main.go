package sendgrid

import (
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	parser "net/mail"
	"os"
)

// TODO: change URL address
const CURRENTHOST = "http://209.38.200.216:8080"

type Email struct {
	From    string
	To      string
	Subject string
	Body    string
}

type ConfirmationEmail struct {
	To             string
	NewsletterName string
}

type NewsletterEmail struct {
	To               string
	Subject          string
	UnsubscribeToken string
	Body             string
}

func createEmail(email Email) (*mail.SGMailV3, error) {
	_, err := parser.ParseAddress(email.To)
	if err != nil {
		return nil, err
	}

	if email.Subject == "" {
		return nil, err
	}

	if email.Body == "" {
		return nil, err
	}

	from := mail.NewEmail("Application", "am7642939@gmail.com")
	to := mail.NewEmail("Destination User", email.To)
	subject := email.Subject
	body := email.Body

	message := mail.NewSingleEmail(from, subject, to, body, body)
	return message, nil
}

func createNewSendClient() *sendgrid.Client {
	return sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
}

func SendNewsletterEmail(email NewsletterEmail) (*rest.Response, error) {
	newEmail := Email{
		To:      email.To,
		From:    "am7642939@gmail.com",
		Subject: email.Subject,
		Body: email.Body + "<br />" +
			"To unsubscribe click from this newsletter click the following link: <a href=" + CURRENTHOST + "/newsletter-api/v1/newsletter/" + email.UnsubscribeToken + "/unsubscribe\">Unsubscribe</a>",
	}
	message, err := createEmail(newEmail)
	if err != nil {
		return nil, err
	}

	client := createNewSendClient()
	response, err := client.Send(message)
	if err != nil {
		return nil, err
	} else {
		return response, nil
	}
}

func SendSubscribeConfirmationEmail(email ConfirmationEmail, token string) (*rest.Response, error) {

	unsublink := CURRENTHOST + "/newsletter-api/v1/newsletter/" + token + "/unsubscribe"

	newEmail := Email{
		To:      email.To,
		From:    "am7642939@gmail.com",
		Subject: "Welcome to our newsletter: " + email.NewsletterName,
		Body: "You have successfully subscribed to our newsletter: " + email.NewsletterName + ". <br />" +
			"To unsubscribe click on the following link: <a href=" + unsublink + ">Unsubscribe</a>",
	}

	message, err := createEmail(newEmail)
	if err != nil {
		return nil, err
	}

	client := createNewSendClient()
	response, err := client.Send(message)
	if err != nil {
		return nil, err
	} else {
		return response, nil
	}
}

func SendUnsubscribeConfirmationEmail(email ConfirmationEmail) (*rest.Response, error) {
	newEmail := Email{
		To:      email.To,
		From:    "am7642939@gmail.com",
		Subject: "Succesfully Unsubscribed from " + email.NewsletterName,
		Body:    "You have successfully unsubscribed from newsletter: " + email.NewsletterName,
	}

	message, err := createEmail(newEmail)
	if err != nil {
		return nil, err
	}

	client := createNewSendClient()
	response, err := client.Send(message)
	if err != nil {
		return nil, err
	} else {
		return response, nil
	}
}

func SendNewsletterDeletionEmail(email ConfirmationEmail) (*rest.Response, error) {
	newEmail := Email{
		To:      email.To,
		From:    "am7642939@gmail.com",
		Subject: "Newsletter " + email.NewsletterName + " has been deleted",
		Body:    "Newsletter " + email.NewsletterName + " has been deleted",
	}

	message, err := createEmail(newEmail)
	if err != nil {
		return nil, err
	}

	client := createNewSendClient()
	response, err := client.Send(message)
	if err != nil {
		return nil, err
	} else {
		return response, nil
	}
}

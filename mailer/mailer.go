package mailer

import (
	"fmt"
	"strings"
	"net/smtp"
	"log"
)

func SendMail(recipients []string, subject string, messageBody string, from string, username string, password string, host string, port int, bcc string) error {
	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = strings.Join(recipients,",")
	headers["Subject"] = subject

	message := ""
	for k,v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + messageBody

	auth := smtp.PlainAuth(
		"",
		username,
		password,
		host,
	)
	if len(bcc) > 0 {
		recipients = append(recipients, bcc)
	}


	serverAddressAndPort := fmt.Sprint(host , ":" , port)
	err := smtp.SendMail(
		serverAddressAndPort,
		auth,
		from,
		recipients,
		[]byte(message),
	)
	log.Println(err)
	return err
}
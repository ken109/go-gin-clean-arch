package mail

import (
	_ "embed"
	"fmt"
	"sync"

	"gopkg.in/gomail.v2"

	"go-gin-clean-arch/config"
	"go-gin-clean-arch/domain/service/mailf"
	"go-gin-clean-arch/packages/errors"
)

type Sender interface {
	Send(m mailf.Mail) error
	BulkSend(mails []mailf.Mail) (err error)
}

type sender struct{}

func NewMailSender() Sender {
	return &sender{}
}

func (s *sender) Send(m mailf.Mail) error {
	to, err := m.To()
	if err != nil {
		return err
	}
	title, err := m.Title()
	if err != nil {
		return err
	}
	content, err := m.Content()
	if err != nil {
		return err
	}

	message := gomail.NewMessage()

	message.SetBody("text/html", content)

	message.SetHeaders(
		map[string][]string{
			"From":    {message.FormatAddress(config.Env.Mail.From, config.Env.Mail.Name)},
			"To":      {to},
			"Subject": {fmt.Sprintf("【%s】%s", config.Env.App.Name, title)},
		},
	)

	dialer := gomail.NewDialer(config.Env.SMTP.Host, config.Env.SMTP.Port, config.Env.SMTP.User, config.Env.SMTP.Password)
	err = dialer.DialAndSend(message)
	if err != nil {
		return errors.NewUnexpected(err)
	}

	return nil
}

func (s *sender) BulkSend(mails []mailf.Mail) (err error) {
	if len(mails) == 0 {
		return nil
	}

	dialer := gomail.NewDialer(config.Env.SMTP.Host, config.Env.SMTP.Port, config.Env.SMTP.User, config.Env.SMTP.Password)

	errChan := make(chan error, len(mails))

	var wg sync.WaitGroup
	wg.Add(len(mails))

	for _, m := range mails {
		go func(m mailf.Mail) {
			defer func() {
				if r := recover(); r != nil {
					errChan <- errors.NewUnexpected(r.(error), errors.WithUnexpectedPanic{})
				}
			}()
			defer wg.Done()

			to, err := m.To()
			if err != nil {
				errChan <- err
				return
			}
			title, err := m.Title()
			if err != nil {
				errChan <- err
				return
			}
			content, err := m.Content()
			if err != nil {
				errChan <- err
				return
			}

			message := gomail.NewMessage()

			message.SetBody("text/html", content)

			message.SetHeaders(
				map[string][]string{
					"From":    {message.FormatAddress(config.Env.Mail.From, config.Env.Mail.Name)},
					"To":      {to},
					"Subject": {fmt.Sprintf("【%s】%s", config.Env.App.Name, title)},
				},
			)

			err = dialer.DialAndSend(message)
			if err != nil {
				errChan <- errors.NewUnexpected(err)
				return
			}
		}(m)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err = range errChan {
		return err
	}

	return nil
}

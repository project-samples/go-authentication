package mail

import (
	"fmt"
	"github.com/core-go/mail"
)

type MockMailSender struct {
}

func NewMockMailSender() *MockMailSender {
	return &MockMailSender{}
}

func (s *MockMailSender) Send(m mail.SimpleMail) error {
	var contents = make([]mail.Content, len(m.Content))
	for i, content := range m.Content {
		contents[i] = content
	}
	mail := mail.NewMailInit(m.From, m.Subject, m.To, m.Cc, contents...)
	if len(mail.Content) > 0 {
		fmt.Println(mail.Content[0].Value)
	}
	return nil
}

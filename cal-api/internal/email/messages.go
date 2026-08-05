package email

import (
	"bytes"
	"fmt"
	"html/template"

	"gopkg.in/gomail.v2"
)

type GomailMessages struct{}

func (gm GomailMessages) SendOtp(to, otp string) (*gomail.Message, error) {
	m := gomail.NewMessage()

	m.SetHeader("From", "homeshare@wburg.dev")
	m.SetHeader("To", to)

	m.SetHeader("Subject", "Home Share login code")
	m.SetBody("text/plain", fmt.Sprintf("Your login otp is: %s", otp))

	tmplData := map[string]string{
		"Otp": otp,
	}

	t, err := template.New("email").Parse(`
		<html>
			<body>
				<h1>Sign in to Home Share</h1>
				<p>Your login code is:</p>
				<p>{{.Otp}}</p>
			</body>
		</html>
	`)
	if err != nil {
		return m, err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, tmplData); err != nil {
		return m, err
	}

	m.AddAlternative("text/html", buf.String())

	return m, nil
}

func NewGomailMessages() GomailMessages {
	return GomailMessages{}
}

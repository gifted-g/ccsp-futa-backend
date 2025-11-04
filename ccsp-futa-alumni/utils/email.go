package utils

import (
    "gopkg.in/gomail.v2"
    "os"
    "strconv"
)

func SendEmail(to, subject, body string) error {
    port, _ := strconv.Atoi(os.Getenv("EMAIL_PORT"))

    m := gomail.NewMessage()
    m.SetHeader("From", os.Getenv("EMAIL_USER"))
    m.SetHeader("To", to)
    m.SetHeader("Subject", subject)
    m.SetBody("text/plain", body)

    d := gomail.NewDialer(
        os.Getenv("EMAIL_HOST"),
        port,
        os.Getenv("EMAIL_USER"),
        os.Getenv("EMAIL_PASS"),
    )

    return d.DialAndSend(m)
}

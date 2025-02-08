package services

import (
	"fmt"
	"log"
	"math/rand"
	"myproject/pkg/config"
	"net/smtp"

	"strconv"
	"time"
)

type Services interface {
	GenerateOtp(length int) int

	SendEmailWithOTP(email string) (string, error)
}
type MyService struct {
	Config config.Config
}

func (s MyService) GenerateOtp(length int) int {
	rand.Seed(time.Now().UnixNano())

	randomNum := rand.Intn(90000) + 10000

	fmt.Println("Random 5-digit number:", randomNum)
	return randomNum
}

func (s MyService) SendEmailWithOTP(email string) (string, error) {

	otp := strconv.Itoa(s.GenerateOtp(6))

	message := fmt.Sprintf("Subject: OTP for Verification\n\nYour OTP is: %s", otp)
	fmt.Println("this is my email  !!!!!", s.Config.SMTPemail, "this is my email  !!!!!", s.Config.Password)

	SMTPemail := s.Config.SMTPemail
	SMTPpass := s.Config.Password
	auth := smtp.PlainAuth("", "adithyanunni258@gmail.com", SMTPpass, "smtp.gmail.com")

	err := smtp.SendMail("smtp.gmail.com:587", auth, SMTPemail, []string{email}, []byte(message))
	if err != nil {
		log.Println("Error sending email:", err)
		return "", err
	}

	return otp, nil
}

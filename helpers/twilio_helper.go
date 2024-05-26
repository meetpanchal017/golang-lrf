package helper

import (
	"errors"
	"os"

	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/verify/v2"
)

var client *twilio.RestClient = twilio.NewRestClientWithParams(twilio.ClientParams{
	Username: os.Getenv("TWILIO_ACCOUNT_SID"),
	Password: os.Getenv("TWILIO_AUTH_TOKEN"),
})

func TwilioSendOTP(phoneNumber string) (string, error) {
	params := &twilioApi.CreateVerificationParams{}

	params.SetTo(phoneNumber)
	params.SetChannel("sms")
	res, err := client.VerifyV2.CreateVerification(os.Getenv("VERIFY_SERVICE_SID"), params)
	if err != nil {
		return "", err
	}

	return *res.Sid, nil
}

func TwilioVerifyOTP(phoneNumber string, code string) error {
	params := &twilioApi.CreateVerificationCheckParams{}
	params.SetTo(phoneNumber)
	params.SetCode(code)

	resp, err := client.VerifyV2.CreateVerificationCheck(os.Getenv("VERIFY_SERVICE_SID"), params)
	if err != nil {
		return err
	}

	// BREAKING CHANGE IN THE VERIFY API
	// https://www.twilio.com/docs/verify/quickstarts/verify-totp-change-in-api-response-when-authpayload-is-incorrect
	if *resp.Status != "approved" {
		return errors.New("not a valid code")
	}

	return nil
}

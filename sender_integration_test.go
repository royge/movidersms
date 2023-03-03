//go:build integration
// +build integration

package movidersms_test

import (
	"context"
	"os"
	"testing"

	"github.com/royge/movidersms"
)

func checkEnv(t *testing.T, envs ...string) {
	t.Helper()

	for _, env := range envs {
		env := env // shadow copy

		if v := os.Getenv(env); v == "" {
			t.Fatalf("environment variable `%s` is not defined", env)
		}
	}
}

func Test_Sender_SendMessage_Integration(t *testing.T) {
	checkEnv(
		t,
		"MOVIDER_API_KEY",
		"MOVIDER_API_SECRET",
		"ALLOWED_RECIPIENT_NUMBER",
		"NOT_ALLOWED_RECIPIENT_NUMBER",
	)

	var (
		apiKey              = os.Getenv("MOVIDER_API_KEY")
		secretKey           = os.Getenv("MOVIDER_API_SECRET")
		allowedRecipient    = os.Getenv("ALLOWED_RECIPIENT_NUMBER")
		notAllowedRecipient = os.Getenv("NOT_ALLOWED_RECIPIENT_NUMBER")
	)

	allowedRecipients := []string{allowedRecipient}
	sender := movidersms.NewSender(movidersms.Credentials{apiKey, secretKey}, allowedRecipients)

	tests := []struct {
		scenario  string
		recipient string
		check     func(t *testing.T, res *movidersms.SendMessageResponse)
	}{
		{
			scenario:  "send sms to allowed phone number",
			recipient: allowedRecipient,
			check: func(t *testing.T, res *movidersms.SendMessageResponse) {
				t.Helper()
			},
		},
		{
			scenario:  "send sms to not allowed phone number",
			recipient: notAllowedRecipient,
			check: func(t *testing.T, res *movidersms.SendMessageResponse) {
				t.Helper()

				if res != nil {
					t.Errorf("want nil response")
				}
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.scenario, func(t *testing.T) {
			res, err := sender.SendMessage(
				context.Background(),
				[]string{tc.recipient},
				"TEST 1234 QWERTY",
			)
			if err != nil {
				t.Fatalf("unable to send message: %v", err)
			}

			tc.check(t, res)

			t.Logf("Results: %+v\n", res)
		})
	}
}

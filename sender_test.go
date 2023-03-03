package movidersms_test

import (
	"reflect"
	"testing"

	"github.com/royge/movidersms"
)

func Test_NewSender(t *testing.T) {
	creds := movidersms.Credentials{"apikey", "apisecret"}
	allowedRecipients := []string{}
	sender := movidersms.NewSender(creds, allowedRecipients)

	want := "https://api.movider.co/v1/sms"
	if sender.APIURL != want {
		t.Errorf(
			"want API URL to be `%v`, got `%v`",
			want,
			sender.APIURL,
		)
	}
}

func Test_SendMessageRequest_Encode(t *testing.T) {
	creds := movidersms.Credentials{"apikey", "apisecret"}
	credStr := "api_key=apikey&api_secret=apisecret"

	phone1 := "639123456789"
	phone2 := "639123456788"
	tests := []struct {
		name  string
		input *movidersms.SendMessageRequest
		want  string
	}{
		{
			name: "1 word message",
			input: &movidersms.SendMessageRequest{
				To:   []string{phone1},
				Text: "test",
			},
			want: credStr + "&to=639123456789&text=test",
		},
		{
			name: "2 to number",
			input: &movidersms.SendMessageRequest{
				To:   []string{phone1, phone2},
				Text: "test",
			},
			want: credStr + "&to=639123456789,639123456788&text=test",
		},
		{
			name: "2 words message",
			input: &movidersms.SendMessageRequest{
				To:   []string{phone1},
				Text: "test only",
			},
			want: credStr + "&to=639123456789&text=test+only",
		},
		{
			name: "3 words message",
			input: &movidersms.SendMessageRequest{
				To:   []string{phone1},
				Text: "test only now",
			},
			want: credStr + "&to=639123456789&text=test+only+now",
		},
	}

	for _, tc := range tests {
		tc.input.Credentials = creds

		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := tc.input.Encode()
			if tc.want != string(got) {
				t.Errorf("want `%v`, got `%s`", tc.want, got)
			}
		})
	}
}

func Test_SendMessageRequest_Validate(t *testing.T) {
	phone := "639123456789"
	tests := []struct {
		name  string
		input *movidersms.SendMessageRequest
		want  error
	}{
		{
			name: "empty destination",
			input: &movidersms.SendMessageRequest{
				To:   nil,
				Text: "test",
			},
			want: movidersms.ErrNoDestination,
		},
		{
			name: "empty destination",
			input: &movidersms.SendMessageRequest{
				To:   []string{},
				Text: "test",
			},
			want: movidersms.ErrNoDestination,
		},
		{
			name: "empty destination",
			input: &movidersms.SendMessageRequest{
				To:   []string{""},
				Text: "test",
			},
			want: movidersms.ErrNoDestination,
		},
		{
			name: "empty destination",
			input: &movidersms.SendMessageRequest{
				To:   []string{" "},
				Text: "test",
			},
			want: movidersms.ErrNoDestination,
		},
		{
			name: "empty destination",
			input: &movidersms.SendMessageRequest{
				To:   []string{"      "},
				Text: "test",
			},
			want: movidersms.ErrNoDestination,
		},
		{
			name: "invalid destination",
			input: &movidersms.SendMessageRequest{
				To:   []string{"1234567890123456"},
				Text: "test",
			},
			want: movidersms.ErrInvalidDestination,
		},
		{
			name: "empty text",
			input: &movidersms.SendMessageRequest{
				To:   []string{phone},
				Text: "",
			},
			want: movidersms.ErrNoText,
		},
		{
			name: "empty text",
			input: &movidersms.SendMessageRequest{
				To:   []string{phone},
				Text: " ",
			},
			want: movidersms.ErrNoText,
		},
		{
			name: "empty text",
			input: &movidersms.SendMessageRequest{
				To:   []string{phone},
				Text: "      ",
			},
			want: movidersms.ErrNoText,
		},
		{
			name: "valid",
			input: &movidersms.SendMessageRequest{
				To:   []string{phone},
				Text: "test only",
			},
			want: nil,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := tc.input.Validate()
			if tc.want != got {
				t.Errorf("want `%v`, got `%s`", tc.want, got)
			}
		})
	}
}

func Test_MakeValidPhoneNumbers(t *testing.T) {
	phones := []string{
		"092623456789",
		"095823456789",
		"6395823456788",
		"+6395823456787",
	}

	want := []string{
		"6392623456789",
		"6395823456789",
		"6395823456788",
		"6395823456787",
	}

	err := movidersms.MakeValidPhoneNumbers(phones)
	if err != nil {
		t.Fatalf("unable to make valid phone numbers: %v", err)
	}

	if !reflect.DeepEqual(want, phones) {
		t.Errorf("want `%v`, got `%v`", want, phones)
	}
}

package movidersms

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/nyaruka/phonenumbers"
)

const (
	MoviderAPIURL  = "https://api.movider.co/v1"
	MaxPhoneLength = 15
)

var (
	ErrNoDestination      = errors.New("no destination number")
	ErrInvalidDestination = errors.New("recipient number is more than 15 characters")
	ErrNoText             = errors.New("no text message")
)

type Credentials struct {
	APIKey    string `json:"api_key"`
	APISecret string `json:"api_secret"`
}

type Config struct {
	APIURL            string
	Creds             Credentials
	AllowedRecipients []string
}

type Sender struct {
	Config
}

func NewSender(creds Credentials, allowedRecipients []string) *Sender {
	return &Sender{
		Config{
			APIURL:            MoviderAPIURL + "/sms",
			Creds:             creds,
			AllowedRecipients: allowedRecipients,
		},
	}
}

type SendMessageRequest struct {
	Credentials

	To   []string `json:"to"`
	Text string   `json:"text"`
}

func makeValidPhoneNumbers(phones []string) error {
	for i, phone := range phones {
		if strings.HasPrefix(phone, "63") {
			continue
		}

		pn, err := phonenumbers.Parse(phone, "PH")
		if err != nil {
			return err
		}
		phones[i] = strings.Replace(
			phonenumbers.Format(pn, phonenumbers.E164),
			"+",
			"",
			1,
		)
	}

	return nil
}

func (smr *SendMessageRequest) Validate() error {
	if len(smr.To) == 0 {
		return ErrNoDestination
	}
	for _, dest := range smr.To {
		dest := dest
		if strings.TrimSpace(dest) == "" {
			return ErrNoDestination
		}
		if len(dest) > MaxPhoneLength {
			return ErrInvalidDestination
		}
	}
	if strings.TrimSpace(smr.Text) == "" {
		return ErrNoText
	}

	return nil
}

func (smr *SendMessageRequest) Encode() []byte {
	tpl := "api_key=%s&api_secret=%s&to=%s&text=%s"
	text := url.QueryEscape(smr.Text)
	dest := strings.Join(smr.To, ",")
	return []byte(fmt.Sprintf(tpl, smr.APIKey, smr.APISecret, dest, text))
}

type Error struct {
	Code        int    `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (err *Error) Error() string {
	return fmt.Sprintf("code: %v, name: %v, desc: %v", err.Code, err.Name, err.Description)
}

type SendMessageResponse struct {
	RemainingBalance float64 `json:"remaining_balance"`
	TotalSMS         int     `json:"total_sms"`
	Error            *Error  `json:"error"`
}

func (s *Sender) SendMessage(
	ctx context.Context, recipient []string, message string,
) (interface{}, error) {
	smreq := &SendMessageRequest{
		Credentials: s.Creds,
		To:          recipient,
		Text:        message,
	}

	if err := smreq.Validate(); err != nil {
		return nil, err
	}

	if err := makeValidPhoneNumbers(smreq.To); err != nil {
		return nil, err
	}

	if allowed := s.isAllowedPhone(smreq.To[0]); !allowed {
		log.Println("SMS not sent to: ", smreq.To[0])

		return nil, nil
	}

	req, err := http.NewRequest(
		"POST",
		s.APIURL,
		bytes.NewBuffer(smreq.Encode()),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	r, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	smres := &SendMessageResponse{}
	if err := json.Unmarshal(r, smres); err != nil {
		return nil, err
	}

	if smres.Error != nil {
		return nil, smres.Error
	}

	return smres, nil
}

func (s *Sender) isAllowedPhone(phone string) bool {
	if len(s.AllowedRecipients) == 0 {
		return true
	}

	for _, allowed := range s.AllowedRecipients {
		allowed := allowed

		if allowed == phone {
			return true
		}
	}

	return false
}

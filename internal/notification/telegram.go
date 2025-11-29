package notification

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type TelegramService struct {
	botToken string
	client   *http.Client
}

func NewTelegramService(botToken string) *TelegramService {
	return &TelegramService{botToken, &http.Client{}}
}

func (s *TelegramService) SendMessage(ctx context.Context, userID int64, message string) error {
	apiURL := "https://api.telegram.org/bot" + s.botToken + "/sendMessage"
	values := url.Values{}
	values.Set("chat_id", strconv.FormatInt(userID, 10))
	values.Set("text", message)
	values.Set("parse_mode", "HTML")

	req, _ := http.NewRequestWithContext(ctx, "POST", apiURL, strings.NewReader(values.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("telegram API returned %d", resp.StatusCode)
	}
	return nil
}

package notification

import (
	"context"

	"github.com/alibekkenny/simpengine/internal/user"
)

// type of nofication
type NotificationChanel string

const (
	ChannelTelegram NotificationChanel = "telegram"
)

type NotificationService struct {
	telegram *TelegramService
}

func NewNotificationService(telegram *TelegramService) *NotificationService {
	return &NotificationService{telegram: telegram}
}

func (s *NotificationService) Send(ctx context.Context, user user.User, channel NotificationChanel, message string) error {
	switch channel {
	case ChannelTelegram:
		if s.telegram == nil || user.TelegramChatID == nil || *user.TelegramChatID == 0 {
			return nil
		}
		return s.telegram.SendMessage(ctx, *user.TelegramChatID, message)
	}

	return nil
}

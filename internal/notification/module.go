package notification

type Module struct {
	Service *NotificationService
}

func NewModule(botToken string) *Module {
	tgService := NewTelegramService(botToken)
	service := NewNotificationService(tgService)

	return &Module{Service: service}
}

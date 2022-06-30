package core

import (
	"github.com/rendau/limelog/internal/adapters/notification"
)

type NotificationProviderSt struct {
	Id       string
	Levels   []string
	Provider notification.Notification
}

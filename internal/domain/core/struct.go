package core

import (
	"github.com/rendau/limelog/internal/interfaces"
)

type NotificationProviderSt struct {
	Id       string
	Levels   []string
	Provider interfaces.Notification
}

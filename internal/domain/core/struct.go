package core

import (
	"github.com/mechta-market/limelog/internal/interfaces"
)

type NotificationProviderSt struct {
	Id       string
	Levels   []string
	Provider interfaces.Notification
}

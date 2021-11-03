package core

import (
	"github.com/mechta-market/limelog/internal/interfaces"
)

type NotificationProviderSt struct {
	Levels   []string
	Provider interfaces.Notification
}

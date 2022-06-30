package notification

type Notification interface {
	Send(msg map[string]any)
}

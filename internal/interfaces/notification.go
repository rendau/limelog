package interfaces

type Notification interface {
	Send(msg map[string]interface{})
}

package core

type Notification struct {
	r *St
}

func NewNotification(r *St) *Notification {
	return &Notification{r: r}
}

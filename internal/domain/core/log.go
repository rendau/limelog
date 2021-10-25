package core

type Log struct {
	r *St
}

func NewLog(r *St) *Log {
	return &Log{r: r}
}

func (c *Log) HandleMsg(msg interface{}) {
}

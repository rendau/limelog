package entities

import (
	"time"
)

type LogListParsSt struct {
	PaginationParams

	FilterObj map[string]interface{} `json:"filter_obj"`
}

type LogRemoveParsSt struct {
	TsLt *time.Time
}

package entities

import (
	"time"
)

type LogListParsSt struct {
	PaginationParams

	FilterObj map[string]interface{} `json:"filter_obj"`
}

type LogRemoveParsSt struct {
	Tag  *string
	TsLt *time.Time
}

package entities

import (
	"time"

	"github.com/rendau/dop/dopTypes"
)

type LogListParsSt struct {
	dopTypes.ListParams

	FilterObj map[string]any `json:"filter_obj"`
}

type LogListRepSt struct {
	dopTypes.PaginatedListRep
	Results []map[string]any `json:"results"`
}

type LogRemoveParsSt struct {
	Tag  *string
	TsLt *time.Time
}

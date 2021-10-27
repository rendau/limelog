package entities

type LogListParsSt struct {
	PaginationParams

	FilterObj map[string]interface{} `json:"filter_obj"`
}

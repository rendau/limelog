package rest

// swagger:response errRep
type docErrRepSt struct {
	// in:body
	Body ErrRepSt
}

type ErrRepSt struct {
	// Код ошибки
	ErrorCode string `json:"error_code"`
}

type PaginatedListRepSt struct {
	Page       int64       `json:"page"`
	PageSize   int64       `json:"page_size"`
	TotalCount int64       `json:"total_count"`
	Results    interface{} `json:"results"`
}

type DocPaginatedListRepSt struct {
	Page       int64 `json:"page"`
	PageSize   int64 `json:"page_size"`
	TotalCount int64 `json:"total_count"`
}

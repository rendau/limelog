package entities

import "time"

type PaginationParams struct {
	Page     int64 `json:"page"`
	PageSize int64 `json:"page_size"`
}

type PeriodFilterPars struct {
	TsGTE *time.Time
	TsLTE *time.Time
}

type ChartVByTime struct {
	Ts time.Time `json:"ts"`
	V  int64     `json:"v"`
}

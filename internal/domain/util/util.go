package util

import (
	"strconv"
	"time"

	"github.com/mechta-market/limelog/internal/cns"
	"github.com/mechta-market/limelog/internal/domain/entities"
	"github.com/mechta-market/limelog/internal/domain/errs"
)

func RequirePageSize(pars entities.PaginationParams, allowedPageSize int64) error {
	if allowedPageSize == 0 {
		allowedPageSize = cns.MaxPageSize
	}

	if pars.Limit == 0 || pars.Limit > allowedPageSize {
		return errs.IncorrectPageSize
	}

	return nil
}

func CoalesceInt64(v *int64, nv int64) int64 {
	if v == nil {
		return nv
	}

	return *v
}

func TimeInAppLocation(v *time.Time) {
	if v != nil {
		*v = (*v).In(cns.AppTimeLocation)
	}
}

func NewInt(v int) *int {
	return &v
}

func NewInt64(v int64) *int64 {
	return &v
}

func NewFloat64(v float64) *float64 {
	return &v
}

func NewString(v string) *string {
	return &v
}

func NewBool(v bool) *bool {
	return &v
}

func NewTime(v time.Time) *time.Time {
	return &v
}

func NewSliceInt64(v ...int64) *[]int64 {
	res := make([]int64, 0, len(v))
	res = append(res, v...)
	return &res
}

func NewSliceString(v ...string) *[]string {
	res := make([]string, 0, len(v))
	res = append(res, v...)
	return &res
}

func Int64SliceToString(src []int64, delimiter, emptyV string) string {
	if len(src) == 0 {
		return emptyV
	}

	res := ""

	for _, v := range src {
		if res != "" {
			res += delimiter
		}
		res += strconv.FormatInt(v, 10)
	}

	return res
}

func Int64SliceHasValue(sl []int64, v int64) bool {
	for _, x := range sl {
		if x == v {
			return true
		}
	}

	return false
}

func Int64SlicesAreSame(a, b []int64) bool {
	for _, x := range a {
		if !Int64SliceHasValue(b, x) {
			return false
		}
	}

	for _, x := range b {
		if !Int64SliceHasValue(a, x) {
			return false
		}
	}

	return true
}

func Int64SlicesIntersection(sl1, sl2 []int64) []int64 {
	result := make([]int64, 0)

	if len(sl1) == 0 || len(sl2) == 0 {
		return result
	}

	for _, x := range sl1 {
		if Int64SliceHasValue(sl2, x) {
			result = append(result, x)
		}
	}

	return result
}

func Int64SliceExcludeValues(sl, vs []int64) []int64 {
	result := make([]int64, 0, len(sl))

	for _, x := range sl {
		if !Int64SliceHasValue(vs, x) {
			result = append(result, x)
		}
	}

	return result
}

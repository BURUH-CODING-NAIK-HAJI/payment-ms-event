package requestentity

import (
	"net/url"
	"strconv"

	"github.com/rizface/golang-api-template/app/errorgroup"
)

type Pagination struct {
	Limit     int `json:"limit"`
	Match     int `json:"match"`
	Current   int `json:"current"`
	TotalRows int `json:"totalRows"`
	TotalPage int `json:"totalPage"`
}

// func (p *Pagination) Validate() err {

// }

func NewPagination(queries url.Values) *Pagination {
	limit, err := strconv.Atoi(queries.Get("limit"))
	if err != nil {
		invalidLimit := errorgroup.INVALID_PAGINATION_PARAMETER
		invalidLimit.Message = "limit must be number"
		panic(invalidLimit)
	}
	if limit == 0 {
		limit = 1
	}

	page, err := strconv.Atoi(queries.Get("page"))
	if err != nil {
		invalidPage := errorgroup.INVALID_PAGINATION_PARAMETER
		invalidPage.Message = "page must be number"
		panic(err)
	}
	if page == 0 {
		page = 1
	}
	return &Pagination{
		Limit:   limit,
		Current: page,
	}
}

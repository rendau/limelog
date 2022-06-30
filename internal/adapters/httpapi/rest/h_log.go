package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dopHttps "github.com/rendau/dop/adapters/server/https"
	"github.com/rendau/dop/dopTypes"
	"github.com/rendau/limelog/internal/domain/entities"
)

// @Router   /log/list [post]
// @Tags     log
// @Param    body  body  entities.LogListParsSt  false  "body"
// @Success  200  {object}  entities.LogListRepSt
// @Failure  400  {object}  dopTypes.ErrRep
func (o *St) hLogList(c *gin.Context) {
	pars := &entities.LogListParsSt{}
	if !dopHttps.BindJSON(c, pars) {
		return
	}

	result, tCount, err := o.ucs.LogList(o.getRequestContext(c), pars)
	if dopHttps.Error(c, err) {
		return
	}

	c.JSON(http.StatusOK, entities.LogListRepSt{
		PaginatedListRep: dopTypes.PaginatedListRep{
			Page:       pars.Page,
			PageSize:   pars.PageSize,
			TotalCount: tCount,
		},
		Results: result,
	})
}

package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dopHttps "github.com/rendau/dop/adapters/server/https"
	"github.com/rendau/limelog/internal/domain/entities"
)

// @Router   /config [get]
// @Tags     config
// @Produce  json
// @Success  200  {object}  entities.ConfigSt
// @Failure  400  {object}  dopTypes.ErrRep
func (o *St) hConfigGet(c *gin.Context) {
	result, err := o.ucs.ConfigGet(o.getRequestContext(c))
	if dopHttps.Error(c, err) {
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Router   /config [put]
// @Tags     config
// @Summary  Set config
// @Accept   json
// @Param    body  body  entities.ConfigSt  false  "body"
// @Success  200
// @Failure  400  {object}  dopTypes.ErrRep
func (o *St) hConfigSet(c *gin.Context) {
	reqObj := &entities.ConfigSt{}
	if !dopHttps.BindJSON(c, reqObj) {
		return
	}

	dopHttps.Error(c, o.ucs.ConfigSet(o.getRequestContext(c), reqObj))
}

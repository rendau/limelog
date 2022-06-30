package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dopHttps "github.com/rendau/dop/adapters/server/https"
)

// @Router   /tag [get]
// @Tags     tag
// @Produce  json
// @Success  200  {array}   string
// @Failure  400  {object}  dopTypes.ErrRep
func (o *St) hTagList(c *gin.Context) {
	result, err := o.ucs.TagList(o.getRequestContext(c))
	if dopHttps.Error(c, err) {
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Router   /tag/:id [delete]
// @Tags     tag
// @Param    id  path  string  true  "id"
// @Produce  json
// @Success  200
// @Failure  400  {object}  dopTypes.ErrRep
func (o *St) hTagRemove(c *gin.Context) {
	id := c.Param("id")

	err := o.ucs.TagRemove(o.getRequestContext(c), id)
	if dopHttps.Error(c, err) {
		return
	}
}

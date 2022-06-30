package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dopHttps "github.com/rendau/dop/adapters/server/https"
	"github.com/rendau/limelog/internal/domain/entities"
)

// @Router   /profile [get]
// @Tags     profile
// @Produce  json
// @Success  200  {object}  entities.ProfileSt
// @Failure  400  {object}  dopTypes.ErrRep
func (o *St) hProfileGet(c *gin.Context) {
	repObj, err := o.ucs.ProfileGet(o.getRequestContext(c))
	if dopHttps.Error(c, err) {
		return
	}

	c.JSON(http.StatusOK, repObj)
}

// @Router   /profile/auth [post]
// @Tags     profile
// @Param    body    body  entities.ProfileAuthReqSt  false  "body"
// @Produce  json
// @Success  200  {object}  entities.ProfileAuthRepSt
// @Failure  400  {object}  dopTypes.ErrRep
func (o *St) hProfileAuth(c *gin.Context) {
	reqObj := &entities.ProfileAuthReqSt{}
	if !dopHttps.BindJSON(c, reqObj) {
		return
	}

	token, err := o.ucs.ProfileAuth(o.getRequestContext(c), reqObj.Password)
	if dopHttps.Error(c, err) {
		return
	}

	c.JSON(http.StatusOK, entities.ProfileAuthRepSt{Token: token})
}

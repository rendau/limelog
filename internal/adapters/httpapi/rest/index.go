package rest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rendau/dop/adapters/logger"
	dopHttps "github.com/rendau/dop/adapters/server/https"
	"github.com/rendau/limelog/internal/domain/usecases"
	swagFiles "github.com/swaggo/files"
	ginSwag "github.com/swaggo/gin-swagger"
)

type St struct {
	lg  logger.Lite
	ucs *usecases.St
}

func GetHandler(lg logger.Lite, ucs *usecases.St, withCors bool) http.Handler {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	// middlewares

	r.Use(dopHttps.MwRecovery(lg, nil))
	if withCors {
		r.Use(dopHttps.MwCors())
	}

	// handlers

	// doc
	r.GET("/doc/*any", ginSwag.WrapHandler(swagFiles.Handler, func(c *ginSwag.Config) {
		c.DefaultModelsExpandDepth = 0
		c.DocExpansion = "none"
	}))

	s := &St{lg: lg, ucs: ucs}

	// healthcheck
	r.GET("/healthcheck", func(c *gin.Context) { c.Status(http.StatusOK) })

	// config
	r.GET("/config", s.hConfigGet)
	r.PUT("/config", s.hConfigSet)

	// profile
	r.GET("/profile", s.hProfileGet)
	r.POST("/profile/auth", s.hProfileAuth)

	// tag
	r.GET("/tag", s.hTagList)
	r.DELETE("/tag/:id", s.hTagRemove)

	// log
	r.POST("/log/list", s.hLogList)

	return r
}

func (o *St) getRequestContext(c *gin.Context) context.Context {
	return o.ucs.SessionSetToContextByToken(nil, dopHttps.GetAuthToken(c))
}

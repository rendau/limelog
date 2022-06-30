package htp

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rendau/dop/adapters/logger"
	dopHttps "github.com/rendau/dop/adapters/server/https"
	"github.com/rendau/limelog/internal/domain/usecases"
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

	s := &St{lg: lg, ucs: ucs}

	// log
	r.POST("/log", s.hLog)

	return r
}

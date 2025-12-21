package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moth13/finance_tracker/internal/util"
	"github.com/moth13/finance_tracker/internal/web/static/views"
)

func (server *Server) aboutPageHandler(ctx *gin.Context) {
	err := render(ctx, http.StatusOK, views.About())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}
}

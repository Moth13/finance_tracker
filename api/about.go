package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moth13/finance_tracker/views"
)

func (server *Server) aboutPageHandler(ctx *gin.Context) {

	err := server.render(ctx, http.StatusOK, views.Layout(views.About(), "about", "/views/about"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
}

package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moth13/finance_tracker/internal/util"
	"github.com/moth13/finance_tracker/internal/web/static/views"
)

func (server *Server) homePage(ctx *gin.Context) {
	var err error

	if nil == getCurrentUser(ctx) {
		err = render(ctx, http.StatusOK, views.Base(views.Landing()))
	} else {
		err = render(ctx, http.StatusOK, views.Base(views.Main("Jeremie GUERINEL")))
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}
}

func isUser(email string, password string) bool {
	return email == "a@a.com"
}

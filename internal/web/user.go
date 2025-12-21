package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moth13/finance_tracker/internal/util"
	"github.com/moth13/finance_tracker/internal/web/static/views"
)

func (server *Server) loginHandler(ctx *gin.Context) {
	email := ctx.PostForm("email")
	password := ctx.PostForm("password")

	var err error
	if isUser(email, password) {
		if err := setUserSession(ctx, email); err != nil {
			ctx.String(http.StatusInternalServerError, "cannot save session")
			return
		}
		err = render(ctx, http.StatusOK, views.Base(views.Main("Jeremie GUERINEL")))
	} else {
		err = render(ctx, http.StatusOK, views.Base(views.SignUpUser()))
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}
}

func (server *Server) logoutHandler(ctx *gin.Context) {
	_ = clearUserSession(ctx)

	err := render(ctx, http.StatusOK, views.Base(views.Landing()))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}

}

func (server *Server) signupUserPage(ctx *gin.Context) {
	err := render(ctx, http.StatusOK, views.Base(views.SignUpAccounts()))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}
}

func (server *Server) signupAccountPage(ctx *gin.Context) {

	var err error
	if isUser("a@a.com", "aa") {
		if err := setUserSession(ctx, "a@a.com"); err != nil {
			ctx.String(http.StatusInternalServerError, "cannot save session")
			return
		}
		err = render(ctx, http.StatusOK, views.Base(views.Main("Jeremie GUERINEL")))
	} else {
		err = render(ctx, http.StatusOK, views.Base(views.Landing()))
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}
}

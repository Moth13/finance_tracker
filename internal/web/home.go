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

func (server *Server) isUser(ctx *gin.Context, email string, password string) (*User, error) {
	user, err := server.Store.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if err = util.CheckPassword(password, user.HashedPassword); err != nil {
		return nil, err
	}

	return &User{Email: user.Email, Username: user.Username}, nil
}

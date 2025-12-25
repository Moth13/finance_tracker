package web

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/moth13/finance_tracker/internal/util"
	"github.com/moth13/finance_tracker/internal/web/static/views"
)

func getCurrentUser(ctx *gin.Context) *User {
	sess := sessions.Default(ctx)
	email, ok := sess.Get("email").(string)
	if !ok || email == "" {
		return nil
	}
	username, ok := sess.Get("username").(string)
	if !ok || username == "" {
		return nil
	}
	return &User{Email: email, Username: username}
}

func setUserSession(ctx *gin.Context, email string, username string) error {
	sess := sessions.Default(ctx)
	sess.Set("email", email)
	sess.Set("username", username)
	return sess.Save()
}

func clearUserSession(ctx *gin.Context) error {
	sess := sessions.Default(ctx)
	sess.Clear()
	return sess.Save()
}

func AuthRequired() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if getCurrentUser(ctx) == nil {
			isHTMX := ctx.GetHeader("HX-Request") == "true"
			if isHTMX {
				ctx.Status(http.StatusUnauthorized)
				err := render(ctx, http.StatusOK, views.Base(views.SignUpUser()))
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
					return
				}
			} else {
				ctx.Redirect(http.StatusSeeOther, "/login")
			}
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

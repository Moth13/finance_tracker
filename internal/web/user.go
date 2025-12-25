package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/moth13/finance_tracker/internal/db/sqlc"
	"github.com/moth13/finance_tracker/internal/util"
	"github.com/moth13/finance_tracker/internal/web/static/views"
	"github.com/shopspring/decimal"
)

func (server *Server) loginHandler(ctx *gin.Context) {
	email := ctx.PostForm("email")
	password := ctx.PostForm("password")

	if user, err := server.isUser(ctx, email, password); err == nil {
		if err = setUserSession(ctx, user.Email, user.Username); err != nil {
			ctx.String(http.StatusInternalServerError, "cannot save session")
			return
		}
		err = render(ctx, http.StatusOK, views.Base(views.Main("Jeremie GUERINEL")))

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
			return
		}
	} else {
		err = render(ctx, http.StatusOK, views.Base(views.SignUpUser()))

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
			return
		}
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

type createUserFormRequest struct {
	Username string `form:"username" binding:"required,alphanum"`
	Password string `form:"password" binding:"required,min=6"`
	FullName string `form:"fullname" binding:"required"`
	Email    string `form:"email" binding:"required,email"`
	Currency string `form:"currency" binding:"required,len=3"`
}

func (server *Server) signupUserPage(ctx *gin.Context) {
	var req createUserFormRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	hashedPassword, err := util.HashedPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
		Currency:       req.Currency,
	}

	user, err := server.Store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, util.ErrorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}

	if err := setUserSession(ctx, user.Email, user.Username); err != nil {
		ctx.String(http.StatusInternalServerError, "cannot save session")
		return
	}

	err = render(ctx, http.StatusOK, views.Base(views.SignUpAccounts()))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}
}

type createAccountFormRequest struct {
	Title       string          `form:"title" binding:"required"`
	Description string          `form:"description" binding:"required"`
	InitBalance decimal.Decimal `form:"init_balance" binding:"required"`
}

func (server *Server) signupAccountPage(ctx *gin.Context) {
	var req createAccountFormRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	user := getCurrentUser(ctx)

	arg := db.CreateAccountParams{
		Owner:       user.Username,
		Title:       req.Title,
		Description: req.Description,
		InitBalance: req.InitBalance,
	}

	_, err := server.Store.CreateAccount(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, util.ErrorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}

	err = render(ctx, http.StatusOK, views.Base(views.Main("Jeremie GUERINEL")))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}
}

package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/moth13/finance_tracker/db/sqlc"
	"github.com/moth13/finance_tracker/views"
	"github.com/shopspring/decimal"
)

func (server *Server) getViewLinePage(ctx *gin.Context) {
	err := server.render(ctx, http.StatusOK, views.Layout(views.CreateLine(), "", "/"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
}

type createLineFormRequest struct {
	Title        string          `form:"title" binding:"required"`
	AccountName  string          `form:"account_name" binding:"required"`
	MonthName    string          `form:"month_name" binding:"required"`
	CategoryName string          `form:"category_name" binding:"required"`
	Amount       decimal.Decimal `form:"amount" binding:"required"`
	Checked      bool            `form:"checked"`
	Description  string          `form:"description"`
	DueDate      string          `form:"due_date" binding:"required"`
}

func (server *Server) postViewLine(ctx *gin.Context) {
	var req createLineFormRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	due_date, err := time.Parse("2006-01-02", req.DueDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	arg := db.AddLineTxParams{
		// Owner:    authPayload.Username,
		Owner:       "jose",
		Title:       req.Title,
		Description: req.Description,
		Amount:      req.Amount,
		AccountID:   1,
		MonthID:     1,
		YearID:      1,
		CategoryID:  1,
		DueDate:     due_date,
	}

	_, err = server.store.AddLineTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	server.homePage(ctx)
}

type deleteViewLineRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteViewLine(ctx *gin.Context) {
	var req deleteViewLineRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.DeleteLineTxParams{
		ID: req.ID,
	}

	_, err := server.store.DeleteLineTx(ctx, arg)
	fmt.Println(err)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	server.homePage(ctx)
}

package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/moth13/finance_tracker/db/sqlc"
	"github.com/moth13/finance_tracker/views"
	"github.com/moth13/finance_tracker/views/components"
	"github.com/shopspring/decimal"
)

type getViewLineRequest struct {
	ID int64 `uri:"id"`
}

func (server *Server) getViewLinePage(ctx *gin.Context) {
	var req getViewLineRequest
	var vLine components.Line

	if err := ctx.ShouldBindUri(&req); err == nil && req.ID > 0 {
		line, err := server.store.GetExpliciteLine(ctx, req.ID)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		vLine = components.Line{
			Id:          line.Title,
			DbID:        line.ID,
			Description: line.Description,
			Title:       line.Title,
			Amount:      line.Amount,
			DueDate:     line.DueDate,
			Checked:     line.Checked,
			Account:     line.Account,
			Month:       line.Month,
			Category:    line.Category,
		}
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err := server.render(ctx, http.StatusOK, views.Layout(views.CreateLine(vLine), "", "/"))
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

type updateViewLineIDRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type updateViewLineRequest struct {
	Title        *string             `form:"title"`
	AccountName  *string             `form:"account_name"`
	MonthName    *string             `form:"month_name"`
	CategoryName *string             `form:"category_name"`
	AccountID    *int64              `form:"account_id"`
	Amount       decimal.NullDecimal `form:"amount"`
	Checked      *bool               `form:"checked"`
	Description  *string             `form:"description"`
	DueDate      *string             `form:"due_date"`
}

func (server *Server) updateViewLine(ctx *gin.Context) {
	var reqURI updateViewLineIDRequest
	if err := ctx.ShouldBindUri(&reqURI); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var req updateViewLineRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateLineTxParams{
		ID:          reqURI.ID,
		Title:       req.Title,
		Amount:      req.Amount,
		Checked:     req.Checked,
		Description: req.Description,
	}

	if req.DueDate != nil {
		due_date, err := time.Parse("2006-01-02", *req.DueDate)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		arg.DueDate = &due_date
	}

	_, err := server.store.UpdateLineTx(ctx, arg)
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

package web

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/moth13/finance_tracker/internal/db/sqlc"
	"github.com/moth13/finance_tracker/internal/util"
	"github.com/moth13/finance_tracker/internal/web/static/views"
	"github.com/moth13/finance_tracker/internal/web/static/views/components"
	"github.com/shopspring/decimal"
)

type getViewLineRequest struct {
	ID int64 `uri:"id"`
}

func (server *Server) viewLinesPage(ctx *gin.Context) {
	user := getCurrentUser(ctx)
	if user == nil {
		ctx.Redirect(http.StatusSeeOther, "/login")
		return
	}

	arg := db.ListExplicitLinesParams{
		Limit:  10,
		Offset: 0,
		Owner:  user.Username,
	}

	var viewInfos views.Infos

	lines, err := server.Store.ListExplicitLines(ctx, arg)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}

	account, err := server.Store.GetAccount(ctx, 1)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	} else {
		viewInfos.Balance = account.Balance
		viewInfos.FinalBalance = account.FinalBalance
	}

	for _, line := range lines {
		viewsTodo := &components.Line{
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
		viewInfos.Lines = append(viewInfos.Lines, viewsTodo)
	}

	err = render(ctx, http.StatusOK, views.Line(viewInfos))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}
}

func (server *Server) newViewLinePage(ctx *gin.Context) {
	vLine := components.Line{
		DueDate: time.Now().Local().UTC(),
	}

	err := render(ctx, http.StatusOK, views.NewLine(vLine))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}
}

func (server *Server) getViewLinePage(ctx *gin.Context) {
	var req getViewLineRequest
	var vLine components.Line

	if err := ctx.ShouldBindUri(&req); err == nil && req.ID > 0 {
		line, err := server.Store.GetExpliciteLine(ctx, req.ID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				ctx.JSON(http.StatusNotFound, util.ErrorResponse(err))
				return
			}
			ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
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
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}

	err := render(ctx, http.StatusOK, views.NewLine(vLine))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
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
	user := getCurrentUser(ctx)
	if user == nil {
		ctx.Redirect(http.StatusSeeOther, "/login")
		return
	}

	var req createLineFormRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	due_date, err := time.Parse("2006-01-02", req.DueDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}
	arg := db.AddLineTxParams{
		Owner:       user.Username,
		Title:       req.Title,
		Description: req.Description,
		Amount:      req.Amount,
		AccountID:   1,
		MonthID:     1,
		YearID:      1,
		CategoryID:  1,
		DueDate:     due_date,
	}

	_, err = server.Store.AddLineTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}

	if ctx.GetHeader("HX-Request") == "true" {
		server.viewLinesPage(ctx)
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
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	arg := db.DeleteLineTxParams{
		ID: req.ID,
	}

	_, err := server.Store.DeleteLineTx(ctx, arg)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, util.ErrorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}

	if ctx.GetHeader("HX-Request") == "true" {
		ctx.Status(http.StatusOK)
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
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	var req updateViewLineRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
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
			ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
			return
		}
		arg.DueDate = &due_date
	}

	_, err := server.Store.UpdateLineTx(ctx, arg)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, util.ErrorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}

	if ctx.GetHeader("HX-Request") == "true" {
		server.viewLinesPage(ctx)
		return
	}

	server.homePage(ctx)
}

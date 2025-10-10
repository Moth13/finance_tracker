package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/moth13/finance_tracker/db/sqlc"
	"github.com/moth13/finance_tracker/token"
	decimal "github.com/shopspring/decimal"
)

type createLineRequest struct {
	Title       string          `json:"title" binding:"required"`
	AccountID   int64           `json:"account_id" binding:"required"`
	MonthID     int64           `json:"month_id" binding:"required"`
	YearID      int64           `json:"year_id" binding:"required"`
	CategoryID  int64           `json:"category_id" binding:"required"`
	Amount      decimal.Decimal `json:"amount" binding:"required"`
	Checked     *bool           `json:"checked" binding:"required"`
	Description string          `json:"description" binding:"required"`
	DueDate     time.Time       `json:"due_date" binding:"required"`
}

func (server *Server) createLine(ctx *gin.Context) {
	var req createLineRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.AddLineTxParams{
		Owner:       authPayload.Username,
		Title:       req.Title,
		Description: req.Description,
		Checked:     *req.Checked,
		Amount:      req.Amount,
		AccountID:   req.AccountID,
		MonthID:     req.MonthID,
		YearID:      req.YearID,
		CategoryID:  req.CategoryID,
		DueDate:     req.DueDate,
	}

	result, err := server.store.AddLineTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

type getLineRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getLine(ctx *gin.Context) {
	var req getLineRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	line, err := server.store.GetLine(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if line.Owner != authPayload.Username {
		err := errors.New("line doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
	}

	ctx.JSON(http.StatusOK, line)
}

type listLinesRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listLines(ctx *gin.Context) {
	var req listLinesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.ListLinesParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
		Owner:  authPayload.Username,
	}

	lines, err := server.store.ListLines(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, lines)
}

type deleteLineRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteLine(ctx *gin.Context) {
	var req deleteLineRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.DeleteLineTxParams{
		ID: req.ID,
	}

	result, err := server.store.DeleteLineTx(ctx, arg)
	fmt.Println(err)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

type updateLineIDRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type updateLineJSONRequest struct {
	Title       *string             `json:"title"`
	AccountID   *int64              `json:"account_id"`
	MonthID     *int64              `json:"month_id"`
	YearID      *int64              `json:"year_id"`
	CategoryID  *int64              `json:"category_id"`
	Amount      decimal.NullDecimal `json:"amount"`
	Checked     *bool               `json:"checked"`
	Description *string             `json:"description"`
	DueDate     *time.Time          `json:"due_date"`
}

func (server *Server) updateLine(ctx *gin.Context) {
	var reqURI updateLineIDRequest
	if err := ctx.ShouldBindUri(&reqURI); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var reqJSON updateLineJSONRequest
	if err := ctx.ShouldBindJSON(&reqJSON); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateLineTxParams{
		ID: reqURI.ID,
	}

	result, err := server.store.UpdateLineTx(ctx, arg)
	fmt.Println(err)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

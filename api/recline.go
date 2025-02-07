package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/moth13/finance_tracker/db/sqlc"
	decimal "github.com/shopspring/decimal"
)

type createRecLineRequest struct {
	Title       string          `json:"title" binding:"required"`
	AccountID   int64           `json:"account_id" binding:"required"`
	Amount      decimal.Decimal `json:"amount" binding:"required"`
	CategoryID  int64           `json:"category_id" binding:"required"`
	Description string          `json:"description" binding:"required"`
	Recurrency  string          `json:"recurrency" binding:"required"`
	DueDate     time.Time       `json:"due_date" binding:"required"`
}

func (server *Server) createRecLine(ctx *gin.Context) {
	var req createRecLineRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		fmt.Println(err)
		fmt.Print(req)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.CreateRecLineParams{
		// Owner:    authPayload.Username,
		Owner:       "jose",
		Title:       req.Title,
		Description: req.Description,
		Amount:      req.Amount,
		AccountID:   req.AccountID,
		DueDate:     req.DueDate,
		CategoryID:  req.CategoryID,
		Recurrency:  req.Recurrency,
	}

	recline, err := server.store.CreateRecLine(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, recline)
}

type getRecLineRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getRecLine(ctx *gin.Context) {
	var req getRecLineRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	recline, err := server.store.GetRecLine(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	// if recline.Owner != authPayload.Username {
	// 	err := errors.New("recline doesn't belong to the authenticated user")
	// 	ctx.JSON(http.StatusUnauthorized, errorResponse(err))
	// }

	ctx.JSON(http.StatusOK, recline)
}

type listRecLinesRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listRecLines(ctx *gin.Context) {
	fmt.Println("aaa")
	var req listRecLinesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListRecLinesParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
		Owner:  "jose",
	}

	reclines, err := server.store.ListRecLines(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, reclines)
}

type deleteRecLineRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteRecLine(ctx *gin.Context) {
	var req deleteRecLineRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteRecLine(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Line %d has been deleted", req.ID)})
}

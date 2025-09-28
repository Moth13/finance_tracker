package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/moth13/finance_tracker/db/sqlc"
	"github.com/moth13/finance_tracker/token"
)

type createMonthRequest struct {
	Owner       string    `json:"owner" binding:"required"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description" binding:"required"`
	StartDate   time.Time `json:"start_date" binding:"required"`
	EndDate     time.Time `json:"end_date" binding:"required"`
	YearID      int64     `json:"year_id"  binding:"required,min=1"`
}

func (server *Server) createMonth(ctx *gin.Context) {
	var req createMonthRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.CreateMonthParams{
		// Owner:    authPayload.Username,
		Owner:       req.Owner,
		Title:       req.Title,
		Description: req.Description,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		YearID:      req.YearID,
	}

	if _, valid := server.validYear(ctx, arg.YearID); !valid {
		return
	}

	month, err := server.store.CreateMonth(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, month)
}

type getMonthRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getMonth(ctx *gin.Context) {
	var req getMonthRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	month, err := server.store.GetMonth(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	// if month.Owner != authPayload.Username {
	// 	err := errors.New("month doesn't belong to the authenticated user")
	// 	ctx.JSON(http.StatusUnauthorized, errorResponse(err))
	// }

	ctx.JSON(http.StatusOK, month)
}

type listMonthRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listMonths(ctx *gin.Context) {
	var req listMonthRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.ListMonthsParams{
		Owner:  authPayload.Username,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	months, err := server.store.ListMonths(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, months)
}

type deleteMonthRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteMonth(ctx *gin.Context) {
	var req deleteMonthRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteMonth(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Month %d has been deleted", req.ID)})
}

type updateMonthIDRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type updateMonthJSONRequest struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	YearID      *int64     `json:"year_id"`
}

func (server *Server) updateMonth(ctx *gin.Context) {
	var reqURI updateMonthIDRequest
	if err := ctx.ShouldBindUri(&reqURI); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var reqJSON updateMonthJSONRequest
	if err := ctx.ShouldBindJSON(&reqJSON); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Get month
	month, err := server.store.GetMonth(ctx, reqURI.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.UpdateMonthParams{
		ID:          month.ID,
		Title:       month.Title,
		Description: month.Description,
		StartDate:   month.StartDate,
		EndDate:     month.EndDate,
		YearID:      month.YearID,
	}

	// Overload when needs it
	if reqJSON.Title != nil {
		arg.Title = *reqJSON.Title
	}

	if reqJSON.Description != nil {
		arg.Description = *reqJSON.Description
	}

	if reqJSON.StartDate != nil {
		arg.StartDate = *reqJSON.StartDate
	}

	if reqJSON.EndDate != nil {
		arg.EndDate = *reqJSON.EndDate
	}

	if reqJSON.YearID != nil {
		arg.YearID = *reqJSON.YearID
	}

	month, err = server.store.UpdateMonth(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, month)
}

// func (server *Server) validMonth(ctx *gin.Context, monthID int64) (db.Month, bool) {
// 	month, err := server.store.GetMonth(ctx, monthID)
// 	if err != nil {
// 		if err == pgx.ErrNoRows {
// 			ctx.JSON(http.StatusNotFound, errorResponse(err))
// 			return month, false
// 		}
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return month, false
// 	}

// 	return month, true
// }

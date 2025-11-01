package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	db "github.com/moth13/finance_tracker/db/sqlc"
	"github.com/moth13/finance_tracker/token"
)

type createYearRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description" binding:"required"`
	StartDate   time.Time `json:"start_date" binding:"required"`
	EndDate     time.Time `json:"end_date" binding:"required"`
}

func (server *Server) createYear(ctx *gin.Context) {
	var req createYearRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.CreateYearParams{
		Owner:       authPayload.Username,
		Title:       req.Title,
		Description: req.Description,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	}

	year, err := server.store.CreateYear(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, year)
}

type getYearRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getYear(ctx *gin.Context) {
	var req getYearRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	year, err := server.store.GetYear(ctx, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if year.Owner != authPayload.Username {
		err := errors.New("year doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
	}

	ctx.JSON(http.StatusOK, year)
}

type listYearRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listYears(ctx *gin.Context) {
	var req listYearRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.ListYearsParams{
		Owner:  authPayload.Username,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	years, err := server.store.ListYears(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, years)
}

type deleteYearRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteYear(ctx *gin.Context) {
	var req deleteYearRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteYear(ctx, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Year %d has been deleted", req.ID)})
}

type updateYearIDRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type updateYearJSONRequest struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
}

func (server *Server) updateYear(ctx *gin.Context) {
	var reqURI updateYearIDRequest
	if err := ctx.ShouldBindUri(&reqURI); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var reqJSON updateYearJSONRequest
	if err := ctx.ShouldBindJSON(&reqJSON); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Get year
	year, err := server.store.GetYear(ctx, reqURI.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.UpdateYearParams{
		ID:          year.ID,
		Title:       year.Title,
		Description: year.Description,
		StartDate:   year.StartDate,
		EndDate:     year.EndDate,
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

	year, err = server.store.UpdateYear(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, year)
}

func (server *Server) validYear(ctx *gin.Context, yearID int64) (db.Year, bool) {
	year, err := server.store.GetYear(ctx, yearID)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return year, false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return year, false
	}

	return year, true
}

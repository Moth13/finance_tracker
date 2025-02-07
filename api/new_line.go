package api

import (
	"net/http"
	"time"

	"log"

	"github.com/gin-gonic/gin"
	db "github.com/moth13/finance_tracker/db/sqlc"
	"github.com/moth13/finance_tracker/views"
	decimal "github.com/shopspring/decimal"
)

func (server *Server) newLinePageHandler(ctx *gin.Context) {
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

func (server *Server) postNewLineHandler(ctx *gin.Context) {
	var req createLineFormRequest
	if err := ctx.ShouldBind(&req); err != nil {
		log.Println(req)
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)

	// authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.AddLineTxParams{
		// Owner:    authPayload.Username,
		Owner:       "jose",
		Title:       req.Title,
		Description: req.Description,
		Checked:     req.Checked,
		Amount:      req.Amount,
		AccountID:   1,
		MonthID:     1,
		YearID:      1,
		CategoryID:  1,
		DueDate:     time.Now().UTC(),
	}

	_, err := server.store.AddLineTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	server.homePageHandler(ctx)
}

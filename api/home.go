package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/moth13/finance_tracker/db/sqlc"
	"github.com/moth13/finance_tracker/views"
	"github.com/moth13/finance_tracker/views/components"
)

func (server *Server) homePageHandler(ctx *gin.Context) {

	arg := db.ListExplicitLinesParams{
		Limit:  10,
		Offset: 0,
		Owner:  "jose",
	}

	lines, err := server.store.ListExplicitLines(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	account, err := server.store.GetAccount(ctx, 1)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var viewInfos views.Infos
	viewInfos.Balance = account.Balance
	viewInfos.FinalBalance = account.FinalBalance

	for _, line := range lines {
		viewsTodo := &components.Line{
			Id:          line.Title,
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

	err = server.render(ctx, http.StatusOK, views.Layout(views.Line(viewInfos), "home", "/"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
}

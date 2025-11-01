package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/moth13/finance_tracker/db/sqlc"
	"github.com/moth13/finance_tracker/token"
	decimal "github.com/shopspring/decimal"
)

type createAccountRequest struct {
	Title       string          `json:"title" binding:"required"`
	Description string          `json:"description" binding:"required"`
	InitBalance decimal.Decimal `json:"init_balance" binding:"required"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.CreateAccountParams{
		Owner:       authPayload.Username,
		Title:       req.Title,
		Description: req.Description,
		InitBalance: req.InitBalance,
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if account.Owner != authPayload.Username {
		err := errors.New("account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
	}

	ctx.JSON(http.StatusOK, account)
}

type listAccountRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listAccounts(ctx *gin.Context) {
	var req listAccountRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.ListAccountsParams{
		Owner:  authPayload.Username,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	accounts, err := server.store.ListAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}

type deleteAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteAccount(ctx *gin.Context) {
	var req deleteAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteAccount(ctx, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Account %d has been deleted", req.ID)})
}

type updateAccountIDRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type updateAccountJSONRequest struct {
	Title       *string             `json:"title"`
	Description *string             `json:"description"`
	InitBalance decimal.NullDecimal `json:"init_balance"`
}

func (server *Server) updateAccount(ctx *gin.Context) {
	var reqURI updateAccountIDRequest
	if err := ctx.ShouldBindUri(&reqURI); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var reqJSON updateAccountJSONRequest
	if err := ctx.ShouldBindJSON(&reqJSON); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Get account
	account, err := server.store.GetAccount(ctx, reqURI.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.UpdateAccountParams{
		ID:          account.ID,
		Title:       account.Title,
		Description: account.Description,
		InitBalance: account.InitBalance,
	}

	// Overload when needs it
	if reqJSON.Title != nil {
		arg.Title = *reqJSON.Title
	}

	if reqJSON.Description != nil {
		arg.Description = *reqJSON.Description
	}

	if reqJSON.InitBalance.Valid {
		arg.InitBalance = reqJSON.InitBalance.Decimal
	}

	account, err = server.store.UpdateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

// func (server *Server) validAccount(ctx *gin.Context, accountID int64) (db.Account, bool) {
// 	account, err := server.store.GetAccount(ctx, accountID)
// 	if err != nil {
// 		if err == pgx.ErrNoRows {
// 			ctx.JSON(http.StatusNotFound, errorResponse(err))
// 			return account, false
// 		}
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return account, false
// 	}

// 	return account, true
// }

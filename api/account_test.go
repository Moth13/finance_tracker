package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5"
	mockdb "github.com/moth13/finance_tracker/db/mock"
	db "github.com/moth13/finance_tracker/db/sqlc"
	"github.com/moth13/finance_tracker/util"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestCreateAccountAPI(t *testing.T) {
	user, _ := randomUser(t)
	user.Username = "jose"
	account := randomAccount(user.Username)

	// Test cases definition
	testCases := []struct {
		name          string
		body          createAccountRequest
		buildStubds   func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: createAccountRequest{
				Owner:       user.Username,
				Title:       account.Title,
				Description: account.Description,
				InitBalance: account.InitBalance,
			},
			buildStubds: func(store *mockdb.MockStore) {
				arg := db.CreateAccountParams{
					Owner:       user.Username,
					Title:       account.Title,
					Description: account.Description,
					InitBalance: account.InitBalance,
				}

				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(arg.Owner)).
					Times(1).
					Return(user, nil)

				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name: "InvalidUser",
			body: createAccountRequest{
				Owner:       "josa",
				Title:       account.Title,
				Description: account.Description,
				InitBalance: account.InitBalance,
			},
			buildStubds: func(store *mockdb.MockStore) {
				arg := db.CreateAccountParams{
					Owner:       "josa",
				}

				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(arg.Owner)).Times(1).
					Return(db.User{}, pgx.ErrNoRows)

				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InvalidUserError",
			body: createAccountRequest{
				Owner:       "josa",
				Title:       account.Title,
				Description: account.Description,
				InitBalance: account.InitBalance,
			},
			buildStubds: func(store *mockdb.MockStore) {
				arg := db.CreateAccountParams{
					Owner:       "josa",
				}

				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(arg.Owner)).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)

				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidTitle",
			body: createAccountRequest{
				Owner:       user.Username,
				Title:       "",
				Description: account.Description,
				InitBalance: account.InitBalance,
			},
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidDescription",
			body: createAccountRequest{
				Owner:       user.Username,
				Title:       account.Owner,
				Description: "",
				InitBalance: account.InitBalance,
			},
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	// Checking cases
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubds(store)

			// start test server and send request
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/accounts"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestDeleteAccountAPI(t *testing.T) {
	user, _ := randomUser(t)
	user.Username = "jose"
	account := randomAccount(user.Username)

	// Test cases definition
	testCases := []struct {
		name          string
		accountID     int64
		buildStubds   func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:      "InvalidID",
			accountID: -1,
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "NotFound",
			accountID: account.ID,
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalServerError",
			accountID: account.ID,
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	// Checking cases
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubds(store)

			// start test server and send request
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestGetAccountAPI(t *testing.T) {
	user, _ := randomUser(t)
	user.Username = "jose"
	account := randomAccount(user.Username)

	// Test cases definition
	testCases := []struct {
		name          string
		accountID     int64
		buildStubds   func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:      "NotFound",
			accountID: account.ID,
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalServerError",
			accountID: account.ID,
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "InvalidAccountID",
			accountID: 0,
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	// Checking cases
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubds(store)

			// start test server and send request
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestListAccountsAPI(t *testing.T) {
	user, _ := randomUser(t)
	user.Username = "jose"
	n := 5
	accounts := make([]db.Account, n)
	for i := 0; i < n; i++ {
		accounts[i] = randomAccount(user.Username)
	}

	// Test cases definition
	testCases := []struct {
		name          string
		query         listAccountRequest
		buildStubds   func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: listAccountRequest{
				PageID:   1,
				PageSize: int32(n),
			},
			buildStubds: func(store *mockdb.MockStore) {
				arg := db.ListAccountsParams{
					Owner:  user.Username,
					Offset: 0,
					Limit:  int32(n),
				}
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(accounts, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccounts(t, recorder.Body, accounts)
			},
		},
		{
			name: "InternalError",
			query: listAccountRequest{
				PageID:   1,
				PageSize: int32(n),
			},
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidPageSize",
			query: listAccountRequest{
				PageID:   1,
				PageSize: 2,
			},
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidPageIndex",
			query: listAccountRequest{
				PageID:   -1,
				PageSize: int32(n),
			},
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	// Checking cases
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubds(store)

			// start test server and send request
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			request, err := http.NewRequest(http.MethodGet, "/accounts", nil)
			require.NoError(t, err)

			q := request.URL.Query()
			q.Add("page_id", fmt.Sprintf("%d", tc.query.PageID))
			q.Add("page_size", fmt.Sprintf("%d", tc.query.PageSize))
			request.URL.RawQuery = q.Encode()

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestUpdateAccountAPI(t *testing.T) {
	user, _ := randomUser(t)
	user.Username = "jose"
	account1 := randomAccount(user.Username)
	account2 := randomAccount(user.Username)

	account2.ID = account1.ID
	account2.Owner = account1.Owner

	account3 := db.Account{
		ID:           account2.ID,
		Owner:        account2.Owner,
		Title:        account1.Title,
		Description:  account2.Description,
		InitBalance:  account2.InitBalance,
		Balance:      account2.Balance,
		FinalBalance: account2.FinalBalance,
	}

	account4 := db.Account{
		ID:           account2.ID,
		Owner:        account2.Owner,
		Title:        account2.Title,
		Description:  account1.Description,
		InitBalance:  account2.InitBalance,
		Balance:      account2.Balance,
		FinalBalance: account2.FinalBalance,
	}

	account5 := db.Account{
		ID:           account2.ID,
		Owner:        account2.Owner,
		Title:        account2.Title,
		Description:  account2.Description,
		InitBalance:  account1.InitBalance,
		Balance:      account2.Balance,
		FinalBalance: account2.FinalBalance,
	}

	// Test cases definition
	testCases := []struct {
		name          string
		accountID     int64
		body          updateAccountJSONRequest
		buildStubds   func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account1.ID,
			body: updateAccountJSONRequest{
				Title:       &account2.Title,
				Description: &account2.Description,
				InitBalance: decimal.NewNullDecimal(account2.InitBalance),
			},
			buildStubds: func(store *mockdb.MockStore) {
				arg := db.UpdateAccountParams{
					ID:          account2.ID,
					Title:       account2.Title,
					Description: account2.Description,
					InitBalance: account2.InitBalance,
				}
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account2.ID)).
					Times(1).
					Return(account1, nil)
				store.EXPECT().
					UpdateAccount(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(account2, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account2)
			},
		},
		{
			name:      "NoTitle",
			accountID: account1.ID,
			body: updateAccountJSONRequest{
				Description: &account2.Description,
				InitBalance: decimal.NewNullDecimal(account2.InitBalance),
			},
			buildStubds: func(store *mockdb.MockStore) {

				arg := db.UpdateAccountParams{
					ID:          account2.ID,
					Title:       account1.Title,
					Description: account2.Description,
					InitBalance: account2.InitBalance,
				}
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account2.ID)).
					Times(1).
					Return(account1, nil)
				store.EXPECT().
					UpdateAccount(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(account3, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account3)
			},
		},
		{
			name:      "NoDescription",
			accountID: account1.ID,
			body: updateAccountJSONRequest{
				Title:       &account2.Title,
				InitBalance: decimal.NewNullDecimal(account2.InitBalance),
			},
			buildStubds: func(store *mockdb.MockStore) {

				arg := db.UpdateAccountParams{
					ID:          account2.ID,
					Title:       account2.Title,
					Description: account1.Description,
					InitBalance: account2.InitBalance,
				}
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account2.ID)).
					Times(1).
					Return(account1, nil)
				store.EXPECT().
					UpdateAccount(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(account4, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account4)
			},
		},
		{
			name:      "NoInitBalance",
			accountID: account1.ID,
			body: updateAccountJSONRequest{
				Title:       &account2.Title,
				Description: &account2.Description,
			},
			buildStubds: func(store *mockdb.MockStore) {

				arg := db.UpdateAccountParams{
					ID:          account2.ID,
					Title:       account2.Title,
					Description: account2.Description,
					InitBalance: account1.InitBalance,
				}
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account2.ID)).
					Times(1).
					Return(account1, nil)
				store.EXPECT().
					UpdateAccount(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(account5, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account5)
			},
		},
		{
			name:      "InvalidID",
			accountID: -1,
			body: updateAccountJSONRequest{
				Title:       &account2.Title,
				Description: &account2.Description,
				InitBalance: decimal.NewNullDecimal(account2.InitBalance),
			},
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
				store.EXPECT().
					UpdateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "NotFound",
			accountID: account2.ID,
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalErrorGet",
			accountID: account2.ID,
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "InternalErrorUpdate",
			accountID: account2.ID,
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account2.ID)).
					Times(1).
					Return(account1, nil)
				store.EXPECT().
					UpdateAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	// Checking cases
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubds(store)

			// start test server and send request
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func randomAccount(owner string) db.Account {
	return db.Account{
		ID:           util.RandomInt(1, 1000),
		Owner:        owner,
		Title:        util.RandomTitle(),
		Description:  util.RandomString(14),
		InitBalance:  util.RandomMoney(),
		Balance:      decimal.Zero,
		FinalBalance: decimal.Zero,
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)

	// Have to check each field as decimal.Decimal have to be compared each other directly
	require.Equal(t, account.ID, gotAccount.ID)
	require.Equal(t, account.Owner, gotAccount.Owner)
	require.Equal(t, account.Title, gotAccount.Title)
	require.Equal(t, account.Description, gotAccount.Description)
	require.Equal(t, account.Title, gotAccount.Title)
	require.True(t, account.InitBalance.Equal(gotAccount.InitBalance))
	require.True(t, account.Balance.Equal(gotAccount.Balance))
	require.True(t, account.FinalBalance.Equal(gotAccount.FinalBalance))
}

func requireBodyMatchAccounts(t *testing.T, body *bytes.Buffer, accounts []db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccounts []db.Account
	err = json.Unmarshal(data, &gotAccounts)
	require.NoError(t, err)

	for i, account := range accounts {
		gotAccount := gotAccounts[i]
		// Have to check each field as decimal.Decimal have to be compared each other directly
		require.Equal(t, account.ID, gotAccount.ID)
		require.Equal(t, account.Owner, gotAccount.Owner)
		require.Equal(t, account.Title, gotAccount.Title)
		require.Equal(t, account.Description, gotAccount.Description)
		require.Equal(t, account.Title, gotAccount.Title)
		require.True(t, account.InitBalance.Equal(gotAccount.InitBalance))
		require.True(t, account.Balance.Equal(gotAccount.Balance))
		require.True(t, account.FinalBalance.Equal(gotAccount.FinalBalance))
	}
}

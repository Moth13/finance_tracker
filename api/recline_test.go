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
	mockdb "github.com/moth13/finance_tracker/db/mock"
	db "github.com/moth13/finance_tracker/db/sqlc"
	"github.com/moth13/finance_tracker/util"
	"github.com/stretchr/testify/require"
)

func TestCreateRecLineAPI(t *testing.T) {
	user, _ := randomUser(t)
	user.Username = "jose"
	account := randomAccount(user.Username)
	category := randomCategory(user.Username)

	recline := randomRecLine(user, account, category)

	// Test cases definition
	testCases := []struct {
		name          string
		body          createRecLineRequest
		buildStubds   func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: createRecLineRequest{
				Title:       recline.Title,
				AccountID:   recline.AccountID,
				Amount:      recline.Amount,
				Description: recline.Description,
				DueDate:     recline.DueDate,
				CategoryID:  recline.CategoryID,
				Recurrency:  recline.Recurrency,
			},
			buildStubds: func(store *mockdb.MockStore) {
				arg := db.CreateRecLineParams{
					Owner:       recline.Owner,
					Title:       recline.Title,
					AccountID:   recline.AccountID,
					CategoryID:  recline.CategoryID,
					Amount:      recline.Amount,
					Description: recline.Description,
					DueDate:     recline.DueDate,
					Recurrency:  recline.Recurrency,
				}

				store.EXPECT().
					CreateRecLine(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(recline, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchRecLine(t, recorder.Body, recline)
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

			url := "/api/reclines"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestDeleteRecLineAPI(t *testing.T) {
	user, _ := randomUser(t)
	user.Username = "jose"
	account := randomAccount(user.Username)
	category := randomCategory(user.Username)

	recline := randomRecLine(user, account, category)

	// Test cases definition
	testCases := []struct {
		name          string
		reclineID     int64
		buildStubds   func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			reclineID: recline.ID,
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteRecLine(gomock.Any(), gomock.Eq(recline.ID)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:      "InvalidID",
			reclineID: -1,
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteRecLine(gomock.Any(), gomock.Eq(recline.ID)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "NotFound",
			reclineID: recline.ID,
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteRecLine(gomock.Any(), gomock.Any()).
					Times(1).
					Return(sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalServerError",
			reclineID: recline.ID,
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteRecLine(gomock.Any(), gomock.Any()).
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

			url := fmt.Sprintf("/api/reclines/%d", tc.reclineID)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestGetRecLineAPI(t *testing.T) {
	user, _ := randomUser(t)
	user.Username = "jose"
	account := randomAccount(user.Username)
	category := randomCategory(user.Username)

	recline := randomRecLine(user, account, category)

	// Test cases definition
	testCases := []struct {
		name          string
		reclineID     int64
		buildStubds   func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			reclineID: recline.ID,
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetRecLine(gomock.Any(), gomock.Eq(recline.ID)).
					Times(1).
					Return(recline, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchRecLine(t, recorder.Body, recline)
			},
		},
		{
			name:      "NotFound",
			reclineID: recline.ID,
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetRecLine(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Recline{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalServerError",
			reclineID: recline.ID,
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetRecLine(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Recline{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "InvalidRecLineID",
			reclineID: 0,
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetRecLine(gomock.Any(), gomock.Any()).
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

			url := fmt.Sprintf("/api/reclines/%d", tc.reclineID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestListRecLinesAPI(t *testing.T) {
	user, _ := randomUser(t)
	user.Username = "jose"
	account := randomAccount(user.Username)
	category := randomCategory(user.Username)

	n := 5
	reclines := make([]db.Recline, n)
	for i := 0; i < n; i++ {
		reclines[i] = randomRecLine(user, account, category)
	}

	// Test cases definition
	testCases := []struct {
		name          string
		query         listRecLinesRequest
		buildStubds   func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: listRecLinesRequest{
				PageID:   1,
				PageSize: int32(n),
			},
			buildStubds: func(store *mockdb.MockStore) {
				arg := db.ListRecLinesParams{
					Owner:  user.Username,
					Offset: 0,
					Limit:  int32(n),
				}
				store.EXPECT().
					ListRecLines(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(reclines, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchRecLines(t, recorder.Body, reclines)
			},
		},
		{
			name: "InternalError",
			query: listRecLinesRequest{
				PageID:   1,
				PageSize: int32(n),
			},
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListRecLines(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Recline{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidPageSize",
			query: listRecLinesRequest{
				PageID:   1,
				PageSize: 2,
			},
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListRecLines(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidPageIndex",
			query: listRecLinesRequest{
				PageID:   -1,
				PageSize: int32(n),
			},
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListRecLines(gomock.Any(), gomock.Any()).
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

			request, err := http.NewRequest(http.MethodGet, "/api/reclines", nil)
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

func randomRecLine(user db.User, account db.Account, category db.Category) db.Recline {

	return db.Recline{
		ID:          util.RandomInt(1, 1000),
		Title:       util.RandomTitle(),
		Owner:       user.Username,
		AccountID:   account.ID,
		CategoryID:  category.ID,
		Amount:      util.RandomMoney(),
		DueDate:     util.RandomFutureDate(),
		Description: util.RandomString(14),
		Recurrency:  util.RandomRecurrency(),
	}
}

func requireBodyMatchRecLine(t *testing.T, body *bytes.Buffer, recline db.Recline) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotRecLine db.Recline
	err = json.Unmarshal(data, &gotRecLine)
	require.NoError(t, err)

	// Have to check each field as decimal.Decimal have to be compared each other directly
	require.Equal(t, recline.ID, gotRecLine.ID)
	require.Equal(t, recline.Owner, gotRecLine.Owner)
	require.Equal(t, recline.Title, gotRecLine.Title)
}

func requireBodyMatchRecLines(t *testing.T, body *bytes.Buffer, reclines []db.Recline) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotRecLines []db.Recline
	err = json.Unmarshal(data, &gotRecLines)
	require.NoError(t, err)

	for i, recline := range reclines {
		gotRecLine := gotRecLines[i]
		// Have to check each field as decimal.Decimal have to be compared each other directly
		require.Equal(t, recline.ID, gotRecLine.ID)
		require.Equal(t, recline.Owner, gotRecLine.Owner)
		require.Equal(t, recline.Title, gotRecLine.Title)

	}
}

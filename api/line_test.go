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
	"time"

	"github.com/golang/mock/gomock"
	mockdb "github.com/moth13/finance_tracker/db/mock"
	db "github.com/moth13/finance_tracker/db/sqlc"
	"github.com/moth13/finance_tracker/util"
	"github.com/stretchr/testify/require"
)

func TestCreateLineAPI(t *testing.T) {
	user, _ := randomUser(t)
	user.Username = "jose"
	year := randomYear(user.Username)
	month := randomMonth(user.Username, year)
	account := randomAccount(user.Username)
	category := randomCategory(user.Username)

	line := randomLine(user, month, year, account, category)

	result := db.AddLineTxResult{
		Line: line,
		Balance: util.Balance{
			MonthBalance:        month.Balance,
			MonthFinalBalance:   month.FinalBalance,
			YearBalance:         year.Balance,
			YearFinalBalance:    month.FinalBalance,
			AccountBalance:      year.Balance,
			AccountFinalBalance: month.FinalBalance,
		},
	}

	// Test cases definition
	testCases := []struct {
		name          string
		body          createLineRequest
		buildStubds   func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: createLineRequest{
				Title:       line.Title,
				AccountID:   line.AccountID,
				MonthID:     line.MonthID,
				YearID:      line.YearID,
				CategoryID:  line.CategoryID,
				Amount:      line.Amount,
				Checked:     &line.Checked,
				Description: line.Description,
				DueDate:     line.DueDate,
			},
			buildStubds: func(store *mockdb.MockStore) {
				arg := db.AddLineTxParams{
					Owner:       line.Owner,
					Title:       line.Title,
					AccountID:   line.AccountID,
					MonthID:     line.MonthID,
					YearID:      line.YearID,
					CategoryID:  line.CategoryID,
					Amount:      line.Amount,
					Checked:     line.Checked,
					Description: line.Description,
					DueDate:     line.DueDate,
				}

				store.EXPECT().
					AddLineTx(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(result, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAddingLine(t, recorder.Body, result)
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

			url := "/api/lines"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestDeleteLineAPI(t *testing.T) {
	user, _ := randomUser(t)
	user.Username = "jose"
	year := randomYear(user.Username)
	month := randomMonth(user.Username, year)
	account := randomAccount(user.Username)
	category := randomCategory(user.Username)

	line := randomLine(user, month, year, account, category)

	result := db.DeleteLineTxResult{
		Balance: util.Balance{
			MonthBalance:        month.Balance,
			MonthFinalBalance:   month.FinalBalance,
			YearBalance:         year.Balance,
			YearFinalBalance:    month.FinalBalance,
			AccountBalance:      year.Balance,
			AccountFinalBalance: month.FinalBalance,
		},
	}

	// Test cases definition
	testCases := []struct {
		name          string
		lineID        int64
		buildStubds   func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			lineID: line.ID,
			buildStubds: func(store *mockdb.MockStore) {
				arg := db.DeleteLineTxParams{
					ID: line.ID,
				}

				store.EXPECT().
					DeleteLineTx(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(result, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:   "InvalidID",
			lineID: -1,
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteLineTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "NotFound",
			lineID: 10,
			buildStubds: func(store *mockdb.MockStore) {
				arg := db.DeleteLineTxParams{
					ID: 10,
				}
				store.EXPECT().
					DeleteLineTx(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(result, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:   "InternalServerError",
			lineID: line.ID,
			buildStubds: func(store *mockdb.MockStore) {
				arg := db.DeleteLineTxParams{
					ID: line.ID,
				}
				store.EXPECT().
					DeleteLineTx(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(result, sql.ErrConnDone)
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

			url := fmt.Sprintf("/api/lines/%d", tc.lineID)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestGetLineAPI(t *testing.T) {
	user, _ := randomUser(t)
	user.Username = "jose"
	year := randomYear(user.Username)
	month := randomMonth(user.Username, year)
	account := randomAccount(user.Username)
	category := randomCategory(user.Username)

	line := randomLine(user, month, year, account, category)

	// Test cases definition
	testCases := []struct {
		name          string
		lineID        int64
		buildStubds   func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			lineID: line.ID,
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetLine(gomock.Any(), gomock.Eq(line.ID)).
					Times(1).
					Return(line, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchLine(t, recorder.Body, line)
			},
		},
		{
			name:   "NotFound",
			lineID: line.ID,
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetLine(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Line{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:   "InternalServerError",
			lineID: line.ID,
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetLine(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Line{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:   "InvalidLineID",
			lineID: 0,
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetLine(gomock.Any(), gomock.Any()).
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

			url := fmt.Sprintf("/api/lines/%d", tc.lineID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestListLinesAPI(t *testing.T) {
	user, _ := randomUser(t)
	user.Username = "jose"
	year := randomYear(user.Username)
	month := randomMonth(user.Username, year)
	account := randomAccount(user.Username)
	category := randomCategory(user.Username)

	n := 5
	lines := make([]db.Line, n)
	for i := 0; i < n; i++ {
		lines[i] = randomLine(user, month, year, account, category)
	}

	// Test cases definition
	testCases := []struct {
		name          string
		query         listCategoriesRequest
		buildStubds   func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: listCategoriesRequest{
				PageID:   1,
				PageSize: int32(n),
			},
			buildStubds: func(store *mockdb.MockStore) {
				arg := db.ListLinesParams{
					Owner:  user.Username,
					Offset: 0,
					Limit:  int32(n),
				}
				store.EXPECT().
					ListLines(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(lines, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchLines(t, recorder.Body, lines)
			},
		},
		{
			name: "InternalError",
			query: listCategoriesRequest{
				PageID:   1,
				PageSize: int32(n),
			},
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListLines(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Line{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidPageSize",
			query: listCategoriesRequest{
				PageID:   1,
				PageSize: 2,
			},
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListLines(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidPageIndex",
			query: listCategoriesRequest{
				PageID:   -1,
				PageSize: int32(n),
			},
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListLines(gomock.Any(), gomock.Any()).
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

			request, err := http.NewRequest(http.MethodGet, "/api/lines", nil)
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

func randomLine(user db.User, month db.Month, year db.Year, account db.Account, category db.Category) db.Line {

	return db.Line{
		ID:          util.RandomInt(1, 1000),
		Title:       util.RandomTitle(),
		Owner:       user.Username,
		AccountID:   account.ID,
		MonthID:     month.ID,
		YearID:      year.ID,
		CategoryID:  category.ID,
		Amount:      util.RandomMoney(),
		Checked:     util.RandomBool(),
		DueDate:     util.RandomFutureDate(),
		Description: util.RandomString(14),
	}
}

func requireBodyMatchAddingLine(t *testing.T, body *bytes.Buffer, result db.AddLineTxResult) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotResult db.AddLineTxResult
	err = json.Unmarshal(data, &gotResult)
	require.NoError(t, err)

	checkLine(t, gotResult.Line, result.Line)
	checkBalance(t, gotResult.Balance, result.Balance)
	checkBalance(t, gotResult.Balance, result.Balance)
}

func requireBodyMatchLine(t *testing.T, body *bytes.Buffer, line db.Line) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotLine db.Line
	err = json.Unmarshal(data, &gotLine)
	require.NoError(t, err)

	checkLine(t, line, gotLine)
}

func requireBodyMatchLines(t *testing.T, body *bytes.Buffer, lines []db.Line) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotLines []db.Line
	err = json.Unmarshal(data, &gotLines)
	require.NoError(t, err)

	for i, line := range lines {
		gotLine := gotLines[i]
		checkLine(t, line, gotLine)
	}
}

func checkLine(t *testing.T, line1 db.Line, line2 db.Line) {
	require.Equal(t, line1.ID, line2.ID)
	require.Equal(t, line1.Owner, line2.Owner)
	require.Equal(t, line1.Title, line2.Title)
	require.Equal(t, line1.Owner, line2.Owner)
	require.Equal(t, line1.Title, line2.Title)
	require.Equal(t, line1.AccountID, line2.AccountID)
	require.Equal(t, line1.MonthID, line2.MonthID)
	require.True(t, line1.Amount.Equal(line2.Amount))
	require.Equal(t, line1.Checked, line2.Checked)
	require.Equal(t, line1.Description, line2.Description)
	require.Equal(t, line1.YearID, line2.YearID)
	require.Equal(t, line1.CategoryID, line2.CategoryID)
	require.WithinDuration(t, line1.DueDate, line2.DueDate, time.Second)
}

func checkBalance(t *testing.T, balance1 util.Balance, balance2 util.Balance) {
	require.True(t, balance1.MonthBalance.Equal(balance2.MonthBalance))
	require.True(t, balance1.MonthFinalBalance.Equal(balance2.MonthFinalBalance))
	require.True(t, balance1.AccountBalance.Equal(balance2.AccountBalance))
	require.True(t, balance1.AccountFinalBalance.Equal(balance2.AccountFinalBalance))
	require.True(t, balance1.YearBalance.Equal(balance2.YearBalance))
	require.True(t, balance1.YearFinalBalance.Equal(balance2.YearFinalBalance))
}

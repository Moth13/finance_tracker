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
	"github.com/jackc/pgx/v5"
	mockdb "github.com/moth13/finance_tracker/db/mock"
	db "github.com/moth13/finance_tracker/db/sqlc"
	"github.com/moth13/finance_tracker/token"
	"github.com/moth13/finance_tracker/util"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestCreateMonthAPI(t *testing.T) {
	user, _ := randomUser(t)
	year := randomYear(user.Username)
	month := randomMonth(user.Username, year)

	// Test cases definition
	testCases := []struct {
		name          string
		body          createMonthRequest
		buildStubds   func(store *mockdb.MockStore)
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: createMonthRequest{
				Owner:       user.Username,
				Title:       month.Title,
				Description: month.Description,
				StartDate:   month.StartDate,
				EndDate:     month.EndDate,
				YearID:      month.YearID,
			},
			buildStubds: func(store *mockdb.MockStore) {
				arg := db.CreateMonthParams{
					Owner:       user.Username,
					Title:       month.Title,
					Description: month.Description,
					StartDate:   month.StartDate,
					EndDate:     month.EndDate,
					YearID:      month.YearID,
				}

				store.EXPECT().
					GetYear(gomock.Any(), gomock.Eq(arg.YearID)).
					Times(1).
					Return(year, nil)

				store.EXPECT().
					CreateMonth(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(month, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchMonth(t, recorder.Body, month)
			},
		},
		{
			name: "InvalidTitle",
			body: createMonthRequest{
				Title:       "",
				Owner:       month.Owner,
				Description: month.Description,
				StartDate:   month.StartDate,
				EndDate:     month.EndDate,
				YearID:      month.YearID,
			},
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateMonth(gomock.Any(), gomock.Any()).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidDescription",
			body: createMonthRequest{
				Title:       month.Title,
				Owner:       month.Owner,
				Description: "",
				StartDate:   month.StartDate,
				EndDate:     month.EndDate,
				YearID:      month.YearID,
			},
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateMonth(gomock.Any(), gomock.Any()).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidYearNotFound",
			body: createMonthRequest{
				Owner:       month.Owner,
				Title:       month.Title,
				Description: month.Description,
				StartDate:   month.StartDate,
				EndDate:     month.EndDate,
				YearID:      3,
			},
			buildStubds: func(store *mockdb.MockStore) {
				arg := db.CreateMonthParams{
					YearID: 3,
				}

				store.EXPECT().
					GetYear(gomock.Any(), gomock.Eq(arg.YearID)).
					Times(1).
					Return(db.Year{}, pgx.ErrNoRows)

				store.EXPECT().
					CreateMonth(gomock.Any(), gomock.Any()).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InvalidYearInvalid",
			body: createMonthRequest{
				Owner:       month.Owner,
				Title:       month.Title,
				Description: month.Description,
				StartDate:   month.StartDate,
				EndDate:     month.EndDate,
				YearID:      -1,
			},
			buildStubds: func(store *mockdb.MockStore) {

				store.EXPECT().
					GetYear(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					CreateMonth(gomock.Any(), gomock.Any()).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
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

			url := "/api/months"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestDeleteMonthAPI(t *testing.T) {
	user, _ := randomUser(t)
	year := randomYear(user.Username)
	month := randomMonth(user.Username, year)

	// Test cases definition
	testCases := []struct {
		name          string
		monthID       int64
		buildStubds   func(store *mockdb.MockStore)
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:    "OK",
			monthID: month.ID,
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteMonth(gomock.Any(), gomock.Eq(month.ID)).
					Times(1).
					Return(nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:    "InvalidID",
			monthID: -1,
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteMonth(gomock.Any(), gomock.Eq(month.ID)).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:    "NotFound",
			monthID: month.ID,
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteMonth(gomock.Any(), gomock.Any()).
					Times(1).
					Return(sql.ErrNoRows)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:    "InternalServerError",
			monthID: month.ID,
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteMonth(gomock.Any(), gomock.Any()).
					Times(1).
					Return(sql.ErrConnDone)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
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

			url := fmt.Sprintf("/api/months/%d", tc.monthID)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestGetMonthAPI(t *testing.T) {
	user, _ := randomUser(t)
	year := randomYear(user.Username)
	month := randomMonth(user.Username, year)

	// Test cases definition
	testCases := []struct {
		name          string
		monthID       int64
		buildStubds   func(store *mockdb.MockStore)
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:    "OK",
			monthID: month.ID,
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetMonth(gomock.Any(), gomock.Eq(month.ID)).
					Times(1).
					Return(month, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchMonth(t, recorder.Body, month)
			},
		},
		{
			name:    "NotFound",
			monthID: month.ID,
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetMonth(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Month{}, sql.ErrNoRows)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:    "InternalServerError",
			monthID: month.ID,
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetMonth(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Month{}, sql.ErrConnDone)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:    "InvalidMonthID",
			monthID: 0,
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetMonth(gomock.Any(), gomock.Any()).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
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

			url := fmt.Sprintf("/api/months/%d", tc.monthID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestListMonthsAPI(t *testing.T) {
	user, _ := randomUser(t)
	year := randomYear(user.Username)
	n := 5
	months := make([]db.Month, n)
	for i := 0; i < n; i++ {
		months[i] = randomMonth(user.Username, year)
	}

	// Test cases definition
	testCases := []struct {
		name          string
		query         listMonthRequest
		buildStubds   func(store *mockdb.MockStore)
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: listMonthRequest{
				PageID:   1,
				PageSize: int32(n),
			},
			buildStubds: func(store *mockdb.MockStore) {
				arg := db.ListMonthsParams{
					Owner:  user.Username,
					Offset: 0,
					Limit:  int32(n),
				}
				store.EXPECT().
					ListMonths(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(months, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchMonths(t, recorder.Body, months)
			},
		},
		{
			name: "InternalError",
			query: listMonthRequest{
				PageID:   1,
				PageSize: int32(n),
			},
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListMonths(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Month{}, sql.ErrConnDone)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidPageSize",
			query: listMonthRequest{
				PageID:   1,
				PageSize: 2,
			},
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListMonths(gomock.Any(), gomock.Any()).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidPageIndex",
			query: listMonthRequest{
				PageID:   -1,
				PageSize: int32(n),
			},
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListMonths(gomock.Any(), gomock.Any()).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
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

			request, err := http.NewRequest(http.MethodGet, "/api/months", nil)
			require.NoError(t, err)

			q := request.URL.Query()
			q.Add("page_id", fmt.Sprintf("%d", tc.query.PageID))
			q.Add("page_size", fmt.Sprintf("%d", tc.query.PageSize))
			request.URL.RawQuery = q.Encode()

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestUpdateMonthAPI(t *testing.T) {
	user, _ := randomUser(t)
	year1 := randomYear(user.Username)
	year2 := randomYear(user.Username)
	month1 := randomMonth(user.Username, year1)
	month2 := randomMonth(user.Username, year2)

	month2.ID = month1.ID
	month2.Owner = month1.Owner

	month3 := db.Month{
		ID:           month2.ID,
		Owner:        month2.Owner,
		Title:        month1.Title,
		Description:  month2.Description,
		StartDate:    month2.StartDate,
		EndDate:      month2.EndDate,
		Balance:      month2.Balance,
		FinalBalance: month2.FinalBalance,
		YearID:       month2.YearID,
	}

	month4 := db.Month{
		ID:           month2.ID,
		Owner:        month2.Owner,
		Title:        month2.Title,
		Description:  month1.Description,
		StartDate:    month2.StartDate,
		EndDate:      month2.EndDate,
		Balance:      month2.Balance,
		FinalBalance: month2.FinalBalance,
		YearID:       month2.YearID,
	}

	month5 := db.Month{
		ID:           month2.ID,
		Owner:        month2.Owner,
		Title:        month2.Title,
		Description:  month2.Description,
		StartDate:    month1.StartDate,
		EndDate:      month2.EndDate,
		Balance:      month2.Balance,
		FinalBalance: month2.FinalBalance,
		YearID:       month2.YearID,
	}

	month6 := db.Month{
		ID:           month2.ID,
		Owner:        month2.Owner,
		Title:        month2.Title,
		Description:  month2.Description,
		StartDate:    month2.StartDate,
		EndDate:      month1.EndDate,
		Balance:      month2.Balance,
		FinalBalance: month2.FinalBalance,
		YearID:       month2.YearID,
	}

	month7 := db.Month{
		ID:           month2.ID,
		Owner:        month2.Owner,
		Title:        month2.Title,
		Description:  month2.Description,
		StartDate:    month2.StartDate,
		EndDate:      month2.EndDate,
		Balance:      month2.Balance,
		FinalBalance: month2.FinalBalance,
		YearID:       month1.YearID,
	}

	// Test cases definition
	testCases := []struct {
		name          string
		monthID       int64
		body          updateMonthJSONRequest
		buildStubds   func(store *mockdb.MockStore)
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:    "OK",
			monthID: month1.ID,
			body: updateMonthJSONRequest{
				Title:       &month2.Title,
				Description: &month2.Description,
				StartDate:   &month2.StartDate,
				EndDate:     &month2.EndDate,
				YearID:      &month2.YearID,
			},
			buildStubds: func(store *mockdb.MockStore) {
				arg := db.UpdateMonthParams{
					ID:          month2.ID,
					Title:       month2.Title,
					Description: month2.Description,
					StartDate:   month2.StartDate,
					EndDate:     month2.EndDate,
					YearID:      month2.YearID,
				}
				store.EXPECT().
					GetMonth(gomock.Any(), gomock.Eq(month2.ID)).
					Times(1).
					Return(month1, nil)
				store.EXPECT().
					UpdateMonth(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(month2, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchMonth(t, recorder.Body, month2)
			},
		},
		{
			name:    "NoTitle",
			monthID: month1.ID,
			body: updateMonthJSONRequest{
				Description: &month2.Description,
				StartDate:   &month2.StartDate,
				EndDate:     &month2.EndDate,
				YearID:      &month2.YearID,
			},
			buildStubds: func(store *mockdb.MockStore) {

				arg := db.UpdateMonthParams{
					ID:          month2.ID,
					Title:       month1.Title,
					Description: month2.Description,
					StartDate:   month2.StartDate,
					EndDate:     month2.EndDate,
					YearID:      month2.YearID,
				}
				store.EXPECT().
					GetMonth(gomock.Any(), gomock.Eq(month2.ID)).
					Times(1).
					Return(month1, nil)
				store.EXPECT().
					UpdateMonth(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(month3, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchMonth(t, recorder.Body, month3)
			},
		},
		{
			name:    "NoDescription",
			monthID: month1.ID,
			body: updateMonthJSONRequest{
				Title:     &month2.Title,
				StartDate: &month2.StartDate,
				EndDate:   &month2.EndDate,
				YearID:    &month2.YearID,
			},
			buildStubds: func(store *mockdb.MockStore) {

				arg := db.UpdateMonthParams{
					ID:          month2.ID,
					Title:       month2.Title,
					Description: month1.Description,
					StartDate:   month2.StartDate,
					EndDate:     month2.EndDate,
					YearID:      month2.YearID,
				}
				store.EXPECT().
					GetMonth(gomock.Any(), gomock.Eq(month2.ID)).
					Times(1).
					Return(month1, nil)
				store.EXPECT().
					UpdateMonth(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(month4, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchMonth(t, recorder.Body, month4)
			},
		},
		{
			name:    "NoStartDate",
			monthID: month1.ID,
			body: updateMonthJSONRequest{
				Title:       &month2.Title,
				Description: &month2.Description,
				EndDate:     &month2.EndDate,
				YearID:      &month2.YearID,
			},
			buildStubds: func(store *mockdb.MockStore) {

				arg := db.UpdateMonthParams{
					ID:          month2.ID,
					Title:       month2.Title,
					Description: month2.Description,
					StartDate:   month1.StartDate,
					EndDate:     month2.EndDate,
					YearID:      month2.YearID,
				}
				store.EXPECT().
					GetMonth(gomock.Any(), gomock.Eq(month2.ID)).
					Times(1).
					Return(month1, nil)
				store.EXPECT().
					UpdateMonth(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(month5, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchMonth(t, recorder.Body, month5)
			},
		},
		{
			name:    "NoEndDate",
			monthID: month1.ID,
			body: updateMonthJSONRequest{
				Title:       &month2.Title,
				Description: &month2.Description,
				StartDate:   &month2.StartDate,
				YearID:      &month2.YearID,
			},
			buildStubds: func(store *mockdb.MockStore) {

				arg := db.UpdateMonthParams{
					ID:          month2.ID,
					Title:       month2.Title,
					Description: month2.Description,
					StartDate:   month2.StartDate,
					EndDate:     month1.EndDate,
					YearID:      month2.YearID,
				}
				store.EXPECT().
					GetMonth(gomock.Any(), gomock.Eq(month2.ID)).
					Times(1).
					Return(month1, nil)
				store.EXPECT().
					UpdateMonth(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(month6, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchMonth(t, recorder.Body, month6)
			},
		},
		{
			name:    "NoYear",
			monthID: month1.ID,
			body: updateMonthJSONRequest{
				Title:       &month2.Title,
				Description: &month2.Description,
				StartDate:   &month2.StartDate,
				EndDate:     &month2.EndDate,
			},
			buildStubds: func(store *mockdb.MockStore) {

				arg := db.UpdateMonthParams{
					ID:          month2.ID,
					Title:       month2.Title,
					Description: month2.Description,
					StartDate:   month2.StartDate,
					EndDate:     month2.EndDate,
					YearID:      month1.YearID,
				}
				store.EXPECT().
					GetMonth(gomock.Any(), gomock.Eq(month2.ID)).
					Times(1).
					Return(month1, nil)
				store.EXPECT().
					UpdateMonth(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(month7, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchMonth(t, recorder.Body, month7)
			},
		},
		{
			name:    "InvalidID",
			monthID: -1,
			body: updateMonthJSONRequest{
				Title:       &month2.Title,
				Description: &month2.Description,
				StartDate:   &month2.StartDate,
				EndDate:     &month2.EndDate,
			},
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetMonth(gomock.Any(), gomock.Any()).
					Times(0)
				store.EXPECT().
					UpdateMonth(gomock.Any(), gomock.Any()).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:    "NotFound",
			monthID: month2.ID,
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetMonth(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Month{}, sql.ErrNoRows)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:    "InternalErrorGet",
			monthID: month2.ID,
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetMonth(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Month{}, sql.ErrConnDone)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:    "InternalErrorUpdate",
			monthID: month2.ID,
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetMonth(gomock.Any(), gomock.Eq(month2.ID)).
					Times(1).
					Return(month1, nil)
				store.EXPECT().
					UpdateMonth(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Month{}, sql.ErrConnDone)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
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

			url := fmt.Sprintf("/api/months/%d", tc.monthID)
			request, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func randomMonth(owner string, year db.Year) db.Month {
	start, end := util.RandomMonthDate()

	return db.Month{
		ID:           util.RandomInt(1, 1000),
		Owner:        owner,
		Title:        util.RandomTitle(),
		Description:  util.RandomString(14),
		StartDate:    start,
		EndDate:      end,
		Balance:      decimal.Zero,
		FinalBalance: decimal.Zero,
		YearID:       year.ID,
	}
}

func requireBodyMatchMonth(t *testing.T, body *bytes.Buffer, month db.Month) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotMonth db.Month
	err = json.Unmarshal(data, &gotMonth)
	require.NoError(t, err)

	// Have to check each field as decimal.Decimal have to be compared each other directly
	require.Equal(t, month.ID, gotMonth.ID)
	require.Equal(t, month.Owner, gotMonth.Owner)
	require.Equal(t, month.Title, gotMonth.Title)
	require.Equal(t, month.Description, gotMonth.Description)
	require.Equal(t, month.Title, gotMonth.Title)
	require.WithinDuration(t, month.StartDate, gotMonth.StartDate, time.Second)
	require.WithinDuration(t, month.EndDate, gotMonth.EndDate, time.Second)
}

func requireBodyMatchMonths(t *testing.T, body *bytes.Buffer, months []db.Month) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotMonths []db.Month
	err = json.Unmarshal(data, &gotMonths)
	require.NoError(t, err)

	for i, month := range months {
		gotMonth := gotMonths[i]
		// Have to check each field as decimal.Decimal have to be compared each other directly
		require.Equal(t, month.ID, gotMonth.ID)
		require.Equal(t, month.Owner, gotMonth.Owner)
		require.Equal(t, month.Title, gotMonth.Title)
		require.Equal(t, month.Description, gotMonth.Description)
		require.Equal(t, month.Title, gotMonth.Title)
		require.WithinDuration(t, month.StartDate, gotMonth.StartDate, time.Second)
		require.WithinDuration(t, month.EndDate, gotMonth.EndDate, time.Second)
	}
}

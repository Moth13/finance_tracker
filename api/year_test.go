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
	"github.com/moth13/finance_tracker/token"
	"github.com/moth13/finance_tracker/util"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestCreateYearAPI(t *testing.T) {
	user, _ := randomUser(t)
	year := randomYear(user.Username)

	// Test cases definition
	testCases := []struct {
		name          string
		body          createYearRequest
		buildStubds   func(store *mockdb.MockStore)
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: createYearRequest{
				Title:       year.Title,
				Description: year.Description,
				StartDate:   year.StartDate,
				EndDate:     year.EndDate,
			},
			buildStubds: func(store *mockdb.MockStore) {
				arg := db.CreateYearParams{
					Owner:       year.Owner,
					Title:       year.Title,
					Description: year.Description,
					StartDate:   year.StartDate,
					EndDate:     year.EndDate,
				}

				store.EXPECT().
					CreateYear(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(year, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchYear(t, recorder.Body, year)
			},
		},
		{
			name: "InvalidTitle",
			body: createYearRequest{
				Title:       "",
				Description: year.Description,
				StartDate:   year.StartDate,
				EndDate:     year.EndDate,
			},
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateYear(gomock.Any(), gomock.Any()).
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
			body: createYearRequest{
				Title:       year.Title,
				Description: "",
				StartDate:   year.StartDate,
				EndDate:     year.EndDate,
			},
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateYear(gomock.Any(), gomock.Any()).
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

			url := "/api/years"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestDeleteYearAPI(t *testing.T) {
	user, _ := randomUser(t)
	year := randomYear(user.Username)

	// Test cases definition
	testCases := []struct {
		name          string
		yearID        int64
		buildStubds   func(store *mockdb.MockStore)
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			yearID: year.ID,
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteYear(gomock.Any(), gomock.Eq(year.ID)).
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
			name:   "InvalidID",
			yearID: -1,
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteYear(gomock.Any(), gomock.Eq(year.ID)).
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
			name:   "NotFound",
			yearID: year.ID,
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteYear(gomock.Any(), gomock.Any()).
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
			name:   "InternalServerError",
			yearID: year.ID,
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteYear(gomock.Any(), gomock.Any()).
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

			url := fmt.Sprintf("/api/years/%d", tc.yearID)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestGetYearAPI(t *testing.T) {
	user, _ := randomUser(t)
	year := randomYear(user.Username)

	// Test cases definition
	testCases := []struct {
		name          string
		yearID        int64
		buildStubds   func(store *mockdb.MockStore)
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			yearID: year.ID,
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetYear(gomock.Any(), gomock.Eq(year.ID)).
					Times(1).
					Return(year, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchYear(t, recorder.Body, year)
			},
		},
		{
			name:   "NotFound",
			yearID: year.ID,
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetYear(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Year{}, sql.ErrNoRows)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:   "InternalServerError",
			yearID: year.ID,
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetYear(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Year{}, sql.ErrConnDone)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:   "InvalidYearID",
			yearID: 0,
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetYear(gomock.Any(), gomock.Any()).
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

			url := fmt.Sprintf("/api/years/%d", tc.yearID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestListYearsAPI(t *testing.T) {
	user, _ := randomUser(t)
	n := 5
	years := make([]db.Year, n)
	for i := range n {
		years[i] = randomYear(user.Username)
	}

	// Test cases definition
	testCases := []struct {
		name          string
		query         listYearRequest
		buildStubds   func(store *mockdb.MockStore)
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: listYearRequest{
				PageID:   1,
				PageSize: int32(n),
			},
			buildStubds: func(store *mockdb.MockStore) {
				arg := db.ListYearsParams{
					Owner:  user.Username,
					Offset: 0,
					Limit:  int32(n),
				}
				store.EXPECT().
					ListYears(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(years, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchYears(t, recorder.Body, years)
			},
		},
		{
			name: "InternalError",
			query: listYearRequest{
				PageID:   1,
				PageSize: int32(n),
			},
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListYears(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Year{}, sql.ErrConnDone)
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
			query: listYearRequest{
				PageID:   1,
				PageSize: 2,
			},
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListYears(gomock.Any(), gomock.Any()).
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
			query: listYearRequest{
				PageID:   -1,
				PageSize: int32(n),
			},
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListYears(gomock.Any(), gomock.Any()).
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

			request, err := http.NewRequest(http.MethodGet, "/api/years", nil)
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

func TestUpdateYearAPI(t *testing.T) {
	user, _ := randomUser(t)
	year1 := randomYear(user.Username)
	year2 := randomYear(user.Username)

	year2.ID = year1.ID
	year2.Owner = year1.Owner

	year3 := db.Year{
		ID:           year2.ID,
		Owner:        year2.Owner,
		Title:        year1.Title,
		Description:  year2.Description,
		StartDate:    year2.StartDate,
		EndDate:      year2.EndDate,
		Balance:      year2.Balance,
		FinalBalance: year2.FinalBalance,
	}

	year4 := db.Year{
		ID:           year2.ID,
		Owner:        year2.Owner,
		Title:        year2.Title,
		Description:  year1.Description,
		StartDate:    year2.StartDate,
		EndDate:      year2.EndDate,
		Balance:      year2.Balance,
		FinalBalance: year2.FinalBalance,
	}

	year5 := db.Year{
		ID:           year2.ID,
		Owner:        year2.Owner,
		Title:        year2.Title,
		Description:  year2.Description,
		StartDate:    year1.StartDate,
		EndDate:      year2.EndDate,
		Balance:      year2.Balance,
		FinalBalance: year2.FinalBalance,
	}

	year6 := db.Year{
		ID:           year2.ID,
		Owner:        year2.Owner,
		Title:        year2.Title,
		Description:  year2.Description,
		StartDate:    year2.StartDate,
		EndDate:      year1.EndDate,
		Balance:      year2.Balance,
		FinalBalance: year2.FinalBalance,
	}

	// Test cases definition
	testCases := []struct {
		name          string
		yearID        int64
		body          updateYearJSONRequest
		buildStubds   func(store *mockdb.MockStore)
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			yearID: year1.ID,
			body: updateYearJSONRequest{
				Title:       &year2.Title,
				Description: &year2.Description,
				StartDate:   &year2.StartDate,
				EndDate:     &year2.EndDate,
			},
			buildStubds: func(store *mockdb.MockStore) {
				arg := db.UpdateYearParams{
					ID:          year2.ID,
					Title:       year2.Title,
					Description: year2.Description,
					StartDate:   year2.StartDate,
					EndDate:     year2.EndDate,
				}
				store.EXPECT().
					GetYear(gomock.Any(), gomock.Eq(year2.ID)).
					Times(1).
					Return(year1, nil)
				store.EXPECT().
					UpdateYear(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(year2, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchYear(t, recorder.Body, year2)
			},
		},
		{
			name:   "NoTitle",
			yearID: year1.ID,
			body: updateYearJSONRequest{
				Description: &year2.Description,
				StartDate:   &year2.StartDate,
				EndDate:     &year2.EndDate,
			},
			buildStubds: func(store *mockdb.MockStore) {

				arg := db.UpdateYearParams{
					ID:          year2.ID,
					Title:       year1.Title,
					Description: year2.Description,
					StartDate:   year2.StartDate,
					EndDate:     year2.EndDate,
				}
				store.EXPECT().
					GetYear(gomock.Any(), gomock.Eq(year2.ID)).
					Times(1).
					Return(year1, nil)
				store.EXPECT().
					UpdateYear(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(year3, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchYear(t, recorder.Body, year3)
			},
		},
		{
			name:   "NoDescription",
			yearID: year1.ID,
			body: updateYearJSONRequest{
				Title:     &year2.Title,
				StartDate: &year2.StartDate,
				EndDate:   &year2.EndDate,
			},
			buildStubds: func(store *mockdb.MockStore) {

				arg := db.UpdateYearParams{
					ID:          year2.ID,
					Title:       year2.Title,
					Description: year1.Description,
					StartDate:   year2.StartDate,
					EndDate:     year2.EndDate,
				}
				store.EXPECT().
					GetYear(gomock.Any(), gomock.Eq(year2.ID)).
					Times(1).
					Return(year1, nil)
				store.EXPECT().
					UpdateYear(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(year4, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchYear(t, recorder.Body, year4)
			},
		},
		{
			name:   "NoStartDate",
			yearID: year1.ID,
			body: updateYearJSONRequest{
				Title:       &year2.Title,
				Description: &year2.Description,
				EndDate:     &year2.EndDate,
			},
			buildStubds: func(store *mockdb.MockStore) {

				arg := db.UpdateYearParams{
					ID:          year2.ID,
					Title:       year2.Title,
					Description: year2.Description,
					StartDate:   year1.StartDate,
					EndDate:     year2.EndDate,
				}
				store.EXPECT().
					GetYear(gomock.Any(), gomock.Eq(year2.ID)).
					Times(1).
					Return(year1, nil)
				store.EXPECT().
					UpdateYear(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(year5, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchYear(t, recorder.Body, year5)
			},
		},
		{
			name:   "NoEndDate",
			yearID: year1.ID,
			body: updateYearJSONRequest{
				Title:       &year2.Title,
				Description: &year2.Description,
				StartDate:   &year2.StartDate,
			},
			buildStubds: func(store *mockdb.MockStore) {

				arg := db.UpdateYearParams{
					ID:          year2.ID,
					Title:       year2.Title,
					Description: year2.Description,
					StartDate:   year2.StartDate,
					EndDate:     year1.EndDate,
				}
				store.EXPECT().
					GetYear(gomock.Any(), gomock.Eq(year2.ID)).
					Times(1).
					Return(year1, nil)
				store.EXPECT().
					UpdateYear(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(year6, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchYear(t, recorder.Body, year6)
			},
		},
		{
			name:   "InvalidID",
			yearID: -1,
			body: updateYearJSONRequest{
				Title:       &year2.Title,
				Description: &year2.Description,
				StartDate:   &year2.StartDate,
				EndDate:     &year2.EndDate,
			},
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetYear(gomock.Any(), gomock.Any()).
					Times(0)
				store.EXPECT().
					UpdateYear(gomock.Any(), gomock.Any()).
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
			name:   "NotFound",
			yearID: year2.ID,
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetYear(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Year{}, sql.ErrNoRows)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:   "InternalErrorGet",
			yearID: year2.ID,
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetYear(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Year{}, sql.ErrConnDone)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:   "InternalErrorUpdate",
			yearID: year2.ID,
			buildStubds: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetYear(gomock.Any(), gomock.Eq(year2.ID)).
					Times(1).
					Return(year1, nil)
				store.EXPECT().
					UpdateYear(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Year{}, sql.ErrConnDone)
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

			url := fmt.Sprintf("/api/years/%d", tc.yearID)
			request, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func randomYear(owner string) db.Year {
	start, end := util.RandomYearDate()

	return db.Year{
		ID:           util.RandomInt(1, 1000),
		Owner:        owner,
		Title:        util.RandomTitle(),
		Description:  util.RandomString(14),
		StartDate:    start,
		EndDate:      end,
		Balance:      decimal.Zero,
		FinalBalance: decimal.Zero,
	}
}

func requireBodyMatchYear(t *testing.T, body *bytes.Buffer, year db.Year) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotYear db.Year
	err = json.Unmarshal(data, &gotYear)
	require.NoError(t, err)

	// Have to check each field as decimal.Decimal have to be compared each other directly
	require.Equal(t, year.ID, gotYear.ID)
	require.Equal(t, year.Owner, gotYear.Owner)
	require.Equal(t, year.Title, gotYear.Title)
	require.Equal(t, year.Description, gotYear.Description)
	require.Equal(t, year.Title, gotYear.Title)
	require.WithinDuration(t, year.StartDate, gotYear.StartDate, time.Second)
	require.WithinDuration(t, year.EndDate, gotYear.EndDate, time.Second)
}

func requireBodyMatchYears(t *testing.T, body *bytes.Buffer, years []db.Year) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotYears []db.Year
	err = json.Unmarshal(data, &gotYears)
	require.NoError(t, err)

	for i, year := range years {
		gotYear := gotYears[i]
		// Have to check each field as decimal.Decimal have to be compared each other directly
		require.Equal(t, year.ID, gotYear.ID)
		require.Equal(t, year.Owner, gotYear.Owner)
		require.Equal(t, year.Title, gotYear.Title)
		require.Equal(t, year.Description, gotYear.Description)
		require.Equal(t, year.Title, gotYear.Title)
		require.WithinDuration(t, year.StartDate, gotYear.StartDate, time.Second)
		require.WithinDuration(t, year.EndDate, gotYear.EndDate, time.Second)
	}
}

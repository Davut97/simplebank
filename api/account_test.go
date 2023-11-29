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

	mockdb "github.com/Davut97/simplebank/db/mock"
	db "github.com/Davut97/simplebank/db/sqlc"
	"github.com/Davut97/simplebank/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetAccountAPI(t *testing.T) {
	account := randomAccount()
	testCases := []struct { // testCases is a slice of structs
		name           string                                                  // name is a string
		accountID      int64                                                   // accountID is an int64
		buildStubs     func(store *mockdb.MockStore)                           // buildStubs is a function that takes a store and returns nothing
		checkResponse  func(t *testing.T, recorder *httptest.ResponseRecorder) // checkResponse is a function that takes a testing.T and a recorder and returns nothing
		expectedStatus int                                                     // expectedStatus is an int
	}{
		{
			name:      "OK",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) { // buildStubs is a function that takes a store and returns nothing
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil) // store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) { // checkResponse is a function that takes a testing.T and a recorder and returns nothing
				require.Equal(t, http.StatusOK, recorder.Code)     // require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account) // requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:      "NotFound",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) { // buildStubs is a function that takes a store and returns nothing
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(db.Account{}, sql.ErrNoRows) // store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) { // checkResponse is a function that takes a testing.T and a recorder and returns nothing
				require.Equal(t, http.StatusNotFound, recorder.Code) // require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalError",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) { // buildStubs is a function that takes a store and returns nothing
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(db.Account{}, sql.ErrConnDone) // store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) { // checkResponse is a function that takes a testing.T and a recorder and returns nothing
				require.Equal(t, http.StatusInternalServerError, recorder.Code) // require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "InvalidID",
			accountID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0) // store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
			}, // buildStubs is a function that takes a store and returns nothing
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// checkResponse is a function that takes a testing.T and a recorder and returns nothing

				require.Equal(t, http.StatusBadRequest, recorder.Code) // require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "Error",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(1).Return(db.Account{}, sql.ErrConnDone) // store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(1).Return(db.Account{}, sql.ErrConnDone)
			}, // buildStubs is a function that takes a store and returns nothing
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) { // checkResponse is a function that takes a testing.T and a recorder and returns nothing
				require.Equal(t, http.StatusInternalServerError, recorder.Code) // require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}
	for i := range testCases { // for i := range testCases
		tc := testCases[i]                  // tc is a testCases[i]
		t.Run(tc.name, func(t *testing.T) { // t.Run(tc.name, func(t *testing.T)
			ctrl := gomock.NewController(t)                           // ctrl is a gomock.NewController(t)
			defer ctrl.Finish()                                       // defer ctrl.Finish()
			store := mockdb.NewMockStore(ctrl)                        // store is a mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)                                      // tc.buildStubs(store)
			server := NewServer(store)                                // server is a NewServer(store)
			recorder := httptest.NewRecorder()                        // recorder is a httptest.NewRecorder()
			url := fmt.Sprintf("/accounts/%d", tc.accountID)          // url is a fmt.Sprintf("/accounts/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil) // request, err is a http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)                                   // require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)                // server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)                             // tc.checkResponse(t, recorder)
		})
	}

}
func randomAccount() db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

}
func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)
	var getAccount db.Account
	err = json.Unmarshal(data, &getAccount)
	require.NoError(t, err)
	require.Equal(t, account, getAccount)
}

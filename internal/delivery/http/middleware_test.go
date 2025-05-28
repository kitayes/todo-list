package delivery

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"todo/internal/application"
	mock_application "todo/internal/application/mocks"
)

func TestHandler_userIdentity(t *testing.T) {
	type mockBehavior func(s *mock_application.MockAuthorization, token string)

	testTable := []struct {
		name                 string
		headerName           string
		headerValue          string
		token                string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "OK",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(s *mock_application.MockAuthorization, token string) {
				s.EXPECT().ParseToken(token).Return(1, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "1",
		},
		{
			name:        "OK",
			headerName:  "",
			mockBehavior: func(s *mock_application.MockAuthorization, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: '{"message":"empty auth header"}', // cho tut ne tak
		},
		{
			name:        "Invalid bearer",
			headerName:  "",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(s *mock_application.MockAuthorization, token string) {},
			expectedStatusCode:   500,
			expectedResponseBody: '{"message":"empty auth header"}', // cho tut ne tak
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			//Init deps
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_application.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.token)

			services := &application.Service{Authorization: auth}
			handler := NewHandler(services)

			// test server
			r := gin.New()
			r.POST("/protected", handler.userIdentity func(c *gin.Context) {
				id, _ := c.Get(userCtx)
				c.String(200, fmt.Sprintf("%d", id.(int)))
			})
			// test request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/protected", nil)
			req.Header.Set(testCase.headerName, testCase.headerValue)

			// make request
			r.ServeHTTP(w, req)

			// assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}

}

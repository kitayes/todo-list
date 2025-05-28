package delivery

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"todo/internal/application"
	mock_application "todo/internal/application/mocks"
	"todo/internal/models"
	service "todo/pkg/services"
)

func TestHandler_signUp(t *testing.T) {
	type mockBehavior func(s *mock_application.MockAuthorization.MockAuthorization, user models.User)

	testTable := []struct {
		name string
		inputBody string
		inputUser models.User
		mockBehavior mockBehavior
		expectedStatusCode int
		expectedRequestBody string
	} {
		{
			"OK",
			`{"name":"Test","username":"test","password":"qwerty"}`,
			models.User{
				Name:     "Test",
				Username: "test",
				Password: "qwerty",
			},
			func(s *mock_application.MockAuthorization.MockAuthorization, user models.User) {
				s.EXPECT().CreateUser(user).Return(1, nil)
			},
			200,
			`{"id":1}`,
		}, {
			name:"Empty Fields",
			inputBody: `{"name":"Test","username":"test","password":"qwerty"}`,
			mockBehavior: func(s *mock_application.MockAuthorization.MockAuthorization, user models.User) {},
			expectedStatusCode: 400,
			expectedRequestBody: `{"error":"invalid input body"}`, // maybe tut message vmesto error no vrode net
		},
		{
			name:"Service Failure",
			inputBody: `{"name":"Test","username":"test","password":"qwerty"}`,
			inputUser: models.User{
				Name: "Test",
				Username: "test",
				Password: "qwerty",
			},
			mockBehavior: func(s *mock_application.MockAuthorization.MockAuthorization, user models.User) {
				s.EXPECT().CreateUser(user).Return(1, errors.New("service failure"))
			},
			expectedStatusCode: 500,
			expectedRequestBody: `{"message":"service failure"'}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			//Init deps
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_application.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.inputUser)

			services := &application.Service{Authorization: auth}
			handler := NewHandler(services)

			//Test server
			r := gin.New()
			r.POST("/sign-up", handler.signUp) // тут крч хз че за трабла, ее уже нет но выглядит подозрительно

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-up", bytes.NewBufferString(testCase.inputBody))

			//Perform request
			r.ServeHTTP(w, req)

			// assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

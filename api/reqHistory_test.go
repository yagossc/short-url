package api_test

import (
	"net/http"
	"testing"

	"github.com/yagossc/short-url/app"
)

func TestHistory(t *testing.T) {

	var expected app.HistoryResponse
	expected.Count = 3

	tests := []httpTest{
		{
			testName:       "[GET] Full History",
			url:            "/history",
			body:           `{"url": "http://localhost:8080/yxZ8byjhRui"}`,
			expectedStatus: http.StatusOK,
			expectedBody:   expected,
		},
		{
			testName:       "[GET] Week's History",
			url:            "/history/week",
			body:           `{"url": "http://localhost:8080/yxZ8byjhRui"}`,
			expectedStatus: http.StatusOK,
			expectedBody:   expected,
		},
		{
			testName:       "[GET] Day's History",
			url:            "/history/day",
			body:           `{"url": "http://localhost:8080/yxZ8byjhRui"}`,
			expectedStatus: http.StatusOK,
			expectedBody:   expected,
		},
		{
			testName:       "[GET] Non Existent History",
			url:            "/history/",
			body:           `{"url": "http://localhost:8080/nonExistent"}`,
			expectedStatus: http.StatusNotFound,
		},
	}

	withServer(t, tests, func(t *testing.T, test httpTest) {
		if err := test.run(); err != nil {
			t.Fatal(err)
		}
	})
}

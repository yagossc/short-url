package api_test

import (
	"net/http"
	"testing"
)

func TestShortener(t *testing.T) {
	tests := []httpTest{
		{
			testName:       "[POST] New short URL",
			url:            "/",
			body:           `{"url": "http://www.google.com"}`,
			expectedStatus: http.StatusCreated,
		},
		{
			testName:       "[POST] New short URL",
			url:            "/",
			expectedStatus: http.StatusBadRequest,
		},
	}

	withServer(t, tests, func(t *testing.T, test httpTest) {
		if err := test.run(); err != nil {
			t.Fatal(err)
		}
	})
}

func TestRedirect(t *testing.T) {
	tests := []httpTest{
		{
			testName:       "[GET] Existent URL",
			url:            "/yxZ8byjhRui",
			expectedStatus: http.StatusMovedPermanently,
		},
		{
			testName:       "[GET] NonExistent URL",
			url:            "/nonExistent",
			expectedStatus: http.StatusNotFound,
		},
	}
	withServer(t, tests, func(t *testing.T, test httpTest) {
		if err := test.run(); err != nil {
			t.Fatal(err)
		}
	})
}

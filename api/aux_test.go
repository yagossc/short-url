package api_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/yagossc/short-url/api"
	"github.com/yagossc/short-url/internal/dbtest"
	"github.com/yagossc/short-url/query"
)

type httpTest struct {
	testName string
	method   string
	url      string
	body     string

	expectedStatus int
	expectedBody   interface{}

	actualBody string

	server *api.Server
	db     *query.Executor
}

func withServer(t *testing.T, tests []httpTest, fn func(t *testing.T, test httpTest)) {
	dbtest.WithDB(func(db *query.Executor) {
		if err := dbtest.Mock(db, dbtest.GenDefaultLoad()); err != nil {
			panic(err)
		}

		e := echo.New()
		s := api.NewServer(db, e, "localhost")
		s.Routes()

		for i := range tests {
			test := tests[i]
			test.server = s
			test.db = db

			if test.testName == "" {
				test.testName = fmt.Sprintf("Unnamed #%d", i)
			}

			t.Run(test.testName, func(t *testing.T) {
				fn(t, test)
			})
		}
	})
}

func (tc *httpTest) run() error {
	return tc.runWith(tc.server)
}

func (tc *httpTest) runWith(s *api.Server) error {
	method := tc.method
	if method == "" {
		switch {
		case strings.Contains(tc.testName, "[POST]"):
			method = http.MethodPost
		default:
			method = http.MethodGet
		}
	}

	addr := tc.url
	var payload io.Reader
	if tc.body != "" {
		payload = strings.NewReader(tc.body)
	}

	req := httptest.NewRequest(method, addr, payload)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	s.ServeHTTP(rec, req)

	// check the http code
	if rec.Code != tc.expectedStatus {
		return fmt.Errorf("expected status code of %d, but received %d\nBody: %s", tc.expectedStatus, rec.Code, rec.Body.String())
	}

	// extract the body	content
	var expected string

	// ignore nil pointer
	v := reflect.ValueOf(tc.expectedBody)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}

		v = reflect.Indirect(v)
	}

	// ignore invalid value
	if !v.IsValid() {
		return nil
	}

	// if it is a struct, marshal it
	if v.Kind() == reflect.Struct {
		switch obj := v.Interface().(type) {
		case json.Marshaler:
			b, err := obj.MarshalJSON()
			if err != nil {
				return err
			}
			expected = string(b)

		default:
			b, err := json.Marshal(obj)
			if err != nil {
				return err
			}
			expected = string(b)
		}
	} else {
		switch obj := v.Interface().(type) {
		case string:
			expected = obj
		default:
			return fmt.Errorf("unknown type: %v", obj)
		}
	}

	// get the expected and the actual body
	expected = strings.TrimSpace(expected)
	body := strings.TrimSpace(rec.Body.String())

	tc.actualBody = body

	if body != expected {
		return fmt.Errorf("expected body to be '%s', but it was '%s'", expected, body)
	}

	return nil
}

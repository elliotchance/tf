package tf

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type (
	HTTPCheckFunc   func(t *testing.T, request *http.Request, response *httptest.ResponseRecorder) bool
	HTTPFinallyFunc func(request *http.Request, response *httptest.ResponseRecorder)
	HTTPBeforeFunc  func(request *http.Request, response *httptest.ResponseRecorder)

	HTTPTester interface {
		TestName() string
		Tests() []*HTTPTest
	}

	MultiHTTPTest struct {
		// Name is used as the test name. If it is empty the test name will be based
		// on the Path.
		Name string

		// Before is run before any of the Steps begin.
		Before func()

		Steps []*HTTPTest
	}

	HTTPTest struct {
		// Name is used as the test name. If it is empty the test name will be based
		// on the Path.
		Name string

		// Method is the HTTP request method. If blank then "GET" will be used.
		Method string

		// Path used in the request. If the Path is blank then "/" is used because
		// it is not possible to parse an empty path.
		Path string

		// RequestBody is the body for the request. You use a string as the body
		// with:
		//
		//   RequestBody: strings.NewReader("foo bar")
		//
		RequestBody io.Reader

		// RequestHeaders will add or replace any header on the request.
		RequestHeaders map[string]string

		// ResponseHeaders will be checked from the response. Only the headers in
		// ResponseHeaders will be checked and and their values must be exactly
		// equal.
		//
		// If you need to do more sophisticated checking or headers you should use
		// Check.
		ResponseHeaders map[string]string

		// ResponseBody will check the body of the response. ResponseBody must be
		// not nil for the check to occur.
		//
		// You can check a string with:
		//
		//   ResponseBody: strings.NewReader("foo bar")
		//
		ResponseBody io.Reader

		// Status is the expected response HTTP status code. You can use one of the
		// constants in the http package such as http.StatusOK. If Status is not
		// provided then the response status will not be checked.
		Status int

		// Check is an optional function that is run before any other assertions. It
		// receives the request and response so you can do any custom validation. If
		// Check returns true the built in assertions will continue. Otherwise a
		// return value of false means to stop checking the response because an
		// error has already been found.
		Check HTTPCheckFunc

		// Finally is always called as a last event, even if the test fails. It is
		// useful for guaranteeing cleanup or restoration of environments.
		//
		// The return value is ignored.
		Finally HTTPFinallyFunc

		// Before is run after the request and record is setup but before the
		// request is executed.
		Before HTTPBeforeFunc
	}
)

func safeTestName(s string) string {
	return strings.Replace(s, ":", "", -1)
}

func (ht *HTTPTest) TestName() string {
	if ht.Name != "" {
		return ht.Name
	}

	if ht.Method == "" {
		return "GET " + ht.RealPath()
	}

	return ht.Method + " " + ht.RealPath()
}

func (ht *HTTPTest) Tests() []*HTTPTest {
	return []*HTTPTest{ht}
}

func (ht *HTTPTest) RealPath() string {
	if ht.Path == "" {
		return "/"
	}

	return ht.Path
}

func (ht *MultiHTTPTest) Tests() []*HTTPTest {
	return ht.Steps
}

func (ht *MultiHTTPTest) TestName() string {
	return ht.Name
}

func testSingle(t *testing.T, test *HTTPTest, handlerFunc http.HandlerFunc) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(test.Method, test.RealPath(), test.RequestBody)

	defer func() {
		if test.Finally != nil {
			test.Finally(request, recorder)
		}
	}()

	for k, v := range test.RequestHeaders {
		request.Header.Set(k, v)
	}

	if test.Before != nil {
		test.Before(request, recorder)
	}

	handlerFunc(recorder, request)

	// Check is run before any other assertions and can stop the the
	// test from proceeding if it returns false.
	if test.Check != nil && !test.Check(t, request, recorder) {
		return
	}

	if test.Status != 0 && !assert.Equal(t, test.Status, recorder.Code) {
		return
	}

	for k, v := range test.ResponseHeaders {
		if !assert.Equalf(t, v, recorder.HeaderMap.Get(k), "ResponseHeader[%s]", k) {
			return
		}
	}

	if test.ResponseBody != nil {
		expectedBody, err := ioutil.ReadAll(test.ResponseBody)
		if !assert.NoError(t, err) {
			return
		}

		if !assert.Equal(t, string(expectedBody), recorder.Body.String()) {
			return
		}
	}
}

func ServeHTTP(t *testing.T, handlerFunc http.HandlerFunc) func(HTTPTester) {
	return func(tests HTTPTester) {
		t.Run(safeTestName(tests.TestName()), func(t *testing.T) {
			if multiTest, ok := tests.(*MultiHTTPTest); ok && multiTest.Before != nil {
				multiTest.Before()
			}

			if len(tests.Tests()) > 1 {
				for _, test := range tests.Tests() {
					t.Run(safeTestName(test.TestName()), func(t *testing.T) {
						testSingle(t, test, handlerFunc)
					})
				}
			} else {
				testSingle(t, tests.Tests()[0], handlerFunc)
			}
		})
	}
}

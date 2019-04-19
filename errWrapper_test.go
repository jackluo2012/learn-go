package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func errPanic(writer http.ResponseWriter, request *http.Request) error {
	panic(124)
}

func errUserError(writer http.ResponseWriter, request *http.Request) error {
	return testingUserError("user error	")
}

func errNotFoundError(writer http.ResponseWriter, request *http.Request) error {
	return os.ErrNotExist
}
func errForbiddenError(writer http.ResponseWriter, request *http.Request) error {
	return os.ErrPermission
}
func errNoError(writer http.ResponseWriter, request *http.Request) error {

	fmt.Fprintf(writer, "no error")

	return nil
}

type testingUserError string

func (e testingUserError) Error() string {
	return e.Error()
}

func (e testingUserError) Message() string {
	return string(e)
}

var tests = []struct {
	h       appHandler
	code    int
	message string
}{
	{errPanic, 500, "Internal Server Error"},
	{errUserError, 400, "user error"},
	{errNotFoundError, 404, "user error"},
	{errForbiddenError, 403, "user error"},
	{errNoError, 200, "no error"},
}
/**
 * 假参数去测试函数
 */
func TestErrWrapper(t *testing.T) {

	for _, tt := range tests {
		f := errWrapper(tt.h)
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "http://www.imooc.com", nil)
		f(response, request)
		VerifyResponse(response.Result(), tt.code, tt.message, t)
	}
}

/**
 * 真运行去测试
 */
func TestErrWraperInServer(t *testing.T) {
	for _, tt := range tests {
		f := errWrapper(tt.h)
		server := httptest.NewServer(http.HandlerFunc(f))
		response, _ := http.Get(server.URL)
		VerifyResponse(response, tt.code, tt.message, t)
	}
}

func VerifyResponse(response *http.Response, code int, message string, t *testing.T) {
	b, _ := ioutil.ReadAll(response.Body)
	body := strings.Trim(string(b), "\n")
	if response.StatusCode != code || body != message {
		t.Errorf("expect (%d, %s, %s) ", code, message, body)
	}
}

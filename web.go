package main

import (
	"gopcp.v2/chapter7/filelisting"
	"net/http"
	_ "net/http/pprof" //访问 debug/pprof/
	"os"
)

type appHandler func(writer http.ResponseWriter, request *http.Request) error

func errWrapper(handler appHandler) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		//panic
		defer func() {
			if r := recover(); r != nil {
				//fmt.Printf("Pnaic:%v", r)
				http.Error(writer,
					http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
		}()

		err := handler(writer, request)

		code := http.StatusOK

		if err != nil {
			//fmt.Printf("Error occurred Handling request:%s", err.Error())
			//处理 userError
			if err, ok := err.(userError); ok {
				code = http.StatusInternalServerError
				http.Error(writer, err.Message(), http.StatusBadRequest)
				return
			}
			//system error
			switch {
			case os.IsExist(err):
				code = http.StatusNotFound
			case os.IsPermission(err):
				code = http.StatusForbidden
			default:
				code = http.StatusInternalServerError
			}
			http.Error(writer, http.StatusText(code), code)
		}
	}
}

type userError interface {
	error
	Message() string
}

func main() {
	http.HandleFunc("/", errWrapper(filelisting.HandlerFileListing))
	http.HandleFunc("/list/", errWrapper(filelisting.HandlerFileListing))
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		panic(err)
	}

}

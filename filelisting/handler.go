package filelisting

import (
	"io/ioutil"
	"net/http"
	"os"
)

type userError string

func (e userError) Error() string {
	return e.Error()
}

func (e userError) Message() string {
	return string(e)
}

func HandlerFileListing(w http.ResponseWriter, r *http.Request) error {
	path := r.URL.Path[len("/lists/")-1:]

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	all, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	w.Write(all)
	return nil
}

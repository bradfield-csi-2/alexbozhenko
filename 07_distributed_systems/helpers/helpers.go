package helpers

import (
	"fmt"
	"net/http"
)

func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func ErrorResponse(w *http.ResponseWriter, err error, code int) error {
	if err != nil {
		message := fmt.Sprintf("%s: %s", http.StatusText(code), err)
		(*w).WriteHeader(code)
		(*w).Write([]byte(message))
		return err
	}
	return nil
}

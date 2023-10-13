package utils

import (
	"errors"
	"io/ioutil"
	"net/http"
)

func ReadBody(request *http.Request) (output string, bodySize int, err error) {
	byteBody, err := ioutil.ReadAll(request.Body)
	defer request.Body.Close()
	if err != nil {
		return "", 0, errors.New("BODY_INVALID")
	}
	return string(byteBody), len(byteBody), nil
}

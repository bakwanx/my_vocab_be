package utils

import "net/http"

func ImageUrlProvider(imageName string, r *http.Request) string {
	return r.URL.RawPath + imageName
}

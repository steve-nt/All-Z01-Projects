package utils

import (
	"net/http"
	"strings"
)

// to do: Add error handling for badrequest and other statuses
func FindParamsFromURL(r *http.Request, whichparams string) (params []string) {

	urlparams := r.URL.Query().Get(whichparams)

	if urlparams != "" {
		params = strings.Split(urlparams, "_")
		return params
	}
	return nil
}

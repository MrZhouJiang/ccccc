package common

import "net/http"

func GetUrlParams(req *http.Request) map[string]string {
	req.ParseForm()

	var rs = map[string]string{}
	for k, v := range req.Form {
		rs[k] = v[0]
	}
	return rs
}

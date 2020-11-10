package utils

import (
	"io/ioutil"
	"net/http"
	"strings"
)

func GetUserInfo(url string, token string, isBearer bool) string {
	var (
		err  error
		req  *http.Request
		resp *http.Response
	)
	if isBearer {
		req, err = http.NewRequest("POST", url, nil)
		if err != nil {
			panic(err.Error())
		}
		req.Header.Set("Authorization", "Bearer "+token)
	} else {
		req, err = http.NewRequest("POST", url, strings.NewReader("access_token="+token))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		if err != nil {
			panic(err.Error())
		}
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)
	return string(b)
}

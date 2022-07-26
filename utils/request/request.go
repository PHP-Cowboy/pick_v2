package request

import (
	"io/ioutil"
	"net/http"
	"strings"
)

func Request(url, method string) ([]byte, error) {
	if method == "" {
		method = "post"
	}

	method = strings.ToUpper(method)

	client := &http.Client{}

	rq, err := http.NewRequest(method, url, nil)

	if err != nil {
		return nil, err
	}
	res, err := client.Do(rq)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

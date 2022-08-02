package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"pick_v2/global"
)

func Post(path string, responseData interface{}) ([]byte, error) {
	fmt.Println(responseData)
	cfg := global.ServerConfig

	url := fmt.Sprintf("%s:%d/%s", cfg.GoodsApi.Url, cfg.GoodsApi.Port, path)

	client := &http.Client{}

	jData, err := json.Marshal(responseData)
	if err != nil {
		return nil, err
	}

	rq, err := http.NewRequest("POST", url, bytes.NewReader(jData))

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

func Get(path string) ([]byte, error) {
	cfg := global.ServerConfig

	url := fmt.Sprintf("%s:%d/%s", cfg.GoodsApi.Url, cfg.GoodsApi.Port, path)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, err
}

func TestGet() ([]byte, error) {
	url := "http://121.196.60.92:19090/api/v1/remote/pick/shop/list"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
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

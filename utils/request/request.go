package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"pick_v2/global"
)

func Post(path string, responseData map[string]interface{}) ([]byte, error) {

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

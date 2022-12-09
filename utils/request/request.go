package request

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"

	"pick_v2/global"
	"pick_v2/middlewares"
)

type HttpRsp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func PostTest() {
	url := "http://121.196.60.92:19090/api/v1/remote/get/goods/by/id"
	method := "POST"

	payload := strings.NewReader(`{"order_id":[255]}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}

func Post(path string, responseData interface{}) ([]byte, error) {

	cfg := global.ServerConfig

	url := fmt.Sprintf("%s:%d/%s", cfg.GoodsApi.Url, cfg.GoodsApi.Port, path)

	client := &http.Client{}

	jData, err := json.Marshal(responseData)
	if err != nil {
		return nil, err
	}

	global.Logger["info"].Infof("url:%s", url)
	global.Logger["info"].Infof("params:%s", string(jData))

	rq, err := http.NewRequest("POST", url, bytes.NewReader(jData))

	if err != nil {
		return nil, err
	}

	sign := middlewares.Generate()

	rq.Header.Add("Content-Type", "application/json")
	rq.Header.Add("x-sign", sign)

	res, err := client.Do(rq)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	global.Logger["info"].Infof(string(body))

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

	body, err := io.ReadAll(res.Body)
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

func Call(uri string, params interface{}, res interface{}) (err error) {
	var (
		body []byte
	)

	body, err = Post(uri, params)

	if err != nil {
		return
	}

	err = json.Unmarshal(body, &res)

	if err != nil {
		return
	}

	rspCode, rspMsg := getRspCode(res)

	if rspCode != 200 {
		return errors.New(rspMsg)
	}

	return
}

func getRspCode(rsp interface{}) (code int, msg string) {
	code = -1
	msg = "未知错误"

	if rsp == nil {
		return
	}
	v := reflect.ValueOf(rsp)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		//最多取两层
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
	}
	//kind := v.Kind()
	// 判断是否是结构体
	if v.Kind() != reflect.Struct {
		return
	}
	codeValue := v.FieldByName("Code")
	if !codeValue.IsValid() {
		return
	}

	msgValue := v.FieldByName("Msg")

	if !codeValue.IsValid() {
		return
	}

	code = int(codeValue.Int())

	msg = msgValue.String()

	return
}

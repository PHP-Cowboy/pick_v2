package request

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"pick_v2/global"
	"pick_v2/middlewares"
	"reflect"
)

type HttpRsp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Post(path string, responseData interface{}) ([]byte, error) {

	cfg := global.ServerConfig

	url := fmt.Sprintf("%s:%d/%s", cfg.GoodsApi.Url, cfg.GoodsApi.Port, path)

	client := &http.Client{}

	jData, err := json.Marshal(responseData)
	if err != nil {
		return nil, err
	}

	rq, err := http.NewRequest("POST", url, bytes.NewReader(jData))

	if err != nil {
		global.Logger["err"].Infof("url:%s,params:%s,err:%s", url, string(jData), err.Error())
		return nil, err
	}

	sign := middlewares.Generate()

	rq.Header.Add("Content-Type", "application/json")
	rq.Header.Add("x-sign", sign)

	res, err := client.Do(rq)

	if err != nil {
		global.Logger["err"].Infof("url:%s,params:%s,err:%s", url, string(jData), err.Error())
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		global.Logger["err"].Infof("url:%s,params:%s,err:%s", url, string(jData), err.Error())
		return nil, err
	}

	global.Logger["info"].Infof("url:%s,params:%s,body:%s", url, string(jData), string(body))

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

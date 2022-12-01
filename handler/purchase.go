package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"pick_v2/dao"
	"pick_v2/forms/req"
	"pick_v2/forms/rsp"
	"pick_v2/global"
	"pick_v2/utils/request"
)

func Order(ctx context.Context, messages ...*primitive.MessageExt) (consumer.ConsumeResult, error) {

	var (
		form req.PurchaseOrderForm
		Id   int
	)

	for i := range messages {
		err := json.Unmarshal(messages[i].Body, &Id)
		if err != nil {
			global.Logger["err"].Infof("解析json失败:%s", err.Error())
			return consumer.ConsumeRetryLater, nil
		}
		form.OrderId = append(form.OrderId, Id)
	}

	orderRsp, err := GetOrderInfo(form)

	if err != nil {
		global.Logger["err"].Infof(err.Error())
		return consumer.ConsumeRetryLater, err
	}

	db := global.DB

	for _, info := range orderRsp.Data {

		if info.DistributionType == 6 { //无需出库
			return dao.NoShipping(db, form, info)
		} else {
			return dao.Shipping(db, form, info)
		}
	}

	return consumer.ConsumeRetryLater, errors.New("异常")
}

func GetOrderInfo(responseData interface{}) (rsp.OrderRsp, error) {
	var result rsp.OrderRsp

	body, err := request.Post("api/v1/remote/get/goods/by/id", responseData)

	if err != nil {
		return result, err
	}

	err = json.Unmarshal(body, &result)

	if err != nil {
		return result, err
	}

	if result.Code != 200 {
		return result, errors.New(result.Msg)
	}

	return result, nil
}

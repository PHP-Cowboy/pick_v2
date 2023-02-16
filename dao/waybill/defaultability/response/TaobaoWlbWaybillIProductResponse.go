package response

import (
	"topsdk/defaultability/domain"
)

type TaobaoWlbWaybillIProductResponse struct {

	/*
	   System request id
	*/
	RequestId string `json:"request_id,omitempty" `

	/*
	   System body
	*/
	Body string

	/*
	   产品类型返回
	*/
	ProductTypes []domain.TaobaoWlbWaybillIProductWaybillProductType `json:"product_types,omitempty" `
}

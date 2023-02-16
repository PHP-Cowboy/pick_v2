package response

import (
	"topsdk/defaultability/domain"
)

type TaobaoWlbWaybillIQuerydetailResponse struct {

	/*
	   System request id
	*/
	RequestId string `json:"request_id,omitempty" `

	/*
	   System body
	*/
	Body string

	/*
	   面单查询错误编码
	*/
	ErrorCodes []string `json:"error_codes,omitempty" `
	/*
	   不存在的面单号
	*/
	InexistentWaybillCodes []string `json:"inexistent_waybill_codes,omitempty" `
	/*
	   查询是否成功
	*/
	QuerySuccess bool `json:"query_success,omitempty" `
	/*
	   面单详情
	*/
	WaybillDetails []domain.TaobaoWlbWaybillIQuerydetailWaybillDetailQueryInfo `json:"waybill_details,omitempty" `
}

package response

import (
	"topsdk/defaultability/domain"
)

type TaobaoWlbWaybillIGetResponse struct {

	/*
	   System request id
	*/
	RequestId string `json:"request_id,omitempty" `

	/*
	   System body
	*/
	Body string

	/*
	   面单申请接口返回信息
	*/
	WaybillApplyNewCols []domain.TaobaoWlbWaybillIGetWaybillApplyNewInfo `json:"waybill_apply_new_cols,omitempty" `
}

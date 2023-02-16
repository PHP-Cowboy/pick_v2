package response

import (
	"topsdk/defaultability/domain"
)

type TaobaoWlbWaybillIFullupdateResponse struct {

	/*
	   System request id
	*/
	RequestId string `json:"request_id,omitempty" `

	/*
	   System body
	*/
	Body string

	/*
	   更新接口出参
	*/
	WaybillApplyUpdateInfo domain.TaobaoWlbWaybillIFullupdateWaybillApplyUpdateInfo `json:"waybill_apply_update_info,omitempty" `
}

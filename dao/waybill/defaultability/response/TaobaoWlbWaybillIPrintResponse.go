package response

import (
	"topsdk/defaultability/domain"
)

type TaobaoWlbWaybillIPrintResponse struct {

	/*
	   System request id
	*/
	RequestId string `json:"request_id,omitempty" `

	/*
	   System body
	*/
	Body string

	/*
	   面单打印信息
	*/
	WaybillApplyPrintCheckInfos []domain.TaobaoWlbWaybillIPrintWaybillApplyPrintCheckInfo `json:"waybill_apply_print_check_infos,omitempty" `
}

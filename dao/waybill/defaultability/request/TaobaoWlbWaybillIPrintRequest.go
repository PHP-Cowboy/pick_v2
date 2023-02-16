package request

import (
	"pick_v2/dao/waybill/defaultability/domain"
	"pick_v2/dao/waybill/util"
)

type TaobaoWlbWaybillIPrintRequest struct {
	/*
	   打印请求     */
	WaybillApplyPrintCheckRequest *domain.TaobaoWlbWaybillIPrintWaybillApplyPrintCheckRequest `json:"waybill_apply_print_check_request" required:"true" `
}

func (s *TaobaoWlbWaybillIPrintRequest) SetWaybillApplyPrintCheckRequest(v domain.TaobaoWlbWaybillIPrintWaybillApplyPrintCheckRequest) *TaobaoWlbWaybillIPrintRequest {
	s.WaybillApplyPrintCheckRequest = &v
	return s
}

func (req *TaobaoWlbWaybillIPrintRequest) ToMap() map[string]interface{} {
	paramMap := make(map[string]interface{})
	if req.WaybillApplyPrintCheckRequest != nil {
		paramMap["waybill_apply_print_check_request"] = util.ConvertStruct(*req.WaybillApplyPrintCheckRequest)
	}
	return paramMap
}

func (req *TaobaoWlbWaybillIPrintRequest) ToFileMap() map[string]interface{} {
	fileMap := make(map[string]interface{})
	return fileMap
}

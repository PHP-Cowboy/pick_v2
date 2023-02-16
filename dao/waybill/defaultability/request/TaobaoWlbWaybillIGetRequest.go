package request

import (
	"pick_v2/dao/waybill/defaultability/domain"
	"pick_v2/dao/waybill/util"
)

type TaobaoWlbWaybillIGetRequest struct {
	/*
	   面单申请     */
	WaybillApplyNewRequest *domain.TaobaoWlbWaybillIGetWaybillApplyNewRequest `json:"waybill_apply_new_request" required:"true" `
}

func (s *TaobaoWlbWaybillIGetRequest) SetWaybillApplyNewRequest(v domain.TaobaoWlbWaybillIGetWaybillApplyNewRequest) *TaobaoWlbWaybillIGetRequest {
	s.WaybillApplyNewRequest = &v
	return s
}

func (req *TaobaoWlbWaybillIGetRequest) ToMap() map[string]interface{} {
	paramMap := make(map[string]interface{})
	if req.WaybillApplyNewRequest != nil {
		paramMap["waybill_apply_new_request"] = util.ConvertStruct(*req.WaybillApplyNewRequest)
	}
	return paramMap
}

func (req *TaobaoWlbWaybillIGetRequest) ToFileMap() map[string]interface{} {
	fileMap := make(map[string]interface{})
	return fileMap
}

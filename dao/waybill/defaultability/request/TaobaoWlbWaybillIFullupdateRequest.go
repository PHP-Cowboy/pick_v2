package request

import (
	"pick_v2/dao/waybill/defaultability/domain"
	"pick_v2/dao/waybill/util"
)

type TaobaoWlbWaybillIFullupdateRequest struct {
	/*
	   更新面单信息请求     */
	WaybillApplyFullUpdateRequest *domain.TaobaoWlbWaybillIFullupdateWaybillApplyFullUpdateRequest `json:"waybill_apply_full_update_request" required:"true" `
}

func (s *TaobaoWlbWaybillIFullupdateRequest) SetWaybillApplyFullUpdateRequest(v domain.TaobaoWlbWaybillIFullupdateWaybillApplyFullUpdateRequest) *TaobaoWlbWaybillIFullupdateRequest {
	s.WaybillApplyFullUpdateRequest = &v
	return s
}

func (req *TaobaoWlbWaybillIFullupdateRequest) ToMap() map[string]interface{} {
	paramMap := make(map[string]interface{})
	if req.WaybillApplyFullUpdateRequest != nil {
		paramMap["waybill_apply_full_update_request"] = util.ConvertStruct(*req.WaybillApplyFullUpdateRequest)
	}
	return paramMap
}

func (req *TaobaoWlbWaybillIFullupdateRequest) ToFileMap() map[string]interface{} {
	fileMap := make(map[string]interface{})
	return fileMap
}

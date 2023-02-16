package request

import (
	"pick_v2/dao/waybill/defaultability/domain"
	"pick_v2/dao/waybill/util"
)

type TaobaoWlbWaybillICancelRequest struct {
	/*
	   取消接口入参     */
	WaybillApplyCancelRequest *domain.TaobaoWlbWaybillICancelWaybillApplyCancelRequest `json:"waybill_apply_cancel_request" required:"true" `
}

func (s *TaobaoWlbWaybillICancelRequest) SetWaybillApplyCancelRequest(v domain.TaobaoWlbWaybillICancelWaybillApplyCancelRequest) *TaobaoWlbWaybillICancelRequest {
	s.WaybillApplyCancelRequest = &v
	return s
}

func (req *TaobaoWlbWaybillICancelRequest) ToMap() map[string]interface{} {
	paramMap := make(map[string]interface{})
	if req.WaybillApplyCancelRequest != nil {
		paramMap["waybill_apply_cancel_request"] = util.ConvertStruct(*req.WaybillApplyCancelRequest)
	}
	return paramMap
}

func (req *TaobaoWlbWaybillICancelRequest) ToFileMap() map[string]interface{} {
	fileMap := make(map[string]interface{})
	return fileMap
}

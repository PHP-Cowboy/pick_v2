package request

import (
	"pick_v2/dao/waybill/defaultability/domain"
	"pick_v2/dao/waybill/util"
)

type TaobaoWlbWaybillIQuerydetailRequest struct {
	/*
	   面单查询请求     */
	WaybillDetailQueryRequest *domain.TaobaoWlbWaybillIQuerydetailWaybillDetailQueryRequest `json:"waybill_detail_query_request" required:"true" `
}

func (s *TaobaoWlbWaybillIQuerydetailRequest) SetWaybillDetailQueryRequest(v domain.TaobaoWlbWaybillIQuerydetailWaybillDetailQueryRequest) *TaobaoWlbWaybillIQuerydetailRequest {
	s.WaybillDetailQueryRequest = &v
	return s
}

func (req *TaobaoWlbWaybillIQuerydetailRequest) ToMap() map[string]interface{} {
	paramMap := make(map[string]interface{})
	if req.WaybillDetailQueryRequest != nil {
		paramMap["waybill_detail_query_request"] = util.ConvertStruct(*req.WaybillDetailQueryRequest)
	}
	return paramMap
}

func (req *TaobaoWlbWaybillIQuerydetailRequest) ToFileMap() map[string]interface{} {
	fileMap := make(map[string]interface{})
	return fileMap
}

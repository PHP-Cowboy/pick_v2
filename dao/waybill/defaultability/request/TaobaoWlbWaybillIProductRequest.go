package request

import (
	"pick_v2/dao/waybill/defaultability/domain"
	"pick_v2/dao/waybill/util"
)

type TaobaoWlbWaybillIProductRequest struct {
	/*
	   查询物流商电子面单产品类型入参     */
	WaybillProductTypeRequest *domain.TaobaoWlbWaybillIProductWaybillProductTypeRequest `json:"waybill_product_type_request" required:"true" `
}

func (s *TaobaoWlbWaybillIProductRequest) SetWaybillProductTypeRequest(v domain.TaobaoWlbWaybillIProductWaybillProductTypeRequest) *TaobaoWlbWaybillIProductRequest {
	s.WaybillProductTypeRequest = &v
	return s
}

func (req *TaobaoWlbWaybillIProductRequest) ToMap() map[string]interface{} {
	paramMap := make(map[string]interface{})
	if req.WaybillProductTypeRequest != nil {
		paramMap["waybill_product_type_request"] = util.ConvertStruct(*req.WaybillProductTypeRequest)
	}
	return paramMap
}

func (req *TaobaoWlbWaybillIProductRequest) ToFileMap() map[string]interface{} {
	fileMap := make(map[string]interface{})
	return fileMap
}

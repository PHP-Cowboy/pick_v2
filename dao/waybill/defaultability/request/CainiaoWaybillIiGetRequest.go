package request

import (
	"pick_v2/dao/waybill/defaultability/domain"
	"pick_v2/dao/waybill/util"
)

type CainiaoWaybillIiGetRequest struct {
	/*
	   入参信息     */
	ParamWaybillCloudPrintApplyNewRequest *domain.CainiaoWaybillIiGetWaybillCloudPrintApplyNewRequest `json:"param_waybill_cloud_print_apply_new_request" required:"true" `
}

func (s *CainiaoWaybillIiGetRequest) SetParamWaybillCloudPrintApplyNewRequest(v domain.CainiaoWaybillIiGetWaybillCloudPrintApplyNewRequest) *CainiaoWaybillIiGetRequest {
	s.ParamWaybillCloudPrintApplyNewRequest = &v
	return s
}

func (req *CainiaoWaybillIiGetRequest) ToMap() map[string]interface{} {
	paramMap := make(map[string]interface{})
	if req.ParamWaybillCloudPrintApplyNewRequest != nil {
		paramMap["param_waybill_cloud_print_apply_new_request"] = util.ConvertStruct(*req.ParamWaybillCloudPrintApplyNewRequest)
	}
	return paramMap
}

func (req *CainiaoWaybillIiGetRequest) ToFileMap() map[string]interface{} {
	fileMap := make(map[string]interface{})
	return fileMap
}

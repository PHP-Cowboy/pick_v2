package request

import (
	"pick_v2/dao/waybill/defaultability/domain"
	"pick_v2/dao/waybill/util"
)

type CainiaoWaybillIiUpdateRequest struct {
	/*
	   更新请求信息     */
	ParamWaybillCloudPrintUpdateRequest *domain.CainiaoWaybillIiUpdateWaybillCloudPrintUpdateRequest `json:"param_waybill_cloud_print_update_request" required:"true" `
}

func (s *CainiaoWaybillIiUpdateRequest) SetParamWaybillCloudPrintUpdateRequest(v domain.CainiaoWaybillIiUpdateWaybillCloudPrintUpdateRequest) *CainiaoWaybillIiUpdateRequest {
	s.ParamWaybillCloudPrintUpdateRequest = &v
	return s
}

func (req *CainiaoWaybillIiUpdateRequest) ToMap() map[string]interface{} {
	paramMap := make(map[string]interface{})
	if req.ParamWaybillCloudPrintUpdateRequest != nil {
		paramMap["param_waybill_cloud_print_update_request"] = util.ConvertStruct(*req.ParamWaybillCloudPrintUpdateRequest)
	}
	return paramMap
}

func (req *CainiaoWaybillIiUpdateRequest) ToFileMap() map[string]interface{} {
	fileMap := make(map[string]interface{})
	return fileMap
}

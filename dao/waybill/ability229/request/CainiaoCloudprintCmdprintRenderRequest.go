package request

import (
	"pick_v2/dao/waybill/ability229/domain"
	"pick_v2/dao/waybill/util"
)

type CainiaoCloudprintCmdprintRenderRequest struct {
	/*
	   参数对象     */
	Params *domain.CainiaoCloudprintCmdprintRenderCmdRenderParams `json:"params" required:"true" `
}

func (s *CainiaoCloudprintCmdprintRenderRequest) SetParams(v domain.CainiaoCloudprintCmdprintRenderCmdRenderParams) *CainiaoCloudprintCmdprintRenderRequest {
	s.Params = &v
	return s
}

func (req *CainiaoCloudprintCmdprintRenderRequest) ToMap() map[string]interface{} {
	paramMap := make(map[string]interface{})
	if req.Params != nil {
		paramMap["params"] = util.ConvertStruct(*req.Params)
	}
	return paramMap
}

func (req *CainiaoCloudprintCmdprintRenderRequest) ToFileMap() map[string]interface{} {
	fileMap := make(map[string]interface{})
	return fileMap
}

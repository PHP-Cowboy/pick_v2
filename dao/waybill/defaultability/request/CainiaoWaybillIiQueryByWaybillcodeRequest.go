package request

import (
	"pick_v2/dao/waybill/defaultability/domain"
	"pick_v2/dao/waybill/util"
)

type CainiaoWaybillIiQueryByWaybillcodeRequest struct {
	/*
	   系统自动生成     */
	ParamList *[]domain.CainiaoWaybillIiQueryByWaybillcodeWaybillDetailQueryByWaybillCodeRequest `json:"param_list,omitempty" required:"false" `
}

func (s *CainiaoWaybillIiQueryByWaybillcodeRequest) SetParamList(v []domain.CainiaoWaybillIiQueryByWaybillcodeWaybillDetailQueryByWaybillCodeRequest) *CainiaoWaybillIiQueryByWaybillcodeRequest {
	s.ParamList = &v
	return s
}

func (req *CainiaoWaybillIiQueryByWaybillcodeRequest) ToMap() map[string]interface{} {
	paramMap := make(map[string]interface{})
	if req.ParamList != nil {
		paramMap["param_list"] = util.ConvertStructList(*req.ParamList)
	}
	return paramMap
}

func (req *CainiaoWaybillIiQueryByWaybillcodeRequest) ToFileMap() map[string]interface{} {
	fileMap := make(map[string]interface{})
	return fileMap
}

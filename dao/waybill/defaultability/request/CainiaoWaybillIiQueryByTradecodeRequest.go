package request

import (
	"topsdk/defaultability/domain"
	"topsdk/util"
)

type CainiaoWaybillIiQueryByTradecodeRequest struct {
	/*
	   订单号列表     */
	ParamList *[]domain.CainiaoWaybillIiQueryByTradecodeWaybillDetailQueryByBizSubCodeRequest `json:"param_list,omitempty" required:"false" `
}

func (s *CainiaoWaybillIiQueryByTradecodeRequest) SetParamList(v []domain.CainiaoWaybillIiQueryByTradecodeWaybillDetailQueryByBizSubCodeRequest) *CainiaoWaybillIiQueryByTradecodeRequest {
	s.ParamList = &v
	return s
}

func (req *CainiaoWaybillIiQueryByTradecodeRequest) ToMap() map[string]interface{} {
	paramMap := make(map[string]interface{})
	if req.ParamList != nil {
		paramMap["param_list"] = util.ConvertStructList(*req.ParamList)
	}
	return paramMap
}

func (req *CainiaoWaybillIiQueryByTradecodeRequest) ToFileMap() map[string]interface{} {
	fileMap := make(map[string]interface{})
	return fileMap
}

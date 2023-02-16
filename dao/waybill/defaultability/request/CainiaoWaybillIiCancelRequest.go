package request

type CainiaoWaybillIiCancelRequest struct {
	/*
	   快递公司code     */
	CpCode *string `json:"cp_code" required:"true" `
	/*
	   电子面单号     */
	WaybillCode *string `json:"waybill_code" required:"true" `
}

func (s *CainiaoWaybillIiCancelRequest) SetCpCode(v string) *CainiaoWaybillIiCancelRequest {
	s.CpCode = &v
	return s
}
func (s *CainiaoWaybillIiCancelRequest) SetWaybillCode(v string) *CainiaoWaybillIiCancelRequest {
	s.WaybillCode = &v
	return s
}

func (req *CainiaoWaybillIiCancelRequest) ToMap() map[string]interface{} {
	paramMap := make(map[string]interface{})
	if req.CpCode != nil {
		paramMap["cp_code"] = *req.CpCode
	}
	if req.WaybillCode != nil {
		paramMap["waybill_code"] = *req.WaybillCode
	}
	return paramMap
}

func (req *CainiaoWaybillIiCancelRequest) ToFileMap() map[string]interface{} {
	fileMap := make(map[string]interface{})
	return fileMap
}

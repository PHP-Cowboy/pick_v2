package request

type CainiaoWaybillIiSearchRequest struct {
	/*
	   物流公司code     */
	CpCode *string `json:"cp_code,omitempty" required:"false" `
}

func (s *CainiaoWaybillIiSearchRequest) SetCpCode(v string) *CainiaoWaybillIiSearchRequest {
	s.CpCode = &v
	return s
}

func (req *CainiaoWaybillIiSearchRequest) ToMap() map[string]interface{} {
	paramMap := make(map[string]interface{})
	if req.CpCode != nil {
		paramMap["cp_code"] = *req.CpCode
	}
	return paramMap
}

func (req *CainiaoWaybillIiSearchRequest) ToFileMap() map[string]interface{} {
	fileMap := make(map[string]interface{})
	return fileMap
}

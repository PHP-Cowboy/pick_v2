package request

type CainiaoWaybillIiProductRequest struct {
	/*
	   快递公司code     */
	CpCode *string `json:"cp_code" required:"true" `
}

func (s *CainiaoWaybillIiProductRequest) SetCpCode(v string) *CainiaoWaybillIiProductRequest {
	s.CpCode = &v
	return s
}

func (req *CainiaoWaybillIiProductRequest) ToMap() map[string]interface{} {
	paramMap := make(map[string]interface{})
	if req.CpCode != nil {
		paramMap["cp_code"] = *req.CpCode
	}
	return paramMap
}

func (req *CainiaoWaybillIiProductRequest) ToFileMap() map[string]interface{} {
	fileMap := make(map[string]interface{})
	return fileMap
}

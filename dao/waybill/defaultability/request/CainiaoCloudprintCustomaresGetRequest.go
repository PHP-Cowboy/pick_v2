package request

type CainiaoCloudprintCustomaresGetRequest struct {
	/*
	   用户使用的标准模板id     */
	TemplateId *int64 `json:"template_id" required:"true" `
}

func (s *CainiaoCloudprintCustomaresGetRequest) SetTemplateId(v int64) *CainiaoCloudprintCustomaresGetRequest {
	s.TemplateId = &v
	return s
}

func (req *CainiaoCloudprintCustomaresGetRequest) ToMap() map[string]interface{} {
	paramMap := make(map[string]interface{})
	if req.TemplateId != nil {
		paramMap["template_id"] = *req.TemplateId
	}
	return paramMap
}

func (req *CainiaoCloudprintCustomaresGetRequest) ToFileMap() map[string]interface{} {
	fileMap := make(map[string]interface{})
	return fileMap
}

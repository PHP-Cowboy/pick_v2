package request

type CainiaoCloudprintClientinfoPutRequest struct {
	/*
	   客户端上传json数据     */
	JsonData *string `json:"json_data" required:"true" `
}

func (s *CainiaoCloudprintClientinfoPutRequest) SetJsonData(v string) *CainiaoCloudprintClientinfoPutRequest {
	s.JsonData = &v
	return s
}

func (req *CainiaoCloudprintClientinfoPutRequest) ToMap() map[string]interface{} {
	paramMap := make(map[string]interface{})
	if req.JsonData != nil {
		paramMap["json_data"] = *req.JsonData
	}
	return paramMap
}

func (req *CainiaoCloudprintClientinfoPutRequest) ToFileMap() map[string]interface{} {
	fileMap := make(map[string]interface{})
	return fileMap
}

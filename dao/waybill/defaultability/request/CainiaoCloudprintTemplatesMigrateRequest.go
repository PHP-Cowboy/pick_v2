package request

type CainiaoCloudprintTemplatesMigrateRequest struct {
	/*
	   标准电子面单模板的id     */
	TempalteId *int64 `json:"tempalte_id,omitempty" required:"false" `
	/*
	   自定义区名称     */
	CustomAreaName *string `json:"custom_area_name,omitempty" required:"false" `
	/*
	   自定义区内容     */
	CustomAreaContent *string `json:"custom_area_content,omitempty" required:"false" `
}

func (s *CainiaoCloudprintTemplatesMigrateRequest) SetTempalteId(v int64) *CainiaoCloudprintTemplatesMigrateRequest {
	s.TempalteId = &v
	return s
}
func (s *CainiaoCloudprintTemplatesMigrateRequest) SetCustomAreaName(v string) *CainiaoCloudprintTemplatesMigrateRequest {
	s.CustomAreaName = &v
	return s
}
func (s *CainiaoCloudprintTemplatesMigrateRequest) SetCustomAreaContent(v string) *CainiaoCloudprintTemplatesMigrateRequest {
	s.CustomAreaContent = &v
	return s
}

func (req *CainiaoCloudprintTemplatesMigrateRequest) ToMap() map[string]interface{} {
	paramMap := make(map[string]interface{})
	if req.TempalteId != nil {
		paramMap["tempalte_id"] = *req.TempalteId
	}
	if req.CustomAreaName != nil {
		paramMap["custom_area_name"] = *req.CustomAreaName
	}
	if req.CustomAreaContent != nil {
		paramMap["custom_area_content"] = *req.CustomAreaContent
	}
	return paramMap
}

func (req *CainiaoCloudprintTemplatesMigrateRequest) ToFileMap() map[string]interface{} {
	fileMap := make(map[string]interface{})
	return fileMap
}

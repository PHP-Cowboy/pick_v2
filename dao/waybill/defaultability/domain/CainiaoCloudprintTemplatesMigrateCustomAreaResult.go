package domain

type CainiaoCloudprintTemplatesMigrateCustomAreaResult struct {
	/*
	   自定义区id     */
	CustomAreaId *int64 `json:"custom_area_id,omitempty" `

	/*
	   自定义区名称     */
	CustomAreaName *string `json:"custom_area_name,omitempty" `

	/*
	   自定义区url     */
	CustomAreaUrl *string `json:"custom_area_url,omitempty" `

	/*
	   变量     */
	Keys *[]CainiaoCloudprintTemplatesMigrateKeyResult `json:"keys,omitempty" `

	/*
	   标准模板名称     */
	StandardTemplateId *int64 `json:"standard_template_id,omitempty" `

	/*
	   标准模板url     */
	StandardTemplateUrl *string `json:"standard_template_url,omitempty" `

	/*
	   用户模板id，等同于mystdtemplates.get中返回的用户模板id     */
	UserTemplateId *int64 `json:"user_template_id,omitempty" `
}

func (s *CainiaoCloudprintTemplatesMigrateCustomAreaResult) SetCustomAreaId(v int64) *CainiaoCloudprintTemplatesMigrateCustomAreaResult {
	s.CustomAreaId = &v
	return s
}
func (s *CainiaoCloudprintTemplatesMigrateCustomAreaResult) SetCustomAreaName(v string) *CainiaoCloudprintTemplatesMigrateCustomAreaResult {
	s.CustomAreaName = &v
	return s
}
func (s *CainiaoCloudprintTemplatesMigrateCustomAreaResult) SetCustomAreaUrl(v string) *CainiaoCloudprintTemplatesMigrateCustomAreaResult {
	s.CustomAreaUrl = &v
	return s
}
func (s *CainiaoCloudprintTemplatesMigrateCustomAreaResult) SetKeys(v []CainiaoCloudprintTemplatesMigrateKeyResult) *CainiaoCloudprintTemplatesMigrateCustomAreaResult {
	s.Keys = &v
	return s
}
func (s *CainiaoCloudprintTemplatesMigrateCustomAreaResult) SetStandardTemplateId(v int64) *CainiaoCloudprintTemplatesMigrateCustomAreaResult {
	s.StandardTemplateId = &v
	return s
}
func (s *CainiaoCloudprintTemplatesMigrateCustomAreaResult) SetStandardTemplateUrl(v string) *CainiaoCloudprintTemplatesMigrateCustomAreaResult {
	s.StandardTemplateUrl = &v
	return s
}
func (s *CainiaoCloudprintTemplatesMigrateCustomAreaResult) SetUserTemplateId(v int64) *CainiaoCloudprintTemplatesMigrateCustomAreaResult {
	s.UserTemplateId = &v
	return s
}

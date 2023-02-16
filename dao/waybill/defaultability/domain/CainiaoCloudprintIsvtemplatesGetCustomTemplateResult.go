package domain

type CainiaoCloudprintIsvtemplatesGetCustomTemplateResult struct {
	/*
	   isv模板的id     */
	IsvTemplateId *int64 `json:"isv_template_id,omitempty" `

	/*
	   isv模板的名称     */
	IsvTemplateName *string `json:"isv_template_name,omitempty" `

	/*
	   isv模板的url     */
	IsvTemplateUrl *string `json:"isv_template_url,omitempty" `

	/*
	   模板的keys     */
	Keys *[]CainiaoCloudprintIsvtemplatesGetKeyResult `json:"keys,omitempty" `

	/*
	   版本号     */
	Version *string `json:"version,omitempty" `
}

func (s *CainiaoCloudprintIsvtemplatesGetCustomTemplateResult) SetIsvTemplateId(v int64) *CainiaoCloudprintIsvtemplatesGetCustomTemplateResult {
	s.IsvTemplateId = &v
	return s
}
func (s *CainiaoCloudprintIsvtemplatesGetCustomTemplateResult) SetIsvTemplateName(v string) *CainiaoCloudprintIsvtemplatesGetCustomTemplateResult {
	s.IsvTemplateName = &v
	return s
}
func (s *CainiaoCloudprintIsvtemplatesGetCustomTemplateResult) SetIsvTemplateUrl(v string) *CainiaoCloudprintIsvtemplatesGetCustomTemplateResult {
	s.IsvTemplateUrl = &v
	return s
}
func (s *CainiaoCloudprintIsvtemplatesGetCustomTemplateResult) SetKeys(v []CainiaoCloudprintIsvtemplatesGetKeyResult) *CainiaoCloudprintIsvtemplatesGetCustomTemplateResult {
	s.Keys = &v
	return s
}
func (s *CainiaoCloudprintIsvtemplatesGetCustomTemplateResult) SetVersion(v string) *CainiaoCloudprintIsvtemplatesGetCustomTemplateResult {
	s.Version = &v
	return s
}

package domain

type CainiaoCloudprintStdtemplatesGetStandardTemplateDo struct {
	/*
	   模板id     */
	StandardTemplateId *int64 `json:"standard_template_id,omitempty" `

	/*
	   模板名称     */
	StandardTemplateName *string `json:"standard_template_name,omitempty" `

	/*
	   模板url     */
	StandardTemplateUrl *string `json:"standard_template_url,omitempty" `

	/*
	   1 快递标准面单 ,2 快递三联面单, 3 快递便携式三联单, 4 快运标准面单, 5 快运三联面单, 6 快递一联单     */
	StandardWaybillType *int64 `json:"standard_waybill_type,omitempty" `

	/*
	   如果没有 brandCode,则为 default     */
	BrandCode *string `json:"brand_code,omitempty" `
}

func (s *CainiaoCloudprintStdtemplatesGetStandardTemplateDo) SetStandardTemplateId(v int64) *CainiaoCloudprintStdtemplatesGetStandardTemplateDo {
	s.StandardTemplateId = &v
	return s
}
func (s *CainiaoCloudprintStdtemplatesGetStandardTemplateDo) SetStandardTemplateName(v string) *CainiaoCloudprintStdtemplatesGetStandardTemplateDo {
	s.StandardTemplateName = &v
	return s
}
func (s *CainiaoCloudprintStdtemplatesGetStandardTemplateDo) SetStandardTemplateUrl(v string) *CainiaoCloudprintStdtemplatesGetStandardTemplateDo {
	s.StandardTemplateUrl = &v
	return s
}
func (s *CainiaoCloudprintStdtemplatesGetStandardTemplateDo) SetStandardWaybillType(v int64) *CainiaoCloudprintStdtemplatesGetStandardTemplateDo {
	s.StandardWaybillType = &v
	return s
}
func (s *CainiaoCloudprintStdtemplatesGetStandardTemplateDo) SetBrandCode(v string) *CainiaoCloudprintStdtemplatesGetStandardTemplateDo {
	s.BrandCode = &v
	return s
}

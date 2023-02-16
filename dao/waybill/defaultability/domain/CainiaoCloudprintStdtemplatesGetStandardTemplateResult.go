package domain

type CainiaoCloudprintStdtemplatesGetStandardTemplateResult struct {
	/*
	   cp编码     */
	CpCode *string `json:"cp_code,omitempty" `

	/*
	   该cp的所有标准模板     */
	StandardTemplates *[]CainiaoCloudprintStdtemplatesGetStandardTemplateDo `json:"standard_templates,omitempty" `
}

func (s *CainiaoCloudprintStdtemplatesGetStandardTemplateResult) SetCpCode(v string) *CainiaoCloudprintStdtemplatesGetStandardTemplateResult {
	s.CpCode = &v
	return s
}
func (s *CainiaoCloudprintStdtemplatesGetStandardTemplateResult) SetStandardTemplates(v []CainiaoCloudprintStdtemplatesGetStandardTemplateDo) *CainiaoCloudprintStdtemplatesGetStandardTemplateResult {
	s.StandardTemplates = &v
	return s
}

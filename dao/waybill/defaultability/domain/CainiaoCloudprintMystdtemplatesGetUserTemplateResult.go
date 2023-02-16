package domain

type CainiaoCloudprintMystdtemplatesGetUserTemplateResult struct {
	/*
	   cp编码     */
	CpCode *string `json:"cp_code,omitempty" `

	/*
	   用户使用的模板数据     */
	UserStdTemplates *[]CainiaoCloudprintMystdtemplatesGetUserTemplateDo `json:"user_std_templates,omitempty" `
}

func (s *CainiaoCloudprintMystdtemplatesGetUserTemplateResult) SetCpCode(v string) *CainiaoCloudprintMystdtemplatesGetUserTemplateResult {
	s.CpCode = &v
	return s
}
func (s *CainiaoCloudprintMystdtemplatesGetUserTemplateResult) SetUserStdTemplates(v []CainiaoCloudprintMystdtemplatesGetUserTemplateDo) *CainiaoCloudprintMystdtemplatesGetUserTemplateResult {
	s.UserStdTemplates = &v
	return s
}

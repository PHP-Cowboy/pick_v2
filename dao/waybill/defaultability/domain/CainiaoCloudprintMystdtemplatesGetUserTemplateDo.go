package domain

type CainiaoCloudprintMystdtemplatesGetUserTemplateDo struct {
	/*
	   keys     */
	Keys *[]CainiaoCloudprintMystdtemplatesGetKeyResult `json:"keys,omitempty" `

	/*
	   用户使用模板的url     */
	UserStdTemplateUrl *string `json:"user_std_template_url,omitempty" `

	/*
	   用户使用模板的id     */
	UserStdTemplateId *int64 `json:"user_std_template_id,omitempty" `

	/*
	   用户使用模板名称     */
	UserStdTemplateName *string `json:"user_std_template_name,omitempty" `

	/*
	   品牌 code     */
	BrandCode *string `json:"brand_code,omitempty" `
}

func (s *CainiaoCloudprintMystdtemplatesGetUserTemplateDo) SetKeys(v []CainiaoCloudprintMystdtemplatesGetKeyResult) *CainiaoCloudprintMystdtemplatesGetUserTemplateDo {
	s.Keys = &v
	return s
}
func (s *CainiaoCloudprintMystdtemplatesGetUserTemplateDo) SetUserStdTemplateUrl(v string) *CainiaoCloudprintMystdtemplatesGetUserTemplateDo {
	s.UserStdTemplateUrl = &v
	return s
}
func (s *CainiaoCloudprintMystdtemplatesGetUserTemplateDo) SetUserStdTemplateId(v int64) *CainiaoCloudprintMystdtemplatesGetUserTemplateDo {
	s.UserStdTemplateId = &v
	return s
}
func (s *CainiaoCloudprintMystdtemplatesGetUserTemplateDo) SetUserStdTemplateName(v string) *CainiaoCloudprintMystdtemplatesGetUserTemplateDo {
	s.UserStdTemplateName = &v
	return s
}
func (s *CainiaoCloudprintMystdtemplatesGetUserTemplateDo) SetBrandCode(v string) *CainiaoCloudprintMystdtemplatesGetUserTemplateDo {
	s.BrandCode = &v
	return s
}

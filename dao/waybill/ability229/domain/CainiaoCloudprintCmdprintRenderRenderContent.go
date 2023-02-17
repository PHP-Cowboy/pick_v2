package domain

type CainiaoCloudprintCmdprintRenderRenderContent struct {
	/*
	   打印数据     */
	PrintData *string `json:"print_data,omitempty" `

	/*
	   模板url     */
	TemplateUrl *string `json:"template_url,omitempty" `

	/*
	   是否获取加密数据 defalutValue:false    */
	Encrypted *bool `json:"encrypted,omitempty" `

	/*
	   加密数据使用秘钥版本     */
	Ver *string `json:"ver,omitempty" `

	/*
	   数据签名     */
	Signature *string `json:"signature,omitempty" `

	/*
	   附加数据(用于修改数据)     */
	AddData *string `json:"add_data,omitempty" `
}

func (s *CainiaoCloudprintCmdprintRenderRenderContent) SetPrintData(v string) *CainiaoCloudprintCmdprintRenderRenderContent {
	s.PrintData = &v
	return s
}
func (s *CainiaoCloudprintCmdprintRenderRenderContent) SetTemplateUrl(v string) *CainiaoCloudprintCmdprintRenderRenderContent {
	s.TemplateUrl = &v
	return s
}
func (s *CainiaoCloudprintCmdprintRenderRenderContent) SetEncrypted(v bool) *CainiaoCloudprintCmdprintRenderRenderContent {
	s.Encrypted = &v
	return s
}
func (s *CainiaoCloudprintCmdprintRenderRenderContent) SetVer(v string) *CainiaoCloudprintCmdprintRenderRenderContent {
	s.Ver = &v
	return s
}
func (s *CainiaoCloudprintCmdprintRenderRenderContent) SetSignature(v string) *CainiaoCloudprintCmdprintRenderRenderContent {
	s.Signature = &v
	return s
}
func (s *CainiaoCloudprintCmdprintRenderRenderContent) SetAddData(v string) *CainiaoCloudprintCmdprintRenderRenderContent {
	s.AddData = &v
	return s
}

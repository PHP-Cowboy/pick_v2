package domain

type CainiaoWaybillCloudprintNetprintPrintPrintData struct {
	/*
	   版本     */
	Ver *string `json:"ver,omitempty" `

	/*
	   打印数据     */
	Data *string `json:"data,omitempty" `

	/*
	   是否加密     */
	Encrypted *bool `json:"encrypted,omitempty" `

	/*
	   签名     */
	Signature *string `json:"signature,omitempty" `

	/*
	   追加的 data     */
	AddData *string `json:"add_data,omitempty" `

	/*
	   模板 url     */
	TemplateUrl *string `json:"template_url,omitempty" `
}

func (s *CainiaoWaybillCloudprintNetprintPrintPrintData) SetVer(v string) *CainiaoWaybillCloudprintNetprintPrintPrintData {
	s.Ver = &v
	return s
}
func (s *CainiaoWaybillCloudprintNetprintPrintPrintData) SetData(v string) *CainiaoWaybillCloudprintNetprintPrintPrintData {
	s.Data = &v
	return s
}
func (s *CainiaoWaybillCloudprintNetprintPrintPrintData) SetEncrypted(v bool) *CainiaoWaybillCloudprintNetprintPrintPrintData {
	s.Encrypted = &v
	return s
}
func (s *CainiaoWaybillCloudprintNetprintPrintPrintData) SetSignature(v string) *CainiaoWaybillCloudprintNetprintPrintPrintData {
	s.Signature = &v
	return s
}
func (s *CainiaoWaybillCloudprintNetprintPrintPrintData) SetAddData(v string) *CainiaoWaybillCloudprintNetprintPrintPrintData {
	s.AddData = &v
	return s
}
func (s *CainiaoWaybillCloudprintNetprintPrintPrintData) SetTemplateUrl(v string) *CainiaoWaybillCloudprintNetprintPrintPrintData {
	s.TemplateUrl = &v
	return s
}

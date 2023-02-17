package domain

type CainiaoWaybillCloudprintNetprintPrintCustomData struct {
	/*
	   自定义区数据     */
	Data *string `json:"data,omitempty" `

	/*
	   自定义区链接     */
	TemplateUrl *string `json:"template_url,omitempty" `
}

func (s *CainiaoWaybillCloudprintNetprintPrintCustomData) SetData(v string) *CainiaoWaybillCloudprintNetprintPrintCustomData {
	s.Data = &v
	return s
}
func (s *CainiaoWaybillCloudprintNetprintPrintCustomData) SetTemplateUrl(v string) *CainiaoWaybillCloudprintNetprintPrintCustomData {
	s.TemplateUrl = &v
	return s
}

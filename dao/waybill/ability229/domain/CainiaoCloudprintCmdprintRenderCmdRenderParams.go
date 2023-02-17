package domain

type CainiaoCloudprintCmdprintRenderCmdRenderParams struct {
	/*
	   需要打印的文档，包括模板地址、打印数据     */
	Document *CainiaoCloudprintCmdprintRenderRenderDocument `json:"document,omitempty" `

	/*
	   打印机名称     */
	PrinterName *string `json:"printer_name,omitempty" `

	/*
	   客户端ID     */
	ClientId *string `json:"client_id,omitempty" `

	/*
	   客户端类型：weixin/alipay/native     */
	ClientType *string `json:"client_type,omitempty" `

	/*
	   打印配置     */
	Config *CainiaoCloudprintCmdprintRenderRenderConfig `json:"config,omitempty" `
}

func (s *CainiaoCloudprintCmdprintRenderCmdRenderParams) SetDocument(v CainiaoCloudprintCmdprintRenderRenderDocument) *CainiaoCloudprintCmdprintRenderCmdRenderParams {
	s.Document = &v
	return s
}
func (s *CainiaoCloudprintCmdprintRenderCmdRenderParams) SetPrinterName(v string) *CainiaoCloudprintCmdprintRenderCmdRenderParams {
	s.PrinterName = &v
	return s
}
func (s *CainiaoCloudprintCmdprintRenderCmdRenderParams) SetClientId(v string) *CainiaoCloudprintCmdprintRenderCmdRenderParams {
	s.ClientId = &v
	return s
}
func (s *CainiaoCloudprintCmdprintRenderCmdRenderParams) SetClientType(v string) *CainiaoCloudprintCmdprintRenderCmdRenderParams {
	s.ClientType = &v
	return s
}
func (s *CainiaoCloudprintCmdprintRenderCmdRenderParams) SetConfig(v CainiaoCloudprintCmdprintRenderRenderConfig) *CainiaoCloudprintCmdprintRenderCmdRenderParams {
	s.Config = &v
	return s
}

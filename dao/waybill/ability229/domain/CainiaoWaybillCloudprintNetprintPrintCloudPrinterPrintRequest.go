package domain

type CainiaoWaybillCloudprintNetprintPrintCloudPrinterPrintRequest struct {
	/*
	   共享码     */
	ShareCode *string `json:"share_code,omitempty" `

	/*
	   打印机 id     */
	Uid *string `json:"uid,omitempty" `

	/*
	   自定义内容     */
	CustomData *CainiaoWaybillCloudprintNetprintPrintCustomData `json:"custom_data,omitempty" `

	/*
	   打印数据     */
	PrintData *CainiaoWaybillCloudprintNetprintPrintPrintData `json:"print_data,omitempty" `
}

func (s *CainiaoWaybillCloudprintNetprintPrintCloudPrinterPrintRequest) SetShareCode(v string) *CainiaoWaybillCloudprintNetprintPrintCloudPrinterPrintRequest {
	s.ShareCode = &v
	return s
}
func (s *CainiaoWaybillCloudprintNetprintPrintCloudPrinterPrintRequest) SetUid(v string) *CainiaoWaybillCloudprintNetprintPrintCloudPrinterPrintRequest {
	s.Uid = &v
	return s
}
func (s *CainiaoWaybillCloudprintNetprintPrintCloudPrinterPrintRequest) SetCustomData(v CainiaoWaybillCloudprintNetprintPrintCustomData) *CainiaoWaybillCloudprintNetprintPrintCloudPrinterPrintRequest {
	s.CustomData = &v
	return s
}
func (s *CainiaoWaybillCloudprintNetprintPrintCloudPrinterPrintRequest) SetPrintData(v CainiaoWaybillCloudprintNetprintPrintPrintData) *CainiaoWaybillCloudprintNetprintPrintCloudPrinterPrintRequest {
	s.PrintData = &v
	return s
}

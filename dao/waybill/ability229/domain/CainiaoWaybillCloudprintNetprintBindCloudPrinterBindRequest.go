package domain

type CainiaoWaybillCloudprintNetprintBindCloudPrinterBindRequest struct {
	/*
	   打印机 mac 地址     */
	Uid *string `json:"uid,omitempty" `

	/*
	   验证码     */
	VerifyCode *string `json:"verify_code,omitempty" `
}

func (s *CainiaoWaybillCloudprintNetprintBindCloudPrinterBindRequest) SetUid(v string) *CainiaoWaybillCloudprintNetprintBindCloudPrinterBindRequest {
	s.Uid = &v
	return s
}
func (s *CainiaoWaybillCloudprintNetprintBindCloudPrinterBindRequest) SetVerifyCode(v string) *CainiaoWaybillCloudprintNetprintBindCloudPrinterBindRequest {
	s.VerifyCode = &v
	return s
}

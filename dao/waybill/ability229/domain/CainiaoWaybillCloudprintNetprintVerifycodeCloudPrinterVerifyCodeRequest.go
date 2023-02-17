package domain

type CainiaoWaybillCloudprintNetprintVerifycodeCloudPrinterVerifyCodeRequest struct {
	/*
	   打印机 id     */
	Uid *string `json:"uid,omitempty" `
}

func (s *CainiaoWaybillCloudprintNetprintVerifycodeCloudPrinterVerifyCodeRequest) SetUid(v string) *CainiaoWaybillCloudprintNetprintVerifycodeCloudPrinterVerifyCodeRequest {
	s.Uid = &v
	return s
}

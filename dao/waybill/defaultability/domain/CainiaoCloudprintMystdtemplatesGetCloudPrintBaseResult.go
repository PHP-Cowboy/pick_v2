package domain

type CainiaoCloudprintMystdtemplatesGetCloudPrintBaseResult struct {
	/*
	   所有cp的数据     */
	Datas *[]CainiaoCloudprintMystdtemplatesGetUserTemplateResult `json:"datas,omitempty" `

	/*
	   错误码     */
	ErrorCode *string `json:"error_code,omitempty" `

	/*
	   错误消息     */
	ErrorMessage *string `json:"error_message,omitempty" `

	/*
	   是否成功     */
	Success *bool `json:"success,omitempty" `
}

func (s *CainiaoCloudprintMystdtemplatesGetCloudPrintBaseResult) SetDatas(v []CainiaoCloudprintMystdtemplatesGetUserTemplateResult) *CainiaoCloudprintMystdtemplatesGetCloudPrintBaseResult {
	s.Datas = &v
	return s
}
func (s *CainiaoCloudprintMystdtemplatesGetCloudPrintBaseResult) SetErrorCode(v string) *CainiaoCloudprintMystdtemplatesGetCloudPrintBaseResult {
	s.ErrorCode = &v
	return s
}
func (s *CainiaoCloudprintMystdtemplatesGetCloudPrintBaseResult) SetErrorMessage(v string) *CainiaoCloudprintMystdtemplatesGetCloudPrintBaseResult {
	s.ErrorMessage = &v
	return s
}
func (s *CainiaoCloudprintMystdtemplatesGetCloudPrintBaseResult) SetSuccess(v bool) *CainiaoCloudprintMystdtemplatesGetCloudPrintBaseResult {
	s.Success = &v
	return s
}
